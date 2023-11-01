package chats

import (
	models2 "gochat/controllers/ws"
	"gochat/models"
)

type (
	ChatUser struct {
		ID           int    `json:"id"`
		Presentation string `json:"presentation"`
	}

	ResponseChat struct {
		models.Chat
		Users []*ChatUser `gorm:"-" json:"users"`
	}

	RequestChat struct {
		models.Chat
		Users []int `json:"users"`
	}

	ResponseChats []*ResponseChat

	ResponseMessages []*struct {
		models.Message
		UserName  string `json:"user_name"`
		FileName  string `json:"file_name"`
		FileID    int    `json:"file_id"`
		Extension string `json:"extension"`
		Type      string `json:"type"`
		Guid      string `json:"guid"`
	}
)

func (m ResponseMessages) ToChatMessage() []*models2.ChatMessage {
	res := make([]*models2.ChatMessage, len(m))

	for i, rm := range m {
		res[i] = &models2.ChatMessage{
			ID:   rm.ID,
			Text: rm.Message.Message,
			User: models2.User{
				ID:           rm.UserID,
				Presentation: rm.UserName,
			},
			Chat: models2.Chat{
				ID: rm.ChatID,
			},
			File: models2.File{
				ID:        rm.FileID,
				Name:      rm.FileName,
				Extension: rm.Extension,
				Type:      rm.Type,
				Guid:      rm.Guid,
			},
			Created: rm.Created,
		}
		//if len([]rune(rm.FileName)) > 25 {
		//	res[i].File.Name = rm.FileName[:22] + "..." + rm.Extension
		//}
	}

	return res
}
