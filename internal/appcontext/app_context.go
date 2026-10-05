package appcontext

import (
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/logger"
	dbConfig "github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/repository/database/config"
	serverConfig "github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/server/config"
)

type AppContext struct {
	Logger       logger.Logger
	ServerConfig serverConfig.Provider
	DBConfig     dbConfig.Provider
}

func New(logger logger.Logger, serverConfig serverConfig.Provider, dbConfig dbConfig.Provider) *AppContext {
	return &AppContext{
		Logger:       logger,
		ServerConfig: serverConfig,
		DBConfig:     dbConfig,
	}
}
