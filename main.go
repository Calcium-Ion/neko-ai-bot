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
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			receiveMsg := update.Message.Text
			if strings.HasPrefix(receiveMsg, "/") {
				bot.RunCommand(strings.Split(receiveMsg[1:], " ")[0], *update.Message, *tgBot)
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
				bot.Change(*update.CallbackQuery.Message, *tgBot, index, taskId, "UPSCALE")
			} else if strings.HasPrefix(update.CallbackQuery.Data, "/v") {
				taskId := strings.Split(update.CallbackQuery.Data, " ")[2]
				index, err := strconv.Atoi(strings.Split(update.CallbackQuery.Data, " ")[1])
				if err != nil {
					log.Println(err)
					continue
				}
				log.Printf("VARIATION taskId: %s, index: %d", taskId, index)
				bot.Change(*update.CallbackQuery.Message, *tgBot, index, taskId, "VARIATION")
			}
		}
	}
}
