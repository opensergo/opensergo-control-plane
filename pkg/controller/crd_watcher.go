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
	"github.com/go-logr/logr"
	crdv1alpha1 "github.com/opensergo/opensergo-control-plane/pkg/api/v1alpha1"
	crdv1alpha1traffic "github.com/opensergo/opensergo-control-plane/pkg/api/v1alpha1/traffic"
	"github.com/opensergo/opensergo-control-plane/pkg/model"
	pb "github.com/opensergo/opensergo-control-plane/pkg/proto/fault_tolerance/v1"
	trpb "github.com/opensergo/opensergo-control-plane/pkg/proto/transport/v1"
	"github.com/opensergo/opensergo-control-plane/pkg/util"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	k8sApiError "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"log"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
	"sync"
)

const (
	EXTENSION_ROUTE_FALL_BACK = "envoy.router.cluster_specifier_plugin.cluster_fallback"
)

// CRDWatcher watches a specific kind of CRD.
type CRDWatcher struct {
	kind model.SubscribeKind

	client.Client
	logger logr.Logger
	scheme *runtime.Scheme

	// crdCache represents associated local cache for current kind of CRD.
	crdCache *CRDCache

	// subscribedList consists of all subscribed target of current kind of CRD.
	subscribedList       map[model.SubscribeTarget]bool
	subscribedNamespaces map[string]bool
	subscribedApps       map[string]bool

	crdGenerator    func() client.Object
	sendDataHandler model.DataEntirePushHandler

	updateMux sync.RWMutex
}

const (
	UpdateRule = 201
	DeleteRule = 202
	AddRule    = 203
)

func (r *CRDWatcher) Kind() model.SubscribeKind {
	return r.kind
}

func (r *CRDWatcher) HasSubscribed(target model.SubscribeTarget) bool {
	if target.Kind != r.kind {
		return false
	}

	r.updateMux.RLock()
	defer r.updateMux.RUnlock()

	has, _ := r.subscribedList[target]
	return has
}

func (r *CRDWatcher) AddSubscribeTarget(target model.SubscribeTarget) error {
	// TODO: validate the target
	if target.Kind != r.kind {
		return errors.New("target kind mismatch, expected: " + target.Kind + ", r.kind: " + r.kind)
	}
	r.updateMux.Lock()
	defer r.updateMux.Unlock()

	r.subscribedList[target] = true
	r.subscribedNamespaces[target.Namespace] = true
	r.subscribedApps[target.AppName] = true

	return nil
}

func (r *CRDWatcher) RemoveSubscribeTarget(target model.SubscribeTarget) error {
	// TODO: implement me

	return nil
}

func (r *CRDWatcher) HasAnySubscribedOfNamespace(namespace string) bool {
	r.updateMux.RLock()
	defer r.updateMux.RUnlock()

	_, exist := r.subscribedNamespaces[namespace]
	return exist
}

func (r *CRDWatcher) HasAnySubscribedOfApp(app string) bool {
	r.updateMux.RLock()
	defer r.updateMux.RUnlock()

	_, exist := r.subscribedApps[app]
	return exist
}

