package rabbitmq

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/sirupsen/logrus"
)

type AmqpPubSub struct {
	subscriber *amqp.Subscriber
	publisher  *amqp.Publisher
	router     *message.Router
	config     amqp.ConnectionConfig
}

func NewAmqpPubSub(connConfig amqp.ConnectionConfig) (*AmqpPubSub, error) {
	wmLogger := watermill.NewStdLogger(false, false)

	subscriberConfig := amqp.NewDurableQueueConfig(connConfig.AmqpURI)
	subscriber, err := amqp.NewSubscriber(subscriberConfig, wmLogger)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscriber: %w", err)
	}

	publisherConfig := amqp.NewDurablePubSubConfig(connConfig.AmqpURI, nil)
	publisher, err := amqp.NewPublisher(publisherConfig, wmLogger)
	if err != nil {
		return nil, fmt.Errorf("failed to create publisher: %w", err)
	}

	router, err := message.NewRouter(message.RouterConfig{}, wmLogger)
	if err != nil {
		return nil, fmt.Errorf("failed to create router: %w", err)
	}

	logrus.Info("RabbitMQ connection established successfully")

	return &AmqpPubSub{
		subscriber: subscriber,
		publisher:  publisher,
		router:     router,
		config:     connConfig,
	}, nil
}

func (a *AmqpPubSub) Publish(topic string, msg *message.Message) error {
	return a.publisher.Publish(topic, msg)
}

func (a *AmqpPubSub) RegisterHandler(
	handlerFunc func(ctx context.Context, msg *message.Message) error,
	configs ...ConsumerConfig,
) error {
	if len(configs) == 0 {
		return fmt.Errorf("no consumer config provided")
	}

	cfg := configs[0]

	a.router.AddNoPublisherHandler(
		cfg.HandlerName,
		cfg.Topic,
		a.subscriber,
		func(msg *message.Message) error {
			return handlerFunc(context.Background(), msg)
		},
	)

	return nil
}

func (a *AmqpPubSub) Consume(ctx context.Context) error {
	return a.router.Run(ctx)
}

func (a *AmqpPubSub) CloseConnection() {
	if a.router != nil {
		_ = a.router.Close()
	}
	if a.publisher != nil {
		_ = a.publisher.Close()
	}
	if a.subscriber != nil {
		_ = a.subscriber.Close()
	}
}
