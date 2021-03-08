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
	SkydiveAgentsDaemonSet    = "skydive-agents/daemon-set.yaml"
	SkydiveAnalyzerDeployment = "skydive-analyzer/deployment.yaml"
	SkydiveAnalyzerRoute      = "skydive-analyzer/route.yaml"
	SkydiveAnalyzerService    = "skydive-analyzer/service.yaml"
	SkydiveFlowExporter       = "flow-exporter/deployment.yaml"
)

// SkydiveSuiteReconciler reconciles a SkydiveSuite object
type SkydiveSuiteReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=skydive.example.com,resources=skydivesuites,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=skydive.example.com,resources=skydivesuites/status,verbs=get;update;patch

func (r *SkydiveSuiteReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("skydivesuite", req.NamespacedName)

	skydive_suite := skydivev1.SkydiveSuite{}
	if err := r.Client.Get(ctx, req.NamespacedName, &skydive_suite); err != nil {
		log.Error(err, "failed to get SkydiveSuite resource")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// getting assets and configuration
	assets := manifests.NewAssets("assets")
	config_instance, err := config.GetConfig()
	if err != nil {
		log.Error(err, "Configuration build has failed")
		return ctrl.Result{}, err
	}

	kclient_instance, err := kclient.New(config_instance, "", skydive_suite.Spec.Namespace, "")
	if err != nil {
		log.Error(err, "Kubernets client build failed")
		return ctrl.Result{}, err
	}

	// Creating skydive Analyzers
	if skydive_suite.Spec.Enable.Analyzer {

		dep, err := kclient.NewDeployment(assets.MustNewAssetReader(SkydiveAnalyzerDeployment))
		if err != nil {
			log.Error(err, "initializing skydive-analyzer Deployment failed")
			return ctrl.Result{}, errors.Wrap(err, "initializing skydive-analyzer Deployment failed")
		}
		dep.Namespace = skydive_suite.Spec.Namespace

		err = kclient_instance.CreateOrUpdateDeployment(dep)
		if err != nil {
			return ctrl.Result{}, errors.Wrap(err, "Skydive analyzer deployment creation has failed")
		}

		svc, err := kclient.NewService(assets.MustNewAssetReader(SkydiveAnalyzerService))
		if err != nil {
			return ctrl.Result{}, errors.Wrap(err, "initializing skydive-analyzer Service failed")
		}
		svc.Namespace = skydive_suite.Spec.Namespace

		err = kclient_instance.CreateOrUpdateService(svc)
		if err != nil {
			return ctrl.Result{}, errors.Wrap(err, "Service creation has failed")
		}

		if skydive_suite.Spec.Enable.Route {
			route, err := kclient.NewRoute(assets.MustNewAssetReader(SkydiveAnalyzerRoute))
			if err != nil {
				log.Error(err, "initializing skydive-analyzer Route failed")
			}
			route.Namespace = skydive_suite.Spec.Namespace

			err = kclient_instance.CreateRouteIfNotExists(route)
			if err != nil {
				log.Error(err, "Route creation has failed")
			}

			_, err = kclient_instance.WaitForRouteReady(route)
			if err != nil {
				log.Error(err, "waiting for Skydive Route to become ready failed")
			}
		}

	}

	// Creating skydive Agents
	if skydive_suite.Spec.Enable.Agents {
		ds, err := kclient.NewDaemonSet(assets.MustNewAssetReader(SkydiveAgentsDaemonSet))
		if err != nil {
			return ctrl.Result{}, errors.Wrap(err, "initializing skydive-agents DaemonSet failed")
		}
		ds.Namespace = skydive_suite.Spec.Namespace
		for _, container := range ds.Spec.Template.Spec.Containers {
			for _, env := range container.Env {
				switch env.Name {
				case "SKYDIVE_LOGGING_LEVEL":
					env.Value = skydive_suite.Spec.Logging.Level
				}
			}
		}
		err = kclient_instance.CreateOrUpdateDaemonSet(ds)
		if err != nil {
			return ctrl.Result{}, errors.Wrap(err, "DaemonSet creation failed")
		}
	}

	if skydive_suite.Spec.Enable.FlowExporter {
		dep, err := kclient.NewDeployment(assets.MustNewAssetReader(SkydiveFlowExporter))
		if err != nil {
			log.Error(err, "initializing skydive flow-exporter Deployment failed")
			return ctrl.Result{}, errors.Wrap(err, "initializing skydive flow-exporter Deployment failed")
		}
		dep.Namespace = skydive_suite.Spec.Namespace
		for _, container := range dep.Spec.Template.Spec.Containers {
			for _, env := range container.Env {
				switch env.Name {
				case "SKYDIVE_LOGGING_LEVEL":
					env.Value = skydive_suite.Spec.Logging.Level
				}
			}
		}

		err = kclient_instance.CreateOrUpdateDeployment(dep)
		if err != nil {
			return ctrl.Result{}, errors.Wrap(err, "Flow exporter deployment creation has failed")
		}
	}

	return ctrl.Result{}, nil
}

func (r *SkydiveSuiteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&skydivev1.SkydiveSuite{}).
		Complete(r)
}
