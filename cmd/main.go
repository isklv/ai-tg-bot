package main

import (
	"context"
	"os"
	"time"

	"github.com/isklv/ai-tg-bot/internal/gemini"
	tg "github.com/isklv/ai-tg-bot/internal/telegram"
)

func main() {
	ctx := context.TODO()
	tgSvc := tg.NewTgBotService(os.Getenv("TELEGRAM_BOT_TOKEN"))
	tgSvc.Init()

	gSvc := gemini.NewGeminiService(os.Getenv("GEMINI_API_KEY"))
	gSvc.Create(ctx, gemini.GEMINI_FLASH)

	tgSvc.SetRegexHandler(`^[^/].*`, func(ha tg.HandlerArgs) string {
		return "Сообщение в обработке"
	}, func(ha tg.HandlerArgs) string {
		res, err := gSvc.SendText(ctx, ha.Text)
		if err != nil {
			return "GEMINI error"
		}

		return res
	})

	tgSvc.SetRegexHandler(``, func(ha tg.HandlerArgs) string {
		return "Сообщение в обработке"
	}, func(ha tg.HandlerArgs) string {
		if ha.Image != nil {
			res, err := gSvc.SendImage(ctx, ha.Text, "jpeg", ha.Image)
			if err != nil {
				return "GEMINI error"
			}

			return res
		}
		return "Прикрепи изображение"
	})

	tgSvc.Start()

	time.Sleep(time.Duration(1<<63 - 1))
}
