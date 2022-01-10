package telegram_bot

import (
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"time"
)

func TelegramSender(token, mongoURI, dbName string) {

	//Создаем бота
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}

	for {
		time.Sleep(time.Second * 10)

		messages, err := findMessagesSubscribe(mongoURI, dbName)
		if err == nil && len(messages) > 0 {
			users, _ := findUsersForSubscribe(mongoURI, dbName)
			for _, user := range users {
				msg := tgbotapi.NewMessage(user.ChatID,
					"<b>Новые сообщения!</b>")
				msg.ParseMode = "HTML"
				bot.Send(msg)
				for _, m := range messages {
					msg := tgbotapi.NewMessage(user.ChatID,
						fmt.Sprintf(
							"От: <u>%s</u>\nНомер: <u>%s</u>\nEmail: <u>%s</u>\n"+
								"Сообщение:\n%s", m.Name, m.Phone, m.Email,
							m.Message))
					msg.ParseMode = "HTML"
					bot.Send(msg)
				}
			}
		}
	}
}
