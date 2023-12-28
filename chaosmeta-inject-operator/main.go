/*
Copyright 2023.

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
	"context"
	"flag"
	"fmt"
	initwebhook "github.com/traas-stack/chaosmeta/chaosmeta-common/webhook"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/restclient"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"time"

	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/common"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/config"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/executor/remoteexecutor"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/pkg/selector"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	injectv1alpha1 "github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/api/v1alpha1"
	"github.com/traas-stack/chaosmeta/chaosmeta-inject-operator/controllers"
	"net/http"
	_ "net/http/pprof"
	//+kubebuilder:scaffold:imports
)

const (
	ComponentInject = "inject"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(injectv1alpha1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	//var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	//flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
		//MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "9cb44693.chaosmeta.io",
		// LeaderElectionReleaseOnCancel defines if the leader should step down voluntarily
		// when the Manager ends. This requires the binary to immediately end when the
		// Manager is stopped, otherwise, this setting is unsafe. Setting this significantly
		// speeds up voluntary leader transitions as the new leader don't have to wait
		// LeaseDuration time first.
		//
		// In the default scaffold provided, the program ends immediately after
		// the manager stops, so would be fine to enable this option. However,
		// if you are doing or is intended to do any operation such as perform cleanups
		// after the manager stops then its usage might be unsafe.
		// LeaderElectionReleaseOnCancel: true,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// basic config
	mainConfig, err := config.LoadConfig("./config/chaosmeta-inject.json")
	if err != nil {
		setupLog.Error(err, "load config error")
		os.Exit(1)
	}

	setupLog.Info(fmt.Sprintf("set main config success: %v", mainConfig))

	selector.SetupAnalyzer(mgr.GetClient())
	common.SetGoroutinePool(mainConfig.Worker.PoolCount)
	setupLog.Info(fmt.Sprintf("set goroutine pool success: %d", mainConfig.Worker.PoolCount))
	go func() {
		err = http.ListenAndServe("localhost:8090", nil)
		if err != nil {
			setupLog.Error(err, "failed to start pprof")
		}
	}()
	// create APIServer client
	t := []injectv1alpha1.CloudTargetType{
		injectv1alpha1.PodCloudTarget,
		injectv1alpha1.DeploymentCloudTarget,
		injectv1alpha1.NodeCloudTarget,
		injectv1alpha1.NamespaceCloudTarget,
		injectv1alpha1.JobCloudTarget,
	}

	if err := restclient.SetApiServerClientMap(mgr.GetConfig(), mgr.GetScheme(), t); err != nil {
		setupLog.Error(err, "set APIServer client error")
		os.Exit(1)
	}
	setupLog.Info(fmt.Sprintf("set APIServer for cloud object success: %v", t))
	err = initwebhook.InitCert(setupLog, ComponentInject)
	if err != nil {
		setupLog.Error(err, "init cert failed")
		os.Exit(1)
	}
	// set executor
	if err = remoteexecutor.SetGlobalRemoteExecutor(&mainConfig.Executor, mgr.GetConfig(), mgr.GetScheme()); err != nil {
		setupLog.Error(err, "set remote executor error")
		os.Exit(1)
	}
	setupLog.Info(fmt.Sprintf("set remote executor success: %s", mainConfig.Executor.Mode))

	// start watching
	if err = (&controllers.ExperimentReconciler{
		Client: mgr.GetClient(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Experiment")
		os.Exit(1)
	}

	if err = (&injectv1alpha1.Experiment{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "Experiment")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	// set autoRecoverTicker = config.ticker
	if mainConfig.Ticker.AutoCheckInterval <= 0 {
		setupLog.Error(fmt.Errorf("ticker interval is invalid"), "must provide a positive integer")
		os.Exit(1)
	}
	go autoRecoverChecker(context.Background(), mainConfig.Ticker.AutoCheckInterval, mgr.GetClient())

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func autoRecoverChecker(ctx context.Context, interval int, c client.Client) {
	logger, ticker := log.FromContext(ctx), time.NewTicker(time.Duration(interval)*time.Second)
	defer ticker.Stop()

	logger.Info(fmt.Sprintf("start auto recover checker success, ticker second: %d", interval))
	for {
		<-ticker.C
		autoRecover(ctx, c)
	}
}

func autoRecover(ctx context.Context, c client.Client) {
	logger := log.FromContext(ctx)
	expList, err := selector.GetAnalyzer().GetExperimentListByPhase(ctx, string(injectv1alpha1.InjectPhaseType))
	if err != nil {
		logger.Error(err, fmt.Sprintf("get experiment list of phase[%s] error", injectv1alpha1.InjectPhaseType))
		return
	}

	exp := expList.Items
	for i := range exp {
		if exp[i].Status.Status == injectv1alpha1.CreatedStatusType ||
			exp[i].Status.Status == injectv1alpha1.RunningStatusType ||
			exp[i].Spec.TargetPhase == injectv1alpha1.RecoverPhaseType {
			continue
		}

		needRecover, err := common.IsTimeout(exp[i].Status.CreateTime, exp[i].Spec.Experiment.Duration)
		if err != nil {
			logger.Error(err, fmt.Sprintf("check timeout of experiment[%s] error", exp[i].Name))
			continue
		}

		if needRecover {
			logger.Info(fmt.Sprintf("experiment[%s] is time to auto recover", exp[i].Name))
			exp[i].Spec.TargetPhase = injectv1alpha1.RecoverPhaseType

			if err := c.Update(ctx, &exp[i]); err != nil {
				logger.Error(err, fmt.Sprintf("experiment[%s] is time to recover, but update \"TargetPhase\" error: %s", exp[i].Name, err.Error()))
			}
		}
	}
}
