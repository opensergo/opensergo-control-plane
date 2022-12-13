package traffic

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
type TrafficRouter struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec TrafficRouterSpec `json:"spec,omitempty"`

	Status TrafficRouterStatus `json:"status,omitempty"`
}

// HttpRequestMatchRuleList contains a list of HttpRequestMatchRule.
// +kubebuilder:object:root=true
type TrafficRouterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TrafficRouter `json:"items"`
}

// HttpRequestMatchRuleStatus defines the observed state of HttpRequestMatchRule.
type TrafficRouterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

func init() {
	SchemeBuilder.Register(&TrafficRouter{}, &TrafficRouterList{})
}
