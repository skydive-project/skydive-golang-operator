// Copyright 2018 The Cluster Monitoring Operator Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kclient

import (
	"context"
	"github.com/imdario/mergo"
	routev1 "github.com/openshift/api/route/v1"
	"io"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/rest"
	"reflect"
	"strings"
	"time"

	openshiftrouteclientset "github.com/openshift/client-go/route/clientset/versioned"

	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	deploymentCreateTimeout = 5 * time.Minute
	metadataPrefix          = "monitoring.openshift.io/"
)

type KClient struct {
	version               string
	namespace             string
	userWorkloadNamespace string
	namespaceSelector     string
	kclient               kubernetes.Interface
	osrclient             openshiftrouteclientset.Interface
	eclient               apiextensionsclient.Interface
}

func New(cfg *rest.Config, version string, namespace, namespaceSelector string) (*KClient, error) {
	kclient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "creating kubernetes clientset kclient")
	}

	eclient, err := apiextensionsclient.NewForConfig(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "creating apiextensions kclient")
	}

	osrclient, err := openshiftrouteclientset.NewForConfig(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "creating openshift route client")
	}

	// SCC moved to CRD and CRD does not handle protobuf. Force the SCC kclient to use JSON instead.
	jsonClientConfig := rest.CopyConfig(cfg)
	jsonClientConfig.ContentConfig.AcceptContentTypes = "application/json"
	jsonClientConfig.ContentConfig.ContentType = "application/json"

	return &KClient{
		version:           version,
		namespace:         namespace,
		namespaceSelector: namespaceSelector,
		kclient:           kclient,
		osrclient:         osrclient,
		eclient:           eclient,
	}, nil
}

func (c *KClient) CreateOrUpdateDeployment(dep *appsv1.Deployment) error {
	existing, err := c.kclient.AppsV1().Deployments(dep.GetNamespace()).Get(context.TODO(), dep.GetName(), metav1.GetOptions{})

	if apierrors.IsNotFound(err) {
		err = c.CreateDeployment(dep)
		return errors.Wrap(err, "creating Deployment object failed")
	}
	if err != nil {
		return errors.Wrap(err, "retrieving Deployment object failed")
	}
	if reflect.DeepEqual(dep.Spec, existing.Spec) {
		// Nothing to do, as the currently existing deployment is equivalent to the one that would be applied.
		return nil
	}

	required := dep.DeepCopy()
	mergeMetadata(&required.ObjectMeta, existing.ObjectMeta)

	err = c.UpdateDeployment(required)
	if err != nil {
		uErr, ok := err.(*apierrors.StatusError)
		if ok && uErr.ErrStatus.Code == 422 && uErr.ErrStatus.Reason == metav1.StatusReasonInvalid {
			// try to delete Deployment
			err = c.DeleteDeployment(existing)
			if err != nil {
				return errors.Wrap(err, "deleting Deployment object failed")
			}
			err = c.CreateDeployment(required)
			if err != nil {
				return errors.Wrap(err, "creating Deployment object failed after update failed")
			}
		}
		return errors.Wrap(err, "updating Deployment object failed")
	}
	return nil
}

