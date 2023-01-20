package controller

import (
	"encoding/json"
	"log"

	clusterv3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	"github.com/opensergo/opensergo-control-plane/constant"
	"github.com/opensergo/opensergo-control-plane/pkg/api/v1alpha1/traffic"
)

func BuildClusterByVirtualWorkload(cls *traffic.VirtualWorkload) []*clusterv3.Cluster {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Error on build Cluster By VirtualWorkload: %v", err)
		}
	}()
	// Only support subset and lb policy now
	// TODO: ConnectionPool, OutlierDetection, TLSSettings, LocalityLB, Consistent Hash and so on
	var clusters []*clusterv3.Cluster
	for _, subset := range cls.Spec.Subsets {
		clusters = append(clusters, &clusterv3.Cluster{
			Name:     buildSubsetName(cls.Spec.Host, subset.Name),
			LbPolicy: buildLbPolicy(subset.TrafficPolicy),
		})
	}
	clusters = append(clusters, &clusterv3.Cluster{
		Name:     buildSubsetName(cls.Spec.Host, ""),
		LbPolicy: buildLbPolicy(cls.Spec.TrafficPolicy),
	})
	return clusters
}

func buildLbPolicy(trafficPolicy *traffic.TrafficPolicy) clusterv3.Cluster_LbPolicy {
	if trafficPolicy == nil {
		return clusterv3.Cluster_ROUND_ROBIN
	}
	sample := trafficPolicy.LoadBalancer.Simple

	switch sample {
	case traffic.LoadBalancerSettings_LEAST_REQUEST:
		return clusterv3.Cluster_LEAST_REQUEST
	case traffic.LoadBalancerSettings_RANDOM:
		return clusterv3.Cluster_RANDOM
	case traffic.LoadBalancerSettings_ROUND_ROBIN:
		return clusterv3.Cluster_ROUND_ROBIN
	case traffic.LoadBalancerSettings_PASSTHROUGH:
		return clusterv3.Cluster_CLUSTER_PROVIDED
	}

	consistentHash := trafficPolicy.LoadBalancer.ConsistentHash
	alg := consistentHash.HashAlgorithm

	if _, ok := alg.(*traffic.LoadBalancerSettings_ConsistentHashLB_Maglev); ok {
		return clusterv3.Cluster_MAGLEV
	}

	if _, ok := alg.(*traffic.LoadBalancerSettings_ConsistentHashLB_RingHash_); ok {
		return clusterv3.Cluster_RING_HASH
	}

	// Default algorithm is ROUND_ROBIN
	return clusterv3.Cluster_ROUND_ROBIN
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
