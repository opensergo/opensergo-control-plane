package grpc

import (
	context "context"
	"fmt"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	extension "github.com/envoyproxy/go-control-plane/envoy/service/extension/v3"
	"github.com/opensergo/opensergo-control-plane/pkg/model"
	"github.com/opensergo/opensergo-control-plane/pkg/util"
	uatomic "go.uber.org/atomic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/anypb"
	"log"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var connectionNumber = int64(0)

type DiscoveryServer struct {
	port                 uint32                                          // Port on which the server is running.
	ecdsClients          map[model.ClientIdentifier]*model.XDSConnection // Active gRPC channels for both ADS and EDS.
	ecdsClientsMutex     sync.RWMutex                                    // Mutex for managing concurrent access to ecdsClients.
	grpcServer           *grpc.Server                                    // gRPC server instance.
	XDSConnectionManager *ConnectionManager[*model.XDSConnection]        // Connection manager for XDS connections.
	subscribeHandlers    []model.SubscribeXDsRequestHandler              // List of request handlers for XDS subscription.
	pushVersion          uatomic.Uint64                                  // Atomic counter for push version.
}

func (s *DiscoveryServer) StreamExtensionConfigs(stream model.DiscoveryStream) error {
	return s.Stream(stream)
}

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

func (s *DiscoveryServer) receive(con *model.XDSConnection) {
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
		if err != nil {
			log.Printf("Stream received message failed, err=%s\n", err.Error())
			return
		}

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

func (s *DiscoveryServer) initConnection(node *core.Node, con *model.XDSConnection) error {
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

func (s *DiscoveryServer) addCon(identifier model.ClientIdentifier, con *model.XDSConnection) {
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

func ShouldRespond(con *model.XDSConnection, request *discovery.DiscoveryRequest) (bool, model.ResourceDelta) {

	// NACK
	if request.ErrorDetail != nil {
		//LOG
		return false, emptyResourceDelta
	}

	con.RLock()
	previousInfo := con.WatchedResources[request.TypeUrl]
	con.RUnlock()

	if shouldUnsubscribe(request) {
		con.Lock()
		delete(con.WatchedResources, request.TypeUrl)
		con.Unlock()
		return false, emptyResourceDelta
	}

	// We should always respond with the current resource names.
	if request.ResponseNonce == "" || previousInfo == nil {
		log.Printf("ECDS: INIT/RECONNECT %s %s %s", con.Identifier, request.VersionInfo, request.ResponseNonce)
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

	// log for test
	log.Printf("nonce before %s nonce now %s,", request.ResponseNonce, previousInfo.NonceSent)
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
		log.Println("ack received")
		// this is an ack nonce matched
		return false, emptyResourceDelta
	}

	return true, model.ResourceDelta{
		Subscribed:   added,
		Unsubscribed: removed,
	}
}

func (s *DiscoveryServer) pushXds(con *model.XDSConnection, w *model.WatchedResource, version int64, rules []*anypb.Any) error {

	resp := &discovery.DiscoveryResponse{
		TypeUrl:     w.TypeUrl,
		VersionInfo: strconv.FormatInt(version, 10),
		Nonce:       util.Nonce(),
		Resources:   rules,
	}

	return con.Stream.Send(resp)
}

func (s *DiscoveryServer) AddConnectionToMap(namespace, appname, kind string, con *model.XDSConnection) {
	s.ecdsClientsMutex.Lock()
	defer s.ecdsClientsMutex.Unlock()

	s.XDSConnectionManager.Add(namespace, appname, kind, con, con.Identifier)

}

func (s *DiscoveryServer) RemoveConnectionFromMap(n model.NamespacedApp, kind string, identifier model.ClientIdentifier) error {
	s.ecdsClientsMutex.Lock()
	defer s.ecdsClientsMutex.Unlock()
	// TODO: HANDLE Error
	if err := s.XDSConnectionManager.removeInternal(n, kind, identifier); err != nil {
		return err
	}
	return nil

}

func NewDiscoveryServer(port uint32, subscribeHandlers []model.SubscribeXDsRequestHandler) *DiscoveryServer {
	connectionManager := NewConnectionManager[*model.XDSConnection]()
	return &DiscoveryServer{
		port:                 port,
		ecdsClients:          make(map[model.ClientIdentifier]*model.XDSConnection),
		grpcServer:           grpc.NewServer(),
		XDSConnectionManager: connectionManager,
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
	fmt.Println(listener)
	extension.RegisterExtensionConfigDiscoveryServiceServer(s.grpcServer, s)
	err = s.grpcServer.Serve(listener)

	if err != nil {
		return err
	}
	return nil
}

func (s *DiscoveryServer) NextVersion() string {
	return time.Now().Format(time.RFC3339) + "/" + strconv.FormatUint(s.pushVersion.Inc(), 10)
}
