package types

import (
	"time"
)

type Expression struct {
	WorkerName string    `db:"worker"`
	Username   string    `db:"username"`
	Expression string    `db:"expression"`
	StartTime  time.Time `db:"startTime"`
	EndTime    time.Time `db:"endTime"`
}
