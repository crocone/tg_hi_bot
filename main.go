package main

import (
	"fmt"
	tgbotapi "github.com/crocone/tg-bot"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var bot *tgbotapi.BotAPI
var startMenu = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Скажи привет", "hi"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Скажи пока", "buy"),
	),
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(".env not loaded")
	}

	bot, err = tgbotapi.NewBotAPI(os.Getenv("TG_API_BOT_TOKEN"))
	if err != nil {
		log.Fatalf("Failed to initialize Telegram bot API: %v", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatalf("Failed to start listening for updates %v", err)
	}

	for update := range updates {
		if update.CallbackQuery != nil {
			callbacks(update)
		} else if update.Message.IsCommand() {
			commands(update)
		} else {
			// simply message
		}
	}
}

func callbacks(update tgbotapi.Update) {
	data := update.CallbackQuery.Data
	chatId := update.CallbackQuery.From.ID
	firstName := update.CallbackQuery.From.FirstName
	lastName := update.CallbackQuery.From.LastName
	var text string
	switch data {
	case "hi":
		text = fmt.Sprintf("Привет %v %v", firstName, lastName)
	case "buy":
		text = fmt.Sprintf("Пока %v %v", firstName, lastName)
	default:
		text = "Неизвестная команда"
	}
	msg := tgbotapi.NewMessage(chatId, text)
	sendMessage(msg)
}

func commands(update tgbotapi.Update) {
	command := update.Message.Command()
	switch command {
	case "start":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите действие")
		msg.ReplyMarkup = startMenu
		msg.ParseMode = "Markdown"
		sendMessage(msg)
	}
}

func sendMessage(msg tgbotapi.Chattable) {
	if _, err := bot.Send(msg); err != nil {
		log.Panicf("Send message error: %v", err)
	}
}
