package bots

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (b *Bot) start(_ *tgbotapi.Message) string {
	return startMessage
}

func (b *Bot) help(_ *tgbotapi.Message) string {
	return helpMessage
}

func (b *Bot) stop(message *tgbotapi.Message) string {
	err := b.st.StopSearch(message.Chat.UserName)
	if err != nil {
		log.Println("[start.StartSearchForUser] error:", err)
		trace := ""
		if message.Chat.UserName == "aastashov" {
			trace = "\n\nTrace for developer:\n" + err.Error()
		}
		return "Прости, говнокод сломался" + trace
	}

	return "Я больше не буду отправлять тебе квартиры. Пока :)"
}

func (b *Bot) search(message *tgbotapi.Message) string {
	err := b.st.StartSearch(message.Chat.ID, message.Chat.UserName)
	if err != nil {
		log.Println("[start.StartSearchForUser] error:", err)
		trace := ""
		if message.Chat.UserName == "aastashov" {
			trace = "\n\nTrace for developer:\n" + err.Error()
		}
		return "Прости, говнокод сломался" + trace
	}

	return "Теперь я буду искать для тебя квартиры"
}

func (b *Bot) bookmarks(message *tgbotapi.Message) string {
	offers, chat, err := b.st.Bookmarks(message.Chat.UserName)
	if err != nil {
		log.Println("[start.StartSearchForUser] error:", err)
		trace := ""
		if message.Chat.UserName == "aastashov" {
			trace = "\n\nTrace for developer:\n" + err.Error()
		}
		return "Прости, говнокод сломался" + trace
	}

	if len(offers) <= 0 {
		return "Хм... У тебя пока нет квартир в закладках"
	}

	b.bookmarksMessages(offers, chat)

	return fmt.Sprintf("Список отмеченных квартир %d", len(offers))
}
