package model

import "time"

type ValueWithTimestamp struct {
	Value     float32
	Timestamp time.Time
}
