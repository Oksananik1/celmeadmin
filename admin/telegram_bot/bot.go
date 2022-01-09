package telegram_bot

import (
	"github.com/Syfaro/telegram-bot-api"
	"reflect"
)

func TelegramBot(token, mongoURI, dbName string) {

	//Создаем бота
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}

	//Устанавливаем время обновления
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	//Получаем обновления от бота
	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		user, errUser := findUserUser(update.Message.From.ID, mongoURI, dbName)

		//Проверяем что от пользователья пришло именно текстовое сообщение
		if update.Message.Contact != nil && errUser != nil {
			if update.Message.Contact.UserID == update.Message.From.ID {
				user := User{UserId: update.Message.Contact.UserID,
					ChatID:      update.Message.Chat.ID,
					FirstName:   update.Message.From.FirstName,
					LastName:    update.Message.From.LastName,
					PhoneNumber: update.Message.Contact.PhoneNumber}
				collectUser(user, mongoURI, dbName)
				user, errUser = findUserUser(update.Message.From.ID, mongoURI, dbName)
				hideButtons(bot, update.Message.Chat.ID)
				if user.IsValid {
					if user.ChatID != update.Message.Chat.ID {
						user.ChatID = update.Message.Chat.ID
						collectUser(user, mongoURI, dbName)
					}
					msg := tgbotapi.NewMessage(update.Message.Chat.ID,
						"Привет, "+
							""+user.FirstName+"! Пока нет новых сообщений. "+
							"Но как только появятся, я сразу оповещу тебя.")
					bot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID,
						"Привет, "+
							""+user.FirstName+"! К сожалению я тебя не знаю. "+
							"И дальше общаться больше не хочу. Пока!")
					bot.Send(msg)
				}
			}
		}

		if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {

			switch update.Message.Text {
			case "/start":

				//Отправлем сообщение
				if errUser != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID,
						"Привет! Я Бот, "+
							"который будет оповещать тебя о новых сообщениях оставленых на сайте")
					bot.Send(msg)
					replySendPhoneAndGeo(bot, update.Message.Chat.ID)
				} else {
					if user.IsValid {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID,
							"Привет, "+user.FirstName+"! Пока нет новых сообщений.")
						bot.Send(msg)
					} else {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID,
							"Привет, "+
								""+user.FirstName+"! К сожалению я тебя не знаю. "+
								"И дальше общаться больше не хочу. Пока!")
						bot.Send(msg)
					}
				}

			case "/number_of_users":

				//Отправлем сообщение
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Database not connected, so i can't say you how many peoples used me.")
				bot.Send(msg)

			default:

				msg := tgbotapi.NewMessage(update.Message.Chat.ID,
					"Кто это сказал? Ты??? И это ты сказал?")
				bot.Send(msg)
			}
		}

		//Проходим через срез и отправляем каждый элемент пользователю

	}
}

func replySendPhoneAndGeo(bot *tgbotapi.BotAPI, chat_id int64) {
	msg := tgbotapi.NewMessage(chat_id, "Давай познакомися!")
	var keyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButtonContact(
				"\xF0\x9F\x93\x9E Отправить мой номер"),
		),
	)
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func hideButtons(bot *tgbotapi.BotAPI, chat_id int64) {
	msg := tgbotapi.NewMessage(chat_id, "Давай познакомися!")
	var keyboard = tgbotapi.NewRemoveKeyboard(true)
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}
