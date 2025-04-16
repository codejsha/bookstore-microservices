package support

import (
	"fmt"

	"github.com/conductor-sdk/conductor-go/sdk/worker"

	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/saga"
	"github.com/codejsha/bookstore-microservices/customer/internal/infrastructure/adapter/conductor"
)

func NewConductorWorker(
	workerClient *saga.WorkerClient,
	customerWorker *conductor.CustomerWorker,
) *ConductorWorker {
	return &ConductorWorker{
		taskRunner:     workerClient.TaskRunner,
		customerWorker: customerWorker,
	}
}

type ConductorWorker struct {
	taskRunner     *worker.TaskRunner
	customerWorker *conductor.CustomerWorker
}

func (r *ConductorWorker) Run() {
	tasks := r.customerWorker.WorkerTaskSpecs()
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
