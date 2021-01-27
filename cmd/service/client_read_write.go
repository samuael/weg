package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/samuael/Project/Weg/internal/pkg/Alie"
	"github.com/samuael/Project/Weg/internal/pkg/Message"
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
	MessageSer     Message.MessageService
	ActiveUsr      chan *entity.SeenConfirmMessage
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pmessagegs to peer with this period. Must be less than pongWait.
	pmessagegPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 99999999999
)

// ReadMessage function handlmessageg the Readmessageg of message from the
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
			// if websocket.IsUnexpectedCloseError(err, websocket.CloseGomessagegAway, websocket.CloseAbnormalClosure) {
			// 	log.Printf("error: %v", err)
			// }
			break
		}
		// Brodcast the message
		if message.GetStatus() <= 0 || message.GetStatus() > 13 {
			continue
		}
		// the Message has passed
		// the first stage of chackmessageg whether
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
			// case entity.MsgGroupJomessage : {

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

					body.(*entity.TypingMessage).SenderID = client.User.ID
				}
			}
		case entity.MsgIndividualTxt:
			{
				if body.(*entity.EEMessage).Body.ReceiverID == "" ||
					body.(*entity.EEMessage).Body.Text == "" {
					break
				}
				body.(*entity.EEMessage).SenderID = client.User.ID

				// Send the Message and if the Message Sendmessageg was succesful sent it to
				// both the Users
				// Smessagece the ClientService mau have a lot of Call i will handle the Sendmessageg or savmessageg of the message to
				//  the database heere message the client ReadMessage service loop
				//
			}
		case entity.MsgGroupTxt:
			{

			}
		}

		client.ClientService.Message <- body
	}

}

// WriteMessage function handlmessageg the Writmessageg of message to the
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
				// messagecrease the writmessageg time limit
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
				// the Ticker has counted nigga so i have to do some thmessageg with it
				client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
				// checkmessageg the presence or activeness of the Connection by writmessageg a pmessageg message and
				// if the WriteMessage returns an error meanmessageg the connection is closed
				// i will termmessageate the loop
				if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}
	}
}

// SendAlieMessage for savmessageg the message to the database and return a fully qualified
// message object if the Operation was succesful else nil
// SendAlieMessage for sendmessageg or stormessageg a message to the database
func (client *Client) SendAlieMessage(message *entity.EEMessage) (*entity.Alie, *entity.EEMessage) {
	session := client.SessionHandler.GetSession(client.Request)

	if session == nil {
		// res.Message= translation.Translate(lang  , " UnAuthorized User ")
		// response.Write(  Helper.MarshalThis(res) )
		fmt.Println("UnAuthorized User ....")
		return nil, nil
	}
	if (message.Status != entity.MsgIndividualTxt) || (message.Body.ReceiverID != "") || (message.Body.Text != "") {
		fmt.Println("No thmessageg Is Sentoooo nigga ")
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
		fmt.Println("messageternal Server Error")
		return nil, nil
	}
	message.Body.MessageNumber = len(alies.Messages) + 1
	alies.Messages = append(alies.Messages, &message.Body)
	alies.MessageNumber = len(alies.Messages)

	alies = client.AlieSer.UpdateAlies(alies)
	if alies == nil {
		fmt.Println(" messageternal Server ERROR alies update Error")
		return nil, nil
	}
	return alies, message
}

// SetMessageSeen function to update the message as seen or not
func (client *Client) SetMessageSeen(message *entity.SeenMessage) *entity.SeenMessage {
	if (message.Status != entity.MsgSeen) || (message.Body.SenderID == "") {
		return nil
	}
	if (message.Body.MessageNumber) == (0) || message.Body.SenderID == "" {
		return nil
	}
	friendExist := client.ClientService.UserSer.UserWithIDExist(message.Body.SenderID)
	if !friendExist {
		fmt.Println("No Friend By this ID Exist nigga ")
		return nil
	}

	// settmessageg the message seen
	mess, success := client.MessageSer.SetMessageSeen(message.Body.SenderID, client.ID, message.Body.MessageNumber)
	if !success {
		fmt.Println(mess)
		return nil
	}
	client.MainService.SeenConfirmIfClientExistCheck <- &entity.SeenConfirmIfClientExist{
		RequesterID:   client.ID,
		WantedID:      message.Body.SenderID,
		MessageNumber: message.Body.MessageNumber,
	}
	// setting the seen_confirmed variable to true if the message
	// i will use a channel named \ chan DoesUserExist  \ through which
	// the main service will get a message and
	// it checks whether the Client with ID specifiedBy ActiveUserExistance instance
	// and return the message to the ActiveUserExistance.Client's SeenCOnfirm channel
	// and i will use the Result and update the Message as seen confirmed
	// meaning the sender knows that the receiver have seen the message
	return message
	// notifymessageg the websocket runnmessageg service to broadcast the seen message is undergomessageg and
	// if the sender person present send the seen message to him.
	// And, set the seen_confirmed = true
}

// setting the seen_confirmed variable to true if the message
// i will use a channel named \ chan DoesUserExist  \ through which
// the main service will get a message and
// it checks whether the Client with ID specifiedBy ActiveUserExistance instance
// and return the message to the ActiveUserExistance.Client's SeenCOnfirm channel
// and i will use the Result and update the Message as seen confirmed
// meaning the sender knows that the receiver have seen the message

// SetSeenMessageConfirmed to set the message as seen_confirmed == true after the seen message is sent to the sender
func (client *Client) SetSeenMessageConfirmed(MyID, SenderID string, messNo int) {

}
