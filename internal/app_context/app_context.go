package app_context

import (
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/logger"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/server/config"
)

type AppContext struct {
	Logger       logger.Logger
	ServerConfig *config.Config
}

func New(logger logger.Logger, serverConfig *config.Config) *AppContext {
	return &AppContext{
		Logger:       logger,
		ServerConfig: serverConfig,
	}
}
