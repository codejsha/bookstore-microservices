package conductor

import (
	"fmt"

	"github.com/conductor-sdk/conductor-go/sdk/model"

	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/config"
	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/saga"
)

func NewCustomerWorker(
	conductorCfg *config.ConductorConfig,
) *CustomerWorker {
	customerWorker := &CustomerWorker{
		conductorCfg:    conductorCfg,
		workerTaskSpecs: []saga.WorkerTaskSpec{},
	}
	return customerWorker
}

type CustomerWorker struct {
	conductorCfg    *config.ConductorConfig
	workerTaskSpecs []saga.WorkerTaskSpec
}

func (w *CustomerWorker) WorkerTaskSpecs() []saga.WorkerTaskSpec {
	return w.workerTaskSpecs
}

func Greet(task *model.Task) (interface{}, error) {
	return map[string]interface{}{
		"greetings": "Hello, " + fmt.Sprintf("%v", task.InputData["name"]),
	}, nil
}