func (c *KClient) WaitForDeploymentRollout(dep *appsv1.Deployment) error {
	var lastErr error
	if err := wait.Poll(time.Second, deploymentCreateTimeout, func() (bool, error) {
		d, err := c.kclient.AppsV1().Deployments(dep.GetNamespace()).Get(context.TODO(), dep.GetName(), metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		if d.Generation > d.Status.ObservedGeneration {
			lastErr = errors.Errorf("current generation %d, observed generation %d",
				d.Generation, d.Status.ObservedGeneration)
			return false, nil
		}
		if d.Status.UpdatedReplicas != d.Status.Replicas {
			lastErr = errors.Errorf("expected %d replicas, got %d updated replicas",
				d.Status.Replicas, d.Status.UpdatedReplicas)
			return false, nil
		}
		if d.Status.UnavailableReplicas != 0 {
			lastErr = errors.Errorf("got %d unavailable replicas",
				d.Status.UnavailableReplicas)
			return false, nil
		}
		return true, nil
	}); err != nil {
		if err == wait.ErrWaitTimeout && lastErr != nil {
			err = lastErr
		}
		return errors.Wrapf(err, "waiting for DeploymentRollout of %s/%s", dep.GetNamespace(), dep.GetName())
	}
	return nil
}

// mergeMetadata merges labels and annotations from `existing` map into `required` one where `required` has precedence
// over `existing` keys and values. Additionally function performs filtering of labels and annotations from `exiting` map
// where keys starting from string defined in `metadataPrefix` are deleted. This prevents issues with preserving stale
// metadata defined by the operator
func mergeMetadata(required *metav1.ObjectMeta, existing metav1.ObjectMeta) {
	for k := range existing.Annotations {
		if strings.HasPrefix(k, metadataPrefix) {
			delete(existing.Annotations, k)
		}
	}

	for k := range existing.Labels {
		if strings.HasPrefix(k, metadataPrefix) {
			delete(existing.Labels, k)
		}
	}

	mergo.Merge(&required.Annotations, existing.Annotations)
	mergo.Merge(&required.Labels, existing.Labels)
}

func NewDeployment(manifest io.Reader) (*appsv1.Deployment, error) {
	d := appsv1.Deployment{}
	err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&d)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

func (c *KClient) CreateDeployment(dep *appsv1.Deployment) error {
	d, err := c.kclient.AppsV1().Deployments(dep.GetNamespace()).Create(context.TODO(), dep, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return c.WaitForDeploymentRollout(d)
}

func (c *KClient) UpdateDeployment(dep *appsv1.Deployment) error {
	updated, err := c.kclient.AppsV1().Deployments(dep.GetNamespace()).Update(context.TODO(), dep, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return c.WaitForDeploymentRollout(updated)
}

func (c *KClient) DeleteDeployment(d *appsv1.Deployment) error {
	p := metav1.DeletePropagationForeground
	err := c.kclient.AppsV1().Deployments(d.GetNamespace()).Delete(context.TODO(), d.GetName(), metav1.DeleteOptions{PropagationPolicy: &p})
	if apierrors.IsNotFound(err) {
		return nil
	}

	return err
}

func NewRoute(manifest io.Reader) (*routev1.Route, error) {
	route := routev1.Route{}
	err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&route)
	if err != nil {
		return nil, err
	}

	return &route, nil
}

func (c *KClient) CreateRouteIfNotExists(r *routev1.Route) error {
	rclient := c.osrclient.RouteV1().Routes(r.GetNamespace())
	_, err := rclient.Get(context.TODO(), r.GetName(), metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		_, err := rclient.Create(context.TODO(), r, metav1.CreateOptions{})
		return errors.Wrap(err, "creating Route object failed")
	}
	return nil
}

func (c *KClient) WaitForRouteReady(r *routev1.Route) (string, error) {
	host := ""
	var lastErr error
	if err := wait.Poll(time.Second, deploymentCreateTimeout, func() (bool, error) {
		newRoute, err := c.osrclient.RouteV1().Routes(r.GetNamespace()).Get(context.TODO(), r.GetName(), metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		if len(newRoute.Status.Ingress) == 0 {
			lastErr = errors.New("no status available")
			return false, nil
		}
		for _, c := range newRoute.Status.Ingress[0].Conditions {
			if c.Type == "Admitted" && c.Status == "True" {
				host = newRoute.Spec.Host
				return true, nil
			}
		}
		lastErr = errors.New("route is not yet Admitted")
		return false, nil
	}); err != nil {
		if err == wait.ErrWaitTimeout && lastErr != nil {
			err = lastErr
		}
		return host, errors.Wrapf(err, "waiting for route %s/%s", r.GetNamespace(), r.GetName())
	}
	return host, nil
}

func NewService(manifest io.Reader) (*v1.Service, error) {
	s := v1.Service{}
	err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (c *KClient) CreateOrUpdateService(svc *v1.Service) error {
	sclient := c.kclient.CoreV1().Services(svc.GetNamespace())
	existing, err := sclient.Get(context.TODO(), svc.GetName(), metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		_, err = sclient.Create(context.TODO(), svc, metav1.CreateOptions{})
		return errors.Wrap(err, "creating Service object failed")
	}
	if err != nil {
		return errors.Wrap(err, "retrieving Service object failed")
	}

	required := svc.DeepCopy()
	required.ResourceVersion = existing.ResourceVersion
	if required.Spec.Type == v1.ServiceTypeClusterIP {
		required.Spec.ClusterIP = existing.Spec.ClusterIP
	}

	if reflect.DeepEqual(required.Spec, existing.Spec) {
		return nil
	}

	mergeMetadata(&required.ObjectMeta, existing.ObjectMeta)

	_, err = sclient.Update(context.TODO(), required, metav1.UpdateOptions{})
	return errors.Wrap(err, "updating Service object failed")
}

func NewDaemonSet(manifest io.Reader) (*appsv1.DaemonSet, error) {
	ds := appsv1.DaemonSet{}
	err := yaml.NewYAMLOrJSONDecoder(manifest, 100).Decode(&ds)
	if err != nil {
		return nil, err
	}

	return &ds, nil
}

func (c *KClient) WaitForDaemonSetRollout(ds *appsv1.DaemonSet) error {
	var lastErr error
	if err := wait.Poll(time.Second, deploymentCreateTimeout, func() (bool, error) {
		d, err := c.kclient.AppsV1().DaemonSets(ds.GetNamespace()).Get(context.TODO(), ds.GetName(), metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		if d.Generation > d.Status.ObservedGeneration {
			lastErr = errors.Errorf("current generation %d, observed generation: %d",
				d.Generation, d.Status.ObservedGeneration)
			return false, nil
		}
		if d.Status.UpdatedNumberScheduled != d.Status.DesiredNumberScheduled {
			lastErr = errors.Errorf("expected %d desired scheduled nodes, got %d updated scheduled nodes",
				d.Status.DesiredNumberScheduled, d.Status.UpdatedNumberScheduled)
			return false, nil
		}
		if d.Status.NumberUnavailable != 0 {
			lastErr = errors.Errorf("got %d unavailable nodes",
				d.Status.NumberUnavailable)
			return false, nil
		}
		return true, nil
	}); err != nil {
		if err == wait.ErrWaitTimeout && lastErr != nil {
			err = lastErr
		}
		return errors.Wrapf(err, "waiting for DaemonSetRollout of %s/%s", ds.GetNamespace(), ds.GetName())
	}
	return nil
}

func (c *KClient) CreateDaemonSet(ds *appsv1.DaemonSet) error {
	d, err := c.kclient.AppsV1().DaemonSets(ds.GetNamespace()).Create(context.TODO(), ds, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return c.WaitForDaemonSetRollout(d)
}

func (c *KClient) UpdateDaemonSet(ds *appsv1.DaemonSet) error {
	updated, err := c.kclient.AppsV1().DaemonSets(ds.GetNamespace()).Update(context.TODO(), ds, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return c.WaitForDaemonSetRollout(updated)
}

func (c *KClient) DeleteDaemonSet(d *appsv1.DaemonSet) error {
	orphanDependents := false
	err := c.kclient.AppsV1().DaemonSets(d.GetNamespace()).Delete(context.TODO(), d.GetName(), metav1.DeleteOptions{OrphanDependents: &orphanDependents})
	if apierrors.IsNotFound(err) {
		return nil
	}

	return err
}

func (c *KClient) CreateOrUpdateDaemonSet(ds *appsv1.DaemonSet) error {
	existing, err := c.kclient.AppsV1().DaemonSets(ds.GetNamespace()).Get(context.TODO(), ds.GetName(), metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		err = c.CreateDaemonSet(ds)
		return errors.Wrap(err, "creating DaemonSet object failed")
	}
	if err != nil {
		return errors.Wrap(err, "retrieving DaemonSet object failed")
	}

	required := ds.DeepCopy()
	mergeMetadata(&required.ObjectMeta, existing.ObjectMeta)

	err = c.UpdateDaemonSet(required)
	if err != nil {
		uErr, ok := err.(*apierrors.StatusError)
		if ok && uErr.ErrStatus.Code == 422 && uErr.ErrStatus.Reason == metav1.StatusReasonInvalid {
			// try to delete DaemonSet
			err = c.DeleteDaemonSet(existing)
			if err != nil {
				return errors.Wrap(err, "deleting DaemonSet object failed")
			}
			err = c.CreateDaemonSet(required)
			if err != nil {
				return errors.Wrap(err, "creating DaemonSet object failed after update failed")
			}
		}
		return errors.Wrap(err, "updating DaemonSet object failed")
	}
	return nil
}
