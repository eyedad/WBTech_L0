package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Consumer struct {
	consumer *kafka.Consumer
}

func NewConsumer(brokers string, groupID string, topics []string) (*Consumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	err = c.SubscribeTopics(topics, nil)
	if err != nil {
		return nil, err
	}

	return &Consumer{consumer: c}, nil
}

func (c *Consumer) Poll(timeoutMs int) (*kafka.Message, error) {
	event := c.consumer.Poll(timeoutMs)

	switch e := event.(type) {
	case *kafka.Message:
		return e, nil
	case kafka.Error:
		if e.IsFatal() {
			return nil, e
		}
		return nil, nil
	default:
		return nil, nil
	}
}

func (c *Consumer) Close() {
	c.consumer.Close()
}
