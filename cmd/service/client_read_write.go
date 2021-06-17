package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/samuael/Project/Weg/internal/pkg/Alie"
	"github.com/samuael/Project/Weg/internal/pkg/Message"
	session "github.com/samuael/Project/Weg/internal/pkg/Session"
	"github.com/samuael/Project/Weg/internal/pkg/entity"
	"github.com/samuael/Project/Weg/pkg/Helper"
)

// Client struct
type Client struct {
	User           *entity.User
	Conns          map[string]*entity.ClientConn
	ID             string
	ClientService  *ClientService
	SessionHandler *session.Cookiehandler
	Message        chan entity.EEMBinary
	Request        *http.Request
	MainService    *MainService
	AlieSer        Alie.AlieService
	MessageSer     Message.MessageService
	// this channel variable is used by main service
	// to notify the client about the presence of this client
	// specified in seenConfirmeMsg.AlieID and to modify the
	// message that is exchanged by the two clients which it's
	// message is specified in the instances
	//  SeenConfirmMsg.MessageNumber
	SeenConfirmMsg chan entity.SeenConfirmMessage
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
func (client *Client) ReadMessage(key string) {
	defer func() {
		client.ClientService.MainService.UnRegister<- client
		// println("I Am Leaving nigga")
		(client.Conns[key]).Conn.Close()
		//  close(client.Message)
		client.MainService.DeleteClientConn <- &entity.ClientConnExistance{
			IP: key,
			ID: client.ID,
		}
		// recover if error happened...
		recover()
	}()
	client.Conns[key].Conn.SetReadLimit(maxMessageSize)
	// client.Conns[key].Conn.SetReadDeadline(time.Now().Add(pongWait))
	client.Conns[key].Conn.SetPongHandler(func(string) error { /*client.Conns[key].Conn.SetReadDeadline(time.Now().Add(pongWait)); */return nil })
	for {
		// print("Im running")
		message := &entity.InMess{}
		err := client.Conns[key].Conn.ReadJSON(message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,  websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			// websocket.
			// println(err.Error())
			break
		}
		// println(string(Helper.MarshalThis(message)))
		if message == nil {
			continue
		}
		// Brodcast the message
		if message.GetStatus() <= 0 || message.GetStatus() > 13 {
			// println("Not A Valid Message ...  ")
			continue
		}

		// the Message has passed
		// the first stage of chackmessageg whether
		// they are a valid message or not
		// now i am gonna map the message to their specific Object
		var body entity.XMessage
		body = client.FilterMessage(message)
		if body == nil {
			continue
		}
		switch body.GetStatus() {
		case entity.MsgSeen:
			{

				// print(" Incomming Seen message ................." , body);
				if body.(*entity.SeenMessage).Body.SenderID == "" {
					body = nil
					break
				}
				body.(*entity.SeenMessage).SenderID = client.User.ID
				client.ClientService.Message <- body
			}
		case entity.MsgTyping, entity.MsgStopTyping:
			{
				body.(*entity.TypingMessage).Body.TyperID = client.User.ID
				if body.(*entity.TypingMessage).Body.ReceiverID == "" {
					body = nil
					break

				}
				body.(*entity.TypingMessage).SenderID = client.User.ID
				client.ClientService.Message <- body
			}
		case entity.MsgIndividualTxt:
			{
				if body.(*entity.EEMessage).Body.ReceiverID == "" ||
					body.(*entity.EEMessage).Body.Text == "" {
					body = nil
					break
				}
				body.(*entity.EEMessage).SenderID = client.User.ID
				// Send the Message and if the Message Sendmessageg was succesful sent it to
				// both the Users
				// Since the ClientService may have a lot of Call i will handle the Sending or saving of the message to
				//  the database here message the client ReadMessage service loop
				alie, messaj := client.SendAlieMessage(body.(*entity.EEMessage))
				if alie != nil {
					me := client.ClientService.UserSer.GetUserByID(func() string {
						if alie.A == client.ID {
							return client.ID
						}
						return alie.B
					}())
					if me == nil {
						break
					}
					him := client.ClientService.UserSer.GetUserByID(func() string {
						if alie.A == client.ID {
							return alie.B
						}
						return alie.A
					}())
					if him == nil {
						break
					}
					client.ClientService.Message <- &entity.NewAlie{
						Status: entity.MsgNewAlie,
						Body: entity.NewAlieBody{
							ReceiverID: client.ID,
							User:       him,
						},
						SenderID: client.ID,
					}
					client.ClientService.Message <- &entity.NewAlie{
						Status: entity.MsgNewAlie,
						Body: entity.NewAlieBody{
							ReceiverID: him.ID,
							User:       me,
						},
						SenderID: client.ID,
					}
				}
				// sending the message if the message in not nil
				if messaj != nil {
					// println(" Seending the Message to the Clients \n\n\n")
					messaj.SenderID = client.ID
					client.ClientService.Message <- messaj
				}
			}
		case entity.MsgGroupTxt:
			{

			}
		}
	}
}

// Run loop each client
// this loop will terminate only if the client is null
// the client will be null if and only if all the connected client
// machines become nil >> that is made by the main service class
// And, Here we will chack the presence of this client in each 5 second
// and terminate the loop if the client doesn't exist any more.~!
func (client *Client) Run() {
	ticker := time.NewTicker(time.Second * 3)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case seenconfm := <-client.SeenConfirmMsg:
			{
				client.SetSeenMessageConfirmed(seenconfm.ReceiverID, seenconfm.AlieID, seenconfm.MessageNumber)
			}
		case <-ticker.C:
			{
				// print("I Am Reading ... ")
				if client == nil {
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
		// fmt.Println("UnAuthorized User ....")
		return nil, nil
	}
	if (message.Status != entity.MsgIndividualTxt) || (message.Body.ReceiverID == "") || (message.Body.Text == "") {
		// fmt.Println("No thmessageg Is Sentoooo nigga ")
		return nil, nil
	}
	message.Body.SenderID = session.UserID
	message.Body.Time = time.Now()
	message.Body.Seen = false
	// checke if they are alies or not
	result := client.AlieSer.AreTheyAlies(message.Body.SenderID, message.Body.ReceiverID)
	var alies *entity.Alie
	if !result {
		// fmt.Println("They Are Not Alies and I Am Gonna Create Alies Table  ")
		alies = client.AlieSer.CreateAlieDocument(message.Body.SenderID, message.Body.ReceiverID)
		if alies != nil {
			// fmt.Println(" Alie ID : ", alies.ID)

		} else {
			// fmt.Println(" ALies Response Is Nil ")
		}
	}
	if alies == nil {
		alies = client.AlieSer.GetAlies(message.Body.SenderID, message.Body.ReceiverID)
	}
	if alies == nil {
		// fmt.Println("messageternal Server Error")
		return nil, nil
	}
	message.Body.MessageNumber = len(alies.Messages) + 1
	alies.Messages = append(alies.Messages, &message.Body)
	alies.MessageNumber = len(alies.Messages)

	alies = client.AlieSer.UpdateAlies(alies)
	if alies == nil {
		// fmt.Println(" internal Server ERROR alies update Error")
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
		// fmt.Println("No Friend By this ID Exist nigga ")
		return nil
	}

	// settmessageg the message seen
	mess, success := client.MessageSer.SetMessageSeen(message.Body.SenderID, client.ID, message.Body.MessageNumber)
	if !success {
		// fmt.Println(mess)
		return nil
	}
	// telling the main service if the sender exist in the online list
	// notify me to set that the message as seen_confirmed to true
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
	if client != nil {
		success := client.MessageSer.SetSeenConfirmed(MyID, SenderID, messNo)
		if success {
			return
		}
	}
}

// WriteMessage function handlmessageg the Writmessageg of message to the
// end user client
func (client *Client) WriteMessage(key string) {
	ticker := time.NewTicker(pongWait)
	defer func() {
		ticker.Stop()
		client.Conns[key].Conn.Close()
		client.MainService.UnRegister<-client 
		client.MainService.DeleteClientConn <- &entity.ClientConnExistance{
			IP: key,
			ID: client.ID,
		}
		recover()
	}()
	for {
		select {
		case mess, ok := <-client.Conns[key].Message:
			{
				// messagecrease the writmessageg time limit
				// client.Conns[key].Conn.SetWriteDeadline(time.Now().Add(writeWait))
				// check whether the channel is open if not return and close the loop
				if !ok {
					client.Conns[key].Conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}
				client.Conns[key].Conn.WriteMessage(websocket.TextMessage, mess )
			}
		case <-ticker.C:
			{
				// the Ticker has counted nigga so i have to do some thmessageg with it

				// client.Conns[key].Conn.SetWriteDeadline(time.Now().Add(writeWait))
				
				// checkmessageg the presence or activeness of the Connection by writmessageg a pmessageg message and
				// if the WriteMessage returns an error meanmessageg the connection is closed
				// i will termmessageate the loop
				if err := client.Conns[key].Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}
	}
}
// FilterMessage message nigga 
func (client *Client) FilterMessage(message  *entity.InMess )  entity.XMessage {
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
			// println("The Message Arrived Here")
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
	}
	return body
}
