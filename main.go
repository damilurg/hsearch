package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/comov/gilles_search_kg/background"
	"github.com/comov/gilles_search_kg/bots"
	"github.com/comov/gilles_search_kg/configs"
	"github.com/comov/gilles_search_kg/storage"

	_ "github.com/joho/godotenv/autoload"
)

const (
	BaseURL      = "http://diesel.elcat.kg/index.php?showforum=305&page=%d"
	helpCommands = "gilles_search_kg: '%s' is not a command.\n" +
		"usage: go run main.go [migrate]\n\n" +
		"By default gilles_search_kg run offer manager and telegram" +
		" bot.\nFor example: go run main.go\n\n" +
		"Commands:\n" +
		"\tmigrate - the command for run migration and create DB if not exist\n"
)

func main() {
	cnf, err := configs.GetConf()
	if err != nil {
		log.Fatalln("[main.GetConf] error: ", err)
	}

	db, err := storage.New(cnf)
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
			fmt.Printf(helpCommands, os.Args[1])
		}
		return
	}

	// Telegram bot и Offer manager в дальнейшем нужно запускать как отдельные
	// сервисы, а главный поток оставить следить за ними. Таким образом можно
	// сделать graceful shutdown, reload config да и просто по приколу

	telegramBot := bots.NewTelegramBot(cnf, db)
	go telegramBot.Start()

	omr := background.StartOfferManager(BaseURL, cnf, db, telegramBot)
	omr.Start()
}
