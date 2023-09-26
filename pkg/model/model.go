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

package model

import (
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	extension "github.com/envoyproxy/go-control-plane/envoy/service/extension/v3"
	trpb "github.com/opensergo/opensergo-control-plane/pkg/proto/transport/v1"
	"github.com/opensergo/opensergo-control-plane/pkg/util"
	"google.golang.org/protobuf/types/known/anypb"
)

// Users could control this variable to determine whether use
var GlobalBoolVariable bool = false

type ChooseRulesServer bool

// ClientIdentifier represents a unique identifier for an OpenSergo client.
type ClientIdentifier string

type OpenSergoTransportStream = trpb.OpenSergoUniversalTransportService_SubscribeConfigServer

type SubscribeRequestHandler func(ClientIdentifier, *trpb.SubscribeRequest, OpenSergoTransportStream) error

type SubscribeXDsRequestHandler func(*discovery.DiscoveryRequest, *XDSConnection) error

const ExtensionConfigType = "type.googleapis.com/envoy.config.core.v3.TypedExtensionConfig"

type DataEntirePushHandler func(namespace, app, kind string, dataWithVersion *trpb.DataWithVersion, status *trpb.Status, respId string) error

type XDSPushHandler func(namespace, app, kind string, rules []*anypb.Any) error

type DiscoveryStream = extension.ExtensionConfigDiscoveryService_StreamExtensionConfigsServer

// WatchedResource tracks an active DiscoveryRequest subscription.
type WatchedResource struct {
	// TypeUrl is copied from the DiscoveryRequest.TypeUrl that initiated watching this resource.
	TypeUrl string

	// ResourceNames tracks the list of resources that are actively watched.
	ResourceNames []string

	// NonceSent is the nonce sent in the last sent response. If it is equal with NonceAcked, the
	NonceSent string

	// NonceAcked is the last acked message.
	NonceAcked string
}

// ResourceDelta records the difference in requested resources by an XDS client
type ResourceDelta struct {
	// Subscribed indicates the client requested these additional resources
	Subscribed util.String
	// Unsubscribed indicates the client no longer requires these resources
	Unsubscribed util.String
}

type NamespacedApp struct {
	Namespace string
	App       string
}
