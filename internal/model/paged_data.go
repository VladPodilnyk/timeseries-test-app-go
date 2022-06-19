package model

import (
	"time"
)

type PagedData struct {
	LastTimestamp time.Time            `json:"LastTimestamp"`
	Values        []ValueWithTimestamp `json:"Values"`
}
