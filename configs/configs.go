package configs

import (
	"github.com/caarlos0/env/v6"
)

// Config - структура содержащая все изменяемые конфигурации приложения
type Config struct {
	MaxPage       int    `env:"MAXPAGE"`
	TelegramToken string `env:"TOKEN"`
}

// GetConf - возвращет конфигурацию приложения
func GetConf() (*Config, error) {
	cfg := &Config{
		// здесь можно выставить параметры по умолчанию
	}

	err := env.Parse(cfg)
	return cfg, err
}
