package bot

import (
	"fmt"
	"log"
	"main/db"
	"main/roles/team"
	"main/roles/client"
	"main/roles/admin"
	"main/roles"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)




func HandleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	user := update.Message.From
	db.SaveUser(user)

    // Получаем UserID из обновления
    var UserID int64
    if update.Message != nil {
        UserID = update.Message.From.ID
    }

    // Получаем роль пользователя
    role, err := roles.GetUserRole(UserID)
    if err != nil {
        log.Println("Error getting user role:", err)
        return
    }

	

    // Проверяем, является ли сообщение командой
    if update.Message.IsCommand() {
        switch update.Message.Command() {
		case "start":
			if role == "None" {
				// Если роль None, показываем сообщение о том, что доступ ограничен
				adminContact := fmt.Sprintf("Доступ ограничен.")
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, adminContact)
				bot.Send(msg)
				return
			}
            handleStartCommand(bot, update)
			if role == "Team"{
				team.ShowTeamMenu(bot, update)
			}
			if role == "Admin"{
				admin.ShowAdminMenu(bot, update, UserID)
			}
			if role == "Client"{
				client.ShowClientMenu(bot,update)
			}
        case "admin":
            if role == "Admin" {
                admin.ShowAdminMenu(bot, update, UserID)
            } else {
                bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "У вас нет прав администратора"))
            }
        }
    } else {
        switch update.Message.Text {
        case "Статистика", "Изменить роль":
            if role == "Admin" {
                admin.HandleAdminCommand(bot, update, UserID)
            } else {
                bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "У вас нет прав для этой команды"))
            }
        case "Пользователи", "Нейронные сети":
            admin.HandleStatistics(bot, update, UserID)
		case "Team":
            team.ShowTeamMenu(bot, update)
        case "Сеть 1", "Сеть 2":
            admin.HandleNeuralNetworksSubCommand(bot, update, UserID)
        case "Деньги", "Токены":
            admin.HandleNetwork1SubCommand(bot, update)
		case "Перевести видео",	"Перевести аудио":
			if role == "Team" {
			team.HandleTeamCommand(bot,update)
			}
			if role == "Admin" {
				team.HandleTeamCommand(bot,update)
				}else {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "У вас нет прав для этой команды"))
				}
		case "ElevenLab","Facebook":
			if role == "Team" {
			team.HandleNeiro(bot,update)
			}
			if role == "Admin" {
				team.HandleNeiro(bot,update)
				}else {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "У вас нет прав для этой команды"))
				}
		case "Поделиться контактом":
			if role == "Client" {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "В разработке"))
				}else {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "У вас нет прав для этой команды"))
				}
			case "Сгенерировать ссылку":
				if role == "Client" {
				db.GenerateReferralLink(bot,update)
					}else {
						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "У вас нет прав для этой команды"))
					}
						
        case "Назад":
			switch role {
			case "Admin":
				admin.HandleBackCommand(bot, update, UserID)
			case "Team":
				team.HandleBackCommand(bot, update)
			}
        default:
            bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда"))
        }
    }
	switch update.Message.Text {
	case "Статистика", "Изменить роль":
		if role == "Team" {
			team.HandleTeamCommand(bot, update)
		} else {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "У вас нет прав для этой команды"))
		}
	}
}



func handleStartCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	

	// Проверка на наличие реферального кода
	referralCode := update.Message.CommandArguments()
	if referralCode != "" {
		inviterID, inviterUsername, err := db.GetInviterInfo(referralCode)
		if err == nil && inviterID != 0 {
			inviterName := inviterUsername
			if inviterName == "" {
				inviterName = fmt.Sprintf("пользователь с ID %d", inviterID)
			}

			// Отправка уведомления новому пользователю
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Вас пригласил %s", inviterName))
			bot.Send(msg)
		}
	}

}

/*
func getAdminUsername(adminUserID int64) string {
	// Поиск имени пользователя администратора по его ID
	var username string
	err := db.DB.QueryRow("SELECT username FROM users WHERE user_id = $1", adminUserID).Scan(&username)
	if err != nil {
		log.Printf("Ошибка получения username администратора: %v", err)
		return "admin"
	}
	return username
}
*/

