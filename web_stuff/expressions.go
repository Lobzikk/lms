package httpstuff

import (
	"time"
)

type Expression struct {
	Username   string
	Expression string
	StartTime  time.Time
	EndTime    time.Time
}
