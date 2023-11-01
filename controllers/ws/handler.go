package ws

import (
	"gochat/middleware"
	"gochat/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1 << 10,
	WriteBufferSize: 1 << 10,
	Subprotocols:    []string{"json"},
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Handler() gin.HandlerFunc {
	h := NewHub()
	go h.Run()

	return func(ctx *gin.Context) {
		conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			ctx.AbortWithError(500, err)
			return
		}

		user := ctx.Value(middleware.KeyAuthUser).(*models.User)

		cl := NewClient(h, conn, user)

		h.Register(cl)

		go cl.WritePump(ctx)
		go cl.ReadPump(ctx)
	}
}
