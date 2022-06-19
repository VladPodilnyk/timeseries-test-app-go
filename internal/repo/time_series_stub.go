package repo

import (
	"time"

	"github.com/VladPodilnyk/timeseries-test-app-go/internal/model"
)

type TimeSeriesStub struct {
	storage map[time.Time]float64
}

func MakeStub() *TimeSeriesStub {
	inMemStorage := make(map[time.Time]float64)
	return &TimeSeriesStub{inMemStorage}
}

func (repo *TimeSeriesStub) Submit(timestamp time.Time, value float64) {
	repo.storage[timestamp] = value
}

func (repo *TimeSeriesStub) Fetch(iterator model.UserRequest, pageLimit int) ([]model.ValueWithTimestamp, error) {
	result := make([]model.ValueWithTimestamp, 0)
	for key, element := range repo.storage {
		if key.After(iterator.Start) && key.Before(iterator.End) && len(result) < pageLimit {
			result = append(result, model.ValueWithTimestamp{Value: float32(element), Timestamp: key})
		}
	}
	return result, nil
}
