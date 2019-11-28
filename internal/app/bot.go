package app

import (
	"fmt"
	"strconv"
	"strings"
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

	ok, serverVersion, serverMinimumVersion, err := tc.IsConnected()
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
				_, _ = bot.Send(start(update.Message.Chat.ID))
			}

			if update.Message.Text == "/list" {
				_, _ = bot.Send(list(update.Message.Chat.ID, tc))
			}

			if strings.HasPrefix(update.Message.Text, "/play") {
				tokens := strings.Fields(update.Message.Text)
				if len(tokens) <= 1 {
					_, _ = bot.Send(tgbotapi.MessageConfig{Text:"/play takes 1 arg: Torrent ID."})
				}

				torrentID, err := strconv.Atoi(tokens[1])
				if err != nil {
					_, _ = bot.Send(tgbotapi.MessageConfig{Text:"/play takes an integer argument: Torrent ID."})
				}

					_, _ = bot.Send(play(update.Message.Chat.ID, torrentID,tc))
			}

			if strings.HasPrefix(update.Message.Text, "/pause") {
				tokens := strings.Fields(update.Message.Text)
				if len(tokens) <= 1 {
					_, _ = bot.Send(tgbotapi.MessageConfig{Text:"/pause takes 1 arg: Torrent ID."})
				}

				torrentID, err := strconv.Atoi(tokens[1])
				if err != nil {
					_, _ = bot.Send(tgbotapi.MessageConfig{Text:"/pause takes an integer argument: Torrent ID."})
				}

				_, _ = bot.Send(pause(update.Message.Chat.ID, torrentID,tc))
			}
		}
	}
}