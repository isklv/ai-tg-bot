package tg

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"regexp"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type HandlerArgs struct {
	Text  string
	Image *[]byte
}

type TgBotService struct {
	token string
	b     *bot.Bot
}

func NewTgBotService(token string) TgBotService {
	return TgBotService{
		token: token,
	}
}

func (tbs *TgBotService) Init() {

	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
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

func (tbs *TgBotService) SetRegexHandler(pattern string, handlers ...func(HandlerArgs) string) {

	re := regexp.MustCompile(pattern)

	tbs.b.RegisterHandlerRegexp(bot.HandlerTypeMessageText, re, func(ctx context.Context, b *bot.Bot, update *models.Update) {

		ha := HandlerArgs{}

		ha.Text = update.Message.Text
		image := update.Message.Photo

		if len(image) > 0 {
			f, err := tbs.b.GetFile(ctx, &bot.GetFileParams{
				FileID: image[len(image)-1].FileID,
			})

			if err != nil {
				slog.Error("read file error:", err)
			} else {
				slog.Info("Load image", f)
				url := tbs.b.FileDownloadLink(f)

				b := loadImage(url)

				ha.Image = &b
				ha.Text = update.Message.Caption
			}
		}

		for _, h := range handlers {
			resp := h(ha)

			_, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    update.Message.Chat.ID,
				Text:      resp,
				ParseMode: models.ParseModeMarkdownV1,
			})
			if err != nil {
				slog.Error("send message error:", err)
			}
		}

	})
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	slog.Info("defaultHandler", update.Message)
}

func loadImage(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		slog.Error("download image error:", err)
	}
	defer resp.Body.Close()

	buffer := new(bytes.Buffer)
	reader := bufio.NewReader(resp.Body)
	_, err = io.Copy(buffer, reader)
	if err != nil {
		slog.Error("read image error:", err)
	}

	return buffer.Bytes()
}
