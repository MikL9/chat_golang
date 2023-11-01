package ws

import (
	"context"
	"gochat/middleware"
	"gochat/models"
	"gorm.io/gorm"
	"time"
)

type ChatMessage struct {
	ID      int       `json:"id"`
	Text    string    `json:"text"`
	User    User      `json:"user"`
	Chat    Chat      `json:"chat"`
	File    File      `json:"file"`
	Created time.Time `json:"created"`
}

type User struct {
	ID           int    `json:"id"`
	Presentation string `json:"presentation"`
}

type Chat struct {
	ID int `json:"id"`
}

type File struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Extension string `json:"extension"`
	Type      string `json:"type"`
	Guid      string `json:"guid"`
}

func (m *ChatMessage) SetUser(u *models.User) {
	m.User.ID = u.ID
	m.User.Presentation = u.Presentation
}

func (m *ChatMessage) ToDbModel() *models.Message {
	return &models.Message{
		ID:      m.ID,
		ChatID:  m.Chat.ID,
		UserID:  m.User.ID,
		Message: m.Text,
	}
}

func (m *ChatMessage) Save(ctx context.Context) error {
	db := ctx.Value(middleware.KeyDB).(*gorm.DB)
	dbObj := m.ToDbModel()

	if err := db.Save(dbObj).Error; err != nil {
		return err
	}

	m.ID = dbObj.ID
	m.Created = dbObj.Created

	return nil
}

func (m *ChatMessage) Delete(ctx context.Context) error {
	db := ctx.Value(middleware.KeyDB).(*gorm.DB)
	dbObj := m.ToDbModel()

	if err := db.Delete(dbObj).Error; err != nil {
		return err
	}

	m.ID = dbObj.ID
	m.Created = dbObj.Created

	return nil
}

func (m *ChatMessage) Edit(ctx context.Context) error {
	db := ctx.Value(middleware.KeyDB).(*gorm.DB)
	dbObj := m.ToDbModel()

	if err := db.Updates(dbObj).Error; err != nil {
		return err
	}

	m.ID = dbObj.ID
	m.Created = dbObj.Created

	return nil
}
