package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/app_context"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/logger"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/server"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/server/config"
)

func main() {

	ctx, ctxCancel := context.WithCancel(context.Background())
	logApp := logger.New()

	envSrv, err := config.NewEnv()
	if err != nil {
		logApp.Error("Error initializing environment service", "error", err)
		return
	}
	cfgServer := config.New(envSrv)

	cfgServer.AddressSet("localhost:8080") //TODO get from flag or env
	appContext := app_context.New(logApp, cfgServer)

	var wg sync.WaitGroup

	wg.Add(1)
	// запускаем http сервер
	go func() {
		defer wg.Done()

		serverApp := server.New(appContext)

		err := serverApp.Run(ctx)
		if err != nil {
			logApp.Error("Server.Shutdown", "error", err)
		}
	}()

	wg.Add(1)
	// ожидаем сигнал завершения приложения
	go func() {
		defer wg.Done()
		wait := make(chan os.Signal, 1)
		signal.Notify(wait, syscall.SIGINT, syscall.SIGTERM)
		<-wait
		ctxCancel()
	}()

	wg.Wait()
}
