package bot

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"log"
	"neko-ai-bot/api"
	"neko-ai-bot/conf"
	"neko-ai-bot/model"
	"neko-ai-bot/util"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func RunCommand(user *model.User, cmdText string, message tgbotapi.Message, tgBog tgbotapi.BotAPI) {
	log.Println("run command: ", cmdText)

	switch cmdText {
	case "start":
		Start(message, tgBog)
	case "imagine":
		Imagine(user, message, tgBog)
	case "test":
		test(message, tgBog)
	case "status":
		ShowStatus(message, tgBog)
	case "unlimited":
		Unlimited(user, message, tgBog)
	}
}

func IsAdmin(user *model.User) bool {
	for _, admin := range conf.Conf.AdminUsername {
		if admin == user.Username {
			return true
		}
	}
	return false
}

func test(message tgbotapi.Message, tgBog tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "test")
	msg.ReplyMarkup = GetChangeKeyboard("545641541", "", "TEST")
	_, err := tgBog.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func Unlimited(user *model.User, message tgbotapi.Message, tgBog tgbotapi.BotAPI) {
	args := strings.Split(message.Text, " ")
	if len(args) < 2 {
		//查询自己
		unlimited, err := model.GetUnlimitedByIDOrUsername(user.Id, &user.Username)
		if err != nil {
			msg := tgbotapi.NewMessage(message.Chat.ID, "数据库错误")
			_, _ = tgBog.Send(msg)
			return
		}
		if unlimited == nil {
			msg := tgbotapi.NewMessage(message.Chat.ID, "未开通包天套餐，请联系管理员")
			_, _ = tgBog.Send(msg)
			return
		}
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("包天GPT套餐\nKey: %s\nTPM限制: %d\nRPM限制: %d", unlimited.Key, unlimited.TokenLimit, unlimited.RateLimit))
		_, _ = tgBog.Send(msg)
		return
	}
	if !IsAdmin(user) {
		msg := tgbotapi.NewMessage(message.Chat.ID, "无权限")
		_, _ = tgBog.Send(msg)
		return
	}
	cmd := args[1]
	if cmd == "list" {

		unlimited, err := model.GetAllUnlimited()
		if err != nil {
			msg := tgbotapi.NewMessage(message.Chat.ID, "数据库错误")
			_, _ = tgBog.Send(msg)
			return
		}
		var text string
		for _, user := range unlimited {
			text += fmt.Sprintf("%s %s %d %d\n", user.Username, user.Key, user.TokenLimit, user.RateLimit)
		}
		msg := tgbotapi.NewMessage(message.Chat.ID, text)
		_, _ = tgBog.Send(msg)
	} else {
		if len(args) < 3 {
			msg := tgbotapi.NewMessage(message.Chat.ID, "参数错误: /unlimited add username [tokenLimit] [rateLimit]")
			_, _ = tgBog.Send(msg)
			return
		}
		username := args[2]
		if cmd == "add" {
			user, err := model.GetUnlimitedByUsername(username)
			if err != nil {
				msg := tgbotapi.NewMessage(message.Chat.ID, "数据库错误")
				_, _ = tgBog.Send(msg)
				return
			}
			if user != nil {
				msg := tgbotapi.NewMessage(message.Chat.ID, "用户已存在")
				_, _ = tgBog.Send(msg)
				return
			}
			tokenLimit := 10000
			rateLimit := 500
			if len(args) > 3 {
				tokenLimit, _ = strconv.Atoi(args[3])
				rateLimit, _ = strconv.Atoi(args[4])
			}
			user = &model.Unlimited{
				Username:   username,
				UserID:     0,
				Key:        "sk-" + util.RandomString(48),
				TokenLimit: tokenLimit,
				RateLimit:  rateLimit,
			}
			err = user.Insert()
			if err != nil {
				msg := tgbotapi.NewMessage(message.Chat.ID, "数据库错误")
				_, _ = tgBog.Send(msg)
				return
			}
			msg := tgbotapi.NewMessage(message.Chat.ID, "添加成功")
			_, _ = tgBog.Send(msg)
		}
	}
}

