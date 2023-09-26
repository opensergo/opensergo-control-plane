package grpc

import (
	"context"
	"fmt"
	v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	extension "github.com/envoyproxy/go-control-plane/envoy/service/extension/v3"
	"github.com/opensergo/opensergo-control-plane/pkg/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

const delimiter = "/"

type XDSClient struct {
	address string
	id      string
	// store the rules subscribed
	subscribedRules []string
	// response nonce
	nonce string

	allResponses []*discovery.DiscoveryResponse
}

// Getter method for the 'nonce' field
func (c *XDSClient) GetNonce() string {
	return c.nonce
}

// Setter method for the 'nonce' field
func (c *XDSClient) SetNonce(nonce string) {
	c.nonce = nonce
}

// Setter method for 'subscribedRules'
func (x *XDSClient) SetSubscribedRules(rules []string) {
	x.subscribedRules = rules
}

// Getter method for 'subscribedRules'
func (x *XDSClient) GetSubscribedRules() []string {
	return x.subscribedRules
}

func NewxDSClient(address string, id string, subscribedRules []string) *XDSClient {
	return &XDSClient{
		address:         address,
		id:              id,
		subscribedRules: subscribedRules,
	}
}

// InitGRPCConnection initializes a gRPC connection to the server.
func (c *XDSClient) InitGRPCConnection() (*grpc.ClientConn, error) {
	ipaddr := c.address
	conn, err := grpc.Dial(ipaddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed to connect to the server: ", err)
		return nil, err
	}
	return conn, nil
}

// CreateDiscoveryRequest creates a DiscoveryRequest.
func (c *XDSClient) CreateDiscoveryRequest(namespace string, app string) *discovery.DiscoveryRequest {
	resourceNames := make([]string, 0)
	for _, rule := range c.subscribedRules {
		resourceName := rule + delimiter + app + delimiter + namespace
		resourceNames = append(resourceNames, resourceName)
	}
	return &discovery.DiscoveryRequest{
		TypeUrl: model.ExtensionConfigType,
		Node: &v3.Node{
			Id: c.id,
		},
		ResponseNonce: c.nonce,
		ResourceNames: resourceNames,
	}
}

// createDiscoveryClient creates a DiscoveryServiceClient using the provided gRPC ClientConn.
func (c *XDSClient) CreateStreamExtensionClient(conn *grpc.ClientConn) (extension.ExtensionConfigDiscoveryService_StreamExtensionConfigsClient, error) {
	discoveryClient := extension.NewExtensionConfigDiscoveryServiceClient(conn)
	streamExtensionClients, err := discoveryClient.StreamExtensionConfigs(context.Background())
	if err != nil {
		log.Fatal("Failed to create streamClient: ", err)
		return nil, err
	}
	return streamExtensionClients, nil
}

// sendDiscoveryRequest sends a DiscoveryRequest to the server using the provided client and request.
func (c *XDSClient) SendDiscoveryRequest(client extension.ExtensionConfigDiscoveryService_StreamExtensionConfigsClient, request *discovery.DiscoveryRequest) error {
	if err := client.Send(request); err != nil {
		return err
	}
	return nil
}

// receiveDiscoveryResponses receives and processes DiscoveryResponses from the server.
func (c *XDSClient) ReceiveDiscoveryResponses(client extension.ExtensionConfigDiscoveryService_StreamExtensionConfigsClient) error {
	for {
		discoveryResponse, err := client.Recv()
		if err != nil {
			return err
		}
		c.SetNonce(discoveryResponse.Nonce)
		fmt.Println(discoveryResponse)
	}
}
