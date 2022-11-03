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

	handlers := []model.SubscribeRequestHandler{
		cp.handleSubscribeRequest,
		cp.handleUnSubscribeRequest,
	}
	cp.server = transport.NewServer(uint32(10246), handlers)
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

func (c *ControlPlane) sendAckToStream(stream model.OpenSergoTransportStream, ack string, status *trpb.Status, respId string) error {
	if stream == nil {
		return nil
	}
	return stream.SendMsg(&trpb.SubscribeResponse{
		Status:       status,
		Ack:          ack,
		ControlPlane: c.protoDesc,
		ResponseId:   respId,
	})
}

func (c *ControlPlane) handleSubscribeRequest(clientIdentifier model.ClientIdentifier, request *trpb.SubscribeRequest, stream model.OpenSergoTransportStream) error {

	if trpb.SubscribeOpType_SUBSCRIBE != request.OpType {
		return nil
	}

	//var labels []model.LabelKV
	//if request.Target.Labels != nil {
	//	for _, label := range request.Target.Labels {
	//		labels = append(labels, model.LabelKV{
	//			Key:   label.Key,
	//			Value: label.Value,
	//		})
	//	}
	//}
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

// handleUnSubscribeRequest handle the UnSubscribeRequest request from OpenSergo SDK.
//
// 1.use ConnectionManager to remove from connectionMap for SubscribeTarget
// 2.use operator to UnRegisterWatcher for SubscribeTarget which will remove SubscribeTarget, delete crdCache, and remove CrdWatcher
func (c *ControlPlane) handleUnSubscribeRequest(clientIdentifier model.ClientIdentifier, request *trpb.SubscribeRequest, stream model.OpenSergoTransportStream) error {

	if trpb.SubscribeOpType_UNSUBSCRIBE != request.OpType {
		return nil
	}

	for _, kind := range request.Target.Kinds {
		namespacedApp := model.NamespacedApp{
			Namespace: request.Target.Namespace,
			App:       request.Target.App,
		}
		// remove the relation of Connection and SubscribeTarget from local cache
		err := c.server.ConnectionManager().RemoveWithIdentifier(namespacedApp, kind, clientIdentifier)
		if err != nil {
			log.Printf("Remove map of Connection-SubscribeTarget failed, err=%s\n", err.Error())
			status := &trpb.Status{
				// TODO: defined a new errorCode
				Code:    transport.RegisterWatcherError,
				Message: "Remove from watcher error",
				Details: nil,
			}
			err = c.sendMessageToStream(stream, request.Target.Namespace, request.Target.App, kind, nil, status, request.RequestId)
			if err != nil {
				// TODO: log here
				log.Printf("sendMessageToStream failed, err=%s\n", err.Error())
			}
			continue
		}

		// UnRegisterWatcher for SubscribeTarget
		err = c.operator.UnRegisterWatcher(model.SubscribeTarget{
			Namespace: request.Target.Namespace,
			AppName:   request.Target.App,
			Kind:      kind,
		})
		if err != nil {
			log.Printf("UnRegisterWatcher failed, err=%s\n", err.Error())
		}
	}

	return nil
}
