package ratelimit_plugin

import (
	"context"

	v1 "github.com/opensergo/opensergo-control-plane/pkg/plugin/proto/rate_limit/v1"
)

type GRPCClient struct {
	client v1.RateLimitServiceClient
}

func (g *GRPCClient) RateLimit(t int64) (int64, error) {
	resp, err := g.client.RateLimit(context.Background(), &v1.RateLimitRequest{
		Threshold: t,
	})
	if err != nil {
		return 0, err
	}
	return resp.Threshold, nil
}

type GRPCServer struct {
	v1.UnimplementedRateLimitServiceServer
	Impl RateLimit
}

func (g *GRPCServer) RateLimit(ctx context.Context, req *v1.RateLimitRequest) (*v1.RateLimitResponse, error) {
	resp, err := g.Impl.RateLimit(req.Threshold)
	if err != nil {
		return nil, err
	}
	return &v1.RateLimitResponse{
		Threshold: resp,
	}, nil
}
