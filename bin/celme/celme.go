package main

import (
	"celme/admin"
	"celme/admin/telegram_bot"
	"celme/blank"
	"celme/config"
	"celme/contacts"
	"celme/products"
	"celme/simplePage"
	"log"
	"net/http"
)

func main() {
	var conf config.Config

	conf.Env()
	mux := http.DefaultServeMux
	contacts.Register(conf, mux)
	simplePage.Register(conf, mux)
	admin.Register(conf, mux)
	blank.Register(conf, mux)
	products.Register(conf, mux)
	go telegram_bot.TelegramSender(conf.TelegramToken, conf.MongoURI,
		conf.DBName)
	log.Fatal(http.ListenAndServe(conf.Port, nil))
}
