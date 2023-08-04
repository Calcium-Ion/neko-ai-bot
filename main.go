package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"neko-ai-bot/bot"
	"neko-ai-bot/conf"
	"neko-ai-bot/model"
	"strconv"
	"strings"
)

func GetUser(chatId int64, username string, userId int64, tgBot tgbotapi.BotAPI) *model.User {
	user, err, init := model.GetUserOrInit(chatId, username, userId)
	if err != nil || user == nil || user.Id == 0 {
		if err != nil {
			log.Println(err)
		}
		message := tgbotapi.NewMessage(chatId, "获取用户信息失败，请联系管理员")
		_, _ = tgBot.Send(message)
		return nil
	}
	if init {
		message := tgbotapi.NewMessage(chatId, "初始化用户成功，获赠20积分")
		message.ReplyMarkup = bot.GetMainKeyboard()
		_, _ = tgBot.Send(message)
	}
	return user
}

func main() {
	conf.Setup()
	model.Setup()
	tgBot, err := tgbotapi.NewBotAPI(conf.Conf.BotToken)
	if err != nil {
		log.Panic(err)
	}
	tgBot.Debug = true

	log.Printf("Authorized on account %s", tgBot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := tgBot.GetUpdatesChan(u)

	for update := range updates {
		var user *model.User
		if update.Message != nil {
			user = GetUser(update.Message.Chat.ID, update.Message.From.UserName, update.Message.From.ID, *tgBot)
		} else if update.CallbackQuery != nil {
			user = GetUser(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.From.UserName, update.CallbackQuery.From.ID, *tgBot)
		}
		if user == nil || user.Id == 0 {
			continue
		}
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			receiveMsg := update.Message.Text
			//mainKeyboard := tgbotapi.NewReplyKeyboard(bot.GetMainKeyboard())
			//mainKeyboard.OneTimeKeyboard = false

			if receiveMsg == "查看帮助" {
				bot.Start(*update.Message, *tgBot)
			} else if receiveMsg == "个人信息" {
				bot.UserInfo(user, *update.Message, *tgBot)
			} else if receiveMsg == "签到" {
				bot.Sign(user, *update.Message, *tgBot)
			} else if strings.HasPrefix(receiveMsg, "/") {
				bot.RunCommand(user, strings.Split(receiveMsg[1:], " ")[0], *update.Message, *tgBot)
			}
		} else if update.CallbackQuery != nil {
			//callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			//if _, err := tgBot.Request(callback); err != nil {
			//	log.Printf("callback error: %v", err)
			//	continue
			//}
			// And finally, send a message containing the data received.
			if strings.HasPrefix(update.CallbackQuery.Data, "/u") {
				taskId := strings.Split(update.CallbackQuery.Data, " ")[2]
				index, err := strconv.Atoi(strings.Split(update.CallbackQuery.Data, " ")[1])
				if err != nil {
					log.Println(err)
					continue
				}
				log.Printf("UPSCALE taskId: %s, index: %d", taskId, index)
				bot.Change(user, *update.CallbackQuery.Message, *tgBot, index, taskId, "UPSCALE")
			} else if strings.HasPrefix(update.CallbackQuery.Data, "/v") {
				taskId := strings.Split(update.CallbackQuery.Data, " ")[2]
				index, err := strconv.Atoi(strings.Split(update.CallbackQuery.Data, " ")[1])
				if err != nil {
					log.Println(err)
					continue
				}
				log.Printf("VARIATION taskId: %s, index: %d", taskId, index)
				bot.Change(user, *update.CallbackQuery.Message, *tgBot, index, taskId, "VARIATION")
			}
		}
	}
}
