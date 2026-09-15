package support

import (
	sharedconfig "github.com/codejsha/shared-library-go/pkg/config"
	"github.com/codejsha/shared-library-go/pkg/message"

	"github.com/codejsha/bookstore-microservices/inventory/internal/infrastructure/adapter/restcontroller"
)

func NewKafkaPublisher(kafkaCfg *sharedconfig.KafkaConfig) (*message.KafkaAsyncPublisher, error) {
	pub, err := message.NewKafkaAsyncPublisher(kafkaCfg.Broker.BootstrapServers, message.NewKafkaAsyncProducerConfig())
	if err != nil {
		return nil, err
	}
	return pub, nil
}

func RegisterEventPublisher(pub *message.KafkaAsyncPublisher) {
	restcontroller.SetEventPublisher(pub)
}
