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

	tgSvc.SetRegexHandler(`^[^/].*`, func(s string) string {
		res, err := gSvc.SendText(ctx, s)
		if err != nil {
			return "GEMINI error"
		}

		return res
	})

	tgSvc.Start()

	time.Sleep(time.Duration(1<<63 - 1))
}
