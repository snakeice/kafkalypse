package config

import (
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var (
	ConfigDir = os.ExpandEnv("$HOME/.kafkalypse")
)

type Configuration struct {
	ApiVersion string              `yaml:"apiVersion"`
	Kind       string              `yaml:"kind"`
	Contexts   map[string]*Context `yaml:"contexts"`

	CurrentContext string `yaml:"currentContext"`

	RefreshInterval time.Duration `yaml:"refreshInterval"`
}

func LoadConfig() (*Configuration, error) {
	sync.OnceFunc(func() {
		viper.SetDefault("contexts", map[string]Context{})
		viper.SetDefault("currentContext", "")
		viper.SetDefault("refreshInterval", 10*time.Second)
		viper.SetDefault("apiVersion", "v1")
		viper.SetDefault("kind", "KafkalypseConfig")
	})

	viper.SetConfigName("config")
	viper.AddConfigPath(ConfigDir)
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {

		if err := os.MkdirAll(ConfigDir, 0755); err != nil {
			return nil, errors.Wrap(err, "failed to create config directory")
		}

		err = viper.SafeWriteConfig()
		if err != nil {
			return nil, errors.Wrap(err, "failed to create config file")
		}
	}

	var config Configuration
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Configuration) GetCurrentContext() *Context {
	return c.Contexts[c.CurrentContext]
}

func (c *Configuration) SaveConfig() error {
	return viper.WriteConfig()
}
