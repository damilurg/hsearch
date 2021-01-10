package main

import (
	"context"
	"fmt"
	"github.com/comov/hsearch/bot"
	"github.com/comov/hsearch/storage"
	"github.com/getsentry/sentry-go"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/comov/hsearch/configs"
)

func getSeparatedAlbums(images []string) [][]string {
	maxImages := 10
	albums := make([][]string, 0, (len(images)+maxImages-1)/maxImages)

	for maxImages < len(images) {
		images, albums = images[maxImages:], append(albums, images[0:maxImages:maxImages])
	}
	return append(albums, images)
}

func main() {
	ChatID := int64(-1001171598722)
	cnf, err := configs.GetConf()
	if err != nil {
		log.Fatalln("[main.GetConf] error: ", err)
	}
	fmt.Printf("Release: %s\n", cnf.Release)

	ctx := context.Background()
	db, err := storage.New(ctx, cnf)
	if err != nil {
		log.Fatalln("[main.storage.New] error: ", err)
	}

	defer db.Close()

	telegramBot := bot.NewTelegramBot(cnf, db)

	var images = []string{
		"https://cdn.house.kg/house/images/5/2/3/523ccfef945657ed6c61669cf6fa91fe_1200x900.jpg",
		"https://cdn.house.kg/house/images/2/5/3/253346b47d082bf031d75ccba2a86a9a_1200x900.jpg",
		"https://cdn.house.kg/house/images/e/4/e/e4e74652c0d7194c58bfbf1820fdb048_1200x900.jpg",
		"https://cdn.house.kg/house/images/1/f/a/1fad6dda1e40961f5cf7d15ff2fd3811_1200x900.jpg",
		"https://cdn.house.kg/house/images/0/c/f/0cf822e1ea803e34417ace1428f84da5_1200x900.jpg",
		"https://cdn.house.kg/house/images/7/b/2/7b2af5b7aa6e280027662b5f21da0a73_1200x900.jpg",
		"https://cdn.house.kg/house/images/e/8/0/e804b7e35ae48555e9edb3b7e3ed13ea_1200x900.jpg",
		"https://cdn.house.kg/house/images/2/f/7/2f7aa62c75b1d94ddd12b990de1b3d6f_1200x900.jpg",
		"https://cdn.house.kg/house/images/a/0/7/a071f14adaaf38ca9b6ad864ea20ce2d_1200x900.jpg",
		"https://cdn.house.kg/house/images/0/4/3/0437f77e365a1e01920387624a31e4d0_1200x900.jpg",
		"https://cdn.house.kg/house/images/3/a/e/3aea5ee0a712ba902ed81e1cd6bc2162_1200x900.jpg",
		"https://cdn.house.kg/house/images/1/9/4/19429310b4cea7b88adbeedab0691ea4_1200x900.jpg",
		"https://cdn.house.kg/house/images/d/4/1/d417e5b0b766246bfd2eb2772a8297d5_1200x900.jpg",
		"https://cdn.house.kg/house/images/1/2/4/12440e2e8f7b8f454dc68630a371d4ee_1200x900.jpg",
		"https://cdn.house.kg/house/images/2/7/e/27e44e59576f57b23cf7038259283f23_1200x900.jpg",
		"https://cdn.house.kg/house/images/f/f/4/ff4e9fa3c4b31f953e660ad2508d1f85_1200x900.jpg",
		"https://cdn.house.kg/house/images/7/5/d/75dd87c0d270c49c9a7f2251d3bd092b_1200x900.jpg",
		"https://cdn.house.kg/house/images/1/5/8/158784a493653c9cb094197894c06565_1200x900.jpg",
		"https://cdn.house.kg/house/images/4/c/e/4cea01ca16070adab53281c07ecc0f11_1200x900.jpg",
		"https://cdn.house.kg/house/images/b/7/2/b7210120f4a7801066709673e8c66092_1200x900.jpg",
	}

	for _, album := range getSeparatedAlbums(images) {
		medias := make([]interface{}, 0)
		for _, img := range album {
			medias = append(medias, tgbotapi.NewInputMediaPhoto(img))
		}

		message := tgbotapi.NewMediaGroup(ChatID, medias)
		message.ReplyToMessageID = 52

		messages, err := telegramBot.SendGroupPhotos(message)
		if err != nil {
			sentry.CaptureException(err)
			log.Println("[photo.Send] sending album error:", err)
		}

		_ = messages

		//for _, msg := range messages {
		//	err = db.SaveMessage(
		//		ctx,
		//		msg.MessageID,
		//		1,
		//		ChatID,
		//		structs.KindPhoto,
		//	)
		//	if err != nil {
		//		sentry.CaptureException(err)
		//		log.Println("[photo.SaveMessage] error:", err)
		//	}
		//}
	}
}
