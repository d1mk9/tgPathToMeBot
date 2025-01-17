package bot

import (
	"fmt"
	"log"
	"sort"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	log.Printf("Получено сообщение от %s: %s", message.From.UserName, message.Text)

	// Приветственное сообщение с кнопкой "Запустить тест"
	if message.IsCommand() && message.Command() == "start" {
		sendWelcomeMessage(bot, message.Chat.ID)
	}
}

func handleCallbackQuery(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	log.Printf("Получен callback от %s: %s", callbackQuery.From.UserName, callbackQuery.Data)

	chatID := callbackQuery.Message.Chat.ID

	switch callbackQuery.Data {
	case "start_test":
		// Инициализация состояния пользователя
		userStates[chatID] = &UserState{
			CurrentQuestion: 0,
			SelectedOptions: make(map[int][]string),
		}
		// Отправляем первый вопрос
		sendQuestionAsImage(bot, chatID, 0)
	default:
		// Обработка выбора вариантов
		if strings.HasPrefix(callbackQuery.Data, "option_") {
			handleOptionSelection(bot, callbackQuery)
		} else if callbackQuery.Data == "finish_selection" {
			handleFinishSelection(bot, callbackQuery)
		} else if callbackQuery.Data == "next_question" {
			// Переход к следующему вопросу
			userState := userStates[chatID]
			userState.CurrentQuestion++

			// Проверяем, был ли это последний вопрос
			if userState.CurrentQuestion >= len(questions) {
				// Завершаем тест
				handleFinishTest(bot, chatID)
			} else {
				// Отправляем следующий вопрос
				sendQuestionAsImage(bot, chatID, userState.CurrentQuestion)
			}
		}
	}

	// Подтверждение обработки callback
	callbackConfig := tgbotapi.NewCallback(callbackQuery.ID, "")
	_, err := bot.Request(callbackConfig)
	if err != nil {
		log.Print("Ошибка при подтверждении callback:", err)
	}
}

func handleOptionSelection(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	chatID := callbackQuery.Message.Chat.ID
	option := callbackQuery.Data // Это строка вида "option_1", "option_2" и т.д.

	// Удаляем префикс "option_"
	optionNumber := strings.TrimPrefix(option, "option_")

	// Получаем текущее состояние пользователя
	userState := userStates[chatID]

	// Проверяем, выбран ли вариант
	if contains(userState.SelectedOptions[userState.CurrentQuestion], optionNumber) {
		// Удаляем вариант из списка
		userState.SelectedOptions[userState.CurrentQuestion] = remove(userState.SelectedOptions[userState.CurrentQuestion], optionNumber)
	} else {
		// Добавляем вариант в список, если он еще не выбран
		userState.SelectedOptions[userState.CurrentQuestion] = append(userState.SelectedOptions[userState.CurrentQuestion], optionNumber)
	}

	// Обновляем сообщение с текущими выбранными вариантами
	updateMessageWithSelectedOptions(bot, callbackQuery, userState.SelectedOptions[userState.CurrentQuestion])
}

func handleFinishSelection(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery) {
	chatID := callbackQuery.Message.Chat.ID
	userState := userStates[chatID]

	// Переход к следующему вопросу
	userState.CurrentQuestion++
	sendQuestion(bot, chatID, userState.CurrentQuestion)
}

func handleFinishTest(bot *tgbotapi.BotAPI, chatID int64) {
	userState := userStates[chatID]
	var finalResponse strings.Builder

	finalResponse.WriteString("Вы завершили тест. Ваши варианты:\n")

	// Создаем срез для хранения ключей (номеров вопросов)
	questionIndices := make([]int, 0, len(userState.SelectedOptions))
	for questionIndex := range userState.SelectedOptions {
		questionIndices = append(questionIndices, questionIndex)
	}

	// Сортируем ключи
	sort.Ints(questionIndices)

	// Проходим по отсортированным ключам и добавляем результаты в итоговое сообщение
	for _, questionIndex := range questionIndices {
		options := userState.SelectedOptions[questionIndex]
		finalResponse.WriteString(fmt.Sprintf("%s: %s\n", questions[questionIndex].Text, strings.Join(options, ", ")))
	}

	// Подсчет выбранных вариантов
	count := countSelectedOptions(userState.SelectedOptions)

	// Определение результата
	resultText := getResultDescription(count)

	// Отправляем текстовое описание
	finalResponse.WriteString("\nРезультат:\n")
	finalResponse.WriteString(resultText)
	msg := tgbotapi.NewMessage(chatID, finalResponse.String())
	_, err := bot.Send(msg)
	if err != nil {
		log.Print("Ошибка при отправке итогового сообщения:", err)
	}

	// Очистка состояния пользователя
	delete(userStates, chatID)
}
