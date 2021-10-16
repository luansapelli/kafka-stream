package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/lovoo/goka"
)

type SaramaConfig struct {
	ApplicationName string
}

// Config populate and returns a *SaramaConfig struct.
func Config(applicationName string) *SaramaConfig {
	return &SaramaConfig{
		ApplicationName: applicationName,
	}
}

// Sarama method serves to set the configs, like authentication, producers, consumers, etc...
func (config *SaramaConfig) Sarama() *sarama.Config {
	saramaConfig := sarama.NewConfig()

	saramaConfig.ClientID = config.ApplicationName
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Partitioner = sarama.NewManualPartitioner

	goka.ReplaceGlobalConfig(saramaConfig)

	return saramaConfig
}
