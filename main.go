package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"neko-ai-bot/bot"
	"neko-ai-bot/conf"
	"neko-ai-bot/model"
	"neko-ai-bot/util"
	"strconv"
	"strings"
	"time"
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
		message.ReplyMarkup = bot.GetMainKeyboard(chatId)
		_, _ = tgBot.Send(message)
	}
	return user
}

func CheckStatus(tgBog tgbotapi.BotAPI) {
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			now := time.Now()
			currentUnixTimestamp := now.Unix()
			respBody, err := util.DoRequest(fmt.Sprintf("/api/log/stat?&type=0&token_name=&model_name=&start_timestamp=%d&end_timestamp=%d", currentUnixTimestamp-60, currentUnixTimestamp))
			if respBody["success"].(bool) {
				if err != nil {
					log.Println(err)
					continue
				}
				if int(respBody["data"].(map[string]interface{})["rpm"].(float64)) > 800 {
					msgS := "NekoAPI高负载警告\n"
					if conf.Conf.AccessToken != "" {
						msg := tgbotapi.NewMessage(6134547155, msgS)
						processMsg, err := tgBog.Send(msg)
						if err != nil {
							log.Println(err)
						}
						if err != nil {
							log.Println(err)
							msgS += "今日RPM：获取失败\n"
						} else {
							if respBody["success"].(bool) {
								log.Println(respBody)
								msgS += fmt.Sprintf("上一分钟RPM：%d\n", int(respBody["data"].(map[string]interface{})["rpm"].(float64)))
								msgS += fmt.Sprintf("上一分钟TPM：%d\n", int(respBody["data"].(map[string]interface{})["tpm"].(float64)))
							}
							msg := tgbotapi.NewEditMessageText(6134547155, processMsg.MessageID, msgS)
							_, _ = tgBog.Send(msg)
						}
					}
				}
			}
		}
	}()
}

func main() {
	conf.Setup()
	model.Setup()

	tgBot, err := tgbotapi.NewBotAPI(conf.Conf.BotToken)
	if err != nil {
		log.Panic(err)
	}
	tgBot.Debug = true

	CheckStatus(*tgBot)

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
			} else if receiveMsg == "个人信息" || receiveMsg == "/me" {
				bot.UserInfo(user, *update.Message, *tgBot)
			} else if receiveMsg == "签到" || receiveMsg == "/sign" {
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
