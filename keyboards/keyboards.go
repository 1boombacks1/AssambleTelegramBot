package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	COMING = "coming"
	LATER  = "later"
	SHOW   = "show"
)

const (
	AssembleText = "Квадрат ОБЩИЙ СБОР!👊"
	ComingText   = "Уже выдвигаюсь!🧑‍🦽"
	LaterText    = "Буду попозжа!🤼"
	ShowText     = "Показать рыцарей ⚔️"
)

var AssembleKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(AssembleText),
	),
)

var InlineArriveKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(ComingText, COMING),
		tgbotapi.NewInlineKeyboardButtonData(LaterText, LATER),
	),
)
var UpdatedInlineArriveKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(ComingText, COMING),
		tgbotapi.NewInlineKeyboardButtonData(LaterText, LATER),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(ShowText, SHOW),
	),
)
