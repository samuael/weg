package service

import (
	"time"

	"github.com/gorilla/websocket"
	session "github.com/samuael/Project/Weg/internal/pkg/Session"
	"github.com/samuael/Project/Weg/internal/pkg/entity"
)

// Client struct
type Client struct {
	User           *entity.User
	Conn           *websocket.Conn
	ID             string
	ClientService  *ClientService
	SessionHandler *session.Cookiehandler
	Message        chan entity.XMessage
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// ReadMessage function handling the Reading of message from the
// end user client
func (client *Client) ReadMessage() {

}

// WriteMessage function handling the Writing of message to the
// end user client
func (client *Client) WriteMessage() {

}
