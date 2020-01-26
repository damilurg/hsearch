package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/aastashov/house_search_assistant/bot"
	"github.com/aastashov/house_search_assistant/configs"
	"github.com/aastashov/house_search_assistant/storage"

	_ "github.com/joho/godotenv/autoload"
)

const (
	BaseURL = "http://diesel.elcat.kg/index.php?showforum=305?page=%d"
	// TODO: описать хелпер
	helpCommands = "[migrate|bgm]"
)

func main() {
	// для отладки
	// parser.TestSelector("http://diesel.elcat.kg/index.php?s=6f9a6608d3dc486d57391576b85ea17d&showtopic=292671275")
	// return

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

	// todo: NewBot должен принимать db для использования storage
	telegramBot := bot.NewTelegramBot(cnf)

	// todo: start должен запускаться в go routine
	telegramBot.Start()

	// todo: после запуска бота, запускаем менеджера, который будет искать
	//  квартиры и рассылать через бота
}
