// Copyright 2022, OpenSergo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controller

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/alibaba/sentinel-golang/logging"
	"github.com/alibaba/sentinel-golang/util"
	crdv1alpha1 "github.com/opensergo/opensergo-control-plane/pkg/api/v1alpha1"
	crdv1betanetworking "github.com/opensergo/opensergo-control-plane/pkg/api/v1beta1/networking"
	"github.com/opensergo/opensergo-control-plane/pkg/model"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)

	_ = crdv1alpha1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme

	_ = crdv1betanetworking.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

type CRDType int32

const (
	FaultToleranceRuleCRDType CRDType = iota
	RateLimitStrategyCRDType
)

func (c CRDType) String() string {
	switch c {
	case FaultToleranceRuleCRDType:
		return "fault-tolerance.opensergo.io/v1alpha1/FaultToleranceRule"
	case RateLimitStrategyCRDType:
		return "fault-tolerance.opensergo.io/v1alpha1/RateLimitStrategy"
	default:
		return "Undefined"
	}
}

type KubernetesOperator struct {
	crdManager  ctrl.Manager
	controllers map[string]*CRDWatcher
	ctx         context.Context
	ctxCancel   context.CancelFunc
	started     atomic.Value

	sendDataHandler model.DataEntirePushHandler

	controllerMux sync.RWMutex
}

// NewKubernetesOperator creates a OpenSergo Kubernetes operator.
func NewKubernetesOperator(sendDataHandler model.DataEntirePushHandler) (*KubernetesOperator, error) {
	ctrl.SetLogger(&k8SLogger{
		l:             logging.GetGlobalLogger(),
		level:         logging.GetGlobalLoggerLevel(),
		names:         make([]string, 0),
		keysAndValues: make([]interface{}, 0),
	})
	k8sConfig, err := ctrl.GetConfig()
	if err != nil {
		return nil, err
	}
	mgr, err := ctrl.NewManager(k8sConfig, ctrl.Options{
		Scheme: scheme,
		// disable metric server
		MetricsBindAddress:     "0",
		HealthProbeBindAddress: "0",
		LeaderElection:         false,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	k := &KubernetesOperator{
		crdManager:      mgr,
		controllers:     make(map[string]*CRDWatcher),
		ctx:             ctx,
		ctxCancel:       cancel,
		sendDataHandler: sendDataHandler,
	}
	return k, nil
}

func (k *KubernetesOperator) RegisterControllersAndStart(info model.SubscribeTarget) error {
	_, err := k.RegisterWatcher(info)
	if err != nil {
		return err
	}
	return k.Run()
}

// RegisterWatcher registers given CRD type and CRD name.
// For each CRD type, it can be registered only once.
func (k *KubernetesOperator) RegisterWatcher(target model.SubscribeTarget) (*CRDWatcher, error) {
	k.controllerMux.Lock()
	defer k.controllerMux.Unlock()

	var err error

	existingWatcher, exists := k.controllers[target.Kind]
	if exists {
		if existingWatcher.HasSubscribed(target) {
			// Target has been subscribed
			return existingWatcher, nil
		} else {
			// Add subscribe to existing watcher
			err = existingWatcher.AddSubscribeTarget(target)
			if err != nil {
				return nil, err
			}
		}
	} else {
		crdMetadata, crdSupports := GetCrdMetadata(target.Kind)
		if !crdSupports {
			return nil, errors.New("CRD not supported: " + target.Kind)
		}
		// This kind of CRD has never been watched.
		crdWatcher := NewCRDWatcher(k.crdManager, target.Kind, crdMetadata.Generator(), k.sendDataHandler)
		err = crdWatcher.AddSubscribeTarget(target)
		if err != nil {
			return nil, err
		}
		err = crdWatcher.SetupWithManager(k.crdManager)
		if err != nil {
			return nil, err
		}
		k.controllers[target.Kind] = crdWatcher
	}
	setupLog.Info("OpenSergo CRD watcher has been registered successfully", "kind", target.Kind, "namespace", target.Namespace, "app", target.AppName)
	return k.controllers[target.Kind], nil
}

func (k *KubernetesOperator) AddWatcher(target model.SubscribeTarget) error {
	k.controllerMux.Lock()
	defer k.controllerMux.Unlock()

	var err error

	existingWatcher, exists := k.controllers[target.Kind]
	if exists && !existingWatcher.HasSubscribed(target) {
		// TODO: think more about here
		err = existingWatcher.AddSubscribeTarget(target)
		if err != nil {
			return err
		}
	} else {
		crdMetadata, crdSupports := GetCrdMetadata(target.Kind)
		if !crdSupports {
			return errors.New("CRD not supported: " + target.Kind)
		}
		crdWatcher := NewCRDWatcher(k.crdManager, target.Kind, crdMetadata.Generator(), k.sendDataHandler)
		err = crdWatcher.AddSubscribeTarget(target)
		if err != nil {
			return err
		}

		crdRunnable, err := ctrl.NewControllerManagedBy(k.crdManager).For(crdMetadata.Generator()()).Build(crdWatcher)
		if err != nil {
			return err
		}
		err = k.crdManager.Add(crdRunnable)
		if err != nil {
			return err
		}
		//_ = crdRunnable.Start(k.ctx)
		k.controllers[target.Kind] = crdWatcher

	}
	setupLog.Info("OpenSergo CRD watcher has been added successfully")
	return nil
}

// Close exit the K8S KubernetesOperator
func (k *KubernetesOperator) Close() error {
	k.ctxCancel()
	return nil
}

func (k *KubernetesOperator) ComponentName() string {
	return "OpenSergoKubernetesOperator"
}

// Run runs the k8s KubernetesOperator
func (k *KubernetesOperator) Run() error {

	// +kubebuilder:scaffold:builder
	go util.RunWithRecover(func() {
		setupLog.Info("Starting OpenSergo operator")
		if err := k.crdManager.Start(k.ctx); err != nil {
			setupLog.Error(err, "problem running OpenSergo operator")
		}
		setupLog.Info("OpenSergo operator will be closed")
	})
	return nil
}

func (k *KubernetesOperator) GetWatcher(kind string) (*CRDWatcher, bool) {
	k.controllerMux.RLock()
	defer k.controllerMux.RUnlock()
	watcher, exists := k.controllers[kind]
	return watcher, exists
}
