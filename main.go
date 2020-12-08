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

package main

import (
	"flag"
	"os"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/source"

	//"sigs.k8s.io/controller-runtime/pkg/webhook"

	cloudshipv1alpha1 "github.com/ToucanSoftware/cloudship/api/v1alpha1"
	"github.com/ToucanSoftware/cloudship/controllers"
	"github.com/ToucanSoftware/cloudship/pkg/reconcile"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	log.SetLogger(zap.New())
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(cloudshipv1alpha1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.Parse()

	// Setup a Manager
	setupLog.Info("setting up manager")
	mgr, err := manager.New(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		Port:               9443,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "d4bd861b.toucansoft.io",
	})
	if err != nil {
		setupLog.Error(err, "unable to set up overall controller manager")
		os.Exit(1)
	}

	// Setup a new controller to reconcile Deployments
	setupLog.Info("Setting up controller")
	c, err := controller.New("deployment-controller", mgr, controller.Options{
		Reconciler: &reconcile.DeploymentReconciliator{
			Client: mgr.GetClient(),
			Log:    ctrl.Log.WithName("controllers").WithName("CloudShip"),
		},
	})
	if err != nil {
		setupLog.Error(err, "unable to set up individual controller")
		os.Exit(1)
	}

	// Watch Deployments and enqueue ReplicaSet object key
	if err := c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForObject{},
		reconcile.DefaultDeploymentPredicates()); err != nil {
		setupLog.Error(err, "unable to watch Deployments")
		os.Exit(1)
	}

	// Watch Pods and enqueue owning ReplicaSet key
	if err := c.Watch(&source.Kind{Type: &corev1.Pod{}},
		&handler.EnqueueRequestForOwner{OwnerType: &appsv1.Deployment{}, IsController: true}); err != nil {
		setupLog.Error(err, "unable to watch Pods")
		os.Exit(1)
	}

	if os.Getenv("ENABLE_WEBHOOKS") != "false" {
		// Setup webhooks
		setupLog.Info("setting up webhook server")
		hookServer := mgr.GetWebhookServer()

		setupLog.Info("registering webhooks to the webhook server")
		//	hookServer.Register("/mutate-v1-pod", &webhook.Admission{Handler: &podAnnotator{Client: mgr.GetClient()}})
		//	hookServer.Register("/validate-v1-pod", &webhook.Admission{Handler: &podValidator{Client: mgr.GetClient()}})
		_ = hookServer
	}

	if err = (&controllers.CanaryDeployReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("CanaryDeploy"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "CanaryDeploy")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "unable to run manager")
		os.Exit(1)
	}
}
