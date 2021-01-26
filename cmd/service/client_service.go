// service
package service

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/samuael/Project/Weg/internal/pkg/Group"
	session "github.com/samuael/Project/Weg/internal/pkg/Session"
	"github.com/samuael/Project/Weg/internal/pkg/User"
	"github.com/samuael/Project/Weg/internal/pkg/entity"
)

// ClientService struct handling end to end messages and
// assembling them
type ClientService struct {
	MainService *MainService
	Message     chan entity.XMessage

	UserSer        User.UserService
	GroupSer       Group.GroupService
	GroupService   GroupService
	SessionHandler *session.Cookiehandler
}

// NewClientService function Returning ClientService instance
func NewClientService(mainService *MainService, groupService GroupService, userser User.UserService, groupSer Group.GroupService, session *session.Cookiehandler) *ClientService {
	return &ClientService{
		// GroupSer:       nil,
		MainService:    mainService,
		Message:        make(chan entity.XMessage),
		SessionHandler: session,
		UserSer:        userser,
		GroupService:   groupService,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Run method client service
func (clienservice *ClientService) Run() {

}

// ServeHTTP handler method making the ClientService class handler Interface
func (clienservice *ClientService) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	conn, er := upgrader.Upgrade(response, request, nil)
	if er != nil {
		return
	}
	session := clienservice.SessionHandler.GetSession(request)
	if session == nil {
		return
	}
	user := clienservice.UserSer.GetUserByID(session.UserID)
	if user == nil {
		return
	}
	client := &Client{
		ClientService:  clienservice,
		Conn:           conn,
		ID:             user.ID,
		Message:        make(chan entity.XMessage),
		SessionHandler: clienservice.SessionHandler,
		User:           user,
	}
	clienservice.MainService.Register <- client
	// Running the raed and the Write loops to read and write the messages
	//  from and to the Web Socket Server
	go client.ReadMessage()
	go client.WriteMessage()
}