func (r *CRDWatcher) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	if r.HasAnySubscribedOfApp(req.Namespace) {
		// Ignore unmatched namespace
		return ctrl.Result{Requeue: false, RequeueAfter: 0}, nil
	}
	log := r.logger.WithValues("crdNamespace", req.Namespace, "crdName", req.Name, "kind", r.kind)

	// your logic here
	crd := r.crdGenerator()
	if err := r.Get(ctx, req.NamespacedName, crd); err != nil {
		k8sApiErr, ok := err.(*k8sApiError.StatusError)
		if !ok {
			log.Error(err, "Failed to get OpenSergo CRD")
			return ctrl.Result{
				Requeue:      false,
				RequeueAfter: 0,
			}, nil
		}
		if k8sApiErr.Status().Code != http.StatusNotFound {
			log.Error(err, "Failed to get OpenSergo CRD")
			return ctrl.Result{
				Requeue:      false,
				RequeueAfter: 0,
			}, nil
		}

		// cr had been deleted
		crd = nil
	}

	app := ""
	if crd != nil {
		// TODO: bugs here: we need to check for namespace-app group, not only for app.
		// 		 And we may also need to check for namespace change of a CRD.
		var hasAppLabel bool
		app, hasAppLabel = crd.GetLabels()["app"]
		appSubscribed := r.HasAnySubscribedOfApp(app)
		if !hasAppLabel || !appSubscribed {
			if _, prevContains := r.crdCache.GetByNamespacedName(req.NamespacedName); prevContains {
				log.Info("OpenSergo CRD will be deleted because app label has been changed", "newApp", app)
				crd = nil
			} else {
				// Ignore unmatched app label
				return ctrl.Result{
					Requeue:      false,
					RequeueAfter: 0,
				}, nil
			}
		} else {
			log.Info("OpenSergo CRD received", "crd", crd)
		}
		r.crdCache.SetByNamespaceApp(model.NamespacedApp{
			Namespace: req.Namespace,
			App:       app,
		}, crd)
		r.crdCache.SetByNamespacedName(req.NamespacedName, crd)

	} else {
		app, _ = r.crdCache.GetAppByNamespacedName(req.NamespacedName)
		r.crdCache.DeleteByNamespaceApp(model.NamespacedApp{Namespace: req.Namespace, App: app}, req.Name)
		r.crdCache.DeleteByNamespacedName(req.NamespacedName)
		log.Info("OpenSergo CRD will be deleted")
	}

	nsa := model.NamespacedApp{
		Namespace: req.Namespace,
		App:       app,
	}
	// TODO: Now we can do something for the crd object!
	rules, version := r.GetRules(nsa)
	status := &trpb.Status{
		Code:    int32(200),
		Message: "Get and send rule success",
		Details: nil,
	}
	dataWithVersion := &trpb.DataWithVersion{Data: rules, Version: version}
	err := r.sendDataHandler(req.Namespace, app, r.kind, dataWithVersion, status, "")
	if err != nil {
		log.Error(err, "Failed to send rules", "kind", r.kind)
	}
	return ctrl.Result{}, nil
}

func (r *CRDWatcher) GetRules(n model.NamespacedApp) ([]*anypb.Any, int64) {
	var rules []*anypb.Any
	objs, version := r.crdCache.GetByNamespaceApp(n)
	for _, obj := range objs {
		if obj == nil {
			continue
		}
		rule, err := r.translateCrdToProto(obj)
		if err != nil {
			return rules, version
		}
		rules = append(rules, rule)
	}
	return rules, version
}

func (r *CRDWatcher) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).For(r.crdGenerator()).Complete(r)
}

