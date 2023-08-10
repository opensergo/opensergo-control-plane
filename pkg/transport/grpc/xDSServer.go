package grpc

import (
	context "context"
	"fmt"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	extension "github.com/envoyproxy/go-control-plane/envoy/service/extension/v3"
	"github.com/opensergo/opensergo-control-plane/pkg/model"
	"github.com/opensergo/opensergo-control-plane/pkg/util"
	"google.golang.org/grpc"
	"net"

	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/anypb"
	"log"
	"strconv"
	"sync"
	"sync/atomic"
)

var connectionNumber = int64(0)

type DiscoveryServer struct {
	port uint32
	// adsClients reflect active gRPC channels, for both ADS and EDS.
	ecdsClients      map[model.ClientIdentifier]*model.XDsConnection
	ecdsClientsMutex sync.RWMutex
	grpcServer       *grpc.Server
	// cxz
	// map from  to xdsconnection
	// namespace appname and kind
	XDSConnectionManeger *ConnectionManager[*model.XDsConnection]
	subscribeHandlers    []model.SubscribeXDsRequestHandler
}

func (s *DiscoveryServer) StreamExtensionConfigs(stream model.DiscoveryStream) error {
	return s.Stream(stream)
}

//func (s *DiscoveryServer) WaitForRequestLimit(ctx context.Context) error {
//	if s.RequestRateLimit.Limit() == 0 {
//		// Allow opt out when rate limiting is set to 0qps
//		return nil
//	}
//	// Give a bit of time for queue to clear out, but if not fail fast. Client will connect to another
//	// instance in best case, or retry with backoff.
//	wait, cancel := context.WithTimeout(ctx, time.Second)
//	defer cancel()
//	return s.RequestRateLimit.Wait(wait)
//}

func (s *DiscoveryServer) Stream(stream model.DiscoveryStream) error {
	ctx := stream.Context()
	peerAddr := "0.0.0.0"
	if peerInfo, ok := peer.FromContext(ctx); ok {
		peerAddr = peerInfo.Addr.String()
	}

	// TODO: Rate Limit
	//if err := s.WaitForRequestLimit(stream.Context()); err != nil {
	//	//log.Warnf("ADS: %q exceeded rate limit: %v", peerAddr, err)
	//	return status.Errorf(codes.ResourceExhausted, "request rate limit exceeded: %v", err)
	//}

	con := model.NewConnection(peerAddr, stream)

	go s.receive(con)

	// Wait for the proxy to be fully initialized before we start serving traffic. Because
	// initialization doesn't have dependencies that will block, there is no need to add any timeout
	// here. Prior to this explicit wait, we were implicitly waiting by receive() not sending to
	// reqChannel and the connection not being enqueued for pushes to pushChannel until the
	// initialization is complete.
	<-con.Initialized

	for {
		// Go select{} statements are not ordered; the same channel can be chosen many times.
		// For requests, these are higher priority (client may be blocked on startup until these are done)
		// and often very cheap to handle (simple ACK), so we check it first.
		select {
		case req, ok := <-con.ReqChan:
			if ok {
				for _, handler := range s.subscribeHandlers {
					err := handler(req, con)
					if err != nil {
						// TODO: handle error
						log.Printf("Failed to handle SubscribeRequest, err=%s\n", err.Error())
					}
				}
			}
		case <-con.Stop:
			return nil
		default:
		}

	}

	return nil
}

func (s *DiscoveryServer) receive(con *model.XDsConnection) {
	defer func() {

		close(con.ReqChan)
		// Close the initialized channel, if its not already closed, to prevent blocking the stream.
		select {
		case <-con.Initialized:
		default:
			close(con.Initialized)
		}
	}()

	firstRequest := true
	for {
		req, err := con.Stream.Recv()
		log.Printf("Stream received message failed, err=%s\n", err.Error())

		if firstRequest {
			firstRequest = false
			if err := s.initConnection(req.Node, con); err != nil {
				log.Printf("initConnection failed, err=%s\n", err.Error())
				return
			}
		}
		select {
		case con.ReqChan <- req:
		case <-con.Stream.Context().Done():
			//log.Infof("ADS: %q %s terminated with stream closed", con.peerAddr, con.conID)
			return
		}
	}
}

func (s *DiscoveryServer) initConnection(node *core.Node, con *model.XDsConnection) error {
	// First request so initialize connection id and start tracking it.
	con.Identifier = connectionID(node.Id)
	con.Node = node

	// Register the connection. this allows pushes to be triggered for the proxy. Note: the timing of
	// this and initializeProxy important. While registering for pushes *after* initialization is complete seems like
	// a better choice, it introduces a race condition; If we complete initialization of a new push
	// context between initializeProxy and addCon, we would not get any pushes triggered for the new
	// push context, leading the proxy to have a stale state until the next full push.
	s.addCon(con.Identifier, con)

	defer close(con.Initialized)

	return nil
}

func connectionID(node string) model.ClientIdentifier {
	id := atomic.AddInt64(&connectionNumber, 1)
	return model.ClientIdentifier(node + "-" + strconv.FormatInt(id, 10))
}

func (s *DiscoveryServer) addCon(identifier model.ClientIdentifier, con *model.XDsConnection) {
	s.ecdsClientsMutex.Lock()
	defer s.ecdsClientsMutex.Unlock()
	s.ecdsClients[identifier] = con

}

