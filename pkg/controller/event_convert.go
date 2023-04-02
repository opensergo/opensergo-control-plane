package controller

import (
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
			RuntimeStrategy: &v1.EventRuntimeStrategy{
				FaultToleranceRule: ec.crdWatcher.translateFTRCrdToProto(
					&v1alpha1.FaultToleranceRule{
						Spec: ss.RuntimeStrategy.FaultToleranceRule,
					}),
				RateLimitStrategy: ec.crdWatcher.translateRLSCrdToProto(
					&v1alpha1.RateLimitStrategy{
						Spec: ss.RuntimeStrategy.RateLimitStrategy,
					}),
				CircuitBreakerStrategy: ec.crdWatcher.translateCBSCrdToProto(
					&v1alpha1.CircuitBreakerStrategy{
						Spec: ss.RuntimeStrategy.CircuitBreakerStrategy,
					}),
				ConcurrencyLimitStrateg: ec.crdWatcher.translateCLSCrdToProto(
					&v1alpha1.ConcurrencyLimitStrategy{
						Spec: ss.RuntimeStrategy.ConcurrencyLimitStrategy,
					}),
				RetryRule: &v1.RetryRule{
					RetryMax:     ss.RuntimeStrategy.RetryRule.RetryMax,
					BackOffDelay: ss.RuntimeStrategy.RetryRule.BackOffDelay,
					BackOffPolicyType: util.Int2EventRetryBackOffPolicyType(
						ss.RuntimeStrategy.RetryRule.BackOffPolicyType,
					),
				},
			},
		}
		sourceStrategies = append(sourceStrategies, s)
	}
	for _, ts := range crdStrategies.TriggerStrategies {
		s := &v1.EventTriggerStrategy{
			EventTriggerId:    ts.EventTriggerID,
			ReceiveBufferSize: ts.ReceiveBufferSize,
			EnableIdempotence: ts.EnableIdempotence,
			RuntimeStrategy: &v1.EventRuntimeStrategy{
				FaultToleranceRule: ec.crdWatcher.translateFTRCrdToProto(
					&v1alpha1.FaultToleranceRule{
						Spec: ts.RuntimeStrategy.FaultToleranceRule,
					}),
				RateLimitStrategy: ec.crdWatcher.translateRLSCrdToProto(
					&v1alpha1.RateLimitStrategy{
						Spec: ts.RuntimeStrategy.RateLimitStrategy,
					}),
				CircuitBreakerStrategy: ec.crdWatcher.translateCBSCrdToProto(
					&v1alpha1.CircuitBreakerStrategy{
						Spec: ts.RuntimeStrategy.CircuitBreakerStrategy,
					}),
				ConcurrencyLimitStrateg: ec.crdWatcher.translateCLSCrdToProto(
					&v1alpha1.ConcurrencyLimitStrategy{
						Spec: ts.RuntimeStrategy.ConcurrencyLimitStrategy,
					}),
				RetryRule: &v1.RetryRule{
					RetryMax:     ts.RuntimeStrategy.RetryRule.RetryMax,
					BackOffDelay: ts.RuntimeStrategy.RetryRule.BackOffDelay,
					BackOffPolicyType: util.Int2EventRetryBackOffPolicyType(
						ts.RuntimeStrategy.RetryRule.BackOffPolicyType,
					),
				},
			},
			DeadLetterStrategy: &v1.DeadLetterStrategy{
				Enable:                ts.DeadLetterStrategy.Enable,
				RetryTriggerThreshold: ts.DeadLetterStrategy.RetryTriggerThreshold,
				Store:                 nil,
			},
		}
		// first use event channel if url or topic is empty then use persistence
		ch := ts.DeadLetterStrategy.StoreEventChannel
		if ch.URL != "" && ch.MQTopicName != "" {
			s.DeadLetterStrategy.Store = &v1.DeadLetterStrategy_Channel{
				Channel: &v1.EventChannel{
					UniqueId:    ch.UniqueID,
					Url:         ch.URL,
					MqTopicName: ch.MQTopicName,
				},
			}
			continue
		}
		persist := ts.DeadLetterStrategy.StorePersistence
		if persist.PersistenceAddress.Address != "" {
			s.DeadLetterStrategy.Store = &v1.DeadLetterStrategy_Persistence{
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
		triggerStrategies = append(triggerStrategies, s)
	}

	strategies.SourceStrategies = sourceStrategies
	strategies.TriggerStrategies = triggerStrategies

	return strategies
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
