package config

import (
	"github.com/spf13/viper"
)

type ServerConfig struct {
	LogLevel string
	Host     string
	Port     int
}

func LoadServerConfig() (ServerConfig, error) {
	const logLevel = "log_level"
	const host = "host"
	const port = "port"

	viper.SetDefault(logLevel, "info")
	viper.SetDefault(host, "0.0.0.0")
	viper.SetDefault(port, 3000)

	return ServerConfig{
		LogLevel: viper.GetString(logLevel),
		Host:     viper.GetString(host),
		Port:     viper.GetInt(port),
	}, nil
}
