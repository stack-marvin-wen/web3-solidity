package models

type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" `
	Username string `gorm:"unique;not null"`
	UserPWD  string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
}
