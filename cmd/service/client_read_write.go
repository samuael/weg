package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/samuael/Project/Weg/internal/pkg/Alie"
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
	Message        chan entity.EEMBinary
	Request        *http.Request
	MainService    *MainService
	AlieSer        Alie.AlieService
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
	defer func() {
		client.Conn.Close()
		close(client.Message)
	}()

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
						theBytes, er := json.Marshal(message.GetBody())
						if er != nil {
							return val
						}
						decoder := json.NewDecoder(bytes.NewReader(theBytes))
						decoder.Decode(&val)
						return val
					}(),
					SenderID: client.User.ID,
				}
			}
		case entity.MsgTyping:
			{
				body = &entity.TypingMessage{
					Status: message.GetStatus(),
					Body: func() entity.TypingBody {
						val := entity.TypingBody{}
						theBytes, er := json.Marshal(message.GetBody())
						if er != nil {
							return val
						}
						decoder := json.NewDecoder(bytes.NewReader(theBytes))
						decoder.Decode(&val)
						return val
					}(),
					SenderID: client.User.ID,
				}
			}
		case entity.MsgStopTyping:
			{
				body = &entity.TypingMessage{
					Status: message.GetStatus(),
					Body: func() entity.TypingBody {
						val := entity.TypingBody{}
						theBytes, er := json.Marshal(message.GetBody())
						if er != nil {
							return val
						}
						decoder := json.NewDecoder(bytes.NewReader(theBytes))
						decoder.Decode(&val)
						return val
					}(),
					SenderID: client.User.ID,
				}
			}
		case entity.MsgIndividualTxt:
			{
				body = &entity.EEMessage{
					Status: message.GetStatus(),
					Body: func() entity.Message {
						val := entity.Message{}
						theBytes, er := json.Marshal(message.GetBody())
						if er != nil {
							return val
						}
						decoder := json.NewDecoder(bytes.NewReader(theBytes))
						decoder.Decode(&val)
						return val
					}(),
					SenderID: client.User.ID,
				}
			}
		case entity.MsgGroupTxt:
			{
				body = &entity.GMMessage{
					Status: message.GetStatus(),
					Body: func() entity.GroupMessage {
						val := entity.GroupMessage{}
						theBytes, er := json.Marshal(message.GetBody())
						if er != nil {
							return val
						}
						decoder := json.NewDecoder(bytes.NewReader(theBytes))
						decoder.Decode(&val)
						return val
					}(),
					SenderID: client.User.ID,
				}
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
		if body == nil {
			continue
		}
		switch body.GetStatus() {
		case entity.MsgSeen:
			{
				if body.(*entity.SeenMessage).Body.SenderID == "" {
					break
				}
				body.(*entity.SeenMessage).SenderID = client.User.ID
			}
		case entity.MsgTyping, entity.MsgStopTyping:
			{
				body.(*entity.TypingMessage).Body.TyperID = client.User.ID
				if body.(*entity.TypingMessage).Body.ReceiverID == "" {
					break
				}
				body.(*entity.TypingMessage).SenderID = client.User.ID
			}
		case entity.MsgIndividualTxt:
			{
				if body.(*entity.EEMessage).Body.ReceiverID == "" ||
					body.(*entity.EEMessage).Body.Text == "" {
					break
				}
				body.(*entity.EEMessage).SenderID = client.User.ID

				// Send the Message and if the Message Sending was succesful sent it to
				// both the Users
				// Since the ClientService mau have a lot of Call i will handle the Sending or saving of the message to
				//  the database heere in the client ReadMessage service loop
				//
			}
		case entity.MsgGroupTxt:
			{

			}
		}

		client.ClientService.Message <- body
	}

}

// WriteMessage function handling the Writing of message to the
// end user client
func (client *Client) WriteMessage() {
	ticker := time.NewTicker(pongWait)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()
	for {
		select {
		case mess, ok := <-client.Message:
			{
				// increase the writing time limit
				client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
				// check whether the channel is open if not return and close the loop
				if !ok {
					client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}
				client.Conn.WriteMessage(websocket.BinaryMessage, mess.Data)
			}
		case <-ticker.C:
			{
				// the Ticker has counted nigga so i have to do some thing with it
				client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
				// checking the presence or activeness of the Connection by writing a ping message and
				// if the WriteMessage returns an error meaning the connection is closed
				// i will terminate the loop
				if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}
	}
}

// SendAlieMessage for saving the message to the database and return a fully qualified
// message object if the Operation was succesful else nil
// SendAlieMessage for sending or storing a message to the database
func (client *Client) SendAlieMessage(message *entity.EEMessage) (*entity.Alie, *entity.EEMessage) {
	session := client.SessionHandler.GetSession(client.Request)

	if session == nil {
		// res.Message= translation.Translate(lang  , " UnAuthorized User ")
		// response.Write(  Helper.MarshalThis(res) )
		fmt.Println("UnAuthorized User ....")
		return nil, nil
	}
	if (message.Status != entity.MsgIndividualTxt) || (message.Body.ReceiverID != "") || (message.Body.Text != "") {
		fmt.Println("No thing Is Sentoooo nigga ")
		return nil, nil
	}
	message.Body.SenderID = session.UserID
	message.Body.Time = time.Now()
	message.Body.Seen = false
	// checke if they are alies or not
	result := client.AlieSer.AreTheyAlies(message.Body.SenderID, message.Body.ReceiverID)
	var alies *entity.Alie
	if !result {
		fmt.Println("They Are Not Alies and I Am Gonna Create Alies Table  ")
		alies = client.AlieSer.CreateAlieDocument(message.Body.SenderID, message.Body.ReceiverID)
		if alies != nil {
			fmt.Println(" Alie ID : ", alies.ID)

		} else {
			fmt.Println(" ALies Response Is Nil ")
		}
	}
	if alies == nil {
		alies = client.AlieSer.GetAlies(message.Body.SenderID, message.Body.ReceiverID)
	}
	if alies == nil {
		fmt.Println("Internal Server Error")
		return nil, nil
	}
	message.Body.MessageNumber = len(alies.Messages) + 1
	alies.Messages = append(alies.Messages, &message.Body)
	alies.MessageNumber = len(alies.Messages)

	alies = client.AlieSer.UpdateAlies(alies)
	if alies == nil {
		fmt.Println(" Internal Server ERROR alies update Error")
		return nil, nil
	}
	return alies, message
}

// SetMessageSeen function to update the message as seen or not
func (client *Client) SetMessageSeen() {

}
