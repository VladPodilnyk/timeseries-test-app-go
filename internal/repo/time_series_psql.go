package repo

import "time"

type TimeSeriesPostgres struct {
	Storage map[time.Time]float64
}
