package bot

import (
	"fmt"
	"log"
	"strings"

	"github.com/d1mk9/tgPathToMeBot/configs"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Храним выбранные варианты для каждого пользователя
var userStates = make(map[int64]*UserState)

// Состояние пользователя
type UserState struct {
	CurrentQuestion int
	SelectedOptions map[int][]string
}

func StartBot() {
	bot, err := tgbotapi.NewBotAPI(configs.GlobalConfig.BotToken)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Аккаунт %s авторизован", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			handleMessage(bot, update.Message)
		} else if update.CallbackQuery != nil {
			handleCallbackQuery(bot, update.CallbackQuery)
		}
	}
}

func createOptionsKeyboard(options []string) tgbotapi.InlineKeyboardMarkup {
	var rows []tgbotapi.InlineKeyboardButton
	for i, option := range options {
		rows = append(rows, tgbotapi.NewInlineKeyboardButtonData(option, "option_"+fmt.Sprint(i+1)))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardButtonData("Далее", "next_question"))
	return tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(rows...))
}

func updateMessageWithSelectedOptions(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery, selected []string) {
	// Формируем текст с выбранными вариантами
	selectedText := "Выбранные варианты: " + strings.Join(selected, ", ")

	// Получаем текущее состояние пользователя
	userState := userStates[callbackQuery.Message.Chat.ID]

	// Создаем клавиатуру с вариантами ответов
	inlineKeyboard := createOptionsKeyboard(questions[userState.CurrentQuestion].Options)

	// Добавляем кнопку "Далее" в клавиатуру
	inlineKeyboard.InlineKeyboard = append(inlineKeyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Далее", "next_question"),
	))

	// Обновляем подпись (caption) сообщения
	editMsg := tgbotapi.NewEditMessageCaption(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, selectedText)
	editMsg.ReplyMarkup = &inlineKeyboard

	// Отправляем обновленное сообщение
	_, err := bot.Send(editMsg)
	if err != nil {
		log.Print("Ошибка при обновлении сообщения:", err)
	}
}
