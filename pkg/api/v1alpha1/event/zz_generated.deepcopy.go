//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeadLetterStrategy) DeepCopyInto(out *DeadLetterStrategy) {
	*out = *in
	out.StoreEventChannel = in.StoreEventChannel
	out.StorePersistence = in.StorePersistence
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeadLetterStrategy.
func (in *DeadLetterStrategy) DeepCopy() *DeadLetterStrategy {
	if in == nil {
		return nil
	}
	out := new(DeadLetterStrategy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Event) DeepCopyInto(out *Event) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Event.
func (in *Event) DeepCopy() *Event {
	if in == nil {
		return nil
	}
	out := new(Event)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Event) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EventChannel) DeepCopyInto(out *EventChannel) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EventChannel.
func (in *EventChannel) DeepCopy() *EventChannel {
	if in == nil {
		return nil
	}
	out := new(EventChannel)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EventComponents) DeepCopyInto(out *EventComponents) {
	*out = *in
	if in.Channels != nil {
		in, out := &in.Channels, &out.Channels
		*out = make([]EventChannel, len(*in))
		copy(*out, *in)
	}
	if in.Sources != nil {
		in, out := &in.Sources, &out.Sources
		*out = make([]EventSource, len(*in))
		copy(*out, *in)
	}
	if in.Triggers != nil {
		in, out := &in.Triggers, &out.Triggers
		*out = make([]EventTrigger, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EventComponents.
func (in *EventComponents) DeepCopy() *EventComponents {
	if in == nil {
		return nil
	}
	out := new(EventComponents)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EventFilter) DeepCopyInto(out *EventFilter) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EventFilter.
func (in *EventFilter) DeepCopy() *EventFilter {
	if in == nil {
		return nil
	}
	out := new(EventFilter)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EventList) DeepCopyInto(out *EventList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Event, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EventList.
func (in *EventList) DeepCopy() *EventList {
	if in == nil {
		return nil
	}
	out := new(EventList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EventList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EventProcessorRef) DeepCopyInto(out *EventProcessorRef) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EventProcessorRef.
func (in *EventProcessorRef) DeepCopy() *EventProcessorRef {
	if in == nil {
		return nil
	}
	out := new(EventProcessorRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EventRouter) DeepCopyInto(out *EventRouter) {
	*out = *in
	out.Filter = in.Filter
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EventRouter.
func (in *EventRouter) DeepCopy() *EventRouter {
	if in == nil {
		return nil
	}
	out := new(EventRouter)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EventRouterRules) DeepCopyInto(out *EventRouterRules) {
	*out = *in
	if in.RouterRules != nil {
		in, out := &in.RouterRules, &out.RouterRules
		*out = make([]EventRouter, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EventRouterRules.
func (in *EventRouterRules) DeepCopy() *EventRouterRules {
	if in == nil {
		return nil
	}
	out := new(EventRouterRules)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EventRuntimeStrategy) DeepCopyInto(out *EventRuntimeStrategy) {
	*out = *in
	in.FaultToleranceRule.DeepCopyInto(&out.FaultToleranceRule)
	out.RateLimitStrategy = in.RateLimitStrategy
	out.CircuitBreakerStrategy = in.CircuitBreakerStrategy
	out.ConcurrencyLimitStrategy = in.ConcurrencyLimitStrategy
	out.RetryRule = in.RetryRule
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EventRuntimeStrategy.
func (in *EventRuntimeStrategy) DeepCopy() *EventRuntimeStrategy {
	if in == nil {
		return nil
	}
	out := new(EventRuntimeStrategy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EventSource) DeepCopyInto(out *EventSource) {
	*out = *in
	out.Ref = in.Ref
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EventSource.
func (in *EventSource) DeepCopy() *EventSource {
	if in == nil {
		return nil
	}
	out := new(EventSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EventSourceStrategy) DeepCopyInto(out *EventSourceStrategy) {
	*out = *in
	out.FaultTolerantStorage = in.FaultTolerantStorage
	in.RuntimeStrategy.DeepCopyInto(&out.RuntimeStrategy)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EventSourceStrategy.
func (in *EventSourceStrategy) DeepCopy() *EventSourceStrategy {
	if in == nil {
		return nil
	}
	out := new(EventSourceStrategy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EventSpec) DeepCopyInto(out *EventSpec) {
	*out = *in
	in.Components.DeepCopyInto(&out.Components)
	in.Strategies.DeepCopyInto(&out.Strategies)
	in.RouterRules.DeepCopyInto(&out.RouterRules)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EventSpec.
func (in *EventSpec) DeepCopy() *EventSpec {
	if in == nil {
		return nil
	}
	out := new(EventSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EventStatus) DeepCopyInto(out *EventStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EventStatus.
func (in *EventStatus) DeepCopy() *EventStatus {
	if in == nil {
		return nil
	}
	out := new(EventStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EventStrategies) DeepCopyInto(out *EventStrategies) {
	*out = *in
	if in.SourceStrategies != nil {
		in, out := &in.SourceStrategies, &out.SourceStrategies
		*out = make([]EventSourceStrategy, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.TriggerStrategies != nil {
		in, out := &in.TriggerStrategies, &out.TriggerStrategies
		*out = make([]EventTriggerStrategy, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.DefaultSourceRuntimeStrategy.DeepCopyInto(&out.DefaultSourceRuntimeStrategy)
	in.DefaultTriggerRuntimeStrategy.DeepCopyInto(&out.DefaultTriggerRuntimeStrategy)
	out.DefaultDeadLetterStrategy = in.DefaultDeadLetterStrategy
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EventStrategies.
func (in *EventStrategies) DeepCopy() *EventStrategies {
	if in == nil {
		return nil
	}
	out := new(EventStrategies)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EventTrigger) DeepCopyInto(out *EventTrigger) {
	*out = *in
	out.Ref = in.Ref
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EventTrigger.
func (in *EventTrigger) DeepCopy() *EventTrigger {
	if in == nil {
		return nil
	}
	out := new(EventTrigger)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EventTriggerStrategy) DeepCopyInto(out *EventTriggerStrategy) {
	*out = *in
	in.RuntimeStrategy.DeepCopyInto(&out.RuntimeStrategy)
	out.DeadLetterStrategy = in.DeadLetterStrategy
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EventTriggerStrategy.
func (in *EventTriggerStrategy) DeepCopy() *EventTriggerStrategy {
	if in == nil {
		return nil
	}
	out := new(EventTriggerStrategy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Persistence) DeepCopyInto(out *Persistence) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Persistence.
func (in *Persistence) DeepCopy() *Persistence {
	if in == nil {
		return nil
	}
	out := new(Persistence)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PersistenceAddress) DeepCopyInto(out *PersistenceAddress) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PersistenceAddress.
func (in *PersistenceAddress) DeepCopy() *PersistenceAddress {
	if in == nil {
		return nil
	}
	out := new(PersistenceAddress)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RetryRule) DeepCopyInto(out *RetryRule) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RetryRule.
func (in *RetryRule) DeepCopy() *RetryRule {
	if in == nil {
		return nil
	}
	out := new(RetryRule)
	in.DeepCopyInto(out)
	return out
}
