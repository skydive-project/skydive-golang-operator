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
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	skydivev1 "skydive/api/v1"
	"skydive/pkg/config"
	"skydive/pkg/kclient"
	"skydive/pkg/manifests"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	SkydiveFlowExporterDev = "flow-exporter/deployment.yaml"
)

// SkydiveFlowExporterReconciler reconciles a SkydiveFlowExporter object
type SkydiveFlowExporterReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=skydive.example.com,resources=skydiveflowexporters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=skydive.example.com,resources=skydiveflowexporters/status,verbs=get;update;patch

func (r *SkydiveFlowExporterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("skydiveflowexporter", req.NamespacedName)

	skydive_flow_exporter := skydivev1.SkydiveFlowExporter{}
	if err := r.Client.Get(ctx, req.NamespacedName, &skydive_flow_exporter); err != nil {
		log.Error(err, "failed to get Skydive-Flow-Exporter resource")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// getting assets and configuration either from hack folder (for dev) or from normal assets folder
	assetDir := "assets"
	if skydive_flow_exporter.Spec.DeployDevEnv {
		assetDir = "hack/assets"
	}

	assets := manifests.NewAssets(assetDir)
	config_instance, err := config.GetConfig()
	if err != nil {
		log.Error(err, "Configuration build has failed")
		return ctrl.Result{}, err
	}

	kclient_instance, err := kclient.New(config_instance, "", skydive_flow_exporter.Spec.Namespace, "")
	if err != nil {
		log.Error(err, "Kubernets client build failed")
		return ctrl.Result{}, err
	}

	dep, err := kclient.NewDeployment(assets.MustNewAssetReader(SkydiveFlowExporterDev))
	if err != nil {
		log.Error(err, "initializing skydive flow-exporter Deployment failed")
		return ctrl.Result{}, errors.Wrap(err, "initializing skydive flow-exporter Deployment failed")
	}
	dep.Namespace = skydive_flow_exporter.Spec.Namespace
	for index := range dep.Spec.Template.Spec.Containers {
		if dep.Spec.Template.Spec.Containers[index].Name == "skydive-flow-exporter" {
			dep.Spec.Template.Spec.Containers[index].Env = skydive_flow_exporter.Spec.Deployment.Env
			break
		}
	}

	err = kclient_instance.CreateOrUpdateDeployment(dep)
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "Flow exporter deployment creation has failed")
	}

	return ctrl.Result{}, nil
}

func (r *SkydiveFlowExporterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&skydivev1.SkydiveFlowExporter{}).
		Complete(r)
}
