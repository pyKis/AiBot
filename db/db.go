package db

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	"main/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectToDB(databaseURL string) error {
	var err error
	DB, err = sql.Open("postgres", databaseURL)
	if err != nil {
		return err
	}
	return DB.Ping()
}

func CreateTables() error {
	usersTable := `
    CREATE TABLE IF NOT EXISTS users (
        user_id BIGINT PRIMARY KEY,
        username TEXT,
        first_name TEXT,
        last_name TEXT,
        phone_number TEXT,
        role TEXT DEFAULT 'None'
    );`

	referralsTable := `
    CREATE TABLE IF NOT EXISTS referrals (
        referral_code VARCHAR(8) PRIMARY KEY,
        user_id BIGINT REFERENCES users(user_id)
    );`

	_, err := DB.Exec(usersTable)
	if err != nil {
		return fmt.Errorf("ошибка при создании таблицы users: %v", err)
	}

	_, err = DB.Exec(referralsTable)
	if err != nil {
		return fmt.Errorf("ошибка при создании таблицы referrals: %v", err)
	}

	return nil
}

func saveUser(user *models.User) error {
	_, err := DB.Exec(`
        INSERT INTO users (user_id, username, first_name, last_name, role) 
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (user_id) DO NOTHING`,
		user.UserID, user.Username, user.FirstName, user.LastName, user.Role)
	return err
}

func UpdateUserRole(userID int64, role models.Role) error {
	// Логируем переданные данные для отладки
	log.Printf("Обновление роли для user_id: %d, новая роль: %s", userID, role)

	// Выполняем обновление роли пользователя в базе данных
	result, err := DB.Exec(`UPDATE users SET role=$1 WHERE user_id=$2`, role, userID)
	if err != nil {
		log.Printf("Ошибка при выполнении запроса: %v", err)
		return fmt.Errorf("ошибка при обновлении роли для user_id %d: %v", userID, err)
	}

	// Проверяем, были ли изменены строки
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Ошибка при проверке измененных строк: %v", err)
		return fmt.Errorf("ошибка при проверке измененных строк: %v", err)
	}

	if rowsAffected == 0 {
		log.Printf("Роль не была обновлена для user_id: %d, возможно, пользователь не найден", userID)
		return fmt.Errorf("пользователь с user_id %d не найден", userID)
	}

	log.Printf("Роль успешно обновлена для user_id: %d", userID)
	return nil
}


func DeleteUser(userID int64) error {
	_, err := DB.Exec(`DELETE FROM users WHERE user_id = $1`, userID)
	return err
}

func GetUserRole(userID int64) (string, error) {
	var role string
	err := DB.QueryRow("SELECT role FROM users WHERE user_id = $1", userID).Scan(&role)
	if err != nil {
		if err == sql.ErrNoRows {
			role = "None" // Если пользователь не найден, по умолчанию None
		} else {
			log.Println("Ошибка получения роли:", err)
			return "", err
		}
	}
	return role, nil
}
// Получение общей статистики пользователей
func GetTotalUsers() (int, error) {
	var totalUsers int
	err := DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&totalUsers)
	if err != nil {
		log.Println("Ошибка при получении общего количества пользователей:", err)
		return 0, err
	}
	return totalUsers, nil
}

// Получение количества приглашенных пользователей
func GetInvitedUsers() (int, error) {
	var invitedUsers int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE referrer_id IS NOT NULL").Scan(&invitedUsers)
	if err != nil {
		log.Println("Ошибка при получении количества приглашенных пользователей:", err)
		return 0, err
	}
	return invitedUsers, nil
}

// Получение количества пользователей по ролям
func GetUsersByRole(role models.Role) (int, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE role = $1", role).Scan(&count)
	if err != nil {
		log.Println("Ошибка при получении количества пользователей с ролью:", role, err)
		return 0, err
	}
	return count, nil
}

func GetUserIDByUsername(username string, userID *int64) error {
	// Убираем символ '@' в начале, если он есть
	if len(username) > 0 && username[0] == '@' {
		username = username[1:]
	}

	// Запрашиваем user_id по username
	err := DB.QueryRow("SELECT user_id FROM users WHERE username = $1", username).Scan(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Если пользователь не найден, возвращаем ошибку
			return fmt.Errorf("пользователь с логином '%s' не найден", username)
		}
		// Обрабатываем другие возможные ошибки
		return fmt.Errorf("ошибка при получении user_id: %v", err)
	}

	return nil
}


func GetInviterInfo(referralCode string) (int64, string, error) {
	var inviterID int64
	var inviterUsername string

	// Поиск пригласившего пользователя по реферальному коду
	err := DB.QueryRow(`
        SELECT u.user_id, u.username
        FROM referrals r
        JOIN users u ON r.user_id = u.user_id
        WHERE r.referral_code = $1
    `, referralCode).Scan(&inviterID, &inviterUsername)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, "", nil // реферальный код не найден
		}
		return 0, "", err
	}

	return inviterID, inviterUsername, nil
}



func SaveUser(user *tgbotapi.User) {
	_, err := DB.Exec(`
        INSERT INTO users (user_id, username, first_name, last_name, role) 
        VALUES ($1, $2, $3, $4, $5) 
        ON CONFLICT (user_id) DO NOTHING`,
		user.ID, user.UserName, user.FirstName, user.LastName, "None",
	)
	if err != nil {
		log.Println("Ошибка сохранения пользователя:", err)
	}
}

func GenerateReferralLink(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	referralCode := GenerateReferralCode()
	_, err := DB.Exec(`INSERT INTO referrals (referral_code, user_id) VALUES ($1, $2) ON CONFLICT (referral_code) DO NOTHING`, referralCode, update.Message.From.ID)
	if err != nil {
		log.Println("Ошибка генерации реферальной ссылки:", err)
		return
	}

	referralLink := fmt.Sprintf("https://t.me/%s?start=%s", bot.Self.UserName, referralCode)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Ваша реферальная ссылка: %s", referralLink))
	bot.Send(msg)
}

func GenerateReferralCode() string {
	rand.Seed(time.Now().UnixNano())
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	code := make([]rune, 8)
	for i := range code {
		code[i] = letters[rand.Intn(len(letters))]
	}
	return string(code)
}

func HandleContact(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	contact := update.Message.Contact
	_, err := DB.Exec(`UPDATE users SET phone_number=$1 WHERE user_id=$2`, contact.PhoneNumber, update.Message.From.ID)
	if err != nil {
		log.Println("Ошибка сохранения контакта:", err)
	} else {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Контакт успешно сохранён!")
		bot.Send(msg)
	}
}

func CheckUserExists(userID int64) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE user_id = $1)`
	err := DB.QueryRow(query, userID).Scan(&exists)
	if err != nil {
		log.Println("Ошибка проверки существования пользователя:", err)
		return false, err
	}
	return exists, nil
}