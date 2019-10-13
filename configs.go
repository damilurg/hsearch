package main

import (
	"github.com/caarlos0/env/v6"
)

type config struct {
	MaxPage       string `env:"MAXPAGE"`
	TelegramToken string `env:"TOKEN"`
}

func GetConf() (*config, error) {
	cfg := &config{}
	err := env.Parse(cfg)
	return cfg, err
}
