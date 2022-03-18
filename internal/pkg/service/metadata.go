package service

import (
	"context"
	"sync"

	metadata "github.com/opensergo/opensergo-go/pkg/api/proto"
	"go.uber.org/multierr"
)

var proxies []metadata.MetadataServiceClient

type MetadataService struct {
	metadata.UnimplementedMetadataServiceServer
}

func NewMetadataService() (*MetadataService, error) {
	return &MetadataService{}, nil
}

func (m *MetadataService) ReportMetadata(ctx context.Context, req *metadata.ReportMetadataRequest) (*metadata.ReportMetadataReply, error) {
	var wg sync.WaitGroup
	var globalErr error

	wg.Add(len(proxies))
	for _, client := range proxies {
		client2 := client
		go func() {
			defer wg.Done()
			_, err := client2.ReportMetadata(ctx, req)
			globalErr = multierr.Append(globalErr, err)
		}()
	}
	wg.Wait()
	return &metadata.ReportMetadataReply{}, globalErr
}
