package api

import "github.com/VladPodilnyk/timeseries-test-app-go/internal/protobuf"

type GrpcService interface {
	// TODO: async call
	FetchData(iterator *protobuf.Iterator) (*protobuf.QueryResponse, error)
}
