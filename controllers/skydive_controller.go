/*
Copyright 2021.

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
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"skydive-operator/pkg/kclient"
	"skydive-operator/pkg/manifests"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	skydivegroupv1 "skydive-operator/api/v1"
)

var (
	SkydiveAgentsDaemonSet    = "skydive-agents/daemon-set.yaml"
	SkydiveAnalyzerDeployment = "skydive-analyzer/deployment.yaml"
	SkydiveAnalyzerRoute      = "skydive-analyzer/route.yaml"
	SkydiveAnalyzerService    = "skydive-analyzer/service.yaml"

	SkydiveFlowExporterDep = "flow-exporter/deployment.yaml"

	PrometheusConnectorDeployment = "prometheus-connector/deployment.yaml"
	PrometheusConnectorService    = "prometheus-connector/service.yaml"
	PrometheusConnectorRoute      = "prometheus-connector/route.yaml"
	PrometheusConfigMap           = "prometheus/config-map.yaml"
	PrometheusDeployment          = "prometheus/deployment.yaml"
	PrometheusService             = "prometheus/service.yaml"
	PrometheusRoute               = "prometheus/route.yaml"
)

// SkydiveReconciler reconciles a Skydive object
type SkydiveReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=skydive-group.example.com,resources=skydives,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=skydive-group.example.com,resources=skydives/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=skydive-group.example.com,resources=skydives/finalizers,verbs=update

func (r *SkydiveReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("skydive", req.NamespacedName)

	log.Info("Starting Reconciler")

	skydive_suite := skydivegroupv1.Skydive{}
	if err := r.Client.Get(ctx, req.NamespacedName, &skydive_suite); err != nil {
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

	kclient_instance, err := kclient.New(config_instance, "", skydive_suite.Spec.Namespace, "")
	if err != nil {
		log.Error(err, "Kubernets client build failed")
		return ctrl.Result{}, err
	}

	err = kclient_instance.InitializeNamespace("skydive")
	if err != nil {
		log.Error(err, "Namespace Initialization has failed")
		return ctrl.Result{}, err
	}

	dep, err := kclient.NewDeployment(assets.MustNewAssetReader(SkydiveAnalyzerDeployment))
	if err != nil {
		log.Error(err, "initializing skydive-analyzer Deployment failed")
		return ctrl.Result{}, errors.Wrap(err, "initializing skydive-analyzer Deployment failed")
	}
	dep.Namespace = skydive_suite.Spec.Namespace
	for index := range dep.Spec.Template.Spec.Containers {
		if dep.Spec.Template.Spec.Containers[index].Name == "skydive-analyzer" {
			dep.Spec.Template.Spec.Containers[index].Env = skydive_suite.Spec.Analyzer.Deployment.Env
			break
		}
	}

	svc, err := kclient.NewService(assets.MustNewAssetReader(SkydiveAnalyzerService))
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "initializing skydive-analyzer Service failed")
	}
	svc.Namespace = skydive_suite.Spec.Namespace

	route, err := kclient.NewRoute(assets.MustNewAssetReader(SkydiveAnalyzerRoute))
	if err != nil {
		log.Error(err, "initializing skydive-analyzer Route failed")
	}
	route.Namespace = skydive_suite.Spec.Namespace

	// Creating skydive Analyzers
	if skydive_suite.Spec.Enable.Analyzer {

		log.Info("Starting Skydive Analyzers")
		err = kclient_instance.CreateOrUpdateDeployment(dep)
		if err != nil {
			return ctrl.Result{}, errors.Wrap(err, "Skydive analyzer deployment creation has failed")
		}

		err = kclient_instance.CreateNewService(svc)
		if err != nil {
			return ctrl.Result{}, errors.Wrap(err, "Service creation has failed")
		}
		if skydive_suite.Spec.OpenShiftDeployment {
			if skydive_suite.Spec.Enable.Route {
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

	} else {

		log.Info("Deleting Skydive Analyzers")

		err = kclient_instance.DeleteDeployment(dep)
		if err != nil {
			return ctrl.Result{}, errors.Wrap(err, "Skydive analyzer deployment deletion has failed")
		}

		err = kclient_instance.DeleteService(svc)
		if err != nil {
			return ctrl.Result{}, errors.Wrap(err, "Skydive analyzer service deletion has failed")
		}

		if skydive_suite.Spec.OpenShiftDeployment {
			err = kclient_instance.DeleteRoute(route)
			if err != nil {
				return ctrl.Result{}, errors.Wrap(err, "Skydive analyzer route deletion has failed")
			}
		}
	}

	log.Info("Starting Skydive Agents")

	ds, err := kclient.NewDaemonSet(assets.MustNewAssetReader(SkydiveAgentsDaemonSet))
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "initializing skydive-agents DaemonSet failed")
	}
	ds.Namespace = skydive_suite.Spec.Namespace
	for index := range ds.Spec.Template.Spec.Containers {
		if ds.Spec.Template.Spec.Containers[index].Name == "skydive-agent" {
			ds.Spec.Template.Spec.Containers[index].Env = skydive_suite.Spec.Agents.DaemonSet.Env
			break
		}
	}

	// Creating skydive Agents
	if skydive_suite.Spec.Enable.Agents {
		err = kclient_instance.CreateOrUpdateDaemonSet(ds)
		if err != nil {
			return ctrl.Result{}, errors.Wrap(err, "DaemonSet creation failed")
		}
	} else {
		log.Info("Deleting Skydive Agents")

		err = kclient_instance.DeleteDaemonSet(ds)
		if err != nil {
			return ctrl.Result{}, errors.Wrap(err, "Skydive agents DaemonSet deletion has failed")
		}
	}

	dep, err = kclient.NewDeployment(assets.MustNewAssetReader(SkydiveFlowExporterDep))
	if err != nil {
		log.Error(err, "initializing skydive flow-exporter Deployment failed")
		return ctrl.Result{}, errors.Wrap(err, "initializing skydive flow-exporter Deployment failed")
	}
	dep.Namespace = skydive_suite.Spec.Namespace
	for index := range dep.Spec.Template.Spec.Containers {
		if dep.Spec.Template.Spec.Containers[index].Name == "skydive-flow-exporter" {
			dep.Spec.Template.Spec.Containers[index].Env = skydive_suite.Spec.FlowExporter.Deployment.Env
			break
		}
	}

	if skydive_suite.Spec.Enable.FlowExporter {
		log.Info("Starting Skydive FlowExporter")
		err = kclient_instance.CreateOrUpdateDeployment(dep)
		if err != nil {
			return ctrl.Result{}, errors.Wrap(err, "Flow exporter deployment creation has failed")
		}
	} else {
		err = kclient_instance.DeleteDeployment(dep)
		if err != nil {
			return ctrl.Result{}, errors.Wrap(err, "Flow exporter deployment deletion has failed")
		}
	}

	config_map, err := kclient.NewConfigMap(assets.MustNewAssetReader(PrometheusConfigMap))
	if err != nil {
		log.Error(err, "initializing Prometheus ConfigMap failed")
	}
	config_map.Namespace = skydive_suite.Spec.Namespace

	dep, err = kclient.NewDeployment(assets.MustNewAssetReader(PrometheusDeployment))
	if err != nil {
		log.Error(err, "initializing Prometheus Deployment failed")
		return ctrl.Result{}, errors.Wrap(err, "initializing Prometheus Deployment failed")
	}
	dep.Namespace = skydive_suite.Spec.Namespace

	svc, err = kclient.NewService(assets.MustNewAssetReader(PrometheusService))
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "initializing Prometheus Service failed")
	}
	svc.Namespace = skydive_suite.Spec.Namespace

	connector_dep, err := kclient.NewDeployment(assets.MustNewAssetReader(PrometheusConnectorDeployment))
	if err != nil {
		log.Error(err, "initializing Prometheus Connector Deployment failed")
		return ctrl.Result{}, errors.Wrap(err, "initializing Prometheus Deployment failed")
	}

	connector_dep.Namespace = skydive_suite.Spec.Namespace
	for index := range connector_dep.Spec.Template.Spec.Containers {
		if connector_dep.Spec.Template.Spec.Containers[index].Name == "skydive-prometheus-connector" {
			connector_dep.Spec.Template.Spec.Containers[index].Env =
				skydive_suite.Spec.PrometheusConnector.Deployment.Env
			break
		}
	}

	connector_svc, err := kclient.NewService(assets.MustNewAssetReader(PrometheusConnectorService))
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "initializing Prometheus Service failed")
	}
	connector_svc.Namespace = skydive_suite.Spec.Namespace

	connector_route, err := kclient.NewRoute(assets.MustNewAssetReader(PrometheusConnectorRoute))
	if err != nil {
		log.Error(err, "initializing Prometheus Route failed")
	}
	connector_route.Namespace = skydive_suite.Spec.Namespace

	if skydive_suite.Spec.Enable.PrometheusConnector {
		log.Info("Starting PrometheusConnector")

		err = kclient_instance.CreateOrUpdateConfigMap(config_map)
		if err != nil {
			return ctrl.Result{}, errors.Wrap(err, "Prometheus ConfigMap creation has failed")
		}

		err = kclient_instance.CreateOrUpdateDeployment(dep)
		if err != nil {
			return ctrl.Result{}, errors.Wrap(err, "Prometheus deployment creation has failed")
		}

		err = kclient_instance.CreateNewService(svc)
		if err != nil {
			return ctrl.Result{}, errors.Wrap(err, "Prometheus Service creation has failed")
		}
		if skydive_suite.Spec.OpenShiftDeployment {
			route, err := kclient.NewRoute(assets.MustNewAssetReader(PrometheusRoute))
			if err != nil {
				log.Error(err, "initializing Prometheus Route failed")
			}
			route.Namespace = skydive_suite.Spec.Namespace

			err = kclient_instance.CreateRouteIfNotExists(route)
			if err != nil {
				log.Error(err, "Prometheus Route creation has failed")
			}

			_, err = kclient_instance.WaitForRouteReady(route)
			if err != nil {
				log.Error(err, "waiting for Prometheus Route to become ready failed")
			}
		}

		err = kclient_instance.CreateOrUpdateDeployment(connector_dep)
		if err != nil {
			return ctrl.Result{}, errors.Wrap(err, "Prometheus Connector deployment creation has failed")
		}

		err = kclient_instance.CreateNewService(connector_svc)
		if err != nil {
			return ctrl.Result{}, errors.Wrap(err, "Prometheus Service creation has failed")
		}

		if skydive_suite.Spec.OpenShiftDeployment {

			err = kclient_instance.CreateRouteIfNotExists(connector_route)
			if err != nil {
				log.Error(err, "Prometheus Route creation has failed")
			}

			_, err = kclient_instance.WaitForRouteReady(connector_route)
			if err != nil {
				log.Error(err, "waiting for Prometheus Route to become ready failed")
			}
		}

	} else {

		err = kclient_instance.DeleteConfigMap(config_map)
		if err != nil {
			return ctrl.Result{}, errors.Wrap(err, "Prometheus ConfigMap deletion has failed")
		}

		err = kclient_instance.DeleteDeployment(dep)
		if err != nil {
			return ctrl.Result{}, errors.Wrap(err, "Prometheus deployment deletion has failed")
		}

		err = kclient_instance.DeleteService(svc)
		if err != nil {
			return ctrl.Result{}, errors.Wrap(err, "Prometheus Service deletion has failed")
		}
		if skydive_suite.Spec.OpenShiftDeployment {
			route, err := kclient.NewRoute(assets.MustNewAssetReader(PrometheusRoute))
			if err != nil {
				log.Error(err, "initializing Prometheus Route failed")
			}
			route.Namespace = skydive_suite.Spec.Namespace

			err = kclient_instance.CreateRouteIfNotExists(route)
			if err != nil {
				log.Error(err, "Prometheus Route creation has failed")
			}

			_, err = kclient_instance.WaitForRouteReady(route)
			if err != nil {
				log.Error(err, "waiting for Prometheus Route to become ready failed")
			}
		}

		err = kclient_instance.DeleteDeployment(connector_dep)
		if err != nil {
			return ctrl.Result{}, errors.Wrap(err, "Prometheus Connector deployment deletion has failed")
		}

		err = kclient_instance.DeleteService(connector_svc)
		if err != nil {
			return ctrl.Result{}, errors.Wrap(err, "Prometheus Service deletion has failed")
		}

		err = kclient_instance.DeleteRoute(connector_route)
		if err != nil {
			log.Error(err, "Prometheus Route deletion has failed")
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SkydiveReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&skydivegroupv1.Skydive{}).
		Complete(r)
}
