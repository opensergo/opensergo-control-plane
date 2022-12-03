package networking

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
type VirtualService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec VirtualServiceSpec `json:"spec,omitempty"`

	Status VirtualServiceStatus `json:"status,omitempty"`
}

// HttpRequestMatchRuleList contains a list of HttpRequestMatchRule.
// +kubebuilder:object:root=true
type VirtualServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualService `json:"items"`
}

// HttpRequestMatchRuleStatus defines the observed state of HttpRequestMatchRule.
type VirtualServiceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

func init() {
	SchemeBuilder.Register(&VirtualService{}, &VirtualServiceList{})
}
