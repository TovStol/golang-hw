package main

import (
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Level        string
	LogLocation  string
	DbDriverName string
	Dsn          string
	KafkaBrokers []string
	KafkaTopic   string
	ScanInterval duration
}

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

func NewConfig(configFile string) Config {
	var config Config
	config.ScanInterval = duration{10 * time.Second}
	_, err := toml.DecodeFile(configFile, &config)
	if err != nil {
		fmt.Println("config error:", err)
	}
	return config
}
