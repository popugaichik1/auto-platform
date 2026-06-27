package core_kafka

import (
	"fmt"
	"strings"

	"github.com/kelseyhightower/envconfig"
)


type ConsumerCfg struct {
	Brokers []string  `envconfig:"BROKERS" required:"true"`
}


type ProducerConfig struct {
	Brokers []string `envconfig:"BROKERS" required:"true"`
}


func NewConsumerConfig() (ConsumerCfg, error) {
	var config ConsumerCfg

	if err := envconfig.Process("KAFKA", &config); err != nil {
		return ConsumerCfg{}, fmt.Errorf("process envconfig: %w", err)
	}
	return config, nil
}


func NewProducerConfig() (ProducerConfig, error) {
	var config ProducerConfig

	if err := envconfig.Process("KAFKA", &config); err != nil {
		return ProducerConfig{}, fmt.Errorf("process envconfig: %w", err)
	}
	return config, nil
}


func NewConsumerConfigMust() ConsumerCfg {
	config, err := NewConsumerConfig()
	if err != nil {
		err = fmt.Errorf("get Kafka producer config: %w", err)
		panic(err)
	}
	return config
}


func NewProducerConfigMust() ProducerConfig {
	config, err := NewProducerConfig()
	if err != nil {
		err = fmt.Errorf("get kafka producer config: %w", err)
		panic(err)
	}
	return config
}


func (c ConsumerCfg) BrokersString() string {
	return strings.Join(c.Brokers, ",")
}


func (c ProducerConfig) BrokersString() string {
	return strings.Join(c.Brokers, ",")
}

