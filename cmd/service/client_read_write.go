package service

import (
	"encoding/json"
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
	maxMessageSize = 99999999999
)

// ReadMessage function handling the Reading of message from the
// end user client
func (client *Client) ReadMessage() {
	client.Conn.SetReadLimit(maxMessageSize)
	client.Conn.SetReadDeadline(time.Now().Add(pongWait))
	client.Conn.SetPongHandler(func(string) error { client.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		message := new(entity.InMess)
		err := client.Conn.ReadJSON(message)
		if err != nil {
			// if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			// 	log.Printf("error: %v", err)
			// }
			break
		}
		// Brodcast the message
		if message.GetStatus() <= 0 || message.GetStatus() > 13 {
			continue
		}
		// the Message has passed
		// the first stage of chacking whether
		// they are a valid message or not
		// now i am gonna map the message to their specific Object
		var body entity.XMessage
		switch message.GetStatus() {

		case entity.MsgSeen:
			{
				body = &entity.SeenMessage{
					Status: message.GetStatus(),
					Body: func() entity.SeenBody {
						val := entity.SeenBody{}
						_, er := json.Marshal(message.GetBody())
						if er != nil {
							return val
						}
						return val
					}(),
				}
			}
		case entity.MsgTyping:
			{

			}
		case entity.MsgStopTyping:
			{

			}
		case entity.MsgIndividualTxt:
			{

			}

		case entity.MsgGroupTxt:
			{

			}
			// case entity.MsgAlieProfileChange : {

			// }
			// case entity.MsgNewAlie :{

			// }
			// case entity.MsgGroupProfileCHange :{

			// }
			// case entity.MsgGroupJoin : {

			// }
			// case entity.MsgDeleteUser :{

			// }

			// default:
			// 	{
			// 	}
			// }
		}
		client.ClientService.Message <- body
	}

}

// WriteMessage function handling the Writing of message to the
// end user client
func (client *Client) WriteMessage() {

}
