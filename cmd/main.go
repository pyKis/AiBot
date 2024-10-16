package main

import (
	"fmt"
	"log"
	"os"


	"main/bot"
	"main/db"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var userID int64

func main() {
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Ошибка загрузки файла .env")
	}

	// Получение значений переменных окружения
	botToken := os.Getenv("BOT_TOKEN")
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	

	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)

	err := db.ConnectToDB(databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.DB.Close()

	err = db.CreateTables()
	if err != nil {
		log.Fatal(err)
	}

	botAPI, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	botAPI.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := botAPI.GetUpdatesChan(u)

	
	for update := range updates {
		bot.HandleUpdate(botAPI, update)
	}
}
