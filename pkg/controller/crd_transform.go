package controller

import (
	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	routev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	matcherv3 "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	"github.com/golang/protobuf/ptypes/wrappers"
	crdv1beta1 "github.com/opensergo/opensergo-control-plane/pkg/api/v1beta1/networking"
	route "github.com/opensergo/opensergo-control-plane/pkg/proto/router/v1"
	"github.com/opensergo/opensergo-control-plane/pkg/util"
)

// BuildRouteConfiguration for Istio RouteConfiguration
func BuildRouteConfiguration(cls *crdv1beta1.VirtualService) *routev3.RouteConfiguration {
	virtualHost := &routev3.VirtualHost{
		Name:   cls.Name,
		Routes: []*routev3.Route{},
	}
	virtualHost.Routes = append(virtualHost.Routes, buildHTTPRoutes(cls)...)
	for _, domain := range cls.Spec.Hosts {
		virtualHost.Domains = append(virtualHost.Domains, buildFQDN(domain, cls.Namespace))
	}
	rule := &routev3.RouteConfiguration{
		Name:         cls.Name,
		VirtualHosts: []*routev3.VirtualHost{virtualHost},
	}
	return rule
}

func buildHTTPRoutes(vs *crdv1beta1.VirtualService) []*routev3.Route {
	var routes []*routev3.Route
	for _, httpRoute := range vs.Spec.Http {
		r := &routev3.Route{
			Match: &routev3.RouteMatch{
				Headers:         buildHeaderMatchers(httpRoute.Match),
				QueryParameters: buildParamMatchers(httpRoute.Match),
			},
			Action: &routev3.Route_Route{
				Route: buildRouteAction(httpRoute, vs),
			},
		}
		routes = append(routes, r)
	}
	return routes
}

func buildUnweightedRouteAction(destination *crdv1beta1.HTTPRouteDestination, vs *crdv1beta1.VirtualService) *routev3.RouteAction {
	if destination.Destination.Fallback != nil {
		return &routev3.RouteAction{
			ClusterSpecifier: &routev3.RouteAction_InlineClusterSpecifierPlugin{
				InlineClusterSpecifierPlugin: buildClusterSpecifierPlugin(true, buildClusterFallbackConfig(buildRouteActionCluster(destination.Destination.Host, vs.Namespace, destination.Destination.Subset), buildRouterFallbackActionCluster(destination.Destination.Fallback, vs.Namespace))),
			},
		}
	} else {
		return &routev3.RouteAction{
			ClusterSpecifier: &routev3.RouteAction_Cluster{Cluster: buildRouteActionCluster(destination.Destination.Host, vs.Namespace, destination.Destination.Subset)},
		}
	}
}

func buildWeightedRouteAction(destinations []*crdv1beta1.HTTPRouteDestination, vs *crdv1beta1.VirtualService) *routev3.RouteAction {
	return &routev3.RouteAction{
		ClusterSpecifier: &routev3.RouteAction_WeightedClusters{
			WeightedClusters: &routev3.WeightedCluster{
				Clusters: buildWeightedClusters(vs.Namespace, destinations),
			},
		},
	}
}

func buildWeightedClusters(namespace string, destinations []*crdv1beta1.HTTPRouteDestination) []*routev3.WeightedCluster_ClusterWeight {
	var weightedClusters []*routev3.WeightedCluster_ClusterWeight
	for _, destination := range destinations {
		w := &routev3.WeightedCluster_ClusterWeight{
			Name:   buildRouteActionCluster(destination.Destination.Host, namespace, destination.Destination.Subset),
			Weight: &wrappers.UInt32Value{Value: uint32(destination.Weight)},
		}
		weightedClusters = append(weightedClusters, w)
	}
	return weightedClusters
}

func buildRouteAction(httpRoute *crdv1beta1.HTTPRoute, vs *crdv1beta1.VirtualService) *routev3.RouteAction {
	if len(httpRoute.Route) == 1 {
		// unweighted
		return buildUnweightedRouteAction(httpRoute.Route[0], vs)
	} else {
		// weighted
		return buildWeightedRouteAction(httpRoute.Route, vs)
	}
}

