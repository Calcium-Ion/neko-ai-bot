package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"log"
	"neko-ai-bot/api"
	"net/http"
	"strings"
	"time"
)

func RunCommand(cmdText string, message tgbotapi.Message, tgBog tgbotapi.BotAPI) {
	log.Println("run command: ", cmdText)

	switch cmdText {
	case "start":
		Start(message, tgBog)
	case "imagine":
		Imagine(message, tgBog)
	case "test":
		test(message, tgBog)
	}
}

func test(message tgbotapi.Message, tgBog tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "test")
	msg.ReplyMarkup = GetChangeKeyboard("545641541", "", "TEST")
	_, err := tgBog.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func Change(message tgbotapi.Message, tgBog tgbotapi.BotAPI, index int, taskId string, action string) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "请稍等，正在发送绘图请求")
	processMsg, err := tgBog.Send(msg)
	if err != nil {
		log.Println(err)
		return
	}
	go func() {
		result := api.Change(taskId, action, index)
		if result.Code == 1 || result.Code == 22 {
			msg := tgbotapi.NewEditMessageText(message.Chat.ID, processMsg.MessageID, "正在绘图中，当前进度： 排队中")
			_, _ = tgBog.Send(msg)
		} else {
			msg := tgbotapi.NewEditMessageText(message.Chat.ID, processMsg.MessageID, "绘图失败，原因："+result.Description)
			_, _ = tgBog.Send(msg)
			return
		}
		taskId := result.Data.(string)
		log.Printf("result: %+v", result)
		for {
			result = api.Fetch(taskId)
			if result.Data == nil {
				continue
			}
			midjourneyResponse := result.Data.(api.Midjourney)
			if midjourneyResponse.FailReason != "" {
				msg := tgbotapi.NewEditMessageText(message.Chat.ID, processMsg.MessageID, "绘图失败，原因："+midjourneyResponse.FailReason)
				_, _ = tgBog.Send(msg)
				return
			}
			log.Printf("midjourneyResponse: %+v", midjourneyResponse)
			if midjourneyResponse.Status != "SUCCESS" {
				if midjourneyResponse.Status == "IN_PROGRESS" {
					msgStr := "正在绘图中，当前进度：" + midjourneyResponse.Progress
					if msgStr != processMsg.Text {
						msg := tgbotapi.NewEditMessageText(message.Chat.ID, processMsg.MessageID, msgStr)
						_, _ = tgBog.Send(msg)
					}
				}
				time.Sleep(time.Duration(5) * time.Second)
				continue
			} else {
				msg := tgbotapi.NewEditMessageText(message.Chat.ID, processMsg.MessageID, "绘图完成，正在下载图片")
				_, _ = tgBog.Send(msg)
				log.Println("download image")
				resp, err := http.Get(midjourneyResponse.ImageUrl + "?width=868&height=868")
				if err != nil {
					log.Println(err)
					return
				}
				defer func() {
					msg = tgbotapi.NewEditMessageText(message.Chat.ID, processMsg.MessageID, "绘图已完成")
					_, _ = tgBog.Send(msg)
				}()
				defer resp.Body.Close()
				// 读取图片数据
				data, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Println(err)
					return
				}

				// 创建一个新的图片消息
				photo := tgbotapi.FileBytes{Name: "image.jpg", Bytes: data}
				imgMsg := tgbotapi.NewPhoto(message.Chat.ID, photo)
				imgMsg.ReplyMarkup = GetChangeKeyboard(midjourneyResponse.MjId, midjourneyResponse.ImageUrl, action)
				// 发送图片消息
				_, err = tgBog.Send(imgMsg)
				if err != nil {
					log.Println(err)
				}

				break
			}
		}

	}()
}

func Imagine(message tgbotapi.Message, tgBog tgbotapi.BotAPI) {
	prompt := strings.TrimPrefix(message.Text, "/imagine")
	prompt = strings.TrimSpace(prompt)
	if prompt == "" {
		msg := tgbotapi.NewMessage(message.Chat.ID, "请输入绘图内容，例如：\n/imagine 可爱猫猫")
		_, err := tgBog.Send(msg)
		if err != nil {
			log.Println(err)
		}
		return
	}
	log.Printf("prompt: [%s]", prompt)
	msg := tgbotapi.NewMessage(message.Chat.ID, "请稍等，正在发送绘图请求")
	processMsg, err := tgBog.Send(msg)
	if err != nil {
		log.Println(err)
		return
	}
	go func() {
		result := api.Imagine("", prompt)
		if result.Code == 1 || result.Code == 22 {
			msg := tgbotapi.NewEditMessageText(message.Chat.ID, processMsg.MessageID, "正在绘图中，当前进度： 排队中")
			_, _ = tgBog.Send(msg)
		} else {
			msg := tgbotapi.NewEditMessageText(message.Chat.ID, processMsg.MessageID, "绘图失败，原因："+result.Description)
			_, _ = tgBog.Send(msg)
			return
		}
		taskId := result.Data.(string)
		log.Printf("result: %+v", result)
		for {
			result = api.Fetch(taskId)
			if result.Data == nil {
				continue
			}
			midjourneyResponse := result.Data.(api.Midjourney)
			if midjourneyResponse.FailReason != "" {
				msg := tgbotapi.NewEditMessageText(message.Chat.ID, processMsg.MessageID, "绘图失败，原因："+midjourneyResponse.FailReason)
				_, _ = tgBog.Send(msg)
				return
			}
			log.Printf("midjourneyResponse: %+v", midjourneyResponse)
			if midjourneyResponse.Status != "SUCCESS" {
				if midjourneyResponse.Status == "IN_PROGRESS" {
					msgStr := "正在绘图中，当前进度：" + midjourneyResponse.Progress
					if msgStr != processMsg.Text {
						msg := tgbotapi.NewEditMessageText(message.Chat.ID, processMsg.MessageID, msgStr)
						_, _ = tgBog.Send(msg)
					}
				}
				time.Sleep(time.Duration(5) * time.Second)
				continue
			} else {
				msg := tgbotapi.NewEditMessageText(message.Chat.ID, processMsg.MessageID, "绘图完成，正在下载图片")
				_, _ = tgBog.Send(msg)
				log.Println("download image")
				resp, err := http.Get(midjourneyResponse.ImageUrl + "?width=868&height=868")
				if err != nil {
					log.Println(err)
					return
				}
				defer func() {
					msg = tgbotapi.NewEditMessageText(message.Chat.ID, processMsg.MessageID, "绘图已完成")
					_, _ = tgBog.Send(msg)
				}()
				defer resp.Body.Close()
				// 读取图片数据
				data, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Println(err)
					return
				}

				// 创建一个新的图片消息
				photo := tgbotapi.FileBytes{Name: "image.jpg", Bytes: data}
				imgMsg := tgbotapi.NewPhoto(message.Chat.ID, photo)
				imgMsg.ReplyMarkup = GetChangeKeyboard(midjourneyResponse.MjId, midjourneyResponse.ImageUrl, "IMAGINE")
				// 发送图片消息
				_, err = tgBog.Send(imgMsg)
				if err != nil {
					log.Println(err)
				}

				break
			}
		}

	}()
}

func Start(message tgbotapi.Message, tgBog tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "输入绘图内容，例如：\n/imagine 可爱猫猫")
	_, err := tgBog.Send(msg)
	if err != nil {
		log.Println(err)
	}
}
