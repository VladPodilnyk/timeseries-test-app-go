package repo

import (
	"time"

	"github.com/VladPodilnyk/timeseries-test-app-go/internal/model"
)

type TimeSeries interface {
	Submit(timestamp time.Time, value float64)
	Fetch(iterator model.UserRequest, pageLimit int) ([]model.ValueWithTimestamp, error)
}
