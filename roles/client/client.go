package client

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)





func ShowClientMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// Кнопки для действий Team
	buttons := []tgbotapi.KeyboardButton{
		{Text: "Поделиться контактом"},
		{Text: "Сгенерировать ссылку"},
		
	}
	keyboard := tgbotapi.NewReplyKeyboard(buttons)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите действие:")
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func HandleClientCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

	switch update.Message.Text {
	case "Поделиться контактомо":
		
	case "Сгенерировать ссылку":
		
	}
}




