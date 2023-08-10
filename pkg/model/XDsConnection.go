package model

import (
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"sync"
	"time"
)

type XDsConnection struct {
	sync.RWMutex
	// peerAddr is the address of the client, from network layer.
	peerAddr string

	// WatchedResources contains the list of watched resources for the proxy, keyed by the DiscoveryRequest TypeUrl.
	WatchedResources map[string]*WatchedResource

	// Time of connection, for debugging
	connectedAt time.Time

	// conID is the connection conID, used as a key in the connection table.
	// Currently based on the node name and a counter.
	Identifier ClientIdentifier

	// Both ADS and SDS streams implement this interface
	Stream DiscoveryStream

	// Original node metadata, to avoid unmarshal/marshal.
	// This is included in internal events.
	Node *core.Node

	// initialized channel will be closed when proxy is initialized. Pushes, or anything accessing
	// the proxy, should not be started until this channel is closed.
	Initialized chan struct{}

	// stop can be used to end the connection manually via debug endpoints. Only to be used for testing.
	Stop chan struct{}

	// reqChan is used to receive discovery requests for this connection.
	ReqChan chan *discovery.DiscoveryRequest
}

func NewConnection(peerAddr string, stream DiscoveryStream) *XDsConnection {
	return &XDsConnection{

		Initialized: make(chan struct{}),
		Stop:        make(chan struct{}),
		ReqChan:     make(chan *discovery.DiscoveryRequest, 1),
		peerAddr:    peerAddr,
		connectedAt: time.Now(),
		Stream:      stream,
	}
}

func (conn *XDsConnection) Watched(typeUrl string) *WatchedResource {
	conn.RLock()
	defer conn.RUnlock()
	if conn.WatchedResources != nil && conn.WatchedResources[typeUrl] != nil {
		return conn.WatchedResources[typeUrl]
	}
	return nil
}
