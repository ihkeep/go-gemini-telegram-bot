package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"go-gemini-telegram-bot/app"

	"github.com/joho/godotenv"
)

var config Config

type Config struct {
	BotToken       string
	Gemini_API_KEY string
}

func LoadConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, trying to load from environment")
	}

	config := Config{
		BotToken:       getEnv("BOT_TOKEN", ""),
		Gemini_API_KEY: getEnv("Gemini_API_KEY", ""),
	}

	if config.BotToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN must be set in environment variables or .env file")
	}

	return config
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func main() {
	log.Println("Starting bot")
	config = LoadConfig()
	log.Println("Loaded config")

	start_bot()
}

func start_bot() {
	bot, err := tgbotapi.NewBotAPI(config.BotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				switch update.Message.Command() {
				case "start":
					app.StartCommand(update, bot)
				case "new":
					app.NewChatCommand(update, bot)
				default:
					app.DefaultCommand(update, bot)
				}
			} else if update.Message.Text != "" {
				app.HandleText(update, bot)
			} else if update.Message.Photo != nil {
				app.HandlePhoto(update, bot)
			}

		}

	}
}