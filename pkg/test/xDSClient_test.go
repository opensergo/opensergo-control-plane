package test

import (
	"github.com/opensergo/opensergo-control-plane/pkg/controller"
	"github.com/opensergo/opensergo-control-plane/pkg/transport/grpc"
	"log"
	"testing"
)

func TestXDSClient(t *testing.T) {
	// Define the address of the server
	serverAddr := ":8002" // Replace with the actual address

	// Create a new XDSClient instance
	client := grpc.NewxDSClient(serverAddr, "testNode", []string{controller.FaultToleranceRuleKind, controller.ThrottlingStrategyKind})

	// Initialize the gRPC connection to the server
	conn, err := client.InitGRPCConnection()
	if err != nil {
		log.Fatalf("Failed to initialize gRPC connection: %v", err)
	}

	// Create a DiscoveryRequest
	discoveryRequest := client.CreateDiscoveryRequest("default", "foo-app")

	// Create a stream extension client for handling streaming requests
	streamExtensionClient, _ := client.CreateStreamExtensionClient(conn)

	// Send the DiscoveryRequest to the server
	err = client.SendDiscoveryRequest(streamExtensionClient, discoveryRequest)
	if err != nil {
		log.Fatalf("Failed to send DiscoveryRequest: %v", err)
	}

	// Receive and handle DiscoveryResponses from the server
	client.ReceiveDiscoveryResponses(streamExtensionClient)

	// The client will continuously receive and print DiscoveryResponses
	// You can add your own logic to handle the responses as needed
}
