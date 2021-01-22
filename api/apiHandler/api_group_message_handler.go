package apiHandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/samuael/Project/Weg/internal/pkg/Group"
	"github.com/samuael/Project/Weg/internal/pkg/Message"
	session "github.com/samuael/Project/Weg/internal/pkg/Session"
	"github.com/samuael/Project/Weg/internal/pkg/User"
	"github.com/samuael/Project/Weg/internal/pkg/entity"
	"github.com/samuael/Project/Weg/pkg/Helper"
	"github.com/samuael/Project/Weg/pkg/translation"
)

// GroupMessageHandler struct representing the Group And It's
type GroupMessageHandler struct {
	GroupSer Group.GroupService
	UserSer User.UserService 
	Session *session.Cookiehandler
	MessageSer Message.MessageService
}

// NewGroupMessageHandler function returning group handler instance pointer 
func NewGroupMessageHandler(  gser Group.GroupService  , userser User.UserService  , session *session.Cookiehandler , messer Message.MessageService ) *GroupMessageHandler {
	return &GroupMessageHandler{
		GroupSer: gser,
		UserSer: userser,
		Session: session,
		MessageSer: messer,
	}
}
// GetSessionHandler returning session handler to implement DHandler interface
func (gmhandler *GroupMessageHandler)  GetSessionHandler() *session.Cookiehandler {
	return gmhandler.Session
}


// SendGroupMessage method later to be implemented by websocket 
// Method  : POST 
// INPUT  : JSON 
// OUTPUT   : JSON 
// Authorization : being a member of that group and having a valid user session 
/*
// INPUT 
{
	"status" : 5 ,
	"body" : {
		"group_id":"",
		"text" :
	}
}
*/
func (gmhandler *GroupMessageHandler ) SendGroupMessage(response http.ResponseWriter  , request *http.Request ){
	response.Header().Set( "Content-Type"  , "application/json" )
	session := gmhandler.Session.GetSession(request)
	in := &struct{
		Status int `json:"status"`
		Body  *entity.GroupMessage `json:"body"`
	}{}
	jdec   := json.NewDecoder(request.Body )
	er := jdec.Decode(in)
	if er != nil || in.Status != entity.MsgGroupTxt  || in.Body.Text=="" || in.Body.GroupID=="" {
		fmt.Println(" Input Error Nigga ")
		return 
	}
	groupExist := gmhandler.GroupSer.DoesGroupExist( in.Body.GroupID )
	if !groupExist{
		fmt.Printf("there is Group By This ID %s \n" ,  in.Body.GroupID )
		return 
	}
	isGroupMember := gmhandler.UserSer.IsGroupMember(session.UserID  , in.Body.GroupID )

	if !isGroupMember {
		fmt.Println("The User is not a member of the Group")
		return 
	}
	group := gmhandler.GroupSer.GetGroupByID(in.Body.GroupID)
	if group == nil {
		fmt.Println("Group Doesn't Exist ")
		return 
	}
	in.Body.SenderID= session.UserID
	in.Body.Time = time.Now()
	in.Body.MessageNumber = len(group.Messages)+1 
	group.Messages= append(group.Messages, in.Body)
	group.LastMessageNumber= in.Body.MessageNumber
	 group = gmhandler.GroupSer.UpdateGroup(group)
	 if group == nil {
		 fmt.Println("  Internal Server Error ...")
		 return 
	 }
	 fmt.Println(" Message Sent Succesfuly ... ")
}

// GetGroupMessage method returning messages exchanged between group members 
// METHOD : GET 
// Variables : group_id  , offset 
// RESPONSE : JSON 
// Authorization : Beign member of the Group and Appropriately Logged In 
func (  gmhandler *GroupMessageHandler ) GetGroupMessage(response http.ResponseWriter   , request *http.Request) {
	response.Header().Set("Content-Type" , "application/json")
	lang := GetSetLang(gmhandler  , response  , request )
	session := gmhandler.Session.GetSession(request)
	res := &struct{
		Success bool `json:"success"`
		Message string `json:"message"`
		Messages []*entity.GroupMessage `json:"messages"`
		GroupID string `json:"group_id"`
		Offset int `json:"offset"`
	}{
		Success: false ,
		Message: translation.Translate(lang ,  "Invalid Input " ),
	}
	groupid := request.FormValue("group_id")
	offsetString := request.FormValue("offset")
	offset :=0
	offset , er := strconv.Atoi(offsetString)
	if groupid ==""  || er != nil  {
		response.Write(Helper.MarshalThis(res))
		return 
	}
	res.GroupID = groupid
	res.Offset = offset
	groupExists := gmhandler.GroupSer.DoesGroupExist(groupid)
	if !groupExists {
		res.Message = fmt.Sprintf(translation.Translate(lang , "Goup With ID : %s doesn't Exist") , groupid )
		response.Write(Helper.MarshalThis(res))
		return 
	}
	isMember := gmhandler.GroupSer.IsGroupMember(groupid , session.UserID )
	if !isMember {
		res.Message = fmt.Sprintf( translation.Translate(lang ," User with ID %s is Not a member in Group With ID : %s") , session.UserID , groupid )
		response.Write(Helper.MarshalThis(res))
		return 
	}
	gmessages := gmhandler.MessageSer.GetGroupMessages( groupid , offset  )
	if gmessages == nil || len(gmessages)==0 {
		res.Message = translation.Translate(lang , "No Message Record Found ")
		response.Write(Helper.MarshalThis(res))
		return 
	}
	res.Message = fmt.Sprintf(translation.Translate(lang  , "Succesfuly Found %d %s ") , len(gmessages) , Helper.SetPlural(lang , "message" , len(gmessages)))
	res.Success = true 
	res.Messages = gmessages
	response.Write(Helper.MarshalThis(res))
}