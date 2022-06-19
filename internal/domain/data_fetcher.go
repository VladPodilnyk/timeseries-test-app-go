package domain

import (
	"github.com/VladPodilnyk/timeseries-test-app-go/internal/grpc"
	"github.com/VladPodilnyk/timeseries-test-app-go/internal/model"
)

type DataFetcher struct {
	client grpc.GrcpClient
}

func (*DataFetcher) Retrieve(query model.UserRequest) (*model.PagedData, error) {
	return nil, nil
}

func (*DataFetcher) validateRequest(request model.UserRequest) error {
	return nil
}
