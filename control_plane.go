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
	"github.com/avast/retry-go/v4"
	"github.com/opensergo/opensergo-control-plane/pkg/controller"
	"github.com/opensergo/opensergo-control-plane/pkg/model"
	"github.com/opensergo/opensergo-control-plane/pkg/options"
	trpb "github.com/opensergo/opensergo-control-plane/pkg/proto/transport/v1"
	transport "github.com/opensergo/opensergo-control-plane/pkg/transport/grpc"
	"github.com/opensergo/opensergo-control-plane/pkg/util"
	"github.com/pkg/errors"
	"os"
	"sync"
)

type ControlPlane struct {
	operator *controller.KubernetesOperator
	server   *transport.Server

	protoDesc *trpb.ControlPlaneDesc

	opts *options.Options

	mux sync.RWMutex
}

func NewControlPlane(opts *options.Options) (*ControlPlane, error) {
	cp := &ControlPlane{}
	cp.opts = opts

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

func (c *ControlPlane) sendMessageToStream(stream model.OpenSergoTransportStream, namespace, app, kind string, dataWithVersion *trpb.DataWithVersion, trpbStatus *trpb.Status, respId string) error {
	if stream == nil {
		return nil
	}

	return retry.Do(
		func() error {
			error := stream.SendMsg(&trpb.SubscribeResponse{
				Status:          trpbStatus,
				Ack:             "",
				Namespace:       namespace,
				App:             app,
				Kind:            kind,
				DataWithVersion: dataWithVersion,
				ControlPlane:    c.protoDesc,
				ResponseId:      respId,
			})

			return error
		},
		retry.Attempts(uint(c.opts.ConfigPushMaxAttempt)),
		retry.RetryIf(util.TimeoutCondition))
}

func (c *ControlPlane) handleSubscribeRequest(clientIdentifier model.ClientIdentifier, request *trpb.SubscribeRequest, stream model.OpenSergoTransportStream) error {
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
			}
			continue
		}
		_ = c.server.ConnectionManager().Add(request.Target.Namespace, request.Target.App, kind, transport.NewConnection(clientIdentifier, stream))
		// watcher缓存不空就发送
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
			}
		}
	}
	return nil
}
