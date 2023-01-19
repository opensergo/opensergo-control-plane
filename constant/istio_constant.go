package constant

const (
	EXTENSION_ROUTE_FALL_BACK = "envoy.router.cluster_specifier_plugin.cluster_fallback"
	CRD_API_VERSION           = "apiVersion"
	CRD_KIND                  = "kind"
	CRD_METADATA              = "metadata"
	CRD_SPEC                  = "spec"
	CRD_NAME                  = "name"

	META_RESOURCE_VERSION = "resourceVersion"
	META_UID              = "uid"

	VERSION_V1_ALPHA3 = "networking.istio.io/v1alpha3"

	VIRTUAL_SERVICE_KIND       = "VirtualService"
	VIRTUAL_SERVICE_HOST       = "hosts"
	VIRTUAL_SERVICE_HTTP_MATCH = "http"

	DESTINATION_RULE_KIND              = "DestinationRule"
	DESTINATION_RULE_EXPORT_TO         = "exportTo"
	DESTINATION_RULE_HOST              = "host"
	DESTINATION_RULE_SUBSETS           = "subsets"
	DESTINATION_RULE_TRAFFIC_POLICY    = "trafficPolicy"
	DESTINATION_RULE_WORKLOAD_SELECTOR = "workloadSelector"
)
