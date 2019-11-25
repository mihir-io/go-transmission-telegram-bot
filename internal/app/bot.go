package app

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

type BotConfig struct {
	Token string
	Username string
	Password string
	Hostname string
	Port int
}

func StartBot(config *BotConfig, verbose bool){
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {log.Fatal(err)}

	log.Info("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		_, _ = bot.Send(msg)
	}
}