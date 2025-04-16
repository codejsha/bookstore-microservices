package saga

import (
	"github.com/conductor-sdk/conductor-go/sdk/client"
	"github.com/conductor-sdk/conductor-go/sdk/settings"
	"github.com/conductor-sdk/conductor-go/sdk/worker"

	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/config"
)

func NewWorkerClient(
	conductorCfg *config.ConductorConfig,
) *WorkerClient {
	return &WorkerClient{
		TaskRunner: createClient(conductorCfg),
	}
}

type WorkerClient struct {
	TaskRunner *worker.TaskRunner
}

func createClient(conductorCfg *config.ConductorConfig) *worker.TaskRunner {
	httpSetting := settings.NewHttpSettings(conductorCfg.Client.Url)
	apiClient := client.NewAPIClient(nil, httpSetting)
	taskRunner := worker.NewTaskRunnerWithApiClient(apiClient)
	return taskRunner
}
