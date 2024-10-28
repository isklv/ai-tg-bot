package main

import (
    "fmt"
    "log"

    "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
    bot, err := telegram.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
    if err != nil {
        log.Panicln(err)
    }

    bot.OnText(&bot.DefaultFilter, func(msg *tgbotapi.Message) {
        fmt.Printf("Received message: %s\n", msg.Text)
        bot.SendMessage(msg.Chat.ID, "You sent: " + msg.Text)
    })

    err = bot.StartPolling()
    if err != nil {
        log.Panicln(err)
    }
}
