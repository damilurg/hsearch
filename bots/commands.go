package bots

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (b *Bot) start(_ *tgbotapi.Message) string {
	return startMessage
}

func (b *Bot) stop(message *tgbotapi.Message) string {
	err := b.st.StopSearch(message.Chat.UserName)
	if err != nil {
		log.Println("[stop.StartSearchForUser] error:", err)
		return "Прости, говнокод сломался"
	}

	return "Я больше не буду отправлять тебе квартиры. Пока :)"
}

func (b *Bot) search(message *tgbotapi.Message) string {
	err := b.st.StartSearch(message.Chat.ID, message.Chat.UserName)
	if err != nil {
		log.Println("[search.StartSearchForUser] error:", err)
		return "Прости, говнокод сломался"
	}

	return "Теперь я буду искать для тебя квартиры"
}

func (b *Bot) bookmarks(message *tgbotapi.Message) string {
	offers, chat, err := b.st.Bookmarks(message.Chat.UserName)
	if err != nil {
		log.Println("[bookmarks.StartSearchForUser] error:", err)
		return "Прости, говнокод сломался"
	}

	if len(offers) <= 0 {
		return "Хм... У тебя пока нет квартир в закладках"
	}

	b.bookmarksMessages(offers, chat)

	return fmt.Sprintf("Список отмеченных квартир %d", len(offers))
}

func (b *Bot) feedback(message *tgbotapi.Message) string {
	if message.CommandArguments() == "" {
		return "Нужно оставить комментарий в виде:\n /feedback Я тебя найду...."
	}

	err := b.st.Feedback(message.Chat.ID, message.Chat.UserName, message.CommandArguments())
	if err != nil {
		log.Println("[feedback.StartSearchForUser] error:", err)
		return "Прости, даже фидбек может быть сломан"
	}
	return "Понял, предам!"
}
