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
	"log"
	"os"
	"sync"

	"github.com/alibaba/sentinel-golang/util"

	"github.com/opensergo/opensergo-control-plane/pkg/controller"
	"github.com/opensergo/opensergo-control-plane/pkg/model"
	trpb "github.com/opensergo/opensergo-control-plane/pkg/proto/transport/v1"
	transport "github.com/opensergo/opensergo-control-plane/pkg/transport/grpc"
)

type ControlPlane struct {
	operator     *controller.KubernetesOperator
	server       *transport.Server
	secureServer *transport.Server

	protoDesc *trpb.ControlPlaneDesc

	mux sync.RWMutex
	ch  chan error
}

func NewControlPlane() (*ControlPlane, error) {
	cp := &ControlPlane{}

	operator, err := controller.NewKubernetesOperator(cp.sendMessage)
	if err != nil {
		return nil, err
	}

	cp.server = transport.NewServer(uint32(10246), []model.SubscribeRequestHandler{cp.handleSubscribeRequest})
	// On port 10248, it can use tls transport
	cp.secureServer = transport.NewSecureServer(uint32(10248), []model.SubscribeRequestHandler{cp.handleSubscribeRequest})
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

	go util.RunWithRecover(func() {
		// Run the transport server
		log.Println("Starting grpc server on port 10246!")
		err = c.server.Run()
		if err != nil {
			c.ch <- err
			log.Fatal("Failed to run the grpc server")
		}
	})

	go util.RunWithRecover(func() {
		// Run the secure transport server
		log.Println("Starting secure grpc server on port 10248!")
		err = c.secureServer.Run()
		if err != nil {
			c.ch <- err
			log.Fatal("Failed to run the secure grpc server")
		}
	})
	err = <-c.ch
	return err
}

func (c *ControlPlane) sendMessage(namespace, app, kind string, dataWithVersion *trpb.DataWithVersion, status *trpb.Status, respId string) error {
	var connections []*transport.Connection
	var exists bool
	scs, exists := c.secureServer.ConnectionManager().Get(namespace, app, kind)
	if !exists || connections == nil {
		log.Printf("There is no secure connection for app %s kind %s in ns %s", app, kind, namespace)
	} else {
		connections = append(connections, scs...)
	}
	cs, exists := c.server.ConnectionManager().Get(namespace, app, kind)
	if !exists || connections == nil {
		log.Printf("There is no connection for app %s kind %s in ns %s", app, kind, namespace)
	} else {
		connections = append(connections, cs...)
	}
	return c.innerSendMessage(namespace, app, kind, dataWithVersion, status, respId, connections)
}

func (c *ControlPlane) innerSendMessage(namespace, app, kind string, dataWithVersion *trpb.DataWithVersion, status *trpb.Status, respId string, connections []*transport.Connection) error {
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

func (c *ControlPlane) handleSubscribeRequest(clientIdentifier model.ClientIdentifier, request *trpb.SubscribeRequest, stream model.OpenSergoTransportStream, isSecure bool) error {
	for _, kind := range request.Target.Kinds {
		crdWatcher, err := c.operator.RegisterWatcher(model.SubscribeTarget{
			Namespace: request.Target.Namespace,
			AppName:   request.Target.App,
			Kind:      kind,
		}, isSecure)
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

		if isSecure {
			_ = c.secureServer.ConnectionManager().Add(request.Target.Namespace, request.Target.App, kind, transport.NewConnection(clientIdentifier, stream))
		} else {
			_ = c.server.ConnectionManager().Add(request.Target.Namespace, request.Target.App, kind, transport.NewConnection(clientIdentifier, stream))
		}

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
