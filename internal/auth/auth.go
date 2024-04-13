package auth

import (
	"strings"

	"go.uber.org/zap"
)

type Auth interface {
	Authenticate(string) bool
}

type auth struct {
	logger *zap.Logger
	config *Config
}

func NewAuth(logger *zap.Logger, config *Config) Auth {
	return &auth{logger: logger, config: config}
}

func (auth *auth) Authenticate(token string) bool {
	if auth.config.Username == " " || auth.config.Password == " " {
		return true
	}
	parts := strings.Split(token, ":")
	return parts[0] == auth.config.Username && parts[1] == auth.config.Password
}