func (r *CRDWatcher) translateCrdToProto(object client.Object) (*anypb.Any, error) {
	var packRule *anypb.Any
	var err error
	var rule proto.Message
	switch r.kind {
	case FaultToleranceRuleKind:
		ftr := object.(*crdv1alpha1.FaultToleranceRule)
		var targets []*pb.FaultToleranceRule_FaultToleranceRuleTargetRef
		var strategies []*pb.FaultToleranceRule_FaultToleranceStrategyRef
		if ftr != nil {
			for _, target := range ftr.Spec.Targets {
				targets = append(targets, &pb.FaultToleranceRule_FaultToleranceRuleTargetRef{TargetResourceName: target.TargetResourceName})
			}
			for _, strategy := range ftr.Spec.Strategies {
				strategies = append(strategies, &pb.FaultToleranceRule_FaultToleranceStrategyRef{
					Name: strategy.Name,
					Kind: strategy.Kind,
				})
			}
		}
		rule = &pb.FaultToleranceRule{
			Targets:    targets,
			Strategies: strategies,
			Action:     nil,
		}

	case RateLimitStrategyKind:
		rls := object.(*crdv1alpha1.RateLimitStrategy)
		mType, _ := strconv.ParseInt(rls.Spec.MetricType, 10, 32)
		limitMode, _ := strconv.ParseInt(rls.Spec.LimitMode, 10, 32)
		rule = &pb.RateLimitStrategy{
			Name:         rls.Name,
			MetricType:   pb.RateLimitStrategy_MetricType(mType),
			LimitMode:    pb.RateLimitStrategy_LimitMode(limitMode),
			Threshold:    rls.Spec.Threshold,
			StatDuration: rls.Spec.StatDurationSeconds,
			//StatDurationTimeUnit: 2, //todo 2应该为参数
		}
	case ThrottlingStrategyKind:
		ts := object.(*crdv1alpha1.ThrottlingStrategy)
		miMill, err := util.Str2MillSeconds(ts.Spec.MinIntervalOfRequests)
		if err != nil {
			miMill = -1
			log.Println("translate to MinIntervalMillisOfRequests error, ", err)
		}
		qtMill, err := util.Str2MillSeconds(ts.Spec.QueueTimeout)
		if err != nil {
			qtMill = -1
			log.Println("translate to QueueTimeoutMillis error, ", err)
		}
		rule = &pb.ThrottlingStrategy{
			Name:                        ts.Name,
			MinIntervalMillisOfRequests: miMill,
			QueueTimeoutMillis:          qtMill,
		}
	case CircuitBreakerStrategyKind:
		cbs := object.(*crdv1alpha1.CircuitBreakerStrategy)
		tr, err := util.RatioStr2Float(cbs.Spec.TriggerRatio)
		if err != nil {
			tr = -1.0
			log.Println("translate to TriggerRatio error, ", err)
		}
		sdMill, err := util.Str2MillSeconds(cbs.Spec.StatDuration)
		if err != nil {
			sdMill = -1
		}
		rtMill, err := util.Str2MillSeconds(cbs.Spec.RecoveryTimeout)
		if err != nil {
			rtMill = -1
		}
		maMill, err := util.Str2MillSeconds(cbs.Spec.SlowConditions.MaxAllowedRt)
		if err != nil {
			maMill = -1
		}

		rule = &pb.CircuitBreakerStrategy{
			Name:                    cbs.Name,
			Strategy:                util.Str2CBStrategy(cbs.Spec.Strategy),
			TriggerRatio:            tr,
			StatDuration:            sdMill, // todo int64 or int32
			StatDurationTimeUnit:    0,
			RecoveryTimeout:         int32(rtMill),
			RecoveryTimeoutTimeUnit: 0,
			MinRequestAmount:        cbs.Spec.MinRequestAmount,
			SlowCondition:           &pb.CircuitBreakerStrategy_CircuitBreakerSlowCondition{MaxAllowedRtMillis: int32(maMill)},
			ErrorCondition:          nil,
		}
	case ConcurrencyLimitStrategyKind:
		cls := object.(*crdv1alpha1.ConcurrencyLimitStrategy)
		rule = &pb.ConcurrencyLimitStrategy{
			Name:           cls.Name,
			LimitMode:      util.Str2LimitNode(cls.Spec.LimitMode),
			MaxConcurrency: cls.Spec.MaxConcurrencyThreshold,
		}
	case TrafficRouterKind:
		cls := object.(*crdv1alpha1traffic.TrafficRouter)
		rule = BuildRouteConfiguration(cls)
	default:
		return nil, nil
	}
	packRule, err = anypb.New(rule)
	if err != nil {
		log.Println("pack rule error", err)
		return nil, err
	}
	return packRule, nil

}

func NewCRDWatcher(crdManager ctrl.Manager, kind model.SubscribeKind, crdGenerator func() client.Object, sendDataHandler model.DataEntirePushHandler) *CRDWatcher {
	return &CRDWatcher{
		kind:                 kind,
		Client:               crdManager.GetClient(),
		logger:               ctrl.Log.WithName("controller").WithName(kind),
		scheme:               crdManager.GetScheme(),
		subscribedList:       make(map[model.SubscribeTarget]bool, 4),
		subscribedNamespaces: make(map[string]bool),
		subscribedApps:       make(map[string]bool),
		crdGenerator:         crdGenerator,
		crdCache:             NewCRDCache(kind),
		sendDataHandler:      sendDataHandler,
	}
}
