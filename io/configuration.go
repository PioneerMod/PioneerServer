package io

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Server ServerConfig `toml:"server"`
	Logger LoggerConfig `toml:"logger"`
}

type LoggerConfig struct {
	Level      int    `toml:"level"`
	TimeFormat string `toml:"time_format"`
}

type ServerConfig struct {
	IsLocal  bool    `toml:"is_local"`
	Port     int     `toml:"port"`
	Password *string `toml:"password"`
}

var config *Config

func GetConfig() *Config {
	if config != nil {
		return config
	}

	var newConfig = getDefaultConfig()

	if _, err := os.Stat("config.toml"); os.IsNotExist(err) {
		return &newConfig
	}

	_, err := toml.DecodeFile("config.toml", &newConfig)
	if err != nil {
		panic(err)
	}

	config = &newConfig

	return config
}

func getDefaultConfig() Config {
	return Config{
		Server: ServerConfig{
			IsLocal:  true,
			Port:     25319,
			Password: nil,
		},
		Logger: LoggerConfig{
			Level:      1,
			TimeFormat: "15:04:05.000",
		},
	}
}
