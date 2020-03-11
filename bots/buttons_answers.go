package bots

import (
	"log"

	"github.com/comov/hsearch/structs"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// initAnswers - содержит список всех зарегистированных кнопок
func (b *Bot) initAnswers() {
	b.answers["skip"] = authHandler(b.skip)
	b.answers["dislike"] = authHandler(b.dislike)
	b.answers["description"] = authHandler(b.description)
	b.answers["photo"] = authHandler(b.photo)
}

// skip - обраьатывает нажатие на кнопку "Пропустить"
func (b *Bot) skip(user *structs.User, query *tgbotapi.CallbackQuery) {
	err := b.st.Skip(query.Message.MessageID, user)
	if err != nil {
		log.Println("[skip.Skip] error:", err)
		return
	}

	offer, err := b.st.ReadNextOffer(user)
	if err != nil {
		log.Println("[skip.ReadNextOffer] error:", err)
		return
	}

	err = b.SendOffer(offer, user, query, "Покажу позже")
	if err != nil {
		log.Println("[skip.SendOffer] error:", err)
		return
	}
}

// dislike - this button delete order from chat and no more show to user that
// order.
func (b *Bot) dislike(user *structs.User, query *tgbotapi.CallbackQuery) {
	messagesIds, err := b.st.Dislike(query.Message.MessageID, user)
	if err != nil {
		log.Println("[dislike.Dislike] error:", err)
		return
	}

	for _, id := range messagesIds {
		_, _ = b.bot.DeleteMessage(
			tgbotapi.NewDeleteMessage(query.Message.Chat.ID, id),
		)
	}

	err = b.SendOffer(nil, user, query, "Больше никогда не покажу")
	if err != nil {
		log.Println("[dislike.SendOffer] error:", err)
		return
	}
}

// description - return full description about order
func (b *Bot) description(user *structs.User, query *tgbotapi.CallbackQuery) {
	offerId, body, err := b.st.ReadOfferDescription(query.Message.MessageID, user)
	if err != nil {
		log.Println("[description.ReadOfferDescription] error:", err)
		return
	}

	message := tgbotapi.NewMessage(user.Chat, body)
	message.ReplyToMessageID = query.Message.MessageID

	send, err := b.bot.Send(message)
	if err != nil {
		log.Println("[description.Send] error:", err)
	}

	err = b.st.SaveMessage(
		send.MessageID,
		offerId,
		user.Chat,
		structs.KindDescription,
	)
	if err != nil {
		log.Println("[photo.SaveMessage] error:", err)
	}
}

// photo - this button return all orders photos from site
func (b *Bot) photo(user *structs.User, query *tgbotapi.CallbackQuery) {
	offerId, images, err := b.st.ReadOfferImages(query.Message.MessageID, user)
	if err != nil {
		log.Println("[photo.ReadOfferDescription] error:", err)
		return
	}

	waitMessage := tgbotapi.Message{}
	if len(images) != 0 {
		waitMessage, err = b.bot.Send(tgbotapi.NewMessage(
			user.Chat,
			WaitPhotoMessage(len(images)),
		))
		if err != nil {
			log.Println("[photo.Send] error:", err)
		}
	}

	for _, album := range getSeparatedAlbums(images) {
		imgs := make([]interface{}, 0)
		for _, img := range album {
			imgs = append(imgs, tgbotapi.NewInputMediaPhoto(img))
		}

		message := tgbotapi.NewMediaGroup(user.Chat, imgs)
		message.ReplyToMessageID = query.Message.MessageID

		messages, err := b.SendGroupPhotos(message)
		if err != nil {
			log.Println("[photo.Send] sending album error:", err)
		}

		for _, msg := range messages {
			err = b.st.SaveMessage(
				msg.MessageID,
				offerId,
				user.Chat,
				structs.KindPhoto,
			)
			if err != nil {
				log.Println("[photo.SaveMessage] error:", err)
			}
		}
	}

	if len(images) != 0 {
		_, _ = b.bot.DeleteMessage(tgbotapi.NewDeleteMessage(
			query.Message.Chat.ID,
			waitMessage.MessageID,
		))
	}
}

// getSeparatedAlbums - separate images array to 10-items albums. Telegram API
// has limit: `max images in images album is 10`
func getSeparatedAlbums(images []string) [][]string {
	maxImages := 10
	albums := make([][]string, 0, (len(images)+maxImages-1)/maxImages)

	for maxImages < len(images) {
		images, albums = images[maxImages:], append(albums, images[0:maxImages:maxImages])
	}
	return append(albums, images)
}

// authHandler - handler that get the user from query and add to handle func
func authHandler(f func(u *structs.User, q *tgbotapi.CallbackQuery)) func(*tgbotapi.CallbackQuery) {
	return func(query *tgbotapi.CallbackQuery) {
		user := &structs.User{
			Chat:     query.Message.Chat.ID,
			Username: query.Message.Chat.UserName,
		}
		f(user, query)
	}
}
