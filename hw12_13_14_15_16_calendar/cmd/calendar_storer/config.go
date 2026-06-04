package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Level        string
	LogLocation  string
	DbDriverName string
	Dsn          string
	KafkaBrokers []string
	KafkaTopic   string
	KafkaGroupID string
}

func NewConfig(configFile string) Config {
	var config Config
	config.KafkaGroupID = "storer-group"
	_, err := toml.DecodeFile(configFile, &config)
	if err != nil {
		fmt.Println("config error:", err)
	}
	return config
}
