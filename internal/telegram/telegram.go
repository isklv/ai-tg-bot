package tg

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"regexp"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type TgBotService struct {
	token    string
	b        *bot.Bot
	handlers *[]bot.HandlerFunc
}

func NewTgBotService(token string) TgBotService {
	return TgBotService{
		token: token,
	}
}

func (tbs *TgBotService) Init() {

	opts := []bot.Option{
		// bot.WithDefaultHandler(defaultHandler),
	}

	b, err := bot.New(tbs.token, opts...)
	if err != nil {
		slog.Error(err.Error())
	}

	tbs.b = b
}

func (tbs *TgBotService) Start() error {

	if tbs.b == nil {
		return fmt.Errorf("Bot instance is empty")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	tbs.b.Start(ctx)

	return nil
}

func (tbs *TgBotService) SetRegexHandler(pattern string, h func(string) string) {

	re := regexp.MustCompile(pattern)

	tbs.b.RegisterHandlerRegexp(bot.HandlerTypeMessageText, re, func(ctx context.Context, b *bot.Bot, update *models.Update) {
		text := update.Message.Text
		resp := h(text)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   resp,
		})
	})
}

// func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
// 	b.SendMessage(ctx, &bot.SendMessageParams{
// 		ChatID: update.Message.Chat.ID,
// 		Text:   update.Message.Text,
// 	})
// }
