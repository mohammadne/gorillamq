package config

import (
	"github.com/mohammadne/gorillamq/pkg/logger"
)

type Config struct {
	Logger *logger.Config `koanf:"logger"`
}
