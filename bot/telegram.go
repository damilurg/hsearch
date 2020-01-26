package bot

import (
	"log"

	"house_search_assistant/configs"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type (
	Bot struct {
		bot *tgbotapi.BotAPI
	}
)

func NewTelegramBot(cnf *configs.Config) *Bot {
	bot, err := tgbotapi.NewBotAPI(cnf.TelegramToken)
	if err != nil {
		log.Fatalln("[NewBot.NewBotAPI] error: ", err)
		return nil
	}

	return &Bot{bot: bot}
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatalln("[Start.GetUpdatesChan] error: ", err)
		return
	}

	log.Println("Start listen Telegram chanel")
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		//msg.ReplyToMessageID = update.Message.MessageID
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		s := "Hello"
		m := tgbotapi.NewMessage(update.Message.Chat.ID, s)
		_, err := b.bot.Send(m)
		if err != nil {
			log.Println("[Start.Send] error: ", err)
		}
	}
}
