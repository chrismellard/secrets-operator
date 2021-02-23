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
	"time"

	"github.com/chrismellard/secret-operator/pkg/clients/kube"
	"github.com/chrismellard/secret-operator/pkg/secretstores/factory"
	"github.com/chrismellard/secret-operator/pkg/secretstores/health"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	secretoperatorv1alpha1 "github.com/chrismellard/secret-operator/api/v1alpha1"
)

// SecretStoreReconciler reconciles a SecretStore object
type SecretStoreReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=secret-operator.io,resources=secretstores,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=secret-operator.io,resources=secretstores/status,verbs=get;update;patch

func (r *SecretStoreReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("secretstore", req.NamespacedName)

	var store secretoperatorv1alpha1.SecretStore
	if err := r.Get(ctx, req.NamespacedName, &store); err != nil {
		log.Error(err, "unable to fetch CronJob")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	provider, err := factory.CreateSecretStoreProvider(store.Spec.Provider)
	if err != nil {
		return ctrl.Result{Requeue: true, RequeueAfter: 30 * time.Second}, nil
	}

	kubeClient, err := kube.CreateClientSet()
	if err != nil {
		return ctrl.Result{Requeue: true, RequeueAfter: 30 * time.Second}, nil
	}

	manager, err := factory.CreateSecretStoreManager(ctx, store.Spec.Provider, kubeClient)
	healthy, err := health.CheckSecretStoreHealth(manager, provider)

	if !healthy {
		return ctrl.Result{Requeue: true, RequeueAfter: 30 * time.Second}, nil
	}
	return ctrl.Result{}, nil
}

func (r *SecretStoreReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&secretoperatorv1alpha1.SecretStore{}).
		Complete(r)
}
