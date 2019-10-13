package main

import (
	"log"

	_ "github.com/joho/godotenv/autoload"
)

const (
	BaseURL = "http://diesel.elcat.kg/index.php?showforum=305?page=%d"
)

func main() {
	cnf, err := GetConf()
	if err != nil {
		log.Fatal(err)
	}

	bot := NewBot(cnf)
	bot.Start()
}
