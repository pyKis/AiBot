package roles

import (
	"fmt"
	"log"
	"main/db"
	"main/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func AddRole(bot *tgbotapi.BotAPI, update tgbotapi.Update, userID int64, role models.Role) {
	err := db.UpdateUserRole(userID, role)
	if err != nil {
		log.Println("Ошибка назначения роли:", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось назначить роль.")
		bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Роль %s успешно назначена пользователю %d", role, userID))
	bot.Send(msg)
}

func ChangeRole(bot *tgbotapi.BotAPI, update tgbotapi.Update, userID int64, newRole models.Role) {
	err := db.UpdateUserRole(userID, newRole)
	if err != nil {
		log.Println("Ошибка изменения роли:", err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось изменить роль.")
		bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Роль пользователя %d изменена на %s", userID, newRole))
	bot.Send(msg)
}

func HandleChangeRole(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// Запрос информации у администратора
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите логин пользователя, которому хотите изменить роль:")
	msg.ReplyToMessageID = update.Message.MessageID
	bot.Send(msg)

	// Ожидаем ответа от администратора
	updateChannel := bot.GetUpdatesChan(tgbotapi.NewUpdate(0))

	// Получаем логин пользователя
	for userUpdate := range updateChannel {
		

		username := userUpdate.Message.Text
	

		// Проверяем, существует ли пользователь с указанным логином
		var userID int64
		err := db.GetUserIDByUsername(username, &userID)
		if err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пользователь с таким логином не найден. Попробуйте снова:")
			bot.Send(msg)
			continue
		}

		// Запрашиваем новую роль
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Введите новую роль (Admin, Client, Team, None):")
		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)

		// Ожидаем новой роли
		roleUpdate := <-updateChannel
		if roleUpdate.Message == nil { // Игнорируем не сообщения
			continue
		}

		role := roleUpdate.Message.Text
		if role != "Admin" && role != "Client" && role != "Team" && role != "None" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Некорректная роль. Попробуйте снова:")
			bot.Send(msg)
			continue
		}

		// Обновляем роль в базе данных
		err = db.UpdateUserRole(userID, models.Role(role))
		if err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка при изменении роли: "+err.Error())
			bot.Send(msg)
			return
		}

		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Роль пользователя успешно изменена.")
		bot.Send(msg)
		return
	}
}


func GetUserRole(UserID int64) (string, error) {
    var role string
    err := db.DB.QueryRow("SELECT role FROM users WHERE user_id = $1", UserID).Scan(&role)
    if err != nil {
        return "", err
    }
    return role, nil
}