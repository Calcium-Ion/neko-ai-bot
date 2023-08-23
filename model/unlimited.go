package model

import (
	"errors"
	"gorm.io/gorm"
)

type Unlimited struct {
	gorm.Model
	Username   string `json:"username" gorm:"type:varchar(64);unique_index"`
	UserID     int64  `json:"user_id" gorm:"type:int"`
	Key        string `json:"key" gorm:"type:varchar(64);unique_index"`
	TokenLimit int    `json:"token_limit" gorm:"type:int;default:10000"` //-1 means unlimited
	RateLimit  int    `json:"rate_limit" gorm:"type:int;default:500"`    //-1 means unlimited
}

func (unlimited *Unlimited) Insert() error {
	return DB.Create(unlimited).Error
}

func GetAllUnlimited() ([]*Unlimited, error) {
	var unlimiteds []*Unlimited
	err := DB.Find(&unlimiteds).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return unlimiteds, err
}

func GetUnlimitedByIDOrUsername(userID int, username *string) (*Unlimited, error) {
	unlimited := Unlimited{}
	username_ := ""
	if username != nil {
		username_ = *username
	}
	err := DB.Where("id = ? or username = ?", userID, username_).First(&unlimited).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &unlimited, err
}

func GetUnlimitedByUsername(username string) (*Unlimited, error) {
	unlimited := Unlimited{}
	err := DB.Where("username = ?", username).First(&unlimited).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &unlimited, err
}
