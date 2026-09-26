package appcontext

import (
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/logger"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/server/config"
)

type AppContext struct {
	Logger       logger.Logger
	ServerConfig config.Provider
}

func New(logger logger.Logger, serverConfig config.Provider) *AppContext {
	return &AppContext{
		Logger:       logger,
		ServerConfig: serverConfig,
	}
}
