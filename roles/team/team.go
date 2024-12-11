package team

import (
	"fmt"
	"io"
	"log"
	"main/converter"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

func ShowNeiroMenu(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	buttons := []tgbotapi.KeyboardButton{
		{Text: "ElevenLab"},
		{Text: "Facebook"},
		{Text: "Назад"},
	}
	keyboard := tgbotapi.NewReplyKeyboard(buttons)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите нейросеть:")
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func HandleNeiro(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Загрузите файл:")
	bot.Send(msg)
}

func HandleNeiroVideo(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	file := update.Message.Video
	if file == nil {
		log.Println("Не удалось найти видео файл в сообщении.")
		return
	}

	// Скачиваем файл из сообщения
	err := downloadUploadedFile(bot, file.FileID, "downloads")
	if err != nil {
		log.Println("Ошибка при загрузке файла:", err)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Файл загружен, дождитесь обработки.")
	bot.Send(msg)

	converter.SaveFile(file.FileID)

	sendProcessedFiles(bot, update.Message.Chat.ID, file.FileID)
}

func HandleNeiroAudio(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	file := update.Message.Audio
	if file == nil {
		log.Println("Не удалось найти аудио файл в сообщении.")
		return
	}

	// Скачиваем файл из сообщения
	err := downloadUploadedFile(bot, file.FileID, "downloads")
	if err != nil {
		log.Println("Ошибка при загрузке файла:", err)
		return
	}

	converter.SaveFile(file.FileID)

	sendProcessedFiles(bot, update.Message.Chat.ID, file.FileID)
}

// downloadUploadedFile загружает файл по FileID и сохраняет его в указанную директорию
func downloadUploadedFile(bot *tgbotapi.BotAPI, fileID string, downloadDir string) error {
	// Получение файла от Telegram
	fileConfig := tgbotapi.FileConfig{FileID: fileID}
	tgFile, err := bot.GetFile(fileConfig)
	if err != nil {
		return fmt.Errorf("ошибка при получении файла: %w", err)
	}

	// Создаем директорию для загрузки файла, если ее нет
	if _, err := os.Stat(downloadDir); os.IsNotExist(err) {
		os.Mkdir(downloadDir, os.ModePerm)
	}

	// Определяем имя файла и путь
	fileName := strings.ToLower(filepath.Base(tgFile.FilePath))
	filePath := filepath.Join(downloadDir, fileName)

	// Загружаем данные
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(tgFile.Link(bot.Token))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func sendProcessedFiles(bot *tgbotapi.BotAPI, chatID int64, videoFilePath string) {
	extension := filepath.Ext(videoFilePath)
	baseFileName := filepath.Base(videoFilePath[:len(videoFilePath)-len(extension)])

	audioOutputPath := filepath.Join("output/audio", baseFileName+".mp3")
	videoOutputPath := filepath.Join("output/video_no_audio", baseFileName+extension)
	videoOutputTranslation := filepath.Join("output/video_translation_audio", baseFileName+extension)

	audioFile := tgbotapi.NewDocument(chatID, tgbotapi.FilePath(audioOutputPath))
	_, err := bot.Send(audioFile)
	if err != nil {
		log.Println("Ошибка при отправке аудио:", err)
		return
	}

	videoFile := tgbotapi.NewDocument(chatID, tgbotapi.FilePath(videoOutputPath))
	_, err = bot.Send(videoFile)
	if err != nil {
		log.Println("Ошибка при отправке видео:", err)
		return
	}

	videoFileTranslation := tgbotapi.NewDocument(chatID, tgbotapi.FilePath(videoOutputTranslation))
	_, err = bot.Send(videoFileTranslation)
	if err != nil {
		log.Println("Ошибка при отправке видео:", err)
		return
	}

	log.Printf("Аудио и видео отправлены пользователю: %s, %s", audioOutputPath, videoOutputPath)

	err = clearFolders("output/audio", "output/video_no_audio", "downloads", "output/video_translation_audio")
	if err != nil {
		log.Println("Ошибка при очистке папок:", err)
	} else {
		log.Println("Папки успешно очищены.")
	}
}

func clearFolders(folders ...string) error {
	for _, folder := range folders {
		err := os.RemoveAll(folder)
		if err != nil {
			return fmt.Errorf("не удалось очистить папку %s: %w", folder, err)
		}
		// Создаем папку заново, чтобы она существовала для следующих загрузок
		err = os.MkdirAll(folder, 0755)
		if err != nil {
			return fmt.Errorf("не удалось создать папку %s: %w", folder, err)
		}
	}
	return nil
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
