package listing_client

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type GRPCClientConfig struct {
	Addr    string  `envconfig:"ADDR" required:"true"`
}

func NewConfig() (GRPCClientConfig, error) {
	var config GRPCClientConfig

	if err := envconfig.Process("GRPC", &config); err != nil {
		return GRPCClientConfig{}, fmt.Errorf("process envconfig: %w", err)
	}

	return config, nil
}

func NewConfigMust() GRPCClientConfig {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get gRPC client config: %w", err)
		panic(err)
	}

	return config
}
