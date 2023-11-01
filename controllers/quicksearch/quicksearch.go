package quicksearch

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var DB *gorm.DB

func QuickSearchRoutes(r *gin.RouterGroup, db *gorm.DB) {
	DB = db
	r.POST("/", handleSearch)
}

func handleSearch(c *gin.Context) {
	//TODO
	/*	search := c.Query("q")
		messages := make(messages.ResponseMessages, 0)
		chats := make(chats2.ResponseChats, 0)
		users := make(users2.Users, 0)*/
}
