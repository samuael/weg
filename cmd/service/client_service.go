// service
package service

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/samuael/Project/Weg/internal/pkg/Alie"
	"github.com/samuael/Project/Weg/internal/pkg/Group"
	"github.com/samuael/Project/Weg/internal/pkg/Message"
	session "github.com/samuael/Project/Weg/internal/pkg/Session"
	"github.com/samuael/Project/Weg/internal/pkg/User"
	"github.com/samuael/Project/Weg/internal/pkg/entity"
)

// ClientService struct handling end to end messages and
// assembling them
type ClientService struct {
	MainService    *MainService
	Message        chan entity.XMessage
	UserSer        User.UserService
	GroupSer       Group.GroupService
	GroupService   GroupService
	SessionHandler *session.Cookiehandler
	MessageSer     Message.MessageService
	AlieSer        Alie.AlieService
}

// NewClientService function Returning ClientService instance
func NewClientService(
	mainService *MainService,
	messageser Message.MessageService,
	groupService GroupService,
	userser User.UserService,
	groupSer Group.GroupService,
	alieSer Alie.AlieService,
	session *session.Cookiehandler) *ClientService {
	return &ClientService{
		GroupSer:       groupSer,
		MainService:    mainService,
		Message:        make(chan entity.XMessage),
		SessionHandler: session,
		UserSer:        userser,
		GroupService:   groupService,
		MessageSer:     messageser,
		AlieSer:        alieSer,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Run method client service
func (clientservice *ClientService) Run() {

}

// ServeHTTP handler method making the ClientService class handler Interface
func (clientservice *ClientService) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	conn, er := upgrader.Upgrade(response, request, nil)
	if er != nil {
		return
	}
	session := clientservice.SessionHandler.GetSession(request)
	if session == nil {
		return
	}
	user := clientservice.UserSer.GetUserByID(session.UserID)
	if user == nil {
		return
	}
	client := &Client{
		ClientService:  clientservice,
		Conn:           conn,
		ID:             user.ID,
		Message:        make(chan entity.EEMBinary),
		SessionHandler: clientservice.SessionHandler,
		User:           user,
		Request:        request,
		MainService:    clientservice.MainService,
		MessageSer:     clientservice.MessageSer,
		AlieSer:        clientservice.AlieSer,
	}
	clientservice.MainService.Register <- client
	// Running the raed and the Write loops to read and write the messages
	//  from and to the Web Socket Server
	go client.ReadMessage()
	go client.WriteMessage()
}
