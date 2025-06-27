package config

import (
	"github.com/twmb/franz-go/pkg/sasl"
	"github.com/twmb/franz-go/pkg/sasl/scram"
)

type SaslMechanism struct {
	Mechanism string `yaml:"mechanism"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
}

type KafkaContext struct {
	Name             string        `yaml:"name"`
	BootstrapServers string        `yaml:"bootstrapServers"`
	SecurityProtocol string        `yaml:"securityProtocol"`
	Sasl             SaslMechanism `yaml:"sasl"`
}

func (s *SaslMechanism) AsMechanism() sasl.Mechanism {
	auth := scram.Auth{
		User: s.Username,
		Pass: s.Password,
	}

	switch s.Mechanism {
	case "SCRAM-SHA-256":
		return auth.AsSha256Mechanism()
	case "SCRAM-SHA-512":
		return auth.AsSha512Mechanism()
	}

	return nil
}
