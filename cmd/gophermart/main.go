package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/appcontext"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/client"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/logger"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/repository/database"
	configDb "github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/repository/database/config"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/server"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/server/config"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/service/loyalty"
	"golang.org/x/sync/errgroup"
)

type applicationFlags struct {
	runAddress           *string
	databaseURI          *string
	accrualSystemAddress *string
}

func main() {

	ctx, cancelApp := context.WithCancel(context.Background())
	defer cancelApp()

	g, ctx := errgroup.WithContext(ctx)
	logApp := logger.New()

	envSrv, err := config.NewEnv()
	if err != nil {
		logApp.Error("Error initializing environment service", "error", err)
		os.Exit(1)
	}

	envDB, err := configDb.NewEnv()
	if err != nil {
		logApp.Error("Error initializing environment db", "error", err)
		os.Exit(1)
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

	flag.Parse()

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
		os.Exit(1)
	}
	defer db.Close()

	loyaltyStorage := database.NewStorage(db)
	loyaltyManager := loyalty.NewManager(loyaltyStorage)

	// запускаем http сервер
	g.Go(func() error {
		serverApp := server.New(appContext, loyaltyManager)

		return serverApp.Run(ctx)
	})

	g.Go(func() error {
		accrualReceiver := client.NewAccrualReceiver(appContext, loyaltyManager)
		accrualReceiver.Run(ctx)
		return nil
	})

	// Ждём сигнал завершения.
	g.Go(func() error {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-stop:
			logApp.Info("server shutting down")
			cancelApp()
			return nil
		case <-ctx.Done():
			return nil
		}
	})

	if err = g.Wait(); err != nil {
		logApp.Error("application stopped", "error", err)
		os.Exit(1)
	}

	logApp.Info("application stopped gracefully")
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
