package main

import (
	"flag"
	"log"

	"gochat/controllers/chats"
	"gochat/controllers/files"
	"gochat/controllers/messages"
	"gochat/controllers/users"
	"gochat/controllers/ws"
	"gochat/middleware"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var port string

func main() {
	flag.StringVar(&port, "p", "8989", "port to listen on")
	flag.Parse()

	log.SetFlags(log.Lshortfile | log.LstdFlags)

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}

	var (
		db  = initDB()
		rdb = initRedis()
	)

	r := gin.Default()
	r.Use(static.Serve("/files", static.LocalFile("./files", false)))
	r.Use(
		middleware.CORS,
		middleware.Redis(rdb),
		middleware.DB(db),
	)

	// TODO: разгрести эту мусорку
	r.POST("/login", middleware.JWT().LoginHandler)

	chats.ConfigRoutes(r, middleware.JWT().MiddlewareFunc())

	auth := r.Group("", middleware.JWT().MiddlewareFunc())
	{
		auth.GET("/refresh_token", middleware.JWT().RefreshHandler)
		auth.POST("/logout", middleware.JWT().LogoutHandler)

		users.UserRoutes(auth)
		messages.MessageRoutes(auth, db)
		files.FileRoutes(auth, db)

		auth.GET("/ws", ws.Handler())
	}

	log.Fatal(r.Run(":" + port))
}
