package server

import (
	"strings"

	"github.com/spf13/viper"
)

// Config defines server configuration.
type Config struct {
	Host     string `mapstructure:"host" yaml:"host"`
	Port     string `mapstructure:"port" yaml:"port"`
	Mode     string `mapstructure:"mode" yaml:"mode"`
	RootPath string `mapstructure:"rootPath" yaml:"rootPath"`
}

func (c *Config) GetAddr() string {
	return c.Host + ":" + c.Port
}

func DefaultConfig() *Config {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.rootPath", "")

	return &Config{
		Host:     viper.GetString("server.host"),
		Port:     viper.GetString("server.port"),
		Mode:     viper.GetString("server.mode"),
		RootPath: viper.GetString("server.rootPath"),
	}
}
