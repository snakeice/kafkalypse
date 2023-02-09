package config

type Context struct {
	BootstrapServers string `yaml:"bootstrapServers"`
	Username         string `yaml:"username"`
	Password         string `yaml:"password"`

	SchemaRegistryUrl string `yaml:"schemaRegistryUrl"`
	RegistryUsername  string `yaml:"registryUsername"`
	RegistryPassword  string `yaml:"registryPassword"`
}
