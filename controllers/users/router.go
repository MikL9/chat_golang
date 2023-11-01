package users

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"gochat/controllers/files"
	"gochat/middleware"
	"gorm.io/gorm"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gochat/models"
)

func UserRoutes(router *gin.RouterGroup) {

	r := router.Group("/user")
	r.GET("", getAllUsers)
	r.GET("/:id", getUser)
	r.POST("", addUser)
	r.DELETE("/:id", deleteUser)
	r.POST("/avatar/", addAvatar)
	r.POST("/theme/", changeTheme)
	r.POST("/themeColor/", changeThemeColor)
}

func getAllUsers(ctx *gin.Context) {
	var (
		users  = make(ResponseUsers, 0)
		search = ctx.Query("quicksearch")
		offset = ctx.Query("start")
		limit  = ctx.Query("limit")
		db     = ctx.Value(middleware.KeyDB).(*gorm.DB)
	)
	uQuery := db.
		Select("u.*, f.guid, f.name as file_name").
		Table("users u").
		Joins("LEFT JOIN files f ON u.avatar = f.id")
	if search != "" {
		uQuery = uQuery.Where("presentation LIKE '%" + search + "%' " +
			"OR login LIKE '%" + search + "%' " +
			"OR phone LIKE '%" + search + "%' " +
			"OR email LIKE '%" + search + "%' ")
	}
	offsetVal, _ := strconv.Atoi(offset)
	limitVal, _ := strconv.Atoi(limit)
	err := uQuery.Find(&users).Error
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.Writer.Header().Set("X-Total-Count", strconv.Itoa(len(users)))

	err = uQuery.Limit(limitVal).Offset(offsetVal).Find(&users).Error
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, &users)
}

func getUser(ctx *gin.Context) {
	var (
		param   = ctx.Param("id")
		id, err = strconv.Atoi(param)
		db      = ctx.Value(middleware.KeyDB).(*gorm.DB)
	)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	var user = models.User{ID: id}
	if db.First(&user).Error != nil {
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func addUser(ctx *gin.Context) {
	var (
		user = &models.User{}
		err  = ctx.BindJSON(&user)
		db   = ctx.Value(middleware.KeyDB).(*gorm.DB)
	)
	var existedUser = &models.User{}

	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	existsQuery := `
		SELECT EXISTS (
			SELECT *
			FROM users u
			WHERE u.login = ?
		`
	if user.ID != 0 {
		existsQuery = existsQuery + " AND u.id != " + strconv.Itoa(user.ID)
		if user.Password == "" {
			if db.First(&existedUser).Error != nil {
				ctx.JSON(http.StatusNotFound, err)
				return
			}
			user.Password = existedUser.Password
		} else {
			user.Password.Encrypt()
		}
	}
	existsQuery = existsQuery + `) as exist`
	var exists bool
	if err = db.Raw(existsQuery, user.Login).Row().Scan(&exists); err != nil {
		ctx.JSON(http.StatusConflict, err)
		return
	}
	if exists {
		ctx.JSON(http.StatusConflict, gin.H{"message": "Login already exist"})
		return
	}

	if user.ID != 0 {
		if db.Updates(&user).Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Can't save user"})
			return
		}
	} else {
		user.Password.Encrypt()
		if db.Save(user).Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Can't save user"})
			return
		}
	}
	ctx.JSON(http.StatusOK, user)
}

func deleteUser(ctx *gin.Context) {
	var (
		userId, _ = strconv.Atoi(ctx.Param("id"))
		db        = ctx.Value(middleware.KeyDB).(*gorm.DB)
		user      = models.User{ID: userId}
	)

	db.Delete(&user)
}

func DeleteToken(ctx context.Context, userId int, token string) bool {
	var rdb = ctx.Value(middleware.KeyRedis).(*redis.Client)
	_, e := rdb.SRem(ctx, strconv.Itoa(userId), token).Result()
	if e != nil {
		return false
	}
	return true
}

func DeleteAllTokens(ctx context.Context, userId int) bool {
	var rdb = ctx.Value(middleware.KeyRedis).(*redis.Client)
	_, e := rdb.Del(ctx, strconv.Itoa(userId)).Result()
	if e != nil {
		return false
	}
	return true
}

func addAvatar(ctx *gin.Context) {
	var (
		user = ctx.Value(middleware.KeyAuthUser).(*models.User)
		db   = ctx.Value(middleware.KeyDB).(*gorm.DB)
	)

	var avatar struct {
		files.ResponseFile
	}
	if err := ctx.ShouldBind(&avatar); err != nil {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("file read error: %s", err.Error()))
		return
	}
	db.Model(&user).Update("avatar", avatar.FileID)

	ctx.JSON(http.StatusOK, user)
}

func changeTheme(ctx *gin.Context) {
	type getUser struct {
		ID    int `gorm:"column:id;primary_key" json:"id"`
		Theme int `gorm:"column:theme" json:"theme"`
	}
	var (
		user = &getUser{}
		db   = ctx.Value(middleware.KeyDB).(*gorm.DB)
	)
	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("file read error: %s", err.Error()))
		return
	}
	db.Model(&models.User{ID: user.ID}).Update("theme", user.Theme)
	ctx.JSON(http.StatusOK, user.Theme)
	return
}

func changeThemeColor(ctx *gin.Context) {
	type getUser struct {
		ID    int    `gorm:"column:id;primary_key" json:"user_id"`
		Theme string `gorm:"column:theme" json:"theme_color"`
	}
	var (
		user = &getUser{}
		db   = ctx.Value(middleware.KeyDB).(*gorm.DB)
	)
	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("file read error: %s", err.Error()))
		return
	}
	db.Model(&models.User{ID: user.ID}).Update("theme_color", user.Theme)
	ctx.JSON(http.StatusOK, user.Theme)
	return
}
