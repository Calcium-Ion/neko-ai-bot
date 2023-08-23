package model

import (
	"log"
	"neko-ai-bot/conf"
	"time"
)

type User struct {
	Id          int    `json:"id"`
	Username    string `json:"username"`
	UserId      int64  `json:"user_id"`
	ChatId      int64  `json:"chat_id"`
	Balance     int    `json:"balance"`
	SignDate    string `json:"sign_date"`
	Key         string `json:"key"`
	AccessToken string `json:"access_token"`
}

var users = make(map[int64]*User)

func Sign(user *User) (bool, error) {
	if user.SignDate == time.Now().Format("2006-01-02") {
		return false, nil
	} else {
		user.SignDate = time.Now().Format("2006-01-02")
		user.Balance += conf.Conf.SignGiftBalance
		err := DB.Save(user).Error
		if err == nil {
			users[user.ChatId] = user
			return true, nil
		} else {
			return false, err
		}
	}
}

func GetUserOrInit(chatId int64, username string, userId int64) (*User, error, bool) {
	user, _ := GetUserByChatId(chatId)
	log.Printf("GetUserOrInit: %v", user)
	if user.Id == 0 {
		// 用户不存在，创建用户
		user = &User{
			Username: username,
			UserId:   userId,
			ChatId:   chatId,
			Balance:  20,
		}
		err := DB.Create(user).Error
		log.Printf("user not exist, create new user: %v", user)

		if err != nil {
			return nil, err, false
		} else {
			// 插入成功，更新缓存
			users[chatId] = user
			return user, nil, true
		}
	}
	return user, nil, false
}

func GetUserByChatId(chatId int64) (*User, error) {
	var user *User
	var err error
	if users[chatId] != nil {
		// 从缓存中获取
		user = users[chatId]
		log.Printf("get user from cache: %v", user)
	} else {
		// 从数据库中获取
		err = DB.Where("chat_id = ?", chatId).First(&user).Error
		if err == nil && user.Id != 0 {
			users[chatId] = user
			log.Printf("get user from db: %v", user)
		}
	}
	return user, err
}

func DecreaseBalance(user *User, balance int) error {
	user.Balance -= balance
	err := DB.Save(user).Error
	if err == nil {
		users[user.ChatId] = user
	}
	return err
}

func IncreaseBalance(user *User, balance int) error {
	user.Balance += balance
	err := DB.Save(user).Error
	if err == nil {
		users[user.ChatId] = user
	}
	return err
}
