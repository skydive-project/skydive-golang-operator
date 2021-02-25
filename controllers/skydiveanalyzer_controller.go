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
	"k8s.io/client-go/tools/clientcmd"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"skydive/pkg/kclient"
	"skydive/pkg/manifests"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	skydivev1beta1 "skydive/api/v1beta1"
)

const ( //TODO move to a diffrent file
	namespaceName  = "skydive"
	kubeconfigPath = "/Users/orannahoum/.kube/config"
)

var (
	SkydiveAnalyzerDeployment = "skydive-analyzer/deployment.yaml"
	SkydiveAnalyzerRoute      = "skydive-analyzer/route.yaml"
	SkydiveAnalyzerRouteUI    = "skydive-analyzer/route-ui.yaml"
	SkydiveAnalyzerService    = "skydive-analyzer/service.yaml"
	SkydiveAnalyzerServiceUI  = "skydive-analyzer/service-ui.yaml"
)

// SkydiveAnalyzerReconciler reconciles a SkydiveAnalyzer object
type SkydiveAnalyzerReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme

	kclient.KClient
}

// +kubebuilder:rbac:groups=skydive.example.com,resources=skydiveanalyzers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=skydive.example.com,resources=skydiveanalyzers/status,verbs=get;update;patch
func (r *SkydiveAnalyzerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("skydiveanalyzer", req.NamespacedName)

	log.Info("initializing assets for skydive-analyzer")
	assets := manifests.NewAssets("assets")

	log.Info("Building configuration")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		log.Error(err, "Configuration build has failed")
		return ctrl.Result{}, err
	}

	log.Info("Creating kubernetes client")
	kclient_instance, err := kclient.New(config, "", namespaceName, "")
	if err != nil {
		log.Error(err, "Kubernets client build failed")
		return ctrl.Result{}, err
	}

	log.Info("initializing skydive-analyzer Deployment...")
	dep, err := kclient.NewDeployment(assets.MustNewAssetReader(SkydiveAnalyzerDeployment))
	if err != nil {
		log.Error(err, "initializing skydive-analyzer Deployment failed")
		return ctrl.Result{}, errors.Wrap(err, "initializing skydive-analyzer Deployment failed")
	}
	dep.Namespace = namespaceName

	log.Info("initializing skydive-analyzer Route...")
	route, err := kclient.NewRoute(assets.MustNewAssetReader(SkydiveAnalyzerRoute))
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "initializing skydive-analyzer Route failed")
	}
	route.Namespace = namespaceName

	log.Info("initializing skydive-analyzer-UI Route...")
	route_ui, err := kclient.NewRoute(assets.MustNewAssetReader(SkydiveAnalyzerRouteUI))
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "initializing skydive-analyzer-UI Route failed")
	}
	route_ui.Namespace = namespaceName

	log.Info("initializing skydive-analyzer Service...")
	svc, err := kclient.NewService(assets.MustNewAssetReader(SkydiveAnalyzerService))
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "initializing skydive-analyzer Service failed")
	}
	svc.Namespace = namespaceName

	log.Info("initializing skydive-analyzer-UI Service...")
	svc_ui, err := kclient.NewService(assets.MustNewAssetReader(SkydiveAnalyzerServiceUI))
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "initializing skydive-analyzer-UI Service failed")
	}
	svc_ui.Namespace = namespaceName

	log.Info("Creating deployment...")
	err = kclient_instance.CreateOrUpdateDeployment(dep)
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "Deployment creation has failed")
	}

	//log.Info("Creating Service...")
	//err = kclient_instance.CreateOrUpdateService(svc)
	//if err != nil {
	//	return ctrl.Result{}, errors.Wrap(err, "Service creation has failed")
	//}
	//
	//log.Info("Creating UI-Service...")
	//err = kclient_instance.CreateOrUpdateService(svc_ui)
	//if err != nil {
	//	return ctrl.Result{}, errors.Wrap(err, "UI-Service creation has failed")
	//}

	log.Info("Creating route...")
	err = kclient_instance.CreateRouteIfNotExists(route)
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "Route creation has failed")
	}

	_, err = kclient_instance.WaitForRouteReady(route)
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "waiting for Skydive Route to become ready failed")
	}

	log.Info("Creating UI-route...")
	err = kclient_instance.CreateRouteIfNotExists(route_ui)
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "UI-Route creation has failed")
	}

	_, err = kclient_instance.WaitForRouteReady(route_ui)
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "waiting for Skydive-UI Route to become ready failed")
	}

	log.Info("Reconcile finished all tasks successfully")

	return ctrl.Result{}, nil

}

func (r *SkydiveAnalyzerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&skydivev1beta1.SkydiveAnalyzer{}).
		Complete(r)
}
