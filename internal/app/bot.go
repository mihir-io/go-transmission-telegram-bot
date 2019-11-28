package app

import (
	"fmt"
	"transmission-telegram-bot/internal/pkg/rpc"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/hekmon/transmissionrpc"
	log "github.com/sirupsen/logrus"
)

type BotConfig struct {
	Token string
	Username string
	Password string
	Hostname string
	Port int
	HTTPS bool
}

func StartBot(config *BotConfig, verbose bool){
	log.Info(fmt.Sprintf("%+v", config))
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {log.Fatal(err)}
	log.Info(fmt.Sprintf("Authorized on account %s", bot.Self.UserName))
	if err != nil {log.Fatal(err)}

	tc, err := rpc.NewTransmissionConnection(config.Hostname, config.Port, config.Username, config.Password, config.HTTPS)
	if err != nil {log.Fatal(err)}

	ok, err := tc.is
	if err != nil {log.Fatal(err)}
	if !ok {
		log.Fatal(fmt.Sprintf("Remote transmission RPC version (v%d) is incompatible with the transmission library (v%d): remote needs at least v%d",
			serverVersion, transmissionrpc.RPCVersion, serverMinimumVersion))
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		fmt.Println(update.Message.Chat.ID)

		if update.Message.IsCommand() {
			if update.Message.Text == "/start"{
				_, _ = bot.Send(greeting(update.Message.Chat.ID))
			}
		}

		//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		//
		//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		//msg.ReplyToMessageID = update.Message.MessageID

		//_, _ = bot.Send(msg)
	}
}