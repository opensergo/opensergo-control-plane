package controller

import (
	"reflect"

	"github.com/opensergo/opensergo-control-plane/pkg/api/v1alpha1"
	crdv1alpha1event "github.com/opensergo/opensergo-control-plane/pkg/api/v1alpha1/event"
	v1 "github.com/opensergo/opensergo-control-plane/pkg/proto/event/v1"
	"github.com/opensergo/opensergo-control-plane/pkg/util"
)

// EventConvertor transform crd to pb
type EventConvertor struct {
	crdWatcher *CRDWatcher
}

// NewEventTransform new a event transform
func NewEventTransform(r *CRDWatcher) *EventConvertor {
	return &EventConvertor{
		crdWatcher: r,
	}
}

// BuildEventConfiguration build event pb object
func (ec *EventConvertor) BuildEventConfiguration(e *crdv1alpha1event.Event) *v1.Event {
	event := &v1.Event{
		Components:  ec.buildEventComponents(e.Spec.Components),
		Strategies:  ec.buildEventStrategies(e.Spec.Strategies),
		RouterRules: ec.buildEventRouterRules(e.Spec.RouterRules),
	}
	return event
}

func (ec *EventConvertor) buildEventComponents(crdComponents crdv1alpha1event.EventComponents) *v1.EventComponents {
	var (
		components = new(v1.EventComponents)
		channels   []*v1.EventChannel
		sources    []*v1.EventSource
		triggers   []*v1.EventTrigger
	)

	for _, c := range crdComponents.Channels {
		channel := &v1.EventChannel{
			UniqueId:    c.UniqueID,
			Url:         c.URL,
			MqTopicName: c.MQTopicName,
		}
		channels = append(channels, channel)
	}
	for _, s := range crdComponents.Sources {
		source := &v1.EventSource{Ref: &v1.EventProcessorRef{
			UniqueId: s.Ref.UniqueID,
			Kind:     s.Ref.Kind,
			Name:     s.Ref.Name,
		}}
		sources = append(sources, source)
	}
	for _, r := range crdComponents.Triggers {
		trigger := &v1.EventTrigger{Ref: &v1.EventProcessorRef{
			UniqueId: r.Ref.UniqueID,
			Kind:     r.Ref.Kind,
			Name:     r.Ref.Name,
		}}
		triggers = append(triggers, trigger)
	}

	components.Channels = channels
	components.Sources = sources
	components.Triggers = triggers

	return components
}

func (ec *EventConvertor) buildEventStrategies(crdStrategies crdv1alpha1event.EventStrategies) *v1.EventStrategies {
	var (
		strategies        = new(v1.EventStrategies)
		sourceStrategies  []*v1.EventSourceStrategy
		triggerStrategies []*v1.EventTriggerStrategy
	)

	for _, ss := range crdStrategies.SourceStrategies {
		s := &v1.EventSourceStrategy{
			EventSourceId: ss.EventSourceID,
			AsyncSend:     ss.AsyncSend,
			FaultTolerantStorage: &v1.Persistence{
				PersistenceType: util.Int2EventPersistenceType(ss.FaultTolerantStorage.PersistenceType),
				PersistenceSize: ss.FaultTolerantStorage.PersistenceSize,
				FullStrategy:    util.Int2EventPersistenceFullStrategy(ss.FaultTolerantStorage.FullStrategy),
			},
			RuntimeStrategy: ec.buildRuntimeStrategies(
				ss.RuntimeStrategy, crdStrategies.DefaultSourceRuntimeStrategy),
		}
		sourceStrategies = append(sourceStrategies, s)
	}
	for _, ts := range crdStrategies.TriggerStrategies {
		s := &v1.EventTriggerStrategy{
			EventTriggerId:    ts.EventTriggerID,
			ReceiveBufferSize: ts.ReceiveBufferSize,
			EnableIdempotence: ts.EnableIdempotence,
			RuntimeStrategy: ec.buildRuntimeStrategies(
				ts.RuntimeStrategy, crdStrategies.DefaultTriggerRuntimeStrategy),
			DeadLetterStrategy: ec.buildDeadLetterStrategy(
				ts.DeadLetterStrategy, crdStrategies.DefaultDeadLetterStrategy),
		}
		triggerStrategies = append(triggerStrategies, s)
	}

	strategies.SourceStrategies = sourceStrategies
	strategies.TriggerStrategies = triggerStrategies

	return strategies
}