func buildRouterFallbackActionCluster(fallback *crdv1beta1.Fallback, namespace string) string {
	if fallback == nil {
		return ""
	}

	return buildRouteActionCluster(fallback.Host, namespace, fallback.Subset)
}

func buildClusterFallbackConfig(cluster string, fallbackCluster string) *route.ClusterFallbackConfig_ClusterConfig {
	return &route.ClusterFallbackConfig_ClusterConfig{
		RoutingCluster:  cluster,
		FallbackCluster: fallbackCluster,
	}
}

func buildClusterSpecifierPlugin(isSupport bool, config *route.ClusterFallbackConfig_ClusterConfig) *routev3.ClusterSpecifierPlugin {
	if !isSupport || config == nil {
		return nil
	}

	return &routev3.ClusterSpecifierPlugin{
		Extension: &corev3.TypedExtensionConfig{
			Name:        EXTENSION_ROUTE_FALL_BACK,
			TypedConfig: util.MessageToAny(config),
		},
	}
}

func buildRouteActionCluster(serviceName, namespace, version string) string {
	return "outbound|" + "|" + version + "|" + buildFQDN(serviceName, namespace)
}

func buildParamMatchers(matches []*crdv1beta1.HTTPMatchRequest) []*routev3.QueryParameterMatcher {
	var queryParamMatchers []*routev3.QueryParameterMatcher
	for _, match := range matches {
		for _, matcher := range match.QueryParams {
			queryMatcher := &routev3.QueryParameterMatcher{}
			if matcher.GetRegex() != "" {
				queryMatcher.QueryParameterMatchSpecifier = &routev3.QueryParameterMatcher_StringMatch{
					StringMatch: &matcherv3.StringMatcher{
						MatchPattern: &matcherv3.StringMatcher_SafeRegex{
							SafeRegex: &matcherv3.RegexMatcher{
								Regex: matcher.GetRegex(),
							},
						},
					}}
			}
			if matcher.GetPrefix() != "" {
				queryMatcher.QueryParameterMatchSpecifier = &routev3.QueryParameterMatcher_StringMatch{
					StringMatch: &matcherv3.StringMatcher{
						MatchPattern: &matcherv3.StringMatcher_Prefix{
							Prefix: matcher.GetPrefix(),
						},
					}}
			}
			if matcher.GetExact() != "" {
				queryMatcher.QueryParameterMatchSpecifier = &routev3.QueryParameterMatcher_StringMatch{
					StringMatch: &matcherv3.StringMatcher{
						MatchPattern: &matcherv3.StringMatcher_Exact{
							Exact: matcher.GetExact(),
						},
					}}
			}
			queryParamMatchers = append(queryParamMatchers, queryMatcher)
		}
	}
	return queryParamMatchers
}

func buildHeaderMatchers(matches []*crdv1beta1.HTTPMatchRequest) []*routev3.HeaderMatcher {
	var headerMatchers []*routev3.HeaderMatcher
	for _, match := range matches {
		for key, matcher := range match.Headers {
			headerMatcher := &routev3.HeaderMatcher{
				Name: key,
			}
			if matcher.GetRegex() != "" {
				headerMatcher.HeaderMatchSpecifier = &routev3.HeaderMatcher_SafeRegexMatch{SafeRegexMatch: &matcherv3.RegexMatcher{Regex: matcher.GetRegex()}}
			}
			if matcher.GetPrefix() != "" {
				headerMatcher.HeaderMatchSpecifier = &routev3.HeaderMatcher_PrefixMatch{PrefixMatch: matcher.GetPrefix()}
			}
			if matcher.GetExact() != "" {
				headerMatcher.HeaderMatchSpecifier = &routev3.HeaderMatcher_ExactMatch{ExactMatch: matcher.GetExact()}
			}
			headerMatchers = append(headerMatchers, headerMatcher)
		}
	}
	return headerMatchers
}

func buildFQDN(serviceName, namespace string) string {
	return serviceName + "." + namespace + ".svc.cluster.local"
}
