package stream_plugin

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-plugin"
	pb "github.com/opensergo/opensergo-control-plane/pkg/plugin/proto/stream"
	"google.golang.org/grpc"
)

type GRPCClient struct {
	client pb.StreamGreeterClient
	broker *plugin.GRPCBroker
}

func (g *GRPCClient) Greeter(name string, h Hello) (string, error) {
	addHelperServer := &GRPCHelloServer{Impl: h}

	var s *grpc.Server
	serverFunc := func(opts []grpc.ServerOption) *grpc.Server {
		s = grpc.NewServer(opts...)
		pb.RegisterHelloServer(s, addHelperServer)

		return s
	}

	brokerID := g.broker.NextId()
	go g.broker.AcceptAndServe(brokerID, serverFunc)

	resp, err := g.client.Greet(context.Background(), &pb.StreamReq{
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
	pb.UnimplementedStreamGreeterServer
	Impl   Stream
	broker *plugin.GRPCBroker
}

func (g *GRPCServer) Greet(ctx context.Context, req *pb.StreamReq) (*pb.StreamResp, error) {
	conn, err := g.broker.Dial(req.Id)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	a := &GRPCHelloClient{pb.NewHelloClient(conn)}
	resp, err := g.Impl.Greeter(req.Name, a)
	if err != nil {
		return nil, err
	}
	return &pb.StreamResp{
		Greet: resp,
	}, nil
}

type GRPCHelloClient struct {
	client pb.HelloClient
}

func (g *GRPCHelloClient) Say(s string) string {
	resp, err := g.client.Say(context.Background(), &pb.HelloReq{
		Pre: s,
	})
	if err != nil {
		return ""
	}
	return resp.Resp
}

type GRPCHelloServer struct {
	pb.UnimplementedHelloServer
	Impl Hello
}

func (g *GRPCHelloServer) Say(ctx context.Context, req *pb.HelloReq) (*pb.HelloResp, error) {
	resp := g.Impl.Say(fmt.Sprint(req.Pre, " GRPCHelloServer"))
	return &pb.HelloResp{
		Resp: resp,
	}, nil
}
