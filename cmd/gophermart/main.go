package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/appcontext"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/logger"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/repository/database"
	configDb "github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/repository/database/config"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/server"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/server/config"
)

type applicationFlags struct {
	runAddress           *string
	databaseURI          *string
	accrualSystemAddress *string
}

func main() {

	ctx, ctxCancel := context.WithCancel(context.Background())
	logApp := logger.New()

	envSrv, err := config.NewEnv()
	if err != nil {
		logApp.Error("Error initializing environment service", "error", err)
		return
	}

	envDB, err := configDb.NewEnv()
	if err != nil {
		logApp.Error("Error initializing environment db", "error", err)
		return
	}

	/*
		- адрес и порт запуска сервиса: переменная окружения ОС `RUN_ADDRESS` или флаг `-a`
		- адрес подключения к базе данных: переменная окружения ОС `DATABASE_URI` или флаг `-d`
		- адрес системы расчёта начислений: переменная окружения ОС `ACCRUAL_SYSTEM_ADDRESS` или флаг `-r`
	*/
	appFlags := &applicationFlags{}
	appFlags.runAddress = flag.String("a", "localhost:8080", `адрес и порт запуска сервиса`)
	appFlags.databaseURI = flag.String("d", "", `адрес подключения к базе данных`)
	appFlags.accrualSystemAddress = flag.String("r", "localhost:8081", `адрес системы расчёта начислений`)

	cfgServer := config.New(envSrv)
	cfgDB := configDb.New(envDB)

	serverConfigUpdateByFlags(cfgServer, appFlags)

	if cfgDB.DSN() == "" && appFlags.databaseURI != nil {
		cfgDB.DSNSet(*appFlags.databaseURI)
	}

	appContext := appcontext.New(logApp, cfgServer, cfgDB)
	db, err := database.NewDB(ctx, cfgDB, logApp)
	if err != nil {
		logApp.Error("Error initializing database", "error", err)
		return
	}
	loyaltyStorage := database.NewStorage(db)

	var wg sync.WaitGroup

	wg.Add(1)
	// запускаем http сервер
	go func() {
		defer wg.Done()

		serverApp := server.New(appContext, loyaltyStorage)

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

func serverConfigUpdateByFlags(serverConfig *config.Config, appFlags *applicationFlags) {

	if serverConfig == nil || appFlags == nil {
		return
	}

	if serverConfig.Address() == "" && appFlags.runAddress != nil {
		serverConfig.AddressSet(*appFlags.runAddress)
	}

	if serverConfig.AccrualSystemAddress() == "" && appFlags.accrualSystemAddress != nil {
		serverConfig.AccrualSystemAddressSet(*appFlags.accrualSystemAddress)
	}
}
