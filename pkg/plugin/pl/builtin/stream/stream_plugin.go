package stream_plugin

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-plugin"
	pb "github.com/opensergo/opensergo-control-plane/pkg/plugin/proto/stream"
	"google.golang.org/grpc"
)

type StreamPluginServer struct {
}

var _ Stream = (*StreamPluginServer)(nil)

func (s StreamPluginServer) Greeter(name string, h Hello) (string, error) {
	sp := fmt.Sprintf("pre:%s, h:%s\n", name, h.Say("test"))
	return sp, nil
}

type Hello interface {
	Say(s string) string
}

type Stream interface {
	Greeter(name string, h Hello) (string, error)
}

type StreamPlugin struct {
	plugin.Plugin

	impl Stream
}

func NewStreamPluginServiceServer(impl Stream) (*StreamPlugin, error) {
	if impl == nil {
		return nil, fmt.Errorf("empty underlying stream plugin passed in")
	}
	return &StreamPlugin{
		impl: impl,
	}, nil
}

func (h *StreamPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	pb.RegisterStreamGreeterServer(s, &GRPCServer{
		Impl:   h.impl,
		broker: broker,
	})
	return nil
}

func (h *StreamPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (any, error) {
	return &GRPCClient{
		client: pb.NewStreamGreeterClient(c),
		broker: broker,
	}, nil
}
