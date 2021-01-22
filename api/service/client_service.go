package service

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/samuael/Project/Weg/internal/pkg/Message"
	session "github.com/samuael/Project/Weg/internal/pkg/Session"
	"github.com/samuael/Project/Weg/internal/pkg/User"
	"github.com/samuael/Project/Weg/internal/pkg/entity"
)

// ClientService service
// continuously running service handling End To End Messages and broadcasts
//  coming from group service
// through main_service
type ClientService struct {
	Session        session.Cookiehandler
	Message        chan entity.XMessage
	UserSer        User.UserService
	MainSer        *MainService
	GroupSer       *GroupService
	MessageHandler Message.MessageService
}

// Run service loop
func (cs *ClientService) Run() {
	for {
		select {
		case message := <-cs.Message:
			{
				// for each case there will be some functions ans methods to be called
			}
		}
	}
}

// HandlerClientCreation function
func (cs *ClientService) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	go cs.HandleWS(response, request)
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 999999999999
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// HandleWS creating a client and send it to the Main Service as a form of registration
func (cs *ClientService) HandleWS(response http.ResponseWriter, request *http.Request) {

	conn, er := upgrader.Upgrade(response, request, nil)
	if er != nil {
		return
	}
	session := cs.Session.GetSession(request)
	if session == nil {
		return
	}

	user := cs.UserSer.GetUserByID(session.UserID)
	if user == nil {
		return
	}

	client := &Client{
		Conn:      conn,
		request:   request,
		User:      user,
		Session:   session,
		ID:        session.UserID,
		MainSer:   cs.MainSer,
		ClientSer: cs,
		Message:   make(chan []byte),
		GroupSer:  cs.GroupSer,
		// Running  : make(chan bool ) ,
	}
	cs.MainSer.Register <- client
	go client.WriteMessage()
	go client.ReadMessage()

}
