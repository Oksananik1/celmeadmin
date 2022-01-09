package main

import (
	"celme/admin/telegram_bot"
	"celme/config"
	"log"
	"net/http"
)

func main() {
	var conf config.Config

	conf.Env()
	telegram_bot.TelegramBot(conf.TelegramToken, conf.MongoURI, conf.DBName)
	log.Fatal(http.ListenAndServe(conf.Port, nil))
}
