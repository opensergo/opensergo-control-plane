package main

import (
	"context"
	"fmt"
	v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/opensergo/opensergo-control-plane/pkg/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	conn, err := grpc.Dial(":8002", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("服务端出错，连接不上", err)
	}
	discoveryClient := discovery.NewAggregatedDiscoveryServiceClient(conn)

	// build request
	request := discovery.DiscoveryRequest{
		TypeUrl: model.ExtensionConfigType,
		Node: &v3.Node{
			Id: "testNode",
		},
		ResourceNames: []string{"testns/testapp/RateLimitStrategy", "testns/testapp/FaultToleranceRule", "testns/testapp/ConcurrencyLimitStrategy"},
	}

	streamAggregatedResourcesClient, err := discoveryClient.StreamAggregatedResources(context.Background())
	if err != nil {
		log.Fatal("得不到streamAggregatedResourcesClient", err)
	}

	if err := streamAggregatedResourcesClient.Send(&request); err != nil {
		log.Fatal("发送消息错误", err)
	}
	for {
		discoveryResponse, err := streamAggregatedResourcesClient.Recv()
		if err != nil {
			log.Fatal("回收消息错误", err)
		}
		fmt.Println(discoveryResponse)
	}

}
