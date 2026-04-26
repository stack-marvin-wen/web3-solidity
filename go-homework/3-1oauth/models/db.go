package models

import (
	"fmt"
	"oauth/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	dsn := config.GetDBDSN()
	DB, err = gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(fmt.Sprintf("failed to connect database:%v", err))
	}
}

func AutoMigrate() {
	DB.AutoMigrate(&User{})
}
