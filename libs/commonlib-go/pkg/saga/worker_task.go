package saga

import (
	"time"

	"github.com/conductor-sdk/conductor-go/sdk/model"
)

type WorkerTaskSpec struct {
	TaskName        string
	TaskFunction    model.ExecuteTaskFunction
	BatchSize       int
	PollingInterval time.Duration
	Domain          string
}
