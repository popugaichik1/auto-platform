package transport_kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	core_errors "user-service/internal/core/errors"
	core_logger "user-service/internal/core/logger"
	core_kafka "user-service/internal/core/transport/kafka"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Consumer struct {
	consumer    *kafka.Consumer
	service     Service
	log         *core_logger.Logger
	dlqProducer *Producer
}


func NewConsumer(config core_kafka.ConsumerCfg, service Service, topic string, log *core_logger.Logger, dlqProducer *Producer) (*Consumer, error) {
	conf := kafka.ConfigMap{
		"bootstrap.servers":          config.BrokersString(),
		"group.id":                   core_kafka.RegisterUserConsumerGroup,
		"auto.offset.reset":          "earliest",
		"enable.auto.offset.store":   false,
		"queued.max.messages.kbytes": 1000,
		"enable.auto.commit":         true,
		"auto.commit.interval.ms":    5000,
		// FIX
		"session.timeout.ms": 6000,
	}

	if config.SASLEnable {
    	conf["security.protocol"] = "SASL_SSL"
    	conf["sasl.mechanisms"]   = config.SASLMechanism
    	conf["sasl.username"]     = config.SASLUsername
    	conf["sasl.password"]     = config.SASLPassword
	}
	consumer, err := kafka.NewConsumer(&conf)

	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to the topic: %w", err)
	}

	return &Consumer{
		consumer:    consumer,
		service:     service,
		log:         log,
		dlqProducer: dlqProducer,
	}, nil
}

type Service interface {
	SaveUser(
		ctx context.Context,
		ID uuid.UUID,
		username string,
		phoneNumber string,
	) error
}

func (c *Consumer) Run(ctx context.Context) error {
	var event UserRegisterEvent

	defer func() {
		if err := c.consumer.Close(); err != nil {
			c.log.Error("consumer close error: ", zap.Error(err))
		}
	}()

	for {
		select{
		case <-ctx.Done():
			c.log.Info("Kafka consumer stopped...")
			return nil

		default:
		}

		msg, err := c.consumer.ReadMessage(100 * time.Millisecond)
		if err != nil {
			if err.(kafka.Error).Code() == kafka.ErrTimedOut {
				continue
			}
			c.log.Error("read kafka message error", zap.Error(err))
			continue
		}

		// Битый JSON никогда не станет валидным при повторной доставке —
		// в отличие от ошибок БД, повтор тут бессмысленен, сообщение уходит в DLQ.
		if err = json.Unmarshal(msg.Value, &event); err != nil {
			c.log.Error("unmarshal user register event error", zap.Error(err))
			c.sendToDLQ(ctx, msg, err)
			continue
		}

		userID, err := uuid.Parse(event.ID)
		if err != nil {
			c.log.Error("parse user id from event error", zap.String("id", event.ID), zap.Error(err))
			c.sendToDLQ(ctx, msg, err)
			continue
		}

		if err = c.service.SaveUser(ctx, userID, event.Username, event.PhoneNumber); err != nil {
			c.log.Error("save user error", zap.Error(err))
			// ErrInvalidArgument — постоянная ошибка (данные в событии сами
			// по себе невалидны, повтор не поможет) -> DLQ. Всё остальное
			// (например, БД временно недоступна) оставляем без коммита,
			// чтобы Kafka передоставила сообщение позже.
			if errors.Is(err, core_errors.ErrInvalidArgument) {
				c.sendToDLQ(ctx, msg, err)
			}
			continue
		}

		_, err = c.consumer.CommitMessage(
			msg,
		)
		if err != nil {
			c.log.Error("commit kafka message error", zap.Error(err))
		}

	}
}

// sendToDLQ публикует сообщение, которое не удалось обработать без шанса
// на успех при повторной попытке, в DLQ-топик и коммитит исходное сообщение
// (иначе Kafka продолжит передоставлять то, что мы уже признали "ядовитым").
func (c *Consumer) sendToDLQ(ctx context.Context, msg *kafka.Message, cause error) {
	dlqEvent := DeadLetterEvent{
		OriginalTopic:     *msg.TopicPartition.Topic,
		OriginalPartition: msg.TopicPartition.Partition,
		OriginalOffset:    int64(msg.TopicPartition.Offset),
		Error:             cause.Error(),
		FailedAt:          time.Now(),
		Payload:           string(msg.Value),
	}

	if err := c.dlqProducer.Publish(ctx, core_kafka.NewMessage(
		[]byte(core_kafka.TopicUserRegisteredDLQ),
		msg.Key,
		dlqEvent,
	)); err != nil {
		c.log.Error("failed to publish to DLQ", zap.Error(err))
		return
	}

	if _, err := c.consumer.CommitMessage(msg); err != nil {
		c.log.Error("commit kafka message error after DLQ publish", zap.Error(err))
	}
}
