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

type CircuitBreakerStrategy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec CircuitBreakerStrategySpec `json:"spec,omitempty"`

	Status CircuitBreakerStrategyStatus `json:"status,omitempty"`
}

// CircuitBreakerStrategySpec defines the spec of CircuitBreakerStrategy.
type CircuitBreakerStrategySpec struct {
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Enum=SlowRequestRatio;ErrorRequestRatio
	// +kubebuilder:validation:Required
	Strategy string `json:"strategy"`

	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=^([1-9]\d?|100|0)%$
	TriggerRatio string `json:"triggerRatio"`

	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=^[1-9]\d*(s|ms|m|min|minute|h|d)$
	StatDuration string `json:"statDuration"`

	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=^[1-9]\d*(s|ms|m|min|minute|h|d)$
	RecoveryTimeout string `json:"recoveryTimeout"`

	// +kubebuilder:validation:Type=integer
	// +kubebuilder:validation:Format=int32
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Required
	MinRequestAmount int32 `json:"minRequestAmount"`

	SlowConditions SlowConditions `json:"slowConditions,omitempty"`

	ErrorConditions ErrorConditions `json:"errorConditions,omitempty"`
}

type SlowConditions struct {
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=^[1-9]\d*(s|ms|m|min|minute|h|d)$
	MaxAllowedRt string `json:"maxAllowedRt"`
}

type ErrorConditions struct {
}

// CircuitBreakerStrategyStatus defines the observed state of CircuitBreakerStrategy.
type CircuitBreakerStrategyStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// CircuitBreakerStrategyList contains a list of CircuitBreakerStrategy.
type CircuitBreakerStrategyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CircuitBreakerStrategy `json:"items"`
}

// +kubebuilder:rbac:groups=fault-tolerance.opensergo.io,resources=CircuitBreakerStrategy,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=fault-tolerance.opensergo.io,resources=CircuitBreakerStrategy/status,verbs=get;update;patch

func init() {
	SchemeBuilder.Register(&CircuitBreakerStrategy{}, &CircuitBreakerStrategyList{})
}
