package bots

import (
	"fmt"
	"github.com/aastashov/house_search_assistant/structs"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var defaultKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🔜", "skip"),
		tgbotapi.NewInlineKeyboardButtonData("❤️", "book"),
		tgbotapi.NewInlineKeyboardButtonData("❌", "dislike"),
	),
)

var likedKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("💔", "unBook"),
		tgbotapi.NewInlineKeyboardButtonData("❌", "dislike"),
	),
)

func (b *Bot) initAnswers() {
	b.answers["skip"] = b.skip
	b.answers["book"] = b.book
	b.answers["dislike"] = b.dislike
}

func (b *Bot) skip(query *tgbotapi.CallbackQuery) {
	fmt.Println("query.Message.MessageID: ", query.Message.MessageID)
	//return "Покажу позже"
}

func (b *Bot) book(query *tgbotapi.CallbackQuery) {
	//return "Добавленно в избранное"
}

func (b *Bot) dislike(query *tgbotapi.CallbackQuery) {
	user := &structs.User{
		Chat:     query.Message.Chat.ID,
		Username: query.Message.Chat.UserName,
	}

	err := b.st.Dislike(query.Message.MessageID, user)
	if err != nil {
		fmt.Println("[dislike.Dislike] error:", err)
		return
	}

	offer, err := b.st.ReadNextOffer(user)
	if err != nil {
		fmt.Println("[dislike.ReadNextOffer] error:", err)
		return
	}

	err = b.SendOffer(offer, user, query, "Больше никогда не покажу")
	if err != nil {
		fmt.Println("[dislike.SendOffer] error:", err)
		return
	}
}
