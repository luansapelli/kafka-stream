package environment

import (
	"strings"

	"github.com/Netflix/go-env"
)

type Environment struct {
	Application struct {
		Name string `env:"APPLICATION_NAME"`
	}
	AWS struct {
		SNS struct {
			ARN string `env:"AWS_SNS_ARN"`
		}
	}
	Kafka struct {
		Brokers   CommaSeparated `env:"KAFKA_BROKERS"`
		Topic struct {
			Input string `env:"KAFKA_TOPIC_INPUT"`
		}
	}
}

type CommaSeparated []string

func (cs *CommaSeparated) UnmarshalEnvironmentValue(data string) error {
	values := strings.Split(data, ",")

	for i, value := range values {
		values[i] = strings.TrimSpace(value)
	}

	*cs = values

	return nil
}

func Load() (*Environment, error) {
	var environment Environment

	_, err := env.UnmarshalFromEnviron(&environment)
	if err != nil {
		return nil, err
	}

	return &environment, nil
}
