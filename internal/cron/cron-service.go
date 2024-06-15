package cron

import (
	"fmt"
	"github.com/jasonlvhit/gocron"
	"log/slog"
	"rutube/internal/app/bot"
	"time"
)

type CronSerivce struct {
	storage bot.Storage
	logger  *slog.Logger
	bot     *bot.Bot
}

// New returns instance of CronService struct
func New(storage bot.Storage, logger *slog.Logger, bot *bot.Bot) *CronSerivce {
	return &CronSerivce{
		storage: storage,
		logger:  logger,
		bot:     bot,
	}
}

// Start runs daily task to check info about users
func (cs *CronSerivce) Start() {
	cs.logger.Info("Запускаю фоновые задачи")
	cs.logger.Info(time.Now().String())
	gocron.Every(1).Day().At("10:00").Do(cs.checkAndNotifyUsers)

	<-gocron.Start()

}

func (cs *CronSerivce) checkAndNotifyUsers() {
	cs.logger.Info("Начинаю оповещение")

	today := time.Now()

	birthdayUsers := cs.storage.GetUsersByDate(today)
	if birthdayUsers == nil {
		return
	}

	cs.logger.Info(fmt.Sprintf("Сегодня именинников - %d", len(birthdayUsers)))

	for _, usr := range birthdayUsers {
		for _, chatId := range usr.ChatIds {
			cs.bot.NotifyAboutUserBirthday(chatId, usr.Username)
		}
	}
}
