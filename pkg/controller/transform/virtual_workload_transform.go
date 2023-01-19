package controller

import (
	"encoding/json"

	clusterv3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	"github.com/opensergo/opensergo-control-plane/constant"
	"github.com/opensergo/opensergo-control-plane/pkg/api/v1alpha1/traffic"
)

func BuildRouteConfigurationByVirtualWorkload(cls *traffic.VirtualWorkload) *clusterv3.Cluster {
	// Does not need VirtualWorkload now
	return nil
}

func BuildUnstructuredDestinationRule(cls *traffic.VirtualWorkload) map[string]interface{} {
	crdMeta := map[string]interface{}{}
	b, err := json.Marshal(cls.ObjectMeta)
	if err == nil {
		_ = json.Unmarshal(b, &crdMeta)
	}
	return map[string]interface{}{
		constant.CRD_API_VERSION: constant.VERSION_V1_ALPHA3,
		constant.CRD_KIND:        constant.DESTINATION_RULE_KIND,
		constant.CRD_METADATA:    crdMeta,
		constant.CRD_NAME:        cls.Name,
		constant.CRD_SPEC: map[string]interface{}{
			constant.DESTINATION_RULE_EXPORT_TO:         cls.Spec.ExportTo,
			constant.DESTINATION_RULE_HOST:              cls.Spec.Host,
			constant.DESTINATION_RULE_SUBSETS:           cls.Spec.Subsets,
			constant.DESTINATION_RULE_TRAFFIC_POLICY:    cls.Spec.TrafficPolicy,
			constant.DESTINATION_RULE_WORKLOAD_SELECTOR: cls.Spec.WorkloadSelector,
		},
	}
}
