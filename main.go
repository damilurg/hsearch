package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/aastashov/house_search_assistant/background"
	"github.com/aastashov/house_search_assistant/bots"
	"github.com/aastashov/house_search_assistant/configs"
	"github.com/aastashov/house_search_assistant/storage"

	_ "github.com/joho/godotenv/autoload"
)

const (
	BaseURL      = "http://diesel.elcat.kg/index.php?showforum=305?page=%d"
	helpCommands = "[migrate|bgm]"
)

/* TODO: проблемы
    1. протестировать ReadUsersForOrder
    2. бот не понимает какие квартиры отправлял, а какие нет
    3. менеджер так же спамит все
    4. не реализовано skip
    5. не реализовано images
    6. не реализовано description
    7. не реализовано like
    8. не реализовано skip
    9. описать helpCommands
    10. пройтись по всем туду
*/

func main() {
	cnf, err := configs.GetConf()
	if err != nil {
		log.Fatalln("[main.GetConf] error: ", err)
	}

	db, err := storage.New()
	if err != nil {
		log.Fatalln("[main.storage.New] error: ", err)
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "migrate":
			dir, _ := os.Getwd()
			err = db.Migrate(path.Join(dir, "migrations"))
			if err != nil {
				log.Fatalln("[main.storage.Migrate] error: ", err)
			}
		default:
			fmt.Println("Я тебя не понимаю. Вот что могу\n", helpCommands)
		}
		return
	}

	telegramBot := bots.NewTelegramBot(cnf, db)
	go telegramBot.Start()

	background.StartOfferManager(BaseURL, cnf, db, telegramBot)
}
