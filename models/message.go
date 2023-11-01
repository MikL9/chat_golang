package models

import (
	"time"

	"github.com/guregu/null"
)

type Message struct {
	ID        int       `gorm:"column:id;primary_key" json:"id"`
	ChatID    int       `gorm:"NOT NULL;column:chat_id;" json:"-" form:"chat_id"`
	UserID    int       `gorm:"NOT NULL;column:user_id" json:"-"`
	Message   string    `gorm:"column:message" json:"message" form:"message"`
	Created   time.Time `gorm:"autoCreateTime" json:"created"`
	DeletedAt null.Time `gorm:"column:deleted_at" json:"-"`
}

func (m Message) TableName() string {
	return "messages"
}
