package config

import (
	"github.com/mohammadne/gorillamq/pkg/logger"
	"github.com/mohammadne/gorillamq/pkg/tcp"
)

type Config struct {
	Logger *logger.Config `koanf:"logger"`
	TCP    *tcp.Config    `koanf:"tcp"`
}
