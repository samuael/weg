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
		Conns:          map[string]*entity.ClientConn{GetClientIPFromRequest(request): &entity.ClientConn{Conn: conn, IP: GetClientIPFromRequest(request)}},
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
}

// GetClientIPFromRequest returning the ip address of client making the
// request
func GetClientIPFromRequest(request *http.Request) string {
	forwaredeFor := request.Header.Get("X-FORWARDED-FOR")
	if forwaredeFor != "" {
		return forwaredeFor
	}
	// RemoteAddr allows HTTP servers and other software to
	// record the network address that sent the request,
	//  usually for logging. This field is not filled in by
	// ReadRequest and has no defined format.
	// The HTTP server in this package sets RemoteAddr to an
	// "IP:port" address before invoking a handler.
	//  This field is ignored by the HTTP client.
	return request.RemoteAddr
}
