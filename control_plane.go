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

package opensergo

import (
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/opensergo/opensergo-control-plane/pkg/util"
	"google.golang.org/protobuf/types/known/anypb"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/opensergo/opensergo-control-plane/pkg/controller"
	"github.com/opensergo/opensergo-control-plane/pkg/model"
	trpb "github.com/opensergo/opensergo-control-plane/pkg/proto/transport/v1"
	"github.com/opensergo/opensergo-control-plane/pkg/transport/grpc"
	transport "github.com/opensergo/opensergo-control-plane/pkg/transport/grpc"
	"github.com/pkg/errors"
)

const delimiter = "/"

type ControlPlane struct {
	operator  *controller.KubernetesOperator
	server    *transport.Server
	xdsServer *transport.DiscoveryServer
	protoDesc *trpb.ControlPlaneDesc

	mux sync.RWMutex
}

func NewControlPlane() (*ControlPlane, error) {
	cp := &ControlPlane{}

	operator, err := controller.NewKubernetesOperator(cp.sendMessage, cp.pushXds)
	if err != nil {
		return nil, err
	}

	cp.server = transport.NewServer(uint32(10246), []model.SubscribeRequestHandler{cp.handleSubscribeRequest})
	cp.xdsServer = transport.NewDiscoveryServer(uint32(8002), []model.SubscribeXDsRequestHandler{cp.handleXDSSubscribeRequest})
	cp.operator = operator

	hostname, herr := os.Hostname()
	if herr != nil {
		// TODO: log here
		hostname = "unknown-host"
	}
	cp.protoDesc = &trpb.ControlPlaneDesc{Identifier: "osg-" + hostname}

	return cp, nil
}

