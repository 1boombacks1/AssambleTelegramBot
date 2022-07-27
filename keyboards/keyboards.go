package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	COMING = "coming"
	LATER  = "later"
)

var AssembleKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Квадрат ОБЩИЙ СБОР!👊"),
	),
)

var InlineArriveKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Уже выдвигаюсь!🧑‍🦽", COMING),
		tgbotapi.NewInlineKeyboardButtonData("Одалею монстра и подскачу!🤼", LATER),
	),
)