func (s *DiscoveryServer) removeCon(conID model.ClientIdentifier) {
	s.ecdsClientsMutex.Lock()
	defer s.ecdsClientsMutex.Unlock()

	if _, exist := s.ecdsClients[conID]; !exist {
		//log.Errorf("ADS: Removing connection for non-existing node:%v.", conID)

	} else {
		delete(s.ecdsClients, conID)
	}
}

func shouldUnsubscribe(request *discovery.DiscoveryRequest) bool {
	return len(request.ResourceNames) == 0
}

var emptyResourceDelta = model.ResourceDelta{}

func ShouldRespond(con *model.XDsConnection, request *discovery.DiscoveryRequest) (bool, model.ResourceDelta) {

	// NACK
	if request.ErrorDetail != nil {
		//LOG
		return false, emptyResourceDelta
	}

	if shouldUnsubscribe(request) {
		con.Lock()
		delete(con.WatchedResources, request.TypeUrl)
		con.Unlock()
		return false, emptyResourceDelta
	}

	con.RLock()
	previousInfo := con.WatchedResources[request.TypeUrl]
	con.RUnlock()

	// We should always respond with the current resource names.
	if request.ResponseNonce == "" || previousInfo == nil {
		log.Println("ECDS: INIT/RECONNECT %s %s %s", con.Identifier, request.VersionInfo, request.ResponseNonce)
		con.Lock()
		con.WatchedResources[request.TypeUrl] = &model.WatchedResource{TypeUrl: request.TypeUrl, ResourceNames: request.ResourceNames}
		con.Unlock()
		return true, model.ResourceDelta{
			Subscribed:   util.New(request.ResourceNames...),
			Unsubscribed: util.String{},
		}
	}

	// If there is mismatch in the nonce, that is a case of expired/stale nonce.
	// A nonce becomes stale following a newer nonce being sent to Envoy.
	// previousInfo.NonceSent can be empty if we previously had shouldRespond=true but didn't send any resources.
	if request.ResponseNonce != previousInfo.NonceSent {

		log.Println("ECDS: REQ %s Expired nonce received %s, sent %s",
			con.Identifier, request.ResponseNonce, previousInfo.NonceSent)
		return false, emptyResourceDelta
	}

	// If it comes here, that means nonce match.
	con.Lock()
	previousResources := con.WatchedResources[request.TypeUrl].ResourceNames
	con.WatchedResources[request.TypeUrl].NonceAcked = request.ResponseNonce
	con.WatchedResources[request.TypeUrl].ResourceNames = request.ResourceNames
	con.Unlock()

	// Envoy can send two DiscoveryRequests with same version and nonce.
	// when it detects a new resource. We should respond if they change.
	prev := util.New(previousResources...)
	cur := util.New(request.ResourceNames...)
	removed := prev.Difference(cur)
	added := cur.Difference(prev)

	if len(removed) == 0 && len(added) == 0 {
		// this is an ack nonce matched
		return false, emptyResourceDelta
	}

	return true, model.ResourceDelta{
		Subscribed:   added,
		Unsubscribed: removed,
	}
}

func (s *DiscoveryServer) pushXds(con *model.XDsConnection, w *model.WatchedResource, version int64, rules []*anypb.Any) error {

	resp := &discovery.DiscoveryResponse{
		TypeUrl:     w.TypeUrl,
		VersionInfo: strconv.FormatInt(version, 10),
		Nonce:       util.Nonce(),
		Resources:   rules,
	}

	return con.Stream.Send(resp)
}

// cxz
func (s *DiscoveryServer) AddConnectioonToMap(namespace, appname, kind string, con *model.XDsConnection) {
	s.ecdsClientsMutex.Lock()
	defer s.ecdsClientsMutex.Unlock()

	s.XDSConnectionManeger.Add(namespace, appname, kind, con, con.Identifier)

}

func (s *DiscoveryServer) RemoveConnectionFromMap(n model.NamespacedApp, kind string, identifier model.ClientIdentifier) error {
	s.ecdsClientsMutex.Lock()
	defer s.ecdsClientsMutex.Unlock()
	// TODO: HANDLE Error
	if err := s.XDSConnectionManeger.removeInternal(n, kind, identifier); err != nil {
		return err
	}
	return nil

}

func NewDiscoveryServer(port uint32, subscribeHandlers []model.SubscribeXDsRequestHandler) *DiscoveryServer {
	connectionManager := NewConnectionManager[*model.XDsConnection]()
	return &DiscoveryServer{
		port:                 port,
		ecdsClients:          make(map[model.ClientIdentifier]*model.XDsConnection),
		grpcServer:           grpc.NewServer(),
		XDSConnectionManeger: connectionManager,
		subscribeHandlers:    subscribeHandlers,
	}
}

// TODO : Unimplemented
func (s *DiscoveryServer) DeltaExtensionConfigs(stream extension.ExtensionConfigDiscoveryService_DeltaExtensionConfigsServer) error {
	return nil
}

func (s *DiscoveryServer) FetchExtensionConfigs(context.Context, *discovery.DiscoveryRequest) (*discovery.DiscoveryResponse, error) {
	return &discovery.DiscoveryResponse{}, nil
}

func (s *DiscoveryServer) Run() error {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}

	extension.RegisterExtensionConfigDiscoveryServiceServer(s.grpcServer, s)
	err = s.grpcServer.Serve(listener)
	if err != nil {
		return err
	}
	return nil
}
