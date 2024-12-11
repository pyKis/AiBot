// admin/admin.go
package admin

import (
	"fmt"
	"log"
	"main/db"
	_"main/models"
	"main/roles"
	"main/roles/team"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ShowAdminMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update, adminUserID int64) {
	if update.Message.From.ID != adminUserID {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "У вас нет прав администратора.")
		bot.Send(msg)
		return
	}

	// Кнопки для действий администратора
	buttons := []tgbotapi.KeyboardButton{
		{Text: "Статистика"},
		{Text: "Изменить роль"},
		{Text: "Team"},
		
	}
	keyboard := tgbotapi.NewReplyKeyboard(buttons)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите действие:")
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func HandleAdminCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update, adminUserID int64) {
	
	switch update.Message.Text {
	case "Статистика":
		ShowStatisticsMenu(bot, update)
	case "Изменить роль":
		roles.HandleChangeRole(bot, update)
	case "Team":
		team.ShowTeamMenu(bot, update)
	}
}

func ShowStatisticsMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update){
	buttons := []tgbotapi.KeyboardButton{
		{Text: "Пользователи"},
		{Text: "Нейронные сети"},
		{Text: "Назад"},
	}
	keyboard := tgbotapi.NewReplyKeyboard(buttons)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите действие:")
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func HandleStatistics(bot *tgbotapi.BotAPI, update tgbotapi.Update, adminUserID int64) {
	switch update.Message.Text {
	case "Пользователи":
		handleUsers(bot, update)
	case "Нейронные сети":
		handleNeuralNetworks(bot, update)
	
	
	}
}

func handleUsers(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	  // SQL-запросы для получения статистики
	  totalUsersQuery := `SELECT COUNT(*) FROM users;`
	  invitedUsersQuery := `SELECT COUNT(*) FROM referrals WHERE referral_code IS NOT NULL;`
	  adminUsersQuery := `SELECT COUNT(*) FROM users WHERE role = 'Admin';`
	  clientUsersQuery := `SELECT COUNT(*) FROM users WHERE role = 'Client';`
	  teamUsersQuery := `SELECT COUNT(*) FROM users WHERE role = 'Team';`
	  noneUsersQuery := `SELECT COUNT(*) FROM users WHERE role = 'None';`
  
	  var totalUsers, invitedUsers, adminUsers, clientUsers, teamUsers, noneUsers int
  
	  // Выполнение SQL-запросов и получение данных
	  err := db.DB.QueryRow(totalUsersQuery).Scan(&totalUsers)
	  if err != nil {
		  log.Println("Ошибка при получении общего количества пользователей:", err)
		  return
	  }
  
	  err = db.DB.QueryRow(invitedUsersQuery).Scan(&invitedUsers)
	  if err != nil {
		  log.Println("Ошибка при получении количества приглашенных пользователей:", err)
		  return
	  }
  
	  err = db.DB.QueryRow(adminUsersQuery).Scan(&adminUsers)
	  if err != nil {
		  log.Println("Ошибка при получении количества администраторов:", err)
		  return
	  }
  
	  err = db.DB.QueryRow(clientUsersQuery).Scan(&clientUsers)
	  if err != nil {
		  log.Println("Ошибка при получении количества клиентов:", err)
		  return
	  }
  
	  err = db.DB.QueryRow(teamUsersQuery).Scan(&teamUsers)
	  if err != nil {
		  log.Println("Ошибка при получении количества членов команды:", err)
		  return
	  }
  
	  err = db.DB.QueryRow(noneUsersQuery).Scan(&noneUsers)
	  if err != nil {
		  log.Println("Ошибка при получении количества пользователей без роли:", err)
		  return
	  }
  
	  // Формирование сообщения для отправки
	  statsMessage := fmt.Sprintf(
		  "Статистика пользователей:\n\n"+
			  "Общее количество: %d\n"+
			  "Приглашенные по реферальной ссылке: %d\n"+
			  "Администраторы: %d\n"+
			  "Клиенты: %d\n"+
			  "Команда: %d\n"+
			  "Без роли: %d\n",
		  totalUsers, invitedUsers, adminUsers, clientUsers, teamUsers, noneUsers,
	  )
  
	  // Отправка сообщения пользователю
	  msg := tgbotapi.NewMessage(update.Message.Chat.ID, statsMessage)
	  bot.Send(msg)
}

func handleNeuralNetworks(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// Кнопки для выбора между "Сеть 1" и "Сеть 2"
	buttons := []tgbotapi.KeyboardButton{
		{Text: "ElevenLab"},
		{Text: "Facebook"},
		{Text: "Назад"},
	}
	keyboard := tgbotapi.NewReplyKeyboard(buttons)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите сеть:")
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func HandleNeuralNetworksSubCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update, adminUserID int64) {
	switch update.Message.Text {
	case "ElevenLab":
		handleNetwork1(bot, update)
	case "Facebook":
		handleNetwork1(bot, update)
	case "Назад":
		HandleBackCommand(bot, update, adminUserID)
	}
}

func handleNetwork1(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// Кнопки для выбора "Деньги" и "Токены"
	buttons := []tgbotapi.KeyboardButton{
		{Text: "Деньги"},
		{Text: "Токены"},
		{Text: "Назад"},
	}
	keyboard := tgbotapi.NewReplyKeyboard(buttons)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите категорию:")
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func HandleNetwork1SubCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	switch update.Message.Text {
	case "Деньги":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Информация о деньгах...")
		bot.Send(msg)
	case "Токены":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Информация о токенах...")
		bot.Send(msg)
	}
}

func HandleBackCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update, adminUserID int64) {
	// Проверяем, действительно ли команда "Назад"
	if update.Message.Text == "Назад" {
			// Возвращаем пользователя в главное меню администратора
			ShowAdminMenu(bot, update, adminUserID)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			bot.Send(msg)
	}
}


