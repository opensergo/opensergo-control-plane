package gvr

import "k8s.io/apimachinery/pkg/runtime/schema"

var VirtualServiceGVR = schema.GroupVersionResource{Group: "networking.istio.io", Version: "v1alpha3", Resource: "virtualservices"}
