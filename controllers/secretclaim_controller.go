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

	secretoperatorv1alpha1 "github.com/chrismellard/secret-operator/api/v1alpha1"
	"github.com/chrismellard/secret-operator/pkg/claimhandlers/factory"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// SecretClaimReconciler reconciles a SecretClaim object
type SecretClaimReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=secret-operator.io,resources=secretclaims,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=secret-operator.io,resources=secretclaims/status,verbs=get;update;patch

func (r *SecretClaimReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("secretclaim", req.NamespacedName)

	var claim secretoperatorv1alpha1.SecretClaim
	if err := r.Get(ctx, req.NamespacedName, &claim); err != nil {
		log.Error(err, "unable to fetch CronJob")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	handler, err := factory.CreateClaimHandler(claim, ctx)
	if err != nil {
		log.Error(err, "unable to create handler for claim")
		return ctrl.Result{Requeue: true, RequeueAfter: 30}, err
	}
	err = handler.Handle()
	if err != nil {
		log.Error(err, "handler failure")
		return ctrl.Result{Requeue: true, RequeueAfter: 30}, err
	}

	return ctrl.Result{}, nil
}

func (r *SecretClaimReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&secretoperatorv1alpha1.SecretClaim{}).
		Complete(r)
}
