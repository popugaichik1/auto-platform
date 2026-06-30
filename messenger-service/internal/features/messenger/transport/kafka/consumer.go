package transport_kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
	"go.uber.org/zap"

	core_logger "messenger-service/internal/core/logger"
	core_kafka "messenger-service/internal/core/transport/kafka"
	ws "messenger-service/internal/features/messenger/transport/ws"
)

// Hub — то немногое, что консьюмеру нужно от ws.Hub: толкнуть сообщение
// получателю, если он подключён именно к этой реплике. Сигнатура должна
// совпадать с *ws.Hub буквально (Go проверяет реализацию интерфейса по
// точному совпадению сигнатур, не просто по "совместимости" аргументов).
type Hub interface {
	SendTo(userID uuid.UUID, frame ws.ServerFrame) bool
}

type Consumer struct {
	consumer *kafka.Consumer
	hub      Hub
	log      *core_logger.Logger
}

// NewConsumer подписывается со случайным group.id — это превращает обычную
// Kafka-очередь (один консьюмер из группы получает сообщение) в pub/sub:
// КАЖДАЯ реплика messenger-service получает копию каждого события, а не
// делит их между собой. auto.offset.reset=latest, потому что свежезапущенной
// реплике не нужна историческая очередь — старые сообщения уже в Postgres
// и будут получены через REST-историю, а не через фан-аут.
func NewConsumer(config core_kafka.ConsumerCfg, hub Hub, topic string, log *core_logger.Logger) (*Consumer, error) {
	conf := kafka.ConfigMap{
		"bootstrap.servers":        config.BrokersString(),
		"group.id":                 "messenger.fanout." + uuid.New().String(),
		"auto.offset.reset":        "latest",
		"enable.auto.offset.store": false,
		"enable.auto.commit":       true,
		"auto.commit.interval.ms":  5000,
		"session.timeout.ms":       6000,
	}
	if config.SASLEnable {
    	conf["security.protocol"] = "SASL_SSL"
    	conf["sasl.mechanisms"]   = config.SASLMechanism
    	conf["sasl.username"]     = config.SASLUsername
    	conf["sasl.password"]     = config.SASLPassword
}
	consumer, err := kafka.NewConsumer(&conf)
	if err != nil {
		return nil, err
	}

	if err := consumer.SubscribeTopics([]string{topic}, nil); err != nil {
		return nil, err
	}

	return &Consumer{
		consumer: consumer,
		hub:      hub,
		log:      log,
	}, nil
}

func (c *Consumer) Run(ctx context.Context) error {
	defer func() {
		if err := c.consumer.Close(); err != nil {
			c.log.Error("fanout consumer close error", zap.Error(err))
		}
	}()

	for {
		select {
		case <-ctx.Done():
			c.log.Info("messenger fanout consumer stopped...")
			return nil
		default:
		}

		//одни и те же сообщения читаются на всех трех репликах. Если на одном из подов, подключение
		//клиента будет в хабе, то сработает sendTo. В противном случае, сообщение просто сохраниться в БД
		//Sticky-сессия не подходит, так как Sticky-сессия "Васи" не говорит нам о том на каком поде 
		//находиться "Саня".
		msg, err := c.consumer.ReadMessage(100 * time.Millisecond)
		if err != nil {
			if err.(kafka.Error).Code() == kafka.ErrTimedOut {
				continue
			}
			c.log.Error("read kafka message error", zap.Error(err))
			continue
		}

		// Источник истины — Postgres, сообщение туда уже сохранено до
		// публикации этого события. Битое событие здесь просто пропускаем:
		// в худшем случае получатель не увидит сообщение мгновенно и
		// подхватит его при следующем REST-запросе истории. DLQ для этого
		// побочного канала доставки не нужен.
		var event MessageSentEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			c.log.Error("unmarshal message sent event error", zap.Error(err))
			continue
		}

		c.hub.SendTo(event.RecipientID, ws.ServerFrame{
			Type:    "message",
			Payload: event,
		})
	}
}
