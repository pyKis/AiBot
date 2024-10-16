package team

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)


func ShowTeamMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// Кнопки для действий Team
	buttons := []tgbotapi.KeyboardButton{
		{Text: "Перевести видео"},
		{Text: "Перевести аудио"},
		
	}
	keyboard := tgbotapi.NewReplyKeyboard(buttons)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите действие:")
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func HandleTeamCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

	switch update.Message.Text {
	case "Перевести видео":
		ShowNeiroMenu(bot, update)
	case "Перевести аудио":
		ShowNeiroMenu(bot, update)
	}
}

func ShowNeiroMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update){
	buttons := []tgbotapi.KeyboardButton{
		{Text: "ElevenLab"},
		{Text: "Facebook"},
		{Text: "Назад"},
	}
	keyboard := tgbotapi.NewReplyKeyboard(buttons)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите действие:")
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func HandleNeiro(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	switch update.Message.Text {
	case "ElevenLab":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Тут будет магия")
		bot.Send(msg)
	case "Facebook":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Тут будет магия")
		bot.Send(msg)
	}
}

func HandleBackCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// Проверяем, действительно ли команда "Назад"
	if update.Message.Text == "Назад" {
			// Возвращаем пользователя в главное меню администратора
			ShowTeamMenu(bot, update)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			bot.Send(msg)
	}
}