package service

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	// session "github.com/samuael/Project/Weg/internal/pkg/Session"
	"github.com/samuael/Project/Weg/internal/pkg/entity"
)

// Client representing user and user Service related Datas with
//  Continuously running services references
type Client struct {
	Conn *websocket.Conn
	User *entity.User
	Session *entity.Session
	ID string 
	MainSer *MainService
	Message chan []byte 
	ClientSer *ClientService 
	GroupSer *GroupService 
	request *http.Request
	// Running chan bool 
}

// WriteMessage function loop
func (cl *Client )  WriteMessage() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		cl.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-cl.Message:
			cl.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				cl.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := cl.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(cl.Message)
			for i := 0; i < n; i++ {
				w.Write([]byte("\n"))
				w.Write(<-cl.Message)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			cl.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := cl.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ReadMessage loop 
func (cl *Client ) ReadMessage(){
	defer func() {
		cl.MainSer.UnRegister <- cl
		cl.Conn.Close()
	}()
	cl.Conn.SetReadLimit(maxMessageSize)
	cl.Conn.SetReadDeadline(time.Now().Add(pongWait))
	cl.Conn.SetPongHandler(func(string) error { cl.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := cl.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		// fmt.Println(message)

		// After this the Implementation of the Client assembling the Message 
		// and resending the message to the client service or the Group Service will take place here 
		// message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		// cl.MainSer. <- message

		// newMessage :=
		decoder := bytes.NewReader(message)
		jdec := json.NewDecoder(decoder)
		in := &entity.InMess{}
		er := jdec.Decode(in)
		if er != nil {
			continue
		}
		var val entity.XMessage
		switch in.Status {
		case entity.MsgStopTyping , entity.MsgTyping  :  {
			val = &entity.TypingMessage{
				Status :in.Status ,
				Body: in.Body.(entity.TypingBody), 
			}
		}
		case entity.MsgGroupJoin , entity.MsgGroupLeave : {
			val = &entity.JoinLeaveMessage{
				Status: in.Status,
				Body: in.Body.(entity.JoinLeaveBody),
			}
		}
		case entity.MsgSeen : {
			val = &entity.SeenMessage{
				Status: in.Status,
				Body : in.Body.( entity.SeenBody) , 
			}
		}
		case entity.MsgIndividualTxt : {
			val = &entity.EEMessage{
				Status: in.Status,
				Body : in.Body.(entity.Message) , 
			}
		}
		case entity.MsgGroupTxt : {
			val = &entity.GMMessage{
				Status: in.Status,
				Body : in.Body.(entity.GroupMessage) , 
			}
		}
		case entity.MsgAlieProfileChange ,  entity.MsgNewAlie : {
			val = &entity.AlieProfile{
				Status : in.Status , 
				Body : in.Body.(entity.User)  , 
			}
		}
		case entity.MsgGroupProfileChange : {
			val = &entity.GroupProfile{
				Status: in.Status,
				Body : in.Body.( entity.Group ) , 
			}
		}
		default :{
			val = nil 
		}
		}
		
	if val == nil {
		switch val.(entity.XMessage).GetStatus() {
			case entity.MsgGroupTxt , entity.MsgGroupJoin , entity.MsgGroupLeave  , entity.MsgGroupProfileChange : {
				cl.ClientSer.Message <- val
			}
			default : {
				cl.GroupSer.Message <- val 
			}
		}
	}

	}
}
