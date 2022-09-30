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

type ConcurrencyLimitStrategy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ConcurrencyLimitStrategySpec `json:"spec,omitempty"`

	Status ConcurrencyLimitStrategyStatus `json:"status,omitempty"`
}

// ConcurrencyLimitStrategySpec defines the spec of ConcurrencyLimitStrategy.
type ConcurrencyLimitStrategySpec struct {

	// +kubebuilder:validation:Type=integer
	// +kubebuilder:validation:Format=int64
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Required
	MaxConcurrencyThreshold int64 `json:"maxConcurrency"`

	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Enum=Local;Global
	// +kubebuilder:validation:Required
	LimitMode string `json:"limitMode"`
}

// ConcurrencyLimitStrategyStatus defines the observed state of ConcurrencyLimitStrategy.
type ConcurrencyLimitStrategyStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// ConcurrencyLimitStrategyList contains a list of ConcurrencyLimitStrategy.
type ConcurrencyLimitStrategyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ConcurrencyLimitStrategy `json:"items"`
}

// +kubebuilder:rbac:groups=fault-tolerance.opensergo.io,resources=ConcurrencyLimitStrategy,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=fault-tolerance.opensergo.io,resources=ConcurrencyLimitStrategy/status,verbs=get;update;patch

func init() {
	SchemeBuilder.Register(&ConcurrencyLimitStrategy{}, &ConcurrencyLimitStrategyList{})
}
