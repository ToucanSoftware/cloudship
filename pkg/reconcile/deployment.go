/*
Copyright © 2020 ToucanSoftware

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

package reconcile

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// DeploymentReconciliator reconciles Deployment
type DeploymentReconciliator struct {
	// client can be used to retrieve objects from the APIServer.
	Client client.Client
	Log    logr.Logger
}

// Implement reconcile.Reconciler so the controller can reconcile objects
var _ reconcile.Reconciler = &DeploymentReconciliator{}

// Reconcile reconcilates a object
func (r *DeploymentReconciliator) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// set up a convenient log object so we don't have to type request over and over again
	ctx := context.Background()
	log := r.Log.WithValues("deployment", request.NamespacedName)

	// Fetch the Deployment from the cache
	deploy := &appsv1.Deployment{}
	err := r.Client.Get(ctx, request.NamespacedName, deploy)
	if errors.IsNotFound(err) {
		log.Error(nil, "Could not find Deployment")
		return reconcile.Result{}, nil
	}

	if err != nil {
		return reconcile.Result{}, fmt.Errorf("could not fetch Deployment: %+v", err)
	}

	var namespace = &corev1.Namespace{}
	err = r.Client.Get(ctx, types.NamespacedName{Name: deploy.GetNamespace()}, namespace)
	if errors.IsNotFound(err) {
		log.Error(nil, "Could not find Namespace")
		return reconcile.Result{}, nil
	}

	if namespace.Labels["canary-deploy"] != "enabled" {
		log.Info(fmt.Sprintf("namespace %s of the deployment %s not labeled as canary-deploy=enabled", namespace.GetName(), deploy.GetName()))
		return reconcile.Result{}, nil
	}

	// Print the Deployment
	log.Info("Reconciling Deployment", "container name", deploy.Spec.Template.Spec.Containers[0].Name)

	// Set the label if it is missing
	if deploy.Labels == nil {
		deploy.Labels = map[string]string{}
	}
	if deploy.Labels["clouldship"] == "primary" {
		return reconcile.Result{}, nil
	}

	// Update the Deployment
	deploy.Labels["clouldship"] = "primary"
	err = r.Client.Update(ctx, deploy)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("could not write Deployment: %+v", err)
	}

	return reconcile.Result{}, nil
}
