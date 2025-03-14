package config

type KafkaContext struct {
	Name             string `yaml:"name"`
	BootstrapServers string `yaml:"bootstrapServers"`
	SecurityProtocol string `yaml:"securityProtocol"`
	SaslMechanism    string `yaml:"saslMechanism"`
	SaslUsername     string `yaml:"saslUsername"`
	SaslPassword     string `yaml:"saslPassword"`
}
