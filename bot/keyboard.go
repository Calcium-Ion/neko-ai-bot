package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"neko-ai-bot/conf"
)

func GetMainKeyboard(chatId int64) tgbotapi.ReplyKeyboardMarkup {
	if chatId < 0 {
		return tgbotapi.NewReplyKeyboard()
	}
	buttons := make([][]tgbotapi.KeyboardButton, 0)
	buttons = append(buttons, tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("查看帮助"),
		tgbotapi.NewKeyboardButton("个人信息"),
	))
	buttons = append(buttons, tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("签到"),
		//tgbotapi.NewKeyboardButton("更多功能"),
	))
	//buttons = append(buttons, tgbotapi.NewKeyboardButtonRow(
	//))
	return tgbotapi.NewReplyKeyboard(buttons...)
}

func GetStatusKeyboard() tgbotapi.InlineKeyboardMarkup {
	buttons := make([][]tgbotapi.InlineKeyboardButton, 0)

	for baseUrl := range conf.Conf.Sites {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(baseUrl, "/status "+baseUrl),
		))
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		buttons...,
	)
}

func GetChangeKeyboard(taskId string, url string, action string) tgbotapi.InlineKeyboardMarkup {
	buttons := make([][]tgbotapi.InlineKeyboardButton, 0)

	if action != "UPSCALE" {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("放大1", "/u 1 "+taskId),
			tgbotapi.NewInlineKeyboardButtonData("放大2", "/u 2 "+taskId),
			tgbotapi.NewInlineKeyboardButtonData("放大3", "/u 3 "+taskId),
			tgbotapi.NewInlineKeyboardButtonData("放大4", "/u 4 "+taskId),
		))
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("变换1", "/v 1 "+taskId),
			tgbotapi.NewInlineKeyboardButtonData("变换2", "/v 2 "+taskId),
			tgbotapi.NewInlineKeyboardButtonData("变换3", "/v 3 "+taskId),
			tgbotapi.NewInlineKeyboardButtonData("变换4", "/v 4 "+taskId),
		))
	}
	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("点此查看原图", url),
	))
	return tgbotapi.NewInlineKeyboardMarkup(
		buttons...,
	)
}
