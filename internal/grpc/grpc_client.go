package grpc

import (
	"context"
	"math/rand"

	"github.com/VladPodilnyk/timeseries-test-app-go/internal/protobuf"
)

type GrcpClient interface {
	FetchData(ctx context.Context, iterator *protobuf.Iterator) (*protobuf.QueryResponse, error)
}

type GrcpClientImpl struct {
	client protobuf.TimeSeriesClient
}

func (impl *GrcpClientImpl) FetchData(iterator *protobuf.Iterator) (*protobuf.QueryResponse, error) {
	result, err := impl.client.FetchData(context.Background(), iterator)
	if err != nil {
		return nil, err
	}
	return result, nil
}

type GrcpClientStub struct {
}

func (impl *GrcpClientStub) FetchData(iterator *protobuf.Iterator) (*protobuf.QueryResponse, error) {
	randomValue := float32(rand.Intn(10))
	data := []*protobuf.DataWithTimestamp{{Value: &randomValue, ReceivedTime: iterator.Start}}
	result := protobuf.QueryResponse{LastProcessed: iterator.Start, Data: data}
	return &result, nil
}
