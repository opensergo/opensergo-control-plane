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
	stream_plugin "github.com/opensergo/opensergo-control-plane/pkg/plugin/pl/builtin/stream"
	"log"
	"os"
	"sync"

	"github.com/opensergo/opensergo-control-plane/pkg/controller"
	"github.com/opensergo/opensergo-control-plane/pkg/model"
	trpb "github.com/opensergo/opensergo-control-plane/pkg/proto/transport/v1"
	transport "github.com/opensergo/opensergo-control-plane/pkg/transport/grpc"
	"github.com/pkg/errors"
)

type ControlPlane struct {
	operator *controller.KubernetesOperator
	server   *transport.Server

	protoDesc *trpb.ControlPlaneDesc

	mux sync.RWMutex
}

func NewControlPlane() (*ControlPlane, error) {
	cp := &ControlPlane{}

	operator, err := controller.NewKubernetesOperator(cp.sendMessage)
	if err != nil {
		return nil, err
	}

	cp.server = transport.NewServer(uint32(10246), []model.SubscribeRequestHandler{cp.handleSubscribeRequest})
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

type say struct{}

func (s *say) Say(ss string) string {
	return ss + "这是一个后缀v2"
}

func (c *ControlPlane) sendMessageToStream(stream model.OpenSergoTransportStream, namespace, app, kind string, dataWithVersion *trpb.DataWithVersion, status *trpb.Status, respId string) error {
	if stream == nil {
		return nil
	}
	client, err := c.server.PluginServer.GetPluginClient("stream")
	if err != nil {
		log.Printf("Error:%s\n", err.Error())
	}
	raw, ok := client.(stream_plugin.Stream)
	if !ok {
		log.Printf("Error: %s\n", "can't convert rpc plugin to normal wrapper")
	}
	sa := &say{}
	greet, err := raw.Greeter("这是一个前缀", sa)
	if err != nil {
		log.Printf("Error: %s\n", err.Error())
	}
	log.Printf("Greeting: %s\n", greet)

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
		_ = c.server.ConnectionManager().Add(request.Target.Namespace, request.Target.App, kind, transport.NewConnection(clientIdentifier, stream))
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
