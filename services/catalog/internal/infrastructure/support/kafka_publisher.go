package support

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"go.uber.org/fx"

	sharedconfig "github.com/codejsha/shared-library-go/pkg/config"
	"github.com/codejsha/shared-library-go/pkg/message"

	"github.com/codejsha/bookstore-microservices/catalog/internal/infrastructure/adapter/restcontroller"
)

const sideEffectsDrainTimeout = 10 * time.Second

func NewKafkaPublisher(lc fx.Lifecycle, kafkaCfg *sharedconfig.KafkaConfig) (*message.KafkaAsyncPublisher, error) {
	pub, err := message.NewKafkaAsyncPublisher(kafkaCfg.Broker.BootstrapServers, message.NewKafkaAsyncProducerConfig())
	if err != nil {
		return nil, err
	}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return pub.Close()
		},
	})
	return pub, nil
}

func RegisterEventPublisher(lc fx.Lifecycle, pub *message.KafkaAsyncPublisher) {
	restcontroller.SetEventPublisher(pub)
	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			if !restcontroller.WaitForSideEffects(sideEffectsDrainTimeout) {
				logrus.Warn("timed out draining side-effect goroutines before shutdown")
			}
			return nil
		},
	})
}
