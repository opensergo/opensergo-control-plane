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
	trpb "github.com/opensergo/opensergo-control-plane/pkg/proto/transport/v1"
)

type NamespacedApp struct {
	Namespace string
	App       string
}

// ClientIdentifier represents a unique identifier for an OpenSergo client.
type ClientIdentifier string

type OpenSergoTransportStream = trpb.OpenSergoUniversalTransportService_SubscribeConfigServer

type SubscribeRequestHandler func(ClientIdentifier, *trpb.SubscribeRequest, OpenSergoTransportStream) error

type DataEntirePushHandler func(namespace, app, kind string, dataWithVersion *trpb.DataWithVersion, status *trpb.Status, respId string) error

type NotifyPluginHandler func(pluginName string, e any) error
