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
	"strconv"
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
	cp.xdsServer = transport.NewDiscoveryServer(uint32(8002), []model.SubscribeXDsRequestHandler{cp.handleXDsSubscribeRequest})
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
	// Run the transport server
	err = c.server.Run()
	if err != nil {
		return err
	}
	err = c.xdsServer.Run()
	if err != nil {
		return err
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

// cxz
func (c *ControlPlane) handleXDsSubscribeRequest(req *discovery.DiscoveryRequest, con *model.XDsConnection) error {
	if req.TypeUrl != model.ExtensionConfigType {
		return nil
	}
	shouldRespond, delta := grpc.ShouldRespond(con, req)

	subscribed := delta.Subscribed
	unsubscribed := delta.Unsubscribed
	if !shouldRespond {
		return nil
	}
	if len(subscribed) != 0 {
		for resourcename := range subscribed {
			request := strings.Split(resourcename, delimiter)
			crdWatcher, err := c.operator.RegisterWatcher(model.SubscribeTarget{
				Namespace: request[0],
				AppName:   request[1],
				Kind:      request[2],
			})
			// TODO: unhandled err
			if err != nil {
				continue
			}

			c.xdsServer.AddConnectioonToMap(request[0], request[1], request[2], con)

			rules, version := crdWatcher.GetRules(model.NamespacedApp{
				Namespace: request[0],
				App:       request[1],
			})
			if len(rules) > 0 {
				err := c.pushXdsToStream(con, con.Watched(req.TypeUrl), version, rules)
				if err != nil {
					// TODO: log here
					log.Printf("sendMessageToStream failed, err=%s\n", err.Error())
				}
			}

		}
	}

	if len(unsubscribed) != 0 {
		for resourcename := range subscribed {
			request := strings.Split(resourcename, delimiter)
			c.xdsServer.RemoveConnectionFromMap(model.NamespacedApp{request[0], request[1]}, request[2], con.Identifier)
		}
	}

	return nil
}

func (c *ControlPlane) handleSubscribeRequest(clientIdentifier model.ClientIdentifier, request *trpb.SubscribeRequest, stream model.OpenSergoTransportStream) error {
	// var labels []model.LabelKV
	// if request.Target.Labels != nil {
	//	for _, label := range request.Target.Labels {
	//		labels = append(labels, model.LabelKV{
	//			Key:   label.Key,
	//			Value: label.Value,
	//		})
	//	}
	// }
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

// cxz
func (c *ControlPlane) pushXdsToStream(con *model.XDsConnection, w *model.WatchedResource, version int64, rules []*anypb.Any) error {

	res := &discovery.DiscoveryResponse{
		TypeUrl:     w.TypeUrl,
		VersionInfo: strconv.FormatInt(version, 10),

		// TODO: RECORD THE NONCE AND CHECK THE NONCE
		Nonce:     util.Nonce(),
		Resources: rules,
	}
	// set nonce
	con.Lock()
	if con.WatchedResources[model.ExtensionConfigType] == nil {
		con.WatchedResources[res.TypeUrl] = &model.WatchedResource{TypeUrl: res.TypeUrl}
	}
	con.WatchedResources[res.TypeUrl].NonceSent = res.Nonce
	con.Unlock()

	return con.Stream.Send(res)
}

func (c *ControlPlane) pushXds(namespace, app, kind string, rules []*anypb.Any, version int64) error {
	connections, exists := c.xdsServer.XDSConnectionManeger.Get(namespace, app, kind)
	if !exists || connections == nil {
		return errors.New("There is no connection for this kind")
	}

	for _, connection := range connections {
		if connection == nil {
			// TODO: log.Debug
			continue
		}
		err := c.pushXdsToStream(connection, connection.WatchedResources[model.ExtensionConfigType], version, rules)
		if err != nil {
			// TODO: should not short-break here. Handle partial failure here.
			return err
		}
	}

	return nil
}
