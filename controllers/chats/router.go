package chats

import (
	"fmt"
	"gochat/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gochat/middleware"
	"gorm.io/gorm"
)

func ConfigRoutes(r *gin.Engine, mw ...gin.HandlerFunc) {
	g := r.Group("/chats", mw...)
	{
		g.GET("", getUserChats)
		g.GET("/:id/messages", getChatMessages)
		g.GET("/:id/users", getChatUsers)
		g.POST("", addChat)
		g.POST("/:id/leave", leaveChat)
	}
}

func getUserChats(ctx *gin.Context) {
	var (
		chats = make(ResponseChats, 0)
		db    = ctx.Value(middleware.KeyDB).(*gorm.DB)
		user  = ctx.Value(middleware.KeyAuthUser).(*models.User)
	)

	if err := db.Raw(`
		SELECT c.id, c.name, c.is_group, c.logo
		FROM chats c
		WHERE c.deleted_at IS NULL
				AND EXISTS(
						SELECT 1
						FROM members m
						WHERE m.user_id = ?
								AND m.chat_id = c.id
				)
	`, user.ID).Scan(&chats).Error; err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	for key, chat := range chats {
		var users = make([]*ChatUser, 0)
		if err := db.Raw(`
			SELECT u.id, u.presentation FROM users u
			INNER JOIN members m ON m.user_id=u.id
			WHERE m.chat_id = ?
		`, chat.ID).Scan(&users).Error; err != nil {
			continue
		}
		chats[key].Users = users
	}

	ctx.JSON(http.StatusOK, &chats)
}

func getChatMessages(ctx *gin.Context) {
	var (
		messages        = make(ResponseMessages, 0)
		db              = ctx.Value(middleware.KeyDB).(*gorm.DB)
		paramID         = ctx.Param("id")
		qlimit, qoffset = ctx.Query("limit"), ctx.Query("offset")
	)

	chatID, err := strconv.Atoi(paramID)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	limit, err := strconv.Atoi(qlimit)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	offset, err := strconv.Atoi(qoffset)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	nLastQuery := db.
		Select("m.*, f.name as file_name, f.id as file_id, f.extension as extension, f.mtype as type, f.guid as guid, u.presentation as user_name").
		Table("messages m").
		Joins("LEFT JOIN users u ON m.user_id = u.id").
		Joins("LEFT JOIN files f ON m.attachment = f.id").
		Where("m.deleted_at IS NULL").
		Where("m.chat_id = ?", chatID).
		Order("id DESC, created").
		Limit(limit).
		Offset(offset)

	query := db.Select("*").
		Table("(?) n_last", nLastQuery).
		Order("id")

	if err := query.Scan(&messages).Error; err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, messages.ToChatMessage())
}

func getChatUsers(ctx *gin.Context) {

}

func addChat(ctx *gin.Context) {
	var (
		chat = &ResponseChat{}
		err  = ctx.BindJSON(&chat)
		db   = ctx.Value(middleware.KeyDB).(*gorm.DB)
		user = ctx.Value(middleware.KeyAuthUser).(*models.User)
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	if err = db.Table("chats").Create(&chat).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Can't save chat"})
		return
	}
	query := "INSERT INTO members (chat_id, user_id) VALUES "
	values := []string{fmt.Sprintf("(%d, %d)", chat.ID, user.ID)}
	for _, value := range chat.Users {
		if value.ID == user.ID {
			continue
		}
		values = append(values, fmt.Sprintf("(%d, %d)", chat.ID, value.ID))
	}

	db.Exec(query + strings.Join(values, ", "))
	ctx.JSON(http.StatusOK, chat)
}

func leaveChat(ctx *gin.Context) {
	var (
		db      = ctx.Value(middleware.KeyDB).(*gorm.DB)
		user    = ctx.Value(middleware.KeyAuthUser).(*models.User)
		paramID = ctx.Param("id")
	)
	chatID, err := strconv.Atoi(paramID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Incorrect chat ID"})
		return
	}

	if err = db.Exec(`DELETE FROM members WHERE chat_id = ? AND user_id = ?`, chatID, user.ID).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Couldn't leave chat"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": true})
}
