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

package grpc

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/opensergo/opensergo-control-plane/pkg/model"
	trpb "github.com/opensergo/opensergo-control-plane/pkg/proto/transport/v1"
	"github.com/opensergo/opensergo-control-plane/pkg/util"
	"go.uber.org/atomic"
	"google.golang.org/grpc"
)

const (
	ClientIdentifierKey = "OpenSergoClientIdentifier"
)

// Server represents the transport server of OpenSergo universal transport service (OUTS).
type Server struct {
	transportServer *TransportServer
	grpcServer      *grpc.Server

	connectionManager *ConnectionManager

	port    uint32
	started *atomic.Bool
}

func NewServer(port uint32, subscribeHandlers []model.SubscribeRequestHandler) *Server {
	connectionManager := NewConnectionManager()
	return &Server{
		transportServer:   newTransportServer(connectionManager, subscribeHandlers),
		port:              port,
		grpcServer:        grpc.NewServer(),
		started:           atomic.NewBool(false),
		connectionManager: connectionManager,
	}
}

func (s *Server) ConnectionManager() *ConnectionManager {
	return s.connectionManager
}

func (s *Server) ComponentName() string {
	return "OpenSergoUniversalTransportServer"
}

func (s *Server) Run() error {
	if s.started.CAS(false, true) {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
		if err != nil {
			return err
		}

		trpb.RegisterOpenSergoUniversalTransportServiceServer(s.grpcServer, s.transportServer)
		err = s.grpcServer.Serve(listener)
		if err != nil {
			return err
		}
	}
	return nil
}

// TransportServer represents the gRPC server of OpenSergo universal transport service.
type TransportServer struct {
	trpb.OpenSergoUniversalTransportServiceServer

	connectionManager *ConnectionManager

	subscribeHandlers []model.SubscribeRequestHandler
}

const (
	ACKFlag  = "ACK"
	NACKFlag = "NACK"

	Success              = 1
	CheckFormatError     = 4001
	ReqFormatError       = 4002
	RegisterWatcherError = 500
)

func (s *TransportServer) SubscribeConfig(stream trpb.OpenSergoUniversalTransportService_SubscribeConfigServer) error {
	var clientIdentifier model.ClientIdentifier
	for {
		recvData, err := stream.Recv()
		if err == io.EOF {
			// Stream EOF
			_ = s.connectionManager.RemoveByIdentifier(clientIdentifier)
			return nil
		}
		if err != nil {
			//remove stream
			_ = s.connectionManager.RemoveByIdentifier(clientIdentifier)
			return err
		}

		if recvData.ResponseAck == ACKFlag {
			// This indicates the received data is a response of push-success.
			continue
		} else if recvData.ResponseAck == NACKFlag {
			// This indicates the received data is a response of push-failure.
			if recvData.Status.Code == CheckFormatError {
				// TODO: handle here (cannot retry)
				log.Println("Client response CheckFormatError")
			} else {
				// TODO: record error here and do something
				log.Printf("Client response NACK, code=%d\n", recvData.Status.Code)
			}
		} else {
			// This indicates the received data is a SubscribeRequest.
			if clientIdentifier == "" && recvData.Identifier != "" {
				clientIdentifier = model.ClientIdentifier(recvData.Identifier)
			}

			if !util.IsValidReq(recvData) {
				status := &trpb.Status{
					Code:    ReqFormatError,
					Message: "Request is invalid",
					Details: nil,
				}
				_ = stream.Send(&trpb.SubscribeResponse{
					Status:     status,
					Ack:        NACKFlag,
					ResponseId: recvData.RequestId,
				})
				continue
			}

			for _, handler := range s.subscribeHandlers {
				err = handler(clientIdentifier, recvData, stream)
				if err != nil {
					// TODO: handle error
					log.Printf("Failed to handle SubscribeRequest, err=%s\n", err.Error())
				}
			}
		}

	}
}

func newTransportServer(connectionManager *ConnectionManager, subscribeHandlers []model.SubscribeRequestHandler) *TransportServer {
	return &TransportServer{
		connectionManager: connectionManager,
		subscribeHandlers: subscribeHandlers,
	}
}
