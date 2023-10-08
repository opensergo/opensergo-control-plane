package test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	extension "github.com/envoyproxy/go-control-plane/envoy/service/extension/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/opensergo/opensergo-control-plane/pkg/controller"
	"github.com/opensergo/opensergo-control-plane/pkg/model"
)

// test connection init
func TestInitConnection(t *testing.T) {
	conn, err := grpc.Dial(":8002", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("服务端出错，连接不上", err)
	}
	discoveryClient := extension.NewExtensionConfigDiscoveryServiceClient(conn)

	// build request
	request1 := discovery.DiscoveryRequest{
		TypeUrl: model.ExtensionConfigType,
		Node: &v3.Node{
			Id: "testNode",
		},
		ResourceNames: []string{"default/foo-app/" + controller.RateLimitStrategyKind},
	}

	streamExtensionClients, err := discoveryClient.StreamExtensionConfigs(context.Background())

	if err != nil {
		log.Fatal("得不到streamAggregatedResourcesClient", err)
	}

	if err := streamExtensionClients.Send(&request1); err != nil {
		log.Fatal("发送消息错误", err)
	}

	for {
		discoveryResponse, err := streamExtensionClients.Recv()
		if err != nil {
			log.Fatal("回收消息错误", err)
		}
		fmt.Println(discoveryResponse)
	}

}

func TestMultiConnection(t *testing.T) {
	firstconn, err := grpc.Dial(":8002", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("服务端出错，连接不上", err)
	}
	discoveryClient := extension.NewExtensionConfigDiscoveryServiceClient(firstconn)

	// build request
	firstrequest := discovery.DiscoveryRequest{
		TypeUrl: model.ExtensionConfigType,
		Node: &v3.Node{
			Id: "testNode",
		},
		ResourceNames: []string{"default/foo-app/" + controller.ConcurrencyLimitStrategyKind},
	}

	firststreamExtensionClient, err := discoveryClient.StreamExtensionConfigs(context.Background())

	if err != nil {
		log.Fatal("得不到streamAggregatedResourcesClient", err)
	}

	if err := firststreamExtensionClient.Send(&firstrequest); err != nil {
		log.Fatal("发送消息错误", err)
	}

	for {
		discoveryResponse, err := firststreamExtensionClient.Recv()
		if err != nil {
			log.Fatal("回收消息错误", err)
		}
		fmt.Println(discoveryResponse)
	}

}

func TestResponseNonce(t *testing.T) {
	conn, err := grpc.Dial(":8002", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("服务端出错，连接不上", err)
	}
	discoveryClient := extension.NewExtensionConfigDiscoveryServiceClient(conn)

	// build request
	request1 := discovery.DiscoveryRequest{
		TypeUrl: model.ExtensionConfigType,
		Node: &v3.Node{
			Id: "testNode",
		},
		ResourceNames: []string{"default/foo-app/" + controller.RateLimitStrategyKind},
	}

	streamExtensionClients, err := discoveryClient.StreamExtensionConfigs(context.Background())

	if err != nil {
		log.Fatal("得不到streamAggregatedResourcesClient", err)
	}

	if err := streamExtensionClients.Send(&request1); err != nil {
		log.Fatal("发送消息错误", err)
	}

	for {
		discoveryResponse, err := streamExtensionClients.Recv()
		if err != nil {
			log.Fatal("回收消息错误", err)
		}
		fmt.Println(discoveryResponse)
		// build response
		request2 := &discovery.DiscoveryRequest{
			TypeUrl: model.ExtensionConfigType,
			Node: &v3.Node{
				Id: "testNode",
			},
			ResponseNonce: discoveryResponse.Nonce,
			ResourceNames: []string{"default/foo-app/" + controller.RateLimitStrategyKind},
		}
		// send ack
		streamExtensionClients.Send(request2)
	}
}

func TestSubScribetionChange(t *testing.T) {
	conn, err := grpc.Dial(":8002", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("服务端出错，连接不上", err)
	}
	discoveryClient := extension.NewExtensionConfigDiscoveryServiceClient(conn)

	// build request
	request1 := discovery.DiscoveryRequest{
		TypeUrl: model.ExtensionConfigType,
		Node: &v3.Node{
			Id: "testNode",
		},
		ResourceNames: []string{"default/foo-app/" + controller.RateLimitStrategyKind},
	}

	//build request
	request2 := discovery.DiscoveryRequest{
		TypeUrl: model.ExtensionConfigType,
		Node: &v3.Node{
			Id: "testNode",
		},
		// TODO : kind 放到namespace 前面
		ResourceNames: []string{"default/foo-app/" + controller.RateLimitStrategyKind, "default/foo-app/" + controller.ConcurrencyLimitStrategyKind, "default/foo-app/" + controller.ThrottlingStrategyKind},
	}

	streamExtensionClients, err := discoveryClient.StreamExtensionConfigs(context.Background())
	if err != nil {
		log.Fatal("得不到streamAggregatedResourcesClient", err)
	}
	go func() {
		fmt.Println("request1 send")
		if err := streamExtensionClients.Send(&request1); err != nil {
			log.Fatal("发送消息错误", err)
		}
	}()

	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("request2 send")
		if err := streamExtensionClients.Send(&request2); err != nil {
			log.Fatal("发送消息错误", err)
		}
	}()

	for {
		discoveryResponse, err := streamExtensionClients.Recv()
		if err != nil {
			log.Fatal("回收消息错误", err)
		}
		fmt.Println(discoveryResponse)
	}
}
