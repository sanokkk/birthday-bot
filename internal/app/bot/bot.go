package bot

import (
	"context"
	"fmt"
	tgbot "github.com/go-telegram/bot"
	botmodels "github.com/go-telegram/bot/models"
	"github.com/lib/pq"
	"log/slog"
	"os"
	"os/signal"
	"rutube/internal/config"
	"rutube/internal/models"
	"slices"
	"time"
)

const (
	welcomeTemplate        = "Приветствую, @%s, введите свою дату рождения в формате дд-ММ-ГГГГ"
	exitTemplate           = "Пользователь @%s покинул чат"
	congratulationTemplate = "Поздравляем @%s с Днем Рождения!!!"
)

type Storage interface {
	GetUserInfo(id int64) (*models.User, error)
	AddUserInfo(*models.User) (*models.User, error)
	GetUsersByDate(time time.Time) []models.User
	UpdateUser(id int64, user *models.User)
}

type Bot struct {
	storage Storage
	logger  *slog.Logger
	config  *config.Config
	Bot     *tgbot.Bot
}

// New creates Bot instance
func New(logger *slog.Logger, cfg *config.Config, storage Storage) *Bot {
	return &Bot{
		logger:  logger,
		config:  cfg,
		storage: storage,
	}
}

func (b *Bot) Start() {
	const op = "bot:Start"
	log := b.logger.With(slog.String("operation", op))

	log.Info("starting bot")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []tgbot.Option{
		tgbot.WithDefaultHandler(b.defaultHandler),
	}

	bot, err := tgbot.New(b.config.BotKey, opts...)
	if err != nil {
		log.Error("error while creating bot from config", err)
	}

	//bot.Close(ctx)
	b.Bot = bot
	b.Bot.Start(ctx)
}

func (b *Bot) defaultHandler(ctx context.Context, bot *tgbot.Bot, update *botmodels.Update) {
	msg := update.Message

	if msg != nil {
		if len(msg.NewChatMembers) > 0 {
			b.handleNewUsers(ctx, bot, msg)
		}

		if msg.LeftChatMember != nil {
			b.handleLeftMember(ctx, bot, msg)
		}

		if msg.Text != "" {
			b.handleMessage(ctx, msg)
		}
	}
}

func (b *Bot) handleMessage(ctx context.Context, msg *botmodels.Message) {
	usr, err := b.storage.GetUserInfo(msg.From.ID)
	if err != nil {
		b.logger.Warn(err.Error())
	}

	if usr != nil {
		if !slices.Contains(usr.ChatIds, msg.Chat.ID) {
			newChats := append(usr.ChatIds, msg.Chat.ID)
			b.storage.UpdateUser(usr.Id, &models.User{ChatIds: newChats})
		}
	}

	if usr == nil {
		from := msg.From
		newUser, _ := b.storage.AddUserInfo(
			&models.User{
				Id:       from.ID,
				Username: from.Username,
				ChatIds:  pq.Int64Array([]int64{msg.Chat.ID})})
		usr = newUser
	}

	birthday, err := time.Parse("02-01-2006", msg.Text)

	if err != nil {
		b.logger.Info("сообщение не содержит даты рождения")

		if usr.Birthday == nil {
			b.Bot.SendMessage(ctx,
				&tgbot.SendMessageParams{
					Text:   fmt.Sprintf(welcomeTemplate, usr.Username),
					ChatID: msg.Chat.ID,
				})
		}

		return
	}

	usr.Birthday = &birthday
	b.storage.UpdateUser(usr.Id, usr)

	b.logger.Info(fmt.Sprintf("Дата рождения %s - %s", msg.From.Username, birthday.String()))
}

func (b *Bot) handleLeftMember(ctx context.Context, bot *tgbot.Bot, msg *botmodels.Message) {
	_, _ = bot.SendMessage(
		ctx,
		&tgbot.SendMessageParams{
			Text:   fmt.Sprintf(exitTemplate, msg.LeftChatMember.Username),
			ChatID: msg.Chat.ID,
		},
	)
}

func (b *Bot) handleNewUsers(ctx context.Context, bot *tgbot.Bot, msg *botmodels.Message) {
	for _, user := range msg.NewChatMembers {
		_, _ = bot.SendMessage(
			ctx,
			&tgbot.SendMessageParams{
				Text:   fmt.Sprintf(welcomeTemplate, user.Username),
				ChatID: msg.Chat.ID,
			},
		)
	}
}

func (b *Bot) NotifyAboutUserBirthday(chatId int64, username string) {
	b.logger.Info(fmt.Sprintf("Отправляю поздравление для %s", username))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	_, err := b.Bot.SendMessage(ctx, &tgbot.SendMessageParams{
		Text:   fmt.Sprintf(congratulationTemplate, username),
		ChatID: chatId,
	})

	if err != nil {
		b.logger.Error("Ошибка при отправке поздавления", err)
	}
}
