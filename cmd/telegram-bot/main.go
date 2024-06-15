package main

import (
	"fmt"
	"os"
	"rutube/internal/app/bot"
	"rutube/internal/config"
	"rutube/internal/cron"
	"rutube/internal/logging"
	postgres_storage "rutube/internal/storage/postgres-storage"
)

func main() {
	cfg := config.MustGetConfig()

	logger := logging.MustGetLogger(cfg.Env)
	logger.Info(fmt.Sprintf("%+v", cfg))

	storage, err := postgres_storage.New(cfg)
	if err != nil {
		panic(err)
	}

	err = storage.Migrate()
	if err != nil {
		panic(err)
	}

	stop := make(chan os.Signal, 1)

	bot := bot.New(&logger, cfg, storage)
	go bot.Start()

	cronService := cron.New(storage, &logger, bot)
	go cronService.Start()

	for {
		select {
		case <-stop:
			logger.Info("application stopped")
			os.Exit(1)
		}
	}
}
