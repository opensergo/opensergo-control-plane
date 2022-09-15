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

type ThrottlingStrategy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ThrottlingStrategySpec `json:"spec,omitempty"`

	Status ThrottlingStrategyStatus `json:"status,omitempty"`
}

// ThrottlingStrategySpec defines the spec of ThrottlingStrategy.
type ThrottlingStrategySpec struct {
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=^\d+(s|ms|m|min|minute|h|d)$
	MinIntervalOfRequests string `json:"minIntervalOfRequests"`

	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=^[1-9]\d*(s|ms|m|min|minute|h|d)$
	QueueTimeout string `json:"queueTimeout"`
}

// ThrottlingStrategyStatus defines the observed state of ThrottlingStrategy.
type ThrottlingStrategyStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// ThrottlingStrategyList contains a list of ThrottlingStrategy.
type ThrottlingStrategyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ThrottlingStrategy `json:"items"`
}

// +kubebuilder:rbac:groups=fault-tolerance.opensergo.io,resources=ThrottlingStrategy,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=fault-tolerance.opensergo.io,resources=ThrottlingStrategy/status,verbs=get;update;patch

func init() {
	SchemeBuilder.Register(&ThrottlingStrategy{}, &ThrottlingStrategyList{})
}
