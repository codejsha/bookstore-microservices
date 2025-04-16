package conductor

import (
	"fmt"

	"github.com/conductor-sdk/conductor-go/sdk/model"

	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/config"
	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/saga"
)

func NewBookWorker(
	conductorCfg *config.ConductorConfig,
) *BookWorker {
	bookWorker := &BookWorker{
		conductorCfg:    conductorCfg,
		workerTaskSpecs: []saga.WorkerTaskSpec{},
	}
	return bookWorker
}

type BookWorker struct {
	conductorCfg    *config.ConductorConfig
	workerTaskSpecs []saga.WorkerTaskSpec
}

func (w *BookWorker) WorkerTaskSpecs() []saga.WorkerTaskSpec {
	return w.workerTaskSpecs
}

func Greet(task *model.Task) (interface{}, error) {
	return map[string]interface{}{
		"greetings": "Hello, " + fmt.Sprintf("%v", task.InputData["name"]),
	}, nil
}
