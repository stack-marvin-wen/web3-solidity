package models

import (
	"blog/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	dsn := config.GetDBDSN()
	DB, err = gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic("failed to connect database")
	}
}

func AutoMigrate() {
	DB.AutoMigrate(&User{}, &Article{}, &Tag{})
}
