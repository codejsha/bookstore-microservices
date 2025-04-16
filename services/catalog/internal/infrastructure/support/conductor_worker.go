package support

import (
	"fmt"

	"github.com/conductor-sdk/conductor-go/sdk/worker"

	"github.com/codejsha/bookstore-microservices/catalog/internal/infrastructure/adapter/conductor"
	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/saga"
)

func NewConductorWorker(
	workerClient *saga.WorkerClient,
	bookWorker *conductor.BookWorker,
) *ConductorWorker {
	return &ConductorWorker{
		taskRunner: workerClient.TaskRunner,
		bookWorker: bookWorker,
	}
}

type ConductorWorker struct {
	taskRunner *worker.TaskRunner
	bookWorker *conductor.BookWorker
}

func (r *ConductorWorker) Run() {
	tasks := r.bookWorker.WorkerTaskSpecs()
	for _, task := range tasks {
		r.startWorker(task)
	}
}

func (r *ConductorWorker) startWorker(spec saga.WorkerTaskSpec) {
	err := r.taskRunner.StartWorkerWithDomain(
		spec.TaskName,
		spec.TaskFunction,
		spec.BatchSize,
		spec.PollingInterval,
		spec.Domain,
	)
	if err != nil {
		_ = fmt.Errorf("failed to start worker %s: %v", spec.TaskName, err)
	}
}