func (c *ControlPlane) Start() error {
	// Run the Kubernetes operator
	err := c.operator.Run()

	if err != nil {
		return err
	}

	if model.GlobalBoolVariable {
		//Run the transport server
		err = c.server.Run()
		if err != nil {
			return err
		}
	} else {
		//Run the xDS Server
		err = c.xdsServer.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *ControlPlane) sendMessage(namespace, app, kind string, dataWithVersion *trpb.DataWithVersion, status *trpb.Status, respId string) error {
	connections, exists := c.server.ConnectionManager().Get(namespace, app, kind)
	if !exists || connections == nil {
		return errors.New("There is no connection for this kind")
	}
	for _, connection := range connections {
		if connection == nil || !connection.IsValid() {
			// TODO: log.Debug
			continue
		}
		err := c.sendMessageToStream(connection.Stream(), namespace, app, kind, dataWithVersion, status, respId)
		if err != nil {
			// TODO: should not short-break here. Handle partial failure here.
			return err
		}
	}
	return nil
}

func (c *ControlPlane) sendMessageToStream(stream model.OpenSergoTransportStream, namespace, app, kind string, dataWithVersion *trpb.DataWithVersion, status *trpb.Status, respId string) error {
	if stream == nil {
		return nil
	}
	return stream.SendMsg(&trpb.SubscribeResponse{
		Status:          status,
		Ack:             "",
		Namespace:       namespace,
		App:             app,
		Kind:            kind,
		DataWithVersion: dataWithVersion,
		ControlPlane:    c.protoDesc,
		ResponseId:      respId,
	})
}

// handleXDSSubscribeRequest handles the XDS subscription request.
func (c *ControlPlane) handleXDSSubscribeRequest(req *discovery.DiscoveryRequest, con *model.XDSConnection) error {
	// Check if the request is for ExtensionConfigType.
	if req.TypeUrl != model.ExtensionConfigType {
		return nil
	}

	// Determine whether to respond and calculate the delta.
	shouldRespond, delta := grpc.ShouldRespond(con, req)
	subscribed := delta.Subscribed
	unsubscribed := delta.Unsubscribed

	if !shouldRespond {
		// If there's no need to respond, return early.
		return nil
	}

	if len(subscribed) != 0 {
		var rules []*anypb.Any
		for resourcename := range subscribed {
			// Split the resource name into its components.
			request := strings.Split(resourcename, delimiter)

			// Register a watcher for the specified resource.
			crdWatcher, err := c.operator.RegisterWatcher(model.SubscribeTarget{
				Namespace: request[4],
				AppName:   request[3],
				Kind:      request[0] + delimiter + request[1] + delimiter + request[2],
			})

			if err != nil {
				// Log the error and continue to the next resource.
				log.Printf("Error registering watcher for resource %s: %s\n", resourcename, err.Error())
				continue
			}

			// Add the connection to the connection map.
			c.xdsServer.AddConnectionToMap(request[4], request[3], request[0]+"/"+request[1]+"/"+request[2], con)

			// Get the current rules for the resource.
			curRules, _ := crdWatcher.GetRules(model.NamespacedApp{
				Namespace: request[4],
				App:       request[3],
			})

			if len(curRules) > 0 {
				// Append the current rules to the rules slice.
				rules = append(rules, curRules...)
			}
		}

		// Push XDS rules to the connection.
		err := c.pushXdsToStream(con, con.Watched(req.TypeUrl), rules)
		if err != nil {
			// Log the error if pushing XDS rules fails.
			log.Printf("Failed to push XDS rules to connection: %s\n", err.Error())
		}
	}

	if len(unsubscribed) != 0 {
		// Handle unsubscribed resources.
		for resourcename := range unsubscribed {
			// Split the resource name into its components.
			request := strings.Split(resourcename, delimiter)

			// Remove the connection from the connection map.
			c.xdsServer.RemoveConnectionFromMap(model.NamespacedApp{request[0], request[1]}, request[2], con.Identifier)
		}
	}

	return nil
}

func (c *ControlPlane) handleSubscribeRequest(clientIdentifier model.ClientIdentifier, request *trpb.SubscribeRequest, stream model.OpenSergoTransportStream) error {
	for _, kind := range request.Target.Kinds {
		crdWatcher, err := c.operator.RegisterWatcher(model.SubscribeTarget{
			Namespace: request.Target.Namespace,
			AppName:   request.Target.App,
			Kind:      kind,
		})
		if err != nil {
			status := &trpb.Status{
				Code:    transport.RegisterWatcherError,
				Message: "Register watcher error",
				Details: nil,
			}
			err = c.sendMessageToStream(stream, request.Target.Namespace, request.Target.App, kind, nil, status, request.RequestId)
			if err != nil {
				// TODO: log here
				log.Printf("sendMessageToStream failed, err=%s\n", err.Error())
			}
			continue
		}
		curConnection := transport.NewConnection(clientIdentifier, stream)
		_ = c.server.ConnectionManager().Add(request.Target.Namespace, request.Target.App, kind, curConnection, curConnection.Identifier())
		// send if the watcher cache is not empty
		rules, version := crdWatcher.GetRules(model.NamespacedApp{
			Namespace: request.Target.Namespace,
			App:       request.Target.App,
		})
		if len(rules) > 0 {
			status := &trpb.Status{
				Code:    transport.Success,
				Message: "Get and send rule success",
				Details: nil,
			}
			dataWithVersion := &trpb.DataWithVersion{
				Data:    rules,
				Version: version,
			}
			err = c.sendMessageToStream(stream, request.Target.Namespace, request.Target.App, kind, dataWithVersion, status, request.RequestId)
			if err != nil {
				// TODO: log here
				log.Printf("sendMessageToStream failed, err=%s\n", err.Error())
			}
		}
	}
	return nil
}

func (c *ControlPlane) pushXdsToStream(con *model.XDSConnection, w *model.WatchedResource, rules []*anypb.Any) error {
	res := &discovery.DiscoveryResponse{
		TypeUrl:     w.TypeUrl,
		VersionInfo: c.xdsServer.NextVersion(),

		// TODO: RECORD THE NONCE AND CHECK THE NONCE
		Nonce:     util.Nonce(),
		Resources: rules,
	}

	// Set nonce in the XDSConnection's WatchedResource
	con.Lock()
	if con.WatchedResources[model.ExtensionConfigType] == nil {
		con.WatchedResources[res.TypeUrl] = &model.WatchedResource{TypeUrl: res.TypeUrl}
	}
	con.WatchedResources[res.TypeUrl].NonceSent = res.Nonce
	con.Unlock()

	// Send the DiscoveryResponse over the stream
	err := con.Stream.Send(res)
	if err != nil {
		// Handle the error, e.g., log it or return it
		// TODO: You can log the error or handle it as needed.
		log.Println("Failed to send DiscoveryResponse:", err)
		return err
	}

	return nil
}

func (c *ControlPlane) pushXds(namespace, app, kind string, rules []*anypb.Any) error {
	// Retrieve the XDS connections for the specified namespace, app, and kind.
	connections, exists := c.xdsServer.XDSConnectionManager.Get(namespace, app, kind)
	if !exists || connections == nil {
		// Log that there is no connection for this kind.
		// Replace this with your actual logging mechanism.
		log.Println("No XDS connection found for namespace:", namespace, "app:", app, "kind:", kind)
		return errors.New("There is no connection for this kind")
	}

	for _, connection := range connections {
		if connection == nil {
			// Log a debug message for a nil connection.
			// Replace this with your actual logging mechanism.
			log.Println("Encountered a nil XDS connection")
			continue
		}
		err := c.pushXdsToStream(connection, connection.WatchedResources[model.ExtensionConfigType], rules)
		if err != nil {
			// Log an error and return it if there is an error pushing XDS rules.
			// Replace this with your actual logging mechanism.
			log.Println("Failed to push XDS rules to connection:", err)
			// TODO: You might want to consider handling partial failures here.
			return err
		}
	}

	// Return nil to indicate success.
	return nil
}
