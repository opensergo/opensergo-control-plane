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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// +kubebuilder:object:root=true

type RateLimitStrategy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec RateLimitStrategySpec `json:"spec,omitempty"`

	Status RateLimitStrategyStatus `json:"status,omitempty"`
}

const (
	RequestAmountMetricType string = "RequestAmount"

	LocalLimitMode  string = "Local"
	GlobalLimitMode string = "Global"
)

// RateLimitStrategySpec defines the spec of RateLimitStrategy.
type RateLimitStrategySpec struct {
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Enum=RequestAmount
	// +kubebuilder:validation:Required
	MetricType string `json:"metricType"`

	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Enum=Local;Global
	// +kubebuilder:validation:Required
	LimitMode string `json:"limitMode"`

	// +kubebuilder:validation:Type=integer
	// +kubebuilder:validation:Format=int64
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Required
	Threshold int64 `json:"threshold"`

	// +kubebuilder:validation:Type=integer
	// +kubebuilder:validation:Format=int32
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Required
	StatDurationSeconds int32 `json:"statDurationSeconds"`
}

// RateLimitStrategyStatus defines the observed state of RateLimitStrategy.
type RateLimitStrategyStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// RateLimitStrategyList contains a list of RateLimitStrategy.
type RateLimitStrategyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RateLimitStrategy `json:"items"`
}

// +kubebuilder:rbac:groups=fault-tolerance.opensergo.io,resources=RateLimitStrategy,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=fault-tolerance.opensergo.io,resources=RateLimitStrategy/status,verbs=get;update;patch

func init() {
	SchemeBuilder.Register(&RateLimitStrategy{}, &RateLimitStrategyList{})
}
