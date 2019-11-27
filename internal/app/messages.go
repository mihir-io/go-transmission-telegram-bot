package app

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func greeting(id int64) tgbotapi.MessageConfig {
	s := `
		Hello! 👋
		Welcome to the Transmission Telegram Bot!
		For a list of commands you can use, respond
		with /command
`
	return tgbotapi.NewMessage(id, s)

}
