package middleware

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// KeyDB key to extract database client from context
const KeyDB = "database"

func DB(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(KeyDB, db)
		ctx.Next()
	}
}
