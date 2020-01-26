package configs

import (
	"github.com/caarlos0/env/v6"
	"time"
)

// Config - структура содержащая все изменяемые конфигурации приложения
type Config struct {
	MaxPage       int    `env:"MAXPAGE"`
	TelegramToken string `env:"TOKEN"`

	ManagerDelayMin time.Duration
}

// GetConf - возвращет конфигурацию приложения
func GetConf() (*Config, error) {
	cfg := &Config{
		// здесь можно выставить параметры по умолчанию
		MaxPage:         2,
		ManagerDelayMin: 1,
	}

	err := env.Parse(cfg)
	return cfg, err
}
