package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"gochat/models"
	"log"
)

// Client промежуточное звено между хабом и вебсокет соединением
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan Message
	user *models.User
}

type Clients map[*Client]struct{}

func NewClient(hub *Hub, conn *websocket.Conn, user *models.User) *Client {
	return &Client{
		hub:  hub,
		conn: conn,
		send: make(chan Message, 10),
		user: user,
	}
}

func (c *Client) ReadPump(ctx *gin.Context) {
	defer func() {
		c.hub.Unregister(c)
		c.conn.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			return

		default:
			var msg Message

			if err := c.conn.ReadJSON(&msg); err != nil {
				log.Println(err)
				return
			}

			if err := c.handleWsMessage(ctx, &msg); err != nil {
				log.Println(err)
				break
			}

			c.hub.Broadcast(msg)
		}
	}
}

func (c *Client) WritePump(ctx *gin.Context) {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			return

		case msg, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteJSON(&msg); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func (c *Client) handleWsMessage(ctx *gin.Context, msg *Message) (err error) {
	switch msg.Kind {
	case CHAT_MESSAGE:
		var payload ChatMessage

		err = mapstructure.Decode(msg.Payload, &payload)
		if err != nil {
			return
		}

		payload.SetUser(c.user)

		switch msg.ActionType {
		case Send:
			err = payload.Save(ctx)
			if err != nil {
				return
			}
		case Delete:
			err = payload.Delete(ctx)
			if err != nil {
				return
			}
		case Edit:
			err = payload.Edit(ctx)
			if err != nil {
				return
			}
		}
		msg.Payload = payload

		msg.SetStrategy(NewSendToChatIdStrategy(ctx, payload.Chat.ID))
	case CHAT_CREATED:
		var payload ChatCreated

		err = mapstructure.Decode(msg.Payload, &payload)
		if err != nil {
			return
		}

		msg.SetStrategy(NewSendToUsersStrategy(ctx, payload.Users))
	default:
		return ErrUnsupportedMessageKind
	}

	return
}
