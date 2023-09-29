package ratelimit_plugin

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-plugin"
	pb "github.com/opensergo/opensergo-control-plane/pkg/plugin/proto/rate_limit/v1"
	"google.golang.org/grpc"
)

type RateLimitPluginServer struct {
}

var _ RateLimit = (*RateLimitPluginServer)(nil)

func (s RateLimitPluginServer) RateLimit(t int64) (int64, error) {
	return t + 1, nil
}

type RateLimit interface {
	RateLimit(t int64) (int64, error)
}

type RateLimitPlugin struct {
	plugin.Plugin

	impl RateLimit
}

func NewRateLimitPluginServiceServer(impl RateLimit) (*RateLimitPlugin, error) {
	if impl == nil {
		return nil, fmt.Errorf("empty underlying stream plugin passed in")
	}
	return &RateLimitPlugin{
		impl: impl,
	}, nil
}

func (h *RateLimitPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	pb.RegisterRateLimitServiceServer(s, &GRPCServer{
		Impl: h.impl,
	})
	return nil
}

func (h *RateLimitPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (any, error) {
	return &GRPCClient{
		client: pb.NewRateLimitServiceClient(c),
	}, nil
}