func Sign(user *model.User, message tgbotapi.Message, tgBog tgbotapi.BotAPI) {
	result, err := model.Sign(user)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "签到失败，数据库错误")
		_, _ = tgBog.Send(msg)
	}
	if result {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("签到成功，获得%d积分", conf.Conf.SignGiftBalance))
		_, _ = tgBog.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(message.Chat.ID, "今日已签到")
		_, _ = tgBog.Send(msg)
	}
}

func ShowStatus(message tgbotapi.Message, tgBog tgbotapi.BotAPI) {
	msgS := "请选择您要查看的站点\n"
	msg := tgbotapi.NewMessage(message.Chat.ID, msgS)
	msg.ReplyMarkup = GetStatusKeyboard()
	_, _ = tgBog.Send(msg)
}

func Status(message tgbotapi.Message, tgBog tgbotapi.BotAPI, baseUrl string) {
	msgS := baseUrl + " 运行状态\n"
	//msg := tgbotapi.NewMessage(message.Chat.ID, msgS)
	//processMsg, err := tgBog.Send(msg)
	//if err != nil {
	//	log.Println(err)
	//}
	if conf.Conf.Sites[baseUrl] != "" {
		accessToken := conf.Conf.Sites[baseUrl]
		now := time.Now()
		currentUnixTimestamp := now.Unix()
		startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		// 获取今天0点的Unix时间戳（秒）
		startOfDayUnixTimestamp := startOfDay.Unix()
		//uri := fmt.Sprintf("https://nekoapi.com/api/log/stat?&type=0&token_name=&model_name=&start_timestamp=%d&end_timestamp=%d", startOfDayUnixTimestamp, currentUnixTimestamp)
		// resp body {"data":{"quota":0},"message":"","success":true}

		respBody, err := util.DoRequest(baseUrl, fmt.Sprintf("/api/log/stat?&type=0&token_name=&model_name=&start_timestamp=%d&end_timestamp=%d", startOfDayUnixTimestamp, currentUnixTimestamp), accessToken)
		if err != nil {
			log.Println(err)
			msgS += "今日消耗：获取失败\n"
		} else {
			if respBody["success"].(bool) {
				log.Println(respBody)
				msgS += fmt.Sprintf("今日消耗：＄%f\n", respBody["data"].(map[string]interface{})["quota"].(float64)/500000)
			}
			msg := tgbotapi.NewEditMessageText(message.Chat.ID, message.MessageID, msgS)
			_, _ = tgBog.Send(msg)
		}
		respBody, err = util.DoRequest(baseUrl, fmt.Sprintf("/api/log/stat?&type=0&token_name=&model_name=&start_timestamp=%d&end_timestamp=%d", currentUnixTimestamp-60, currentUnixTimestamp), accessToken)
		if err != nil {
			log.Println(err)
			msgS += "今日RPM：获取失败\n"
		} else {
			if respBody["success"].(bool) {
				log.Println(respBody)
				msgS += fmt.Sprintf("上一分钟RPM：%d\n", int(respBody["data"].(map[string]interface{})["rpm"].(float64)))
				msgS += fmt.Sprintf("上一分钟TPM：%d\n", int(respBody["data"].(map[string]interface{})["tpm"].(float64)))
			}
			msg := tgbotapi.NewEditMessageText(message.Chat.ID, message.MessageID, msgS)
			_, _ = tgBog.Send(msg)
		}
	}

}

