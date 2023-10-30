package stream_plugin

import (
	"context"
	"fmt"

	v1 "github.com/opensergo/opensergo-control-plane/pkg/plugin/proto/stream/v1"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

type GRPCClient struct {
	client v1.RateLimitServiceClient
	broker *plugin.GRPCBroker
}

func (g *GRPCClient) Greeter(name string, h Hello) (string, error) {
	addHelperServer := &GRPCHelloServer{Impl: h}

	var s *grpc.Server
	serverFunc := func(opts []grpc.ServerOption) *grpc.Server {
		s = grpc.NewServer(opts...)
		v1.RegisterHelloServer(s, addHelperServer)

		return s
	}

	brokerID := g.broker.NextId()
	go g.broker.AcceptAndServe(brokerID, serverFunc)

	resp, err := g.client.Greet(context.Background(), &v1.StreamReq{
		Id:   brokerID,
		Name: name,
	})
	if err != nil {
		return "", err
	}

	s.Stop()
	return resp.Greet, nil
}

type GRPCServer struct {
	v1.UnimplementedStreamGreeterServer
	Impl   Stream
	broker *plugin.GRPCBroker
}

func (g *GRPCServer) Greet(ctx context.Context, req *v1.StreamReq) (*v1.StreamResp, error) {
	conn, err := g.broker.Dial(req.Id)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	a := &GRPCHelloClient{v1.NewHelloClient(conn)}
	resp, err := g.Impl.Greeter(req.Name, a)
	if err != nil {
		return nil, err
	}
	return &v1.StreamResp{
		Greet: resp,
	}, nil
}

type GRPCHelloClient struct {
	client v1.HelloClient
}

func (g *GRPCHelloClient) Say(s string) string {
	resp, err := g.client.Say(context.Background(), &v1.HelloReq{
		Pre: s,
	})
	if err != nil {
		return ""
	}
	return resp.Resp
}

type GRPCHelloServer struct {
	v1.UnimplementedHelloServer
	Impl Hello
}

func (g *GRPCHelloServer) Say(ctx context.Context, req *v1.HelloReq) (*v1.HelloResp, error) {
	resp := g.Impl.Say(fmt.Sprint(req.Pre, " GRPCHelloServer"))
	return &v1.HelloResp{
		Resp: resp,
	}, nil
}
