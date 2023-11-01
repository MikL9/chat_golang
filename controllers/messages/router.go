package messages

import (
	"encoding/json"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"gochat/middleware"
	"gochat/models"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"strconv"
)

var DB *gorm.DB

func MessageRoutes(router *gin.RouterGroup, db *gorm.DB) {
	DB = db
	r := router.Group("/message")
	r.GET("/select/", getChatMessages)
	r.POST("/", addMessage)
	r.POST("/delete/", deleteMessages)
}

func getChatMessages(c *gin.Context) {
	messages := make(ResponseMessages, 0)
	claims := jwt.ExtractClaims(c)
	userId := int(claims["id"].(float64))
	chatId := c.Query("chat_id")

	err := messages.selectAll(DB, userId, chatId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	c.Writer.Header().Set("X-Total-Count", strconv.Itoa(len(messages)))
	c.JSON(http.StatusOK, messages)
}

func addMessage(c *gin.Context) {
	var message = models.Message{}
	err := c.Bind(&message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	claims := jwt.ExtractClaims(c)
	userId := int(claims["id"].(float64))
	message.UserID = userId
	err = DB.Create(&message).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, message)
}

func deleteMessages(ctx *gin.Context) {
	var (
		arr []int
		db  = ctx.Value(middleware.KeyDB).(*gorm.DB)
	)
	b, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}
	err = json.Unmarshal(b, &arr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}

	for _, messageID := range arr {
		message := models.Message{ID: messageID}
		db.Delete(&message)
	}
}
