package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config содержит все глобальные конфигурационные параметры
type Config struct {
	BotToken string // Токен для доступа к Telegram Bot API
}

// GlobalConfig - глобальная переменная для хранения конфигурации
var GlobalConfig Config

// LoadConfig загружает конфигурацию из файла .env и переменных окружения
func LoadConfig() {
	// Пытаемся загрузить файл .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Файл .env не найден. Используются переменные окружения.")
	}

	// Загружаем токен бота
	GlobalConfig.BotToken = os.Getenv("TELEGRAM_APITOKEN_PATHTOMEBOT")
	if GlobalConfig.BotToken == "" {
		log.Fatal("Переменная окружения TELEGRAM_APITOKEN_PATHTOMEBOT не установлена")
	}

	log.Println("Конфигурация успешно загружена")
}
