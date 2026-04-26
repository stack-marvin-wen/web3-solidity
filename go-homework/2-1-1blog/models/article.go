package models

import (
	"time"

	"gorm.io/gorm"
)

type Article struct {
	ID        uint           `gorm:"primaryKey"`
	Title     string         `gorm:"not null"`
	Content   string         `gorm:"type:text"`
	UserID    uint           `gorm:"index"`
	Author    User           `gorm:"foreignKey:UserID"`
	Tags      []Tag          `gorm:"many2many:article_tags;"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeleteAt  gorm.DeletedAt `gorm:"index"`
}
