package model

import "time"

type UserRequest struct {
	Start time.Time
	End   time.Time
}
