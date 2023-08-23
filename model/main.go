package model

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func Setup() {
	// Use SQLite
	db, err := gorm.Open(sqlite.Open("nekoai-bot.db"), &gorm.Config{
		PrepareStmt: true, // precompile SQL
	})
	log.Println("database connected")
	if err == nil {
		DB = db
		err := db.AutoMigrate(&User{})
		if err != nil {
			panic(err)
		}
		err = db.AutoMigrate(&Unlimited{})
		if err != nil {
			panic(err)
		}
		log.Println("database migrated")
	} else {
		log.Println(err)
	}
	//return err
	//panic(err)
}
