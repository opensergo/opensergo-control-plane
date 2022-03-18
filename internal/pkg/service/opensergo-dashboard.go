package service

import (
	metadata "github.com/opensergo/opensergo-go/pkg/api/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func init() {
	conn, err := grpc.Dial("127.0.0.1:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client := metadata.NewMetadataServiceClient(conn)
	proxies = append(proxies, client)
}
