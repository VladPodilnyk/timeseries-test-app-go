package api

import (
	"errors"

	"github.com/VladPodilnyk/timeseries-test-app-go/internal/config"
	"github.com/VladPodilnyk/timeseries-test-app-go/internal/model"
	"github.com/VladPodilnyk/timeseries-test-app-go/internal/protobuf"
	"github.com/VladPodilnyk/timeseries-test-app-go/internal/repo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GrpcServiceIml struct {
	repo         repo.TimeSeries
	limitsConfig config.LimitsConfig
}

func (service *GrpcServiceIml) FetchData(iterator *protobuf.Iterator) (*protobuf.QueryResponse, error) {
	dataRequest := model.UserRequest{Start: iterator.Start.AsTime(), End: iterator.End.AsTime()}
	result, err := service.repo.Fetch(dataRequest, service.limitsConfig.PageLimit)

	if err != nil {
		// TODO: better error messages
		return nil, errors.New("Couldn't fetch data :(")
	}

	lastTimestamp := timestamppb.New(result[len(result)-1].Timestamp)
	dataWithTimestamps := make([]*protobuf.DataWithTimestamp, len(result))

	for i := 0; i < len(result); i++ {
		value := float32(result[i].Value)
		data := protobuf.DataWithTimestamp{Value: &value, ReceivedTime: timestamppb.New(result[i].Timestamp)}
		dataWithTimestamps[i] = &data
	}

	return &protobuf.QueryResponse{Data: dataWithTimestamps, LastProcessed: lastTimestamp}, nil
}
