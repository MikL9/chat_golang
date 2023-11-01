package ws

import (
	"context"
	"gochat/middleware"
	"gochat/models"
	"gorm.io/gorm"
)

type BroadcastStrategy interface {
	FilterClients(clients Clients) (filtered Clients)
}

type SendToChatStrategy struct {
	ChatID int
	ctx    context.Context
}

func NewSendToChatIdStrategy(ctx context.Context, chatID int) *SendToChatStrategy {
	return &SendToChatStrategy{
		ChatID: chatID,
		ctx:    ctx,
	}
}

func (s *SendToChatStrategy) FilterClients(clients Clients) Clients {
	db := s.ctx.Value(middleware.KeyDB).(*gorm.DB)

	filtered := make(Clients)
	chatUsers := make(map[int]struct{})

	rows, _ := db.Raw(`
		SELECT user_id
		FROM chat.members
		WHERE chat_id = ?
	`, s.ChatID).Rows()

	var userID int
	for rows.Next() {
		_ = rows.Scan(&userID)
		chatUsers[userID] = struct{}{}
	}

	for client := range clients {
		if _, ok := chatUsers[client.user.ID]; ok {
			filtered[client] = struct{}{}
		}
	}

	return filtered
}

type SendToUsersStrategy struct {
	userIDs map[int]struct{}
}

func NewSendToUsersStrategy(ctx context.Context, users []*User) *SendToUsersStrategy {
	owner := ctx.Value(middleware.KeyAuthUser).(*models.User)

	userIDs := make(map[int]struct{})

	for _, u := range users {
		if u.ID != owner.ID {
			userIDs[u.ID] = struct{}{}
		}
	}

	return &SendToUsersStrategy{
		userIDs: userIDs,
	}
}

func (s *SendToUsersStrategy) FilterClients(clients Clients) Clients {
	filtered := make(Clients)

	for client := range clients {
		if _, ok := s.userIDs[client.user.ID]; ok {
			filtered[client] = struct{}{}
		}
	}

	return filtered
}
