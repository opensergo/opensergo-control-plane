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

package controller

import (
	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	routev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	matcherv3 "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/opensergo/opensergo-control-plane/constant"
	"github.com/opensergo/opensergo-control-plane/pkg/api/v1alpha1/traffic"
	route "github.com/opensergo/opensergo-control-plane/pkg/proto/router/v1"
	"github.com/opensergo/opensergo-control-plane/pkg/util"
	"k8s.io/apimachinery/pkg/util/json"
)

// BuildRouteConfiguration for Istio RouteConfiguration
func BuildRouteConfiguration(cls *traffic.TrafficRouter) *routev3.RouteConfiguration {
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

func BuildUnstructuredVirtualService(cls *traffic.TrafficRouter) map[string]interface{} {
	crdMeta := map[string]interface{}{}
	b, err := json.Marshal(cls.ObjectMeta)
	if err == nil {
		_ = json.Unmarshal(b, &crdMeta)
	}
	return map[string]interface{}{
		constant.CRD_API_VERSION: constant.VIRTUAL_SERVICE_V1_ALPHA3,
		constant.CRD_KIND:        constant.VIRTUAL_SERVICE_KIND,
		constant.CRD_METADATA:    crdMeta,
		constant.CRD_NAME:        cls.Name,
		constant.CRD_SPEC: map[string]interface{}{
			constant.VIRTUAL_SERVICE_HOST:       cls.Spec.Hosts,
			constant.VIRTUAL_SERVICE_HTTP_MATCH: cls.Spec.Http,
		},
	}
}

func buildHTTPRoutes(tr *traffic.TrafficRouter) []*routev3.Route {
	var routes []*routev3.Route
	for _, httpRoute := range tr.Spec.Http {
		r := &routev3.Route{
			Match: &routev3.RouteMatch{
				Headers:         buildHeaderMatchers(httpRoute.Match),
				QueryParameters: buildParamMatchers(httpRoute.Match),
			},
			Action: &routev3.Route_Route{
				Route: buildRouteAction(httpRoute, tr),
			},
		}
		routes = append(routes, r)
	}
	return routes
}

func buildUnweightedRouteAction(destination *traffic.HTTPRouteDestination, vs *traffic.TrafficRouter) *routev3.RouteAction {
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

func buildWeightedRouteAction(destinations []*traffic.HTTPRouteDestination, vs *traffic.TrafficRouter) *routev3.RouteAction {
	return &routev3.RouteAction{
		ClusterSpecifier: &routev3.RouteAction_WeightedClusters{
			WeightedClusters: &routev3.WeightedCluster{
				Clusters: buildWeightedClusters(vs.Namespace, destinations),
			},
		},
	}
}

func buildWeightedClusters(namespace string, destinations []*traffic.HTTPRouteDestination) []*routev3.WeightedCluster_ClusterWeight {
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

func buildRouteAction(httpRoute *traffic.HTTPRoute, vs *traffic.TrafficRouter) *routev3.RouteAction {
	if len(httpRoute.Route) == 1 {
		// unweighted
		return buildUnweightedRouteAction(httpRoute.Route[0], vs)
	} else {
		// weighted
		return buildWeightedRouteAction(httpRoute.Route, vs)
	}
}

func buildRouterFallbackActionCluster(fallback *traffic.Fallback, namespace string) string {
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
			Name:        constant.EXTENSION_ROUTE_FALL_BACK,
			TypedConfig: util.MessageToAny(config),
		},
	}
}

func buildRouteActionCluster(serviceName, namespace, version string) string {
	return "outbound|" + "|" + version + "|" + buildFQDN(serviceName, namespace)
}

func buildParamMatchers(matches []*traffic.HTTPMatchRequest) []*routev3.QueryParameterMatcher {
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

func buildHeaderMatchers(matches []*traffic.HTTPMatchRequest) []*routev3.HeaderMatcher {
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