func (ec *EventConvertor) buildRuntimeStrategies(
	runtimeStrategy,
	defaultStrategy crdv1alpha1event.EventRuntimeStrategy,
) *v1.EventRuntimeStrategy {

	var rs crdv1alpha1event.EventRuntimeStrategy
	if reflect.DeepEqual(runtimeStrategy.FaultToleranceRule, v1alpha1.FaultToleranceRuleSpec{}) {
		rs.FaultToleranceRule = defaultStrategy.FaultToleranceRule
	}
	if reflect.DeepEqual(runtimeStrategy.ConcurrencyLimitStrategy, v1alpha1.ConcurrencyLimitStrategySpec{}) {
		rs.ConcurrencyLimitStrategy = defaultStrategy.ConcurrencyLimitStrategy
	}
	if reflect.DeepEqual(runtimeStrategy.RateLimitStrategy, v1alpha1.RateLimitStrategySpec{}) {
		rs.RateLimitStrategy = defaultStrategy.RateLimitStrategy
	}
	if reflect.DeepEqual(runtimeStrategy.CircuitBreakerStrategy, v1alpha1.CircuitBreakerStrategySpec{}) {
		rs.CircuitBreakerStrategy = defaultStrategy.CircuitBreakerStrategy
	}
	if reflect.DeepEqual(runtimeStrategy.RetryRule, crdv1alpha1event.RetryRule{}) {
		rs.RetryRule = defaultStrategy.RetryRule
	}

	rsPB := &v1.EventRuntimeStrategy{
		FaultToleranceRule: ec.crdWatcher.translateFTRCrdToProto(
			&v1alpha1.FaultToleranceRule{
				Spec: rs.FaultToleranceRule,
			}),
		RateLimitStrategy: ec.crdWatcher.translateRLSCrdToProto(
			&v1alpha1.RateLimitStrategy{
				Spec: rs.RateLimitStrategy,
			}),
		CircuitBreakerStrategy: ec.crdWatcher.translateCBSCrdToProto(
			&v1alpha1.CircuitBreakerStrategy{
				Spec: rs.CircuitBreakerStrategy,
			}),
		ConcurrencyLimitStrateg: ec.crdWatcher.translateCLSCrdToProto(
			&v1alpha1.ConcurrencyLimitStrategy{
				Spec: rs.ConcurrencyLimitStrategy,
			}),
		RetryRule: &v1.RetryRule{
			RetryMax:     rs.RetryRule.RetryMax,
			BackOffDelay: rs.RetryRule.BackOffDelay,
			BackOffPolicyType: util.Int2EventRetryBackOffPolicyType(
				rs.RetryRule.BackOffPolicyType,
			),
		},
	}
	return rsPB
}

func (ec *EventConvertor) buildDeadLetterStrategy(
	deadLetterStrategy,
	defaultStrategy crdv1alpha1event.DeadLetterStrategy) *v1.DeadLetterStrategy {

	var dls crdv1alpha1event.DeadLetterStrategy
	if reflect.DeepEqual(deadLetterStrategy.Enable, false) {
		dls.Enable = defaultStrategy.Enable
	}
	if reflect.DeepEqual(deadLetterStrategy.RetryTriggerThreshold, int64(0)) {
		dls.RetryTriggerThreshold = defaultStrategy.RetryTriggerThreshold
	}
	if reflect.DeepEqual(deadLetterStrategy.StoreEventChannel, crdv1alpha1event.EventChannel{}) {
		dls.StoreEventChannel = defaultStrategy.StoreEventChannel
	}
	if reflect.DeepEqual(deadLetterStrategy.StorePersistence, crdv1alpha1event.Persistence{}) {
		dls.StorePersistence = defaultStrategy.StorePersistence
	}

	dlsPB := &v1.DeadLetterStrategy{
		Enable:                dls.Enable,
		RetryTriggerThreshold: dls.RetryTriggerThreshold,
		Store:                 nil,
	}
	// first use event channel if url or topic is empty then use persistence
	ch := dls.StoreEventChannel
	if ch.URL != "" && ch.MQTopicName != "" {
		dlsPB.Store = &v1.DeadLetterStrategy_Channel{
			Channel: &v1.EventChannel{
				UniqueId:    ch.UniqueID,
				Url:         ch.URL,
				MqTopicName: ch.MQTopicName,
			},
		}
		return dlsPB
	}
	// following is that use persistence
	persist := dls.StorePersistence
	if persist.PersistenceAddress.Address != "" {
		dlsPB.Store = &v1.DeadLetterStrategy_Persistence{
			Persistence: &v1.Persistence{
				PersistenceType: util.Int2EventPersistenceType(
					persist.PersistenceType),
				PersistenceAddress: &v1.Persistence_PersistenceAddress{
					Address: persist.PersistenceAddress.Address,
				},
				PersistenceSize: persist.PersistenceSize,
				FullStrategy: util.Int2EventPersistenceFullStrategy(
					persist.FullStrategy),
			},
		}
	}
	return dlsPB
}

func (ec *EventConvertor) buildEventRouterRules(crdRouterRules crdv1alpha1event.EventRouterRules) *v1.EventRouterRules {
	var (
		routerRules = new(v1.EventRouterRules)
		routers     []*v1.EventRouterRules_Router
	)

	for _, rr := range crdRouterRules.RouterRules {
		r := &v1.EventRouterRules_Router{
			SourceId:  rr.SourceID,
			TriggerId: rr.TriggerID,
			ChannelId: rr.ChannelID,
			Filter: &v1.EventRouterRules_Filter{
				Cesql: rr.Filter.CESQL,
			},
		}
		routers = append(routers, r)
	}

	routerRules.Routers = routers
	return routerRules
}
