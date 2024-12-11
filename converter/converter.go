package converter

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)



func SaveFile(filePath string){
	audioOutputFolder := "output/audio"
	videoOutputFolder := "output/video_no_audio"
	videoExportFolder := "output/video_translation_audio"


	// Сохранение аудио
	err := ExtractAudioAndSave(filePath, audioOutputFolder)
	if err != nil {
		log.Fatalf("Ошибка: %v", err)
	}

	// Сохранение видео без звука
	err = RemoveAudioAndSave(filePath, videoOutputFolder)
	if err != nil {
		log.Fatalf("Ошибка: %v", err)
	}

	// Получаем базовое имя файла без расширения
	baseFileName := filepath.Base(filePath)

	// Определяем пути к аудио и видео
	
	videoFilePath := filepath.Join(videoExportFolder, baseFileName)


	MergeAudioVideo(videoFilePath, baseFileName)
	
}

// Ищем файл с указанным именем в заданной директории
func findFileWithPrefix(dir, fileName string) (string, error) {
	fullPath := filepath.Join(dir, fileName)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("файл с именем %s не найден в директории %s", fileName, dir)
	}
	return fullPath, nil
}

// Ищем файл
func findFile(dir string) (string, error) {
	var foundFile string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Проверяем, что это не директория
		if !info.IsDir() {
			foundFile = path
			return fmt.Errorf("file found") // Используем ошибку как сигнал для остановки
		}
		return nil
	})

	if err != nil && err.Error() == "file found" {
		return foundFile, nil
	}

	return "", fmt.Errorf("файлы не найдены в директории %s", dir)
}

// MergeAudioVideo объединяет видео и аудио и сохраняет результат в указанную папку
func MergeAudioVideo(outputFolder, baseFileName string)  {
	// Определяем пути к директориям
	videoDir := filepath.Join("output", "video_no_audio")
	audioDir := filepath.Join("export", "audio")

	// Ищем видеофайл
	videoFilePath, err := findFileWithPrefix(videoDir, baseFileName)
	if err != nil {
		log.Printf("ошибка при поиске видео файла:", err)
	}



	// Ищем аудиофайл
	audioFilePath, err := findFile(audioDir)
	if err != nil {
		log.Printf("ошибка при поиске аудио файла:", err)
	}

	// Выводим найденные пути для диагностики
	log.Printf("Найден видео файл: %s", videoFilePath)
	log.Printf("Найден аудио файл: %s", audioFilePath)

	// Формируем путь для сохранения выходного файла
	//outputFilePath := filepath.Join(outputFolder, fmt.Sprintf("%s_merged.mp4", baseFileName))
	//log.Fatalf(outputFilePath)

// Проверка существования выходной директории
	if _, err := os.Stat("output/video_translation_audio"); os.IsNotExist(err) {
		err := os.MkdirAll("output/video_translation_audio", os.ModePerm)
		if err != nil {
			log.Printf("Ошибка при создании директории: %v", err)
		}
	}

	// Формируем команду ffmpeg
	cmd := exec.Command("ffmpeg",
		"-i", videoFilePath,
		"-i", audioFilePath,
		"-filter_complex", "[0:v][1:a]concat=n=1:v=1:a=1[outv][outa]",
		"-map", "[outv]",
		"-map", "[outa]",
		"-c:v", "libx264",
		"-c:a", "aac",
		"-strict", "experimental",
		outputFolder,
		"-y",
	)

	// Запускаем команду
	err = cmd.Run()
	if err != nil {
		log.Printf("Ошибка при объединении видео и аудио: %v", err)
	}

	log.Printf("Видео с аудио сохранено в: %s\n", outputFolder)
}



// ExtractAudioAndSave сохраняет аудио из видеофайла в указанную папку
func ExtractAudioAndSave(videoFilePath, outputFolder string) error {
	// Получаем расширение файла видео
	extension := filepath.Ext(videoFilePath)
	// Получаем имя файла без расширения
	baseFileName := filepath.Base(videoFilePath[:len(videoFilePath)-len(extension)])
	// Путь для сохранения аудио файла
	audioOutputPath := filepath.Join(outputFolder, baseFileName+".mp3")

	// Проверяем, существует ли выходная папка
	if _, err := os.Stat(outputFolder); os.IsNotExist(err) {
		err := os.MkdirAll(outputFolder, 0755)
		if err != nil {
			return fmt.Errorf("не удалось создать папку: %w", err)
		}
	}

	// Используем ffmpeg-go для извлечения аудио
	err := ffmpeg.Input(videoFilePath).
		Output(audioOutputPath, ffmpeg.KwArgs{"vn": ""}). // "vn" опция убирает видео дорожку
		OverWriteOutput().                                // Перезаписывать выходной файл, если существует
		Run()
	if err != nil {
		return fmt.Errorf("ошибка при разделении видео и аудио: %w", err)
	}

	log.Printf("Аудио сохранено в: %s\n", audioOutputPath)
	return nil
}

// RemoveAudioAndSave сохраняет видео без звука в указанную папку
func RemoveAudioAndSave(videoFilePath, outputFolder string) error {
	// Получаем расширение файла видео
	extension := filepath.Ext(videoFilePath)
	// Получаем имя файла без расширения
	baseFileName := filepath.Base(videoFilePath[:len(videoFilePath)-len(extension)])
	// Путь для сохранения видео без звука
	videoOutputPath := filepath.Join(outputFolder, baseFileName+extension)

	// Проверяем, существует ли выходная папка
	if _, err := os.Stat(outputFolder); os.IsNotExist(err) {
		err := os.MkdirAll(outputFolder, 0755)
		if err != nil {
			return fmt.Errorf("не удалось создать папку: %w", err)
		}
	}

	// Используем ffmpeg-go для удаления аудио дорожки из видео
	err := ffmpeg.Input(videoFilePath).
		Output(videoOutputPath, ffmpeg.KwArgs{"an": ""}). // "an" опция убирает аудио дорожку
		OverWriteOutput().                                // Перезаписывать выходной файл, если существует
		Run()
	if err != nil {
		return fmt.Errorf("ошибка при сохранении видео без звука: %w", err)
	}

	log.Printf("Видео без звука сохранено в: %s\n", videoOutputPath)
	return nil
}