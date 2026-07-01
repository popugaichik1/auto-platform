package core_kafka

import (
	"fmt"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type ProducerConfig struct {
    Brokers          []string `envconfig:"BROKERS"           required:"true"`
    SASLEnable       bool     `envconfig:"SASL_ENABLE"       default:"false"`
    SASLMechanism    string   `envconfig:"SASL_MECHANISM"    default:"SCRAM-SHA-512"`
    SASLUsername     string   `envconfig:"SASL_USERNAME"`
    SASLPassword     string   `envconfig:"SASL_PASSWORD"`
    SecurityProtocol string   `envconfig:"SECURITY_PROTOCOL" default:"SASL_PLAINTEXT"`
}

func NewProducerConfig() (ProducerConfig, error) {
	var config ProducerConfig

	if err := envconfig.Process("KAFKA", &config); err != nil {
		return ProducerConfig{}, fmt.Errorf("process envconfig: %w", err)
	}

	return config, nil
}

func NewProducerConfigMust() ProducerConfig {
	config, err := NewProducerConfig()
	if err != nil {
		err = fmt.Errorf("get Kafka producer config: %w", err)
		panic(err)
	}

	return config
}

func (c ProducerConfig) BrokersString() string {
	return strings.Join(c.Brokers, ",")
}


type ConsumerCfg struct {
    Brokers          []string `envconfig:"BROKERS"           required:"true"`
    SASLEnable       bool     `envconfig:"SASL_ENABLE"       default:"false"`
    SASLMechanism    string   `envconfig:"SASL_MECHANISM"    default:"SCRAM-SHA-512"`
    SASLUsername     string   `envconfig:"SASL_USERNAME"`
    SASLPassword     string   `envconfig:"SASL_PASSWORD"`
    SecurityProtocol string   `envconfig:"SECURITY_PROTOCOL" default:"SASL_PLAINTEXT"`
}

func NewConsumerConfig() (ConsumerCfg, error) {
	var config ConsumerCfg

	if err := envconfig.Process("KAFKA", &config); err != nil {
		return ConsumerCfg{}, fmt.Errorf("process envconfig: %w", err)
	}

	return config, nil
}

func NewConsumerConfigMust() ConsumerCfg {
	config, err := NewConsumerConfig()
	if err != nil {
		err = fmt.Errorf("get Kafka consumer config: %w", err)
		panic(err)
	}

	return config
}

func (c ConsumerCfg) BrokersString() string {
	return strings.Join(c.Brokers, ",")
}
