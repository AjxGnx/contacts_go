package config

import (
	"sync"

	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/gommon/log"
)

type Config struct {
	ServerHost string `required:"true" split_words:"true"`
	ServerPort int    `required:"true" split_words:"true"`
	DBHost     string `required:"true" split_words:"true"`
	DBPort     int    `required:"true" split_words:"true"`
	DBUser     string `required:"true" split_words:"true"`
	DBName     string `required:"true" split_words:"true"`
	DBPass     string `required:"true" split_words:"true"`
}

var once sync.Once
var config Config

func Environments() Config {
	once.Do(func() {
		if err := envconfig.Process("", &config); err != nil {
			log.Panicf("Error parsing environment vars %#v", err)
		}
	})

	return config
}
