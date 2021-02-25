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
	"skydive/pkg/kclient"
	"skydive/pkg/manifests"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	skydivev1beta1 "skydive/api/v1beta1"
)

var (
	SkydiveAgentsDaemonSet = "skydive-agents/daemon-set.yaml"
)

// SkydiveAgentsReconciler reconciles a SkydiveAgents object
type SkydiveAgentsReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=skydive.example.com,resources=skydiveagents,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=skydive.example.com,resources=skydiveagents/status,verbs=get;update;patch
func (r *SkydiveAgentsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("skydiveanalyzer", req.NamespacedName)

	log.Info("initializing assets for skydive-agents")
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

	log.Info("initializing skydive-agents DaemonSet...")
	ds, err := kclient.NewDaemonSet(assets.MustNewAssetReader(SkydiveAgentsDaemonSet))
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "initializing skydive-agents DaemonSet failed")
	}
	ds.Namespace = namespaceName

	log.Info("Creating DaemonSet...")
	err = kclient_instance.CreateOrUpdateDaemonSet(ds)
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "UI-Route creation has failed")
	}

	log.Info("reconciling skydive-analyzer task completed successfully")
	return ctrl.Result{}, nil

}

func (r *SkydiveAgentsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&skydivev1beta1.SkydiveAgents{}).
		Complete(r)
}
