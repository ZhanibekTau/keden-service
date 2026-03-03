package rabbitmq

type ConsumerConfig struct {
	Topic       string
	Exchange    string
	Queue       string
	HandlerName string
	Prefetch    int
}

func NewConsumerTopicDurableConfig(topic, exchange, queue string, prefetch int) []ConsumerConfig {
	return []ConsumerConfig{
		{
			Topic:       topic,
			Exchange:    exchange,
			Queue:       queue,
			HandlerName: queue + ".handler",
			Prefetch:    prefetch,
		},
	}
}
