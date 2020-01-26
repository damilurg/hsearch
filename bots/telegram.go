package bots

import (
	"log"

	"github.com/aastashov/house_search_assistant/configs"
	"github.com/aastashov/house_search_assistant/structs"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type (
	Storage interface {
		StartSearch(chat int64, username string) error
		StopSearch(username string) error
		Bookmarks(username string) ([]*structs.Offer, int64, error)
		SaveMessage(msgId int, offerId uint64, chat int64) error
		Dislike(msgId int, user *structs.User) error
		ReadNextOffer(user *structs.User) (*structs.Offer, error)
	}

	answer map[string]func(query *tgbotapi.CallbackQuery)

	Bot struct {
		bot     *tgbotapi.BotAPI
		st      Storage
		answers answer
	}
)

func NewTelegramBot(cnf *configs.Config, st Storage) *Bot {
	bot, err := tgbotapi.NewBotAPI(cnf.TelegramToken)
	if err != nil {
		log.Fatalln("[NewBot.NewBotAPI] error: ", err)
		return nil
	}

	bb := &Bot{
		bot:     bot,
		st:      st,
		answers: make(answer, 0),
	}

	bb.initAnswers()
	return bb
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
		if update.CallbackQuery != nil {
			b.answers[update.CallbackQuery.Data](update.CallbackQuery)
		}

		if update.Message != nil {
			msg := "Нет среди доступных команд :("
			if update.Message.IsCommand() {
				switch update.Message.Command() {
				case "start":
					msg = b.start(update.Message)
				case "help":
					msg = b.help(update.Message)
				case "stop":
					msg = b.stop(update.Message)
				case "search":
					msg = b.search(update.Message)
				case "bookmarks":
					msg = b.bookmarks(update.Message)
				}
			}

			message := tgbotapi.NewMessage(update.Message.Chat.ID, msg)
			_, err := b.bot.Send(message)
			if err != nil {
				log.Println("[Start.Message.Send] error: ", err)
			}
		}
	}
}

// SendOffer - отправляет offer пользователю, так же региструрует в бд под
// какие номером сообщения было отправленно сообщение и меняет клавиатуру в
// зависимости от offer
func (b *Bot) SendOffer(offer *structs.Offer, user *structs.User, query *tgbotapi.CallbackQuery, answer string) error {
	if query != nil {
		/* TODO: Найти offer по ID и сделать соответствующее действие */

		_, err := b.bot.AnswerCallbackQuery(tgbotapi.NewCallback(query.ID, answer))
		if err != nil {
			log.Println("[SendOffer.AnswerCallbackQuery] error: ", err)
		}
	}

	message := tgbotapi.NewMessage(user.Chat, DefaultMessage(offer))
	message.ReplyMarkup = defaultKeyboard

	send, err := b.bot.Send(message)
	if err != nil {
		return err
	}

	err = b.st.SaveMessage(send.MessageID, offer.Id, user.Chat)
	if err != nil {
		log.Println("[SendOffer.SaveMessage] error:", err)
	}

	return nil
}

func (b *Bot) SendPreviewMessage(offer *structs.Offer, user *structs.User) error {
	return b.SendOffer(offer, user, nil, "")
}
