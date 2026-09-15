package support

import (
	sharedconfig "github.com/codejsha/shared-library-go/pkg/config"
	"github.com/codejsha/shared-library-go/pkg/message"

	"github.com/codejsha/bookstore-microservices/catalog/internal/infrastructure/adapter/restcontroller"
)

func NewKafkaPublisher(kafkaCfg *sharedconfig.KafkaConfig) (*message.KafkaAsyncPublisher, error) {
	return message.NewKafkaAsyncPublisher(kafkaCfg.Broker.BootstrapServers, message.NewKafkaAsyncProducerConfig())
}

func RegisterEventPublisher(pub *message.KafkaAsyncPublisher) {
	restcontroller.SetEventPublisher(pub)
}
