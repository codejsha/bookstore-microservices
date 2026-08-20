package aggregate

import "time"

type PointChangeType string

const (
	POINTCHANGE_EARN   PointChangeType = "EARN"
	POINTCHANGE_SPEND  PointChangeType = "SPEND"
	POINTCHANGE_ADJUST PointChangeType = "ADJUST"
	POINTCHANGE_EXPIRE PointChangeType = "EXPIRE"
)

type PointAggregate struct {
	Uid     string
	UserUid string
	Balance int32
}

type PointHistoryEntry struct {
	Uid        string
	ChangeType PointChangeType
	Amount     int32
	Reason     *string
	CreatedAt  time.Time
}
