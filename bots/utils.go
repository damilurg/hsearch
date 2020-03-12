package bots

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// SendGroupPhotos - sends a group of photos to a chat room, but unlike Send,
// it returns a list of messages because sending a group of photos is sending
// multiple messages.
func (b *Bot) SendGroupPhotos(config tgbotapi.MediaGroupConfig) ([]tgbotapi.Message, error) {
	params, err := buildParams(config)
	if err != nil {
		return []tgbotapi.Message{}, err
	}

	resp, err := b.bot.MakeRequest("sendMediaGroup", params)
	if err != nil {
		return []tgbotapi.Message{}, err
	}

	var messages []tgbotapi.Message
	err = json.Unmarshal(resp.Result, &messages)
	return messages, err
}

func buildParams(config tgbotapi.MediaGroupConfig) (url.Values, error) {
	chat := config.BaseChat
	v := url.Values{}
	if chat.ChannelUsername != "" {
		v.Add("chat_id", chat.ChannelUsername)
	} else {
		v.Add("chat_id", strconv.FormatInt(chat.ChatID, 10))
	}

	if chat.ReplyToMessageID != 0 {
		v.Add("reply_to_message_id", strconv.Itoa(chat.ReplyToMessageID))
	}

	if chat.ReplyMarkup != nil {
		data, err := json.Marshal(chat.ReplyMarkup)
		if err != nil {
			return v, err
		}

		v.Add("reply_markup", string(data))
	}

	v.Add("disable_notification", strconv.FormatBool(chat.DisableNotification))

	data, err := json.Marshal(config.InputMedia)
	if err != nil {
		return v, err
	}

	v.Add("media", string(data))

	return v, nil
}
