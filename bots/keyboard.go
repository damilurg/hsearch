package bots

import (
	"github.com/comov/gilles_search_kg/structs"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var skipButton = tgbotapi.NewInlineKeyboardButtonData("Пропустить", "skip")
var dislikeButton = tgbotapi.NewInlineKeyboardButtonData("Точно нет!", "dislike")
var descriptionButton = tgbotapi.NewInlineKeyboardButtonData("Описание", "description")
var photoButton = tgbotapi.NewInlineKeyboardButtonData("Фото", "photo")

func getKeyboard(offer *structs.Offer) tgbotapi.InlineKeyboardMarkup {
	row1 := tgbotapi.NewInlineKeyboardRow(dislikeButton)
	row2 := tgbotapi.NewInlineKeyboardRow()

	if len(offer.Body) != 0 {
		row2 = append(row2, descriptionButton)
	}

	if offer.Images != 0 {
		row2 = append(row2, photoButton)
	}

	return tgbotapi.NewInlineKeyboardMarkup(row1, row2)
}
