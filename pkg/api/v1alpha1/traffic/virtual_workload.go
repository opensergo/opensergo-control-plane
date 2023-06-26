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

package traffic

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
type VirtualWorkload struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec VirtualWorkloadSpec `json:"spec,omitempty"`

	Status VirtualWorkloadStatus `json:"status,omitempty"`
}

// HttpRequestMatchRuleList contains a list of HttpRequestMatchRule.
// +kubebuilder:object:root=true
type VirtualWorkloadList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualWorkload `json:"items"`
}

// HttpRequestMatchRuleStatus defines the observed state of HttpRequestMatchRule.
type VirtualWorkloadStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

func init() {
	SchemeBuilder.Register(&VirtualWorkload{}, &VirtualWorkloadList{})
}
