package config

import (
	"github.com/mohammadne/gorillamq/internal/network"
	"github.com/mohammadne/gorillamq/pkg/logger"
)

type Config struct {
	Logger  *logger.Config  `koanf:"logger"`
	Network *network.Config `koanf:"network"`
}