func UserInfo(user *model.User, message tgbotapi.Message, tgBog tgbotapi.BotAPI) {
	msgS := fmt.Sprintf("您的积分：%d\n", user.Balance)
	if user.AccessToken != "" {
		now := time.Now()
		currentUnixTimestamp := now.Unix()
		startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		// 获取今天0点的Unix时间戳（秒）
		startOfDayUnixTimestamp := startOfDay.Unix()
		uri := fmt.Sprintf("https://nekoapi.com/api/log/self/stat?&type=0&token_name=&model_name=&start_timestamp=%d&end_timestamp=%d", startOfDayUnixTimestamp, currentUnixTimestamp)
		client := &http.Client{}
		// 创建请求
		req, err := http.NewRequest("GET", uri, nil)
		if err != nil {
			fmt.Println("NewRequest Error:", err)
			return
		}

		// 添加请求头
		req.Header.Add("Authorization", "Bearer "+conf.Conf.Sites["https://nekoapi.com"])

		log.Println(req.Header)

		// 发送请求
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Do Error:", err)
			return
		}

		// 读取响应体
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("ReadAll Error:", err)
			return
		}
		// resp body {"data":{"quota":0},"message":"","success":true}
		respBody := make(map[string]interface{})
		err = json.Unmarshal(body, &respBody)
		if err != nil {
			fmt.Println("Unmarshal Error:", err)
			return
		}
		if respBody["success"].(bool) {
			msgS += fmt.Sprintf("今日消耗额度：＄%f\n", respBody["data"].(map[string]interface{})["quota"].(float64)/500000)
		}
	} else {
		msgS += "您还未绑定NekoAPI账号，输入/bind查看绑定教程\n"
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, msgS)
	_, err := tgBog.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func Change(user *model.User, message tgbotapi.Message, tgBog tgbotapi.BotAPI, index int, taskId string, action string) {
	if user.Balance < conf.Conf.ImaginePrice {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("您的积分不足，当前积分：%d，绘图所需积分：%d", user.Balance, conf.Conf.ImaginePrice))
		_, err := tgBog.Send(msg)
		if err != nil {
			log.Println(err)
		}
		return
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, "请稍等，正在发送绘图请求")
	processMsg, err := tgBog.Send(msg)
	if err != nil {
		log.Println(err)
		return
	}
	go func() {
		err = model.DecreaseBalance(user, conf.Conf.ImaginePrice)
		if err != nil {
			log.Println(err)
			msg := tgbotapi.NewEditMessageText(message.Chat.ID, processMsg.MessageID, "绘图失败，原因：积分扣除失败")
			_, _ = tgBog.Send(msg)
			return
		}
		result := api.Change(taskId, action, index)
		if result.Code == 1 || result.Code == 22 {
			msg := tgbotapi.NewEditMessageText(message.Chat.ID, processMsg.MessageID, fmt.Sprintf("花费%d积分正在绘图中，当前进度： 排队中", conf.Conf.ImaginePrice))
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
				msg := tgbotapi.NewEditMessageText(message.Chat.ID, processMsg.MessageID, "绘图失败（不消耗积分），原因："+midjourneyResponse.FailReason)
				_, _ = tgBog.Send(msg)
				_ = model.IncreaseBalance(user, conf.Conf.ImaginePrice)
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
					msg = tgbotapi.NewEditMessageText(message.Chat.ID, processMsg.MessageID, fmt.Sprintf("绘图完成，花费%d积分", conf.Conf.ImaginePrice))
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

func Imagine(user *model.User, message tgbotapi.Message, tgBog tgbotapi.BotAPI) {

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
	if user.Balance < conf.Conf.ImaginePrice {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("您的积分不足，当前积分：%d，绘图所需积分：%d", user.Balance, conf.Conf.ImaginePrice))
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
		err = model.DecreaseBalance(user, conf.Conf.ImaginePrice)
		if err != nil {
			log.Println(err)
			msg := tgbotapi.NewEditMessageText(message.Chat.ID, processMsg.MessageID, "绘图失败，原因：积分扣除失败")
			_, _ = tgBog.Send(msg)
			return
		}
		result := api.Imagine("", prompt)
		if result.Code == 1 || result.Code == 22 {
			msg := tgbotapi.NewEditMessageText(message.Chat.ID, processMsg.MessageID, fmt.Sprintf("花费%d积分，正在绘图中，当前进度： 排队中", conf.Conf.ImaginePrice))
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
				msg := tgbotapi.NewEditMessageText(message.Chat.ID, processMsg.MessageID, "绘图失败（不消耗积分），原因："+midjourneyResponse.FailReason)
				_ = model.IncreaseBalance(user, conf.Conf.ImaginePrice)
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
					msg = tgbotapi.NewEditMessageText(message.Chat.ID, processMsg.MessageID, fmt.Sprintf("绘图完成，花费%d积分", conf.Conf.ImaginePrice))
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
	msg := tgbotapi.NewMessage(message.Chat.ID, "输入绘图内容\n例如: /imagine 可爱猫猫\n二次元风格可以加上参数\n例如: /imagine 猫娘 --niji 5\n每次绘图或者变换都会消耗10积分，积分可以通过每日签到获取")
	msg.ReplyMarkup = GetMainKeyboard(message.Chat.ID)
	_, err := tgBog.Send(msg)
	if err != nil {
		log.Println(err)
	}
}
