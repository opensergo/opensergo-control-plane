package model

import (
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"sync"
	"time"
)

// XDSConnection represents a connection to the xDS server.
type XDSConnection struct {
	sync.RWMutex

	// peerAddr is the address of the client, from the network layer.
	peerAddr string

	// WatchedResources contains the list of watched resources for the proxy,
	// keyed by the DiscoveryRequest TypeUrl.
	WatchedResources map[string]*WatchedResource

	// connectedAt stores the time of connection, mainly for debugging purposes.
	connectedAt time.Time

	// Identifier represents the unique connection identifier (conID),
	// used as a key in the connection table. Currently based on the node name and a counter.
	Identifier ClientIdentifier

	// Stream represents the DiscoveryStream interface implemented by both ADS and SDS streams.
	Stream DiscoveryStream

	// Node stores the original node metadata, to avoid unmarshal/marshal operations.
	// This information is included in internal events.
	Node *core.Node

	// Initialized channel is closed when the proxy is initialized.
	// Pushes or any other operations accessing the proxy should not start until this channel is closed.
	Initialized chan struct{}

	// Stop channel can be used to manually end the connection, typically via debug endpoints.
	// It should only be used for testing purposes.
	Stop chan struct{}

	// ReqChan is used to receive discovery requests for this connection.
	ReqChan chan *discovery.DiscoveryRequest
}

func NewConnection(peerAddr string, stream DiscoveryStream) *XDSConnection {
	return &XDSConnection{

		Initialized:      make(chan struct{}),
		Stop:             make(chan struct{}),
		WatchedResources: make(map[string]*WatchedResource),
		ReqChan:          make(chan *discovery.DiscoveryRequest, 1),
		peerAddr:         peerAddr,
		connectedAt:      time.Now(),
		Stream:           stream,
	}
}

func (conn *XDSConnection) Watched(typeUrl string) *WatchedResource {
	conn.RLock()
	defer conn.RUnlock()
	if conn.WatchedResources != nil && conn.WatchedResources[typeUrl] != nil {
		return conn.WatchedResources[typeUrl]
	}
	return nil
}
