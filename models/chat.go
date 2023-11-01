package models

import "github.com/guregu/null"

type Chat struct {
	ID        int       `gorm:"column:id;primary_key" json:"id"`
	Name      string    `gorm:"NOT NULL;column:name;" json:"name"`
	IsGroup   bool      `gorm:"NOT NULL;column:is_group" json:"is_group"`
	Logo      int       `gorm:"column:logo" json:"logo"`
	DeletedAt null.Time `gorm:"column:deleted_at" json:"-"`
}

func (c Chat) TableName() string {
	return "chats"
}