package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"gochat/models"
	"gorm.io/gorm"
)

// KeyAuthUser key to extract authorized user from context
const KeyAuthUser = "auth_user"

var (
	jwtMiddleware *jwt.GinJWTMiddleware
	once          sync.Once
)

func userTokensKey(userID int) string {
	return fmt.Sprintf("user:%d:tokens", userID)
}

func JWT() *jwt.GinJWTMiddleware {
	once.Do(func() {
		var err error

		jwtMiddleware, err = jwt.New(&jwt.GinJWTMiddleware{
			Realm:      "test zone", // што это и чем оно должно быть?)
			Key:        []byte(os.Getenv("token_pass")),
			Timeout:    time.Hour * 24,
			MaxRefresh: time.Hour * 24,

			// LoginHandler step 1
			Authenticator: func(ctx *gin.Context) (interface{}, error) {
				var json struct {
					Login    string          `form:"login" binding:"required"`
					Password models.Password `form:"password" binding:"required"`
				}

				if err := ctx.Bind(&json); err != nil {
					return nil, jwt.ErrMissingLoginValues
				}

				json.Password.Encrypt()

				loginUser := &models.User{}

				db := ctx.Value(KeyDB).(*gorm.DB)

				if err = db.Model(models.User{}).
					Where("login = ? AND password = ? AND status = 1", json.Login, json.Password).
					First(loginUser).Error; err != nil {
					return nil, jwt.ErrFailedAuthentication
				}

				ctx.Set(KeyAuthUser, loginUser)

				return loginUser, nil // loginUser передаётся в PayloadFunc(data)
			},

			// LoginHandler step 2
			PayloadFunc: func(data interface{}) jwt.MapClaims {
				u, ok := data.(*models.User)
				if !ok {
					log.Printf("cannot assert data as *users.User")
					return jwt.MapClaims{}
				}

				return jwt.MapClaims{
					"id":           u.ID,
					"presentation": u.Presentation,
					"role":         u.Role,
				}
			},

			// loginHandler step 3
			LoginResponse: func(ctx *gin.Context, code int, token string, expire time.Time) {
				var (
					user = ctx.Value(KeyAuthUser).(*models.User)
					rdb  = ctx.Value(KeyRedis).(*redis.Client)
				)

				if err := rdb.SAdd(ctx, userTokensKey(user.ID), token).Err(); err != nil {
					ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"error": "Can't save token",
					})
					return
				}

				ctx.JSON(code, gin.H{
					"token":  token,
					"expire": expire.Format(time.RFC3339),
					"user":   user,
				})
			},

			// извлекает из контекста данные, по которым будет проводиться авторизация Authorizator
			IdentityHandler: func(ctx *gin.Context) interface{} {
				claims := jwt.ExtractClaims(ctx)

				// возвращаемое значение передаётся в Authorizator(data)
				return &models.User{
					ID:           int(claims["id"].(float64)),
					Presentation: claims["presentation"].(string),
					Role:         int(claims["role"].(float64)),
				}
			},
			Authorizator: func(data interface{}, ctx *gin.Context) bool {
				var (
					rdb = ctx.Value(KeyRedis).(*redis.Client)

					token = jwt.GetToken(ctx)
					user  = data.(*models.User)
				)

				tokenRegistered, err := rdb.SIsMember(ctx, userTokensKey(user.ID), token).Result()
				if err != nil {
					log.Println(err)
				}

				ctx.Set(KeyAuthUser, user)
				// пох на ошибку, в любом случае, если что-то пошло не так, вернётся false
				return tokenRegistered
			},

			Unauthorized: func(ctx *gin.Context, code int, message string) {
				ctx.JSON(code, gin.H{
					"code":    code,
					"message": message,
				})
			},

			RefreshResponse: func(ctx *gin.Context, code int, message string, time time.Time) {
				// TODO: удалить старый токен и записать новый, если надо

				// TODO: на фронте обработать 401-й статус
				// при его получении надо стучаться на рефреш токен со старым токеном
				// и повторно выполнять запрос с новым
				// и тут уже этот RefreshResponse заменить старый токен на новый
			},

			LogoutResponse: func(ctx *gin.Context, code int) {
				var (
					rdb   = ctx.Value(KeyRedis).(*redis.Client)
					user  = ctx.Value(KeyAuthUser).(*models.User)
					token = jwt.GetToken(ctx)
				)

				rdb.SRem(ctx, userTokensKey(user.ID), token)
				ctx.Status(code)
			},

			TokenLookup:   "header: Authorization, query: token, cookie: jwt",
			TokenHeadName: "Bearer",
		})

		if err != nil {
			log.Fatal(err)
		}
	})

	return jwtMiddleware
}
