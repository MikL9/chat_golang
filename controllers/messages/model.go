package messages

import (
	"gochat/models"
	"gorm.io/gorm"
)

type ResponseMessage struct {
	models.Message
	UserName string `gorm:"embedded" json:"user_name"`
}

type ResponseMessages []*ResponseMessage

func baseQuery(db *gorm.DB, userId int) *gorm.DB {
	return db.
		Select(`
			messages.id,
			messages.chat_id,
			messages.user_id,
			messages.message,
			messages.created,
			users.presentation as user_name
		`).
		Table("messages").
		Joins("INNER JOIN chats ON messages.chat_id=chats.id").
		Joins("LEFT JOIN users ON users.id=messages.user_id").
		Where(gorm.Expr("EXISTS (SELECT * FROM members WHERE members.chat_id=chats.id AND members.user_id = ?)", userId))
}

func (m *ResponseMessages) selectAll(db *gorm.DB, userId int, chatId string) error {
	return baseQuery(db, userId).Where("chats.id = ?", chatId).Scan(&m).Error
}
