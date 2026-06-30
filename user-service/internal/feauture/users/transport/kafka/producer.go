package transport_kafka

import (
	"context"
	"encoding/json"
	"fmt"

	core_kafka "user-service/internal/core/transport/kafka"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Producer struct {
	producer *kafka.Producer
}

func NewProducer(config core_kafka.ProducerConfig) (*Producer, error) {
	conf := kafka.ConfigMap{
		"bootstrap.servers": config.BrokersString(),
		"partitioner":       "random",
	}
	
	if config.SASLEnable {
    	conf["security.protocol"] = "SASL_SSL"
    	conf["sasl.mechanisms"]   = config.SASLMechanism
    	conf["sasl.username"]     = config.SASLUsername
    	conf["sasl.password"]     = config.SASLPassword
	}
	
	p, err := kafka.NewProducer(&conf)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}
	return &Producer{
		producer: p,
	}, nil
}

func (p *Producer) Publish(ctx context.Context, message core_kafka.Message) error {
	body, err := json.Marshal(message.Payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	topic := string(message.Topic)

	deliveryChan := make(chan kafka.Event, 1)
	defer close(deliveryChan)

	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: body,
		Key:   message.Key,
	}

	if err := p.producer.Produce(kafkaMsg, deliveryChan); err != nil {
		return fmt.Errorf("produce message: %w", err)
	}

l:
	select {
	case <-ctx.Done():
		break l
	case event := <-deliveryChan:
		msg := event.(*kafka.Message)

		if msg.TopicPartition.Error != nil {
			return msg.TopicPartition.Error
		}
	}

	return nil
}

func (p *Producer) Close() {
	p.producer.Flush(core_kafka.FlushTimeOut)
	p.producer.Close()
}
