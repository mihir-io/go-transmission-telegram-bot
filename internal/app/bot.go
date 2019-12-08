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
	AllowedUsers []string
	Port int
	HTTPS bool
}

func StartBot(config *BotConfig, verbose bool) {
	log.Info(fmt.Sprintf("%+v", config))
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Fatal(err)
	}
	log.Info(fmt.Sprintf("Authorized on account %s", bot.Self.UserName))
	if err != nil {
		log.Fatal(err)
	}

	tc, err := rpc.NewTransmissionConnection(config.Hostname, config.Port, config.Username, config.Password, config.HTTPS)
	if err != nil {
		log.Fatal(err)
	}

	ok, serverVersion, serverMinimumVersion, err := tc.IsConnected()
	if err != nil {
		log.Fatal(err)
	}
	if !ok {
		log.Error(fmt.Sprintf("Remote transmission RPC version (v%d) is incompatible with the transmission library (v%d): remote needs at least v%d",
			serverVersion, transmissionrpc.RPCVersion, serverMinimumVersion))
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {

		id := update.Message.Chat.ID
		text := update.Message.Text

		invalidUser := true // Assume the user is not authorized to use this bot instance...
		for _, au := range config.AllowedUsers {
			if au == update.Message.From.UserName {
				invalidUser = false // unless they're in the bot's allowedUsers arg list
			}
		}

		if invalidUser {
			log.Info(fmt.Sprintf("Message from unauthorized user %s was received.", update.Message.From.UserName))
			_, _ = bot.Send(tgbotapi.NewMessage(id, "You are not authorized to use this bot. This incident will be reported."))
			continue
		}

		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		fmt.Println(update.Message.Chat.ID)

		if update.Message.IsCommand() {
			if text == "/start" {
				_, _ = bot.Send(start(id))
			}

			if strings.HasPrefix(text, "/list") {
				_, _ = bot.Send(list(update.Message.Chat.ID, tc))
			}

			if strings.HasPrefix(text, "/play") {
				tokens := strings.Fields(update.Message.Text)
				if len(tokens) <= 1 {
					_, _ = bot.Send(tgbotapi.NewMessage(id, "/play takes 1 arg: Torrent ID."))
					continue
				}

				torrentID, err := strconv.Atoi(tokens[1])
				if err != nil {
					_, _ = bot.Send(tgbotapi.NewMessage(id, "/play takes an integer argument: Torrent ID."))
					continue
				}

				_, _ = bot.Send(play(id, torrentID, tc))
			}

			if strings.HasPrefix(text, "/pause") {
				tokens := strings.Fields(text)
				if len(tokens) <= 1 {
					_, _ = bot.Send(tgbotapi.NewMessage(id, "/pause takes 1 arg: Torrent ID."))
					continue
				}

				torrentID, err := strconv.Atoi(tokens[1])
				if err != nil {
					_, _ = bot.Send(tgbotapi.NewMessage(id, "/pause takes an integer argument: Torrent ID."))
					continue
				}

				_, _ = bot.Send(pause(id, torrentID, tc))
			}

			if strings.HasPrefix(text, "/add") {
				tokens := strings.Fields(text)
				if len(tokens) <= 1 {
					_, _ = bot.Send(tgbotapi.NewMessage(id, "/add takes 1 arg: Torrent file URL."))
					continue
				}

				torrentFileURL := tokens[1]

				_, _ = bot.Send(add(id, torrentFileURL, tc))
			}

			if strings.HasPrefix(text, "/remove") {
				tokens := strings.Fields(text)
				if len(tokens) <= 1 {
					_, _ = bot.Send(tgbotapi.NewMessage(id, "/remove takes 1 arg: Torrent ID."))
					continue
				}

				torrentID, err := strconv.Atoi(tokens[1])
				if err != nil {
					_, _ = bot.Send(tgbotapi.NewMessage(id, "/remove takes an integer argument: Torrent ID."))
					continue
				}

				// Todo: Don't assume data deletion.
				_, _ = bot.Send(remove(id, torrentID, true, tc))

			}
		}
	}
}