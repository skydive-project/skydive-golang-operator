/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"github.com/pkg/errors"
	"skydive/pkg/config"
	"skydive/pkg/kclient"
	"skydive/pkg/manifests"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	skydivev1 "skydive/api/v1"
)

var (
	PrometheusConnectorDeployment = "prometheus-connector/deployment.yaml"
	PrometheusConnectorService    = "prometheus-connector/service.yaml"
	PrometheusConnectorRoute      = "prometheus-connector/route.yaml"
	PrometheusConfigMap           = "prometheus/config-map.yaml"
	PrometheusDeployment          = "prometheus/deployment.yaml"
	PrometheusService             = "prometheus/service.yaml"
	PrometheusRoute               = "prometheus/route.yaml"
)

// PrometheusConnectorReconciler reconciles a PrometheusConnector object
type PrometheusConnectorReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=skydive.example.com,resources=prometheusconnectors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=skydive.example.com,resources=prometheusconnectors/status,verbs=get;update;patch

func (r *PrometheusConnectorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("prometheusconnector", req.NamespacedName)

	prometheus_connector := skydivev1.PrometheusConnector{}
	if err := r.Client.Get(ctx, req.NamespacedName, &prometheus_connector); err != nil {
		log.Error(err, "failed to get Skydive resource")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// getting assets and configuration
	assets := manifests.NewAssets("assets")
	config_instance, err := config.GetConfig()
	if err != nil {
		log.Error(err, "Configuration build has failed")
		return ctrl.Result{}, err
	}

	kclient_instance, err := kclient.New(config_instance, "", prometheus_connector.Spec.Namespace, "")
	if err != nil {
		log.Error(err, "Kubernets client build failed")
		return ctrl.Result{}, err
	}

	config_map, err := kclient.NewConfigMap(assets.MustNewAssetReader(PrometheusConfigMap))
	if err != nil {
		log.Error(err, "initializing Prometheus ConfigMap failed")
	}
	config_map.Namespace = prometheus_connector.Spec.Namespace

	err = kclient_instance.CreateOrUpdateConfigMap(config_map)
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "Prometheus ConfigMap creation has failed")
	}

	dep, err := kclient.NewDeployment(assets.MustNewAssetReader(PrometheusDeployment))
	if err != nil {
		log.Error(err, "initializing Prometheus Deployment failed")
		return ctrl.Result{}, errors.Wrap(err, "initializing Prometheus Deployment failed")
	}
	dep.Namespace = prometheus_connector.Spec.Namespace
	err = kclient_instance.CreateOrUpdateDeployment(dep)
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "Prometheus deployment creation has failed")
	}

	svc, err := kclient.NewService(assets.MustNewAssetReader(PrometheusService))
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "initializing Prometheus Service failed")
	}
	svc.Namespace = prometheus_connector.Spec.Namespace

	err = kclient_instance.CreateOrUpdateService(svc)
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "Prometheus Service creation has failed")
	}

	route, err := kclient.NewRoute(assets.MustNewAssetReader(PrometheusRoute))
	if err != nil {
		log.Error(err, "initializing Prometheus Route failed")
	}
	route.Namespace = prometheus_connector.Spec.Namespace

	err = kclient_instance.CreateRouteIfNotExists(route)
	if err != nil {
		log.Error(err, "Prometheus Route creation has failed")
	}

	_, err = kclient_instance.WaitForRouteReady(route)
	if err != nil {
		log.Error(err, "waiting for Prometheus Route to become ready failed")
	}

	connector_dep, err := kclient.NewDeployment(assets.MustNewAssetReader(PrometheusConnectorDeployment))
	if err != nil {
		log.Error(err, "initializing Prometheus Connector Deployment failed")
		return ctrl.Result{}, errors.Wrap(err, "initializing Prometheus Deployment failed")
	}
	connector_dep.Namespace = prometheus_connector.Spec.Namespace
	for index := range connector_dep.Spec.Template.Spec.Containers {
		if connector_dep.Spec.Template.Spec.Containers[index].Name == "skydive-prometheus-connector" {
			connector_dep.Spec.Template.Spec.Containers[index].Env = prometheus_connector.Spec.Deployment.Env
			break
		}
	}

	err = kclient_instance.CreateOrUpdateDeployment(connector_dep)
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "Prometheus Connector deployment creation has failed")
	}

	connector_svc, err := kclient.NewService(assets.MustNewAssetReader(PrometheusConnectorService))
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "initializing Prometheus Service failed")
	}
	connector_svc.Namespace = prometheus_connector.Spec.Namespace

	err = kclient_instance.CreateOrUpdateService(connector_svc)
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "Prometheus Service creation has failed")
	}

	connector_route, err := kclient.NewRoute(assets.MustNewAssetReader(PrometheusConnectorRoute))
	if err != nil {
		log.Error(err, "initializing Prometheus Route failed")
	}
	connector_route.Namespace = prometheus_connector.Spec.Namespace

	err = kclient_instance.CreateRouteIfNotExists(connector_route)
	if err != nil {
		log.Error(err, "Prometheus Route creation has failed")
	}

	_, err = kclient_instance.WaitForRouteReady(connector_route)
	if err != nil {
		log.Error(err, "waiting for Prometheus Route to become ready failed")
	}

	return ctrl.Result{}, nil
}

func (r *PrometheusConnectorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&skydivev1.PrometheusConnector{}).
		Complete(r)
}
