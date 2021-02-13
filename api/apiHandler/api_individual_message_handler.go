package apiHandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"strconv"
	"time"

	"github.com/samuael/Project/Weg/internal/pkg/Alie"
	"github.com/samuael/Project/Weg/internal/pkg/Message"
	session "github.com/samuael/Project/Weg/internal/pkg/Session"
	"github.com/samuael/Project/Weg/internal/pkg/User"
	"github.com/samuael/Project/Weg/internal/pkg/entity"
	"github.com/samuael/Project/Weg/pkg/Helper"
	"github.com/samuael/Project/Weg/pkg/translation"
)

// IndividualMessageHandler struct handling thIndividual Message
type IndividualMessageHandler struct {
	Session    *session.Cookiehandler
	MessageSer Message.MessageService
	UserSer    User.UserService
	AlieSer    Alie.AlieService
}

// NewIndvMessageHandler returns individual Message handle instance
func NewIndvMessageHandler(sess *session.Cookiehandler, inms Message.MessageService, userser User.UserService, alieSer Alie.AlieService) *IndividualMessageHandler {
	return &IndividualMessageHandler{
		Session:    sess,
		MessageSer: inms,
		UserSer:    userser,
		AlieSer:    alieSer,
	}
}

// GetSessionHandler () *session.Cookiehandler
func (inms *IndividualMessageHandler) GetSessionHandler() *session.Cookiehandler {
	return inms.Session
}

// SendAlieMessage method
// METHOD : POST
// INPUT JSON
// OUTPUT JSON
// status
func (inms *IndividualMessageHandler) SendAlieMessage(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	// lang := GetSetLang(inms , response  , request)
	session := inms.Session.GetSession(request)

	in := &struct {
		Status int             `json:"status"`
		Body   *entity.Message `json:"body"`
	}{}

	if session == nil {
		// res.Message= translation.Translate(lang  , " UnAuthorized User ")
		// response.Write(  Helper.MarshalThis(res) )
		fmt.Println("UnAuthorized User ....")
		return
	}

	jdec := json.NewDecoder(request.Body)
	er := jdec.Decode(in)
	if er != nil {
		fmt.Println(" Message Decode Error ....  ")
		return
		// res.Message=translation.Translate(lang  , "Invalid Input ")
	}
	if (in.Status != entity.MsgIndividualTxt) || (in.Body.ReceiverID != "") || (in.Body.Text != "") {
		fmt.Println("No thing Is Sentoooo nigga ")
		return
	}
	in.Body.SenderID = session.UserID
	in.Body.Time = time.Now()
	in.Body.Seen = false
	fmt.Println(in.Body.SenderID, in.Body.ReceiverID)
	// checke if they are alies or not
	result := inms.AlieSer.AreTheyAlies(in.Body.SenderID, in.Body.ReceiverID)
	var alies *entity.Alie
	if !result {
		fmt.Println("They Are Not Alies and I Am Gonna Create Alies Table  ")
		alies = inms.AlieSer.CreateAlieDocument(in.Body.SenderID, in.Body.ReceiverID)
		if alies != nil {
			fmt.Println(" Alie ID : ", alies.ID)

		} else {
			fmt.Println(" ALies Response Is Nil ")
		}
	}
	if alies == nil {
		alies = inms.AlieSer.GetAlies(in.Body.SenderID, in.Body.ReceiverID)
	}
	if alies == nil {
		fmt.Println("Internal Server Error")
		return
	}
	in.Body.MessageNumber = len(alies.Messages) + 1
	alies.Messages = append(alies.Messages, in.Body)
	alies.MessageNumber = len(alies.Messages)

	alies = inms.AlieSer.UpdateAlies(alies)
	if alies == nil {
		fmt.Println(" Internal Server ERROR alies update Error")
		return
	}

}

// OurEndToEndMessage method returning the messages that we chated and
// Method : GET
// AUTHENTICATION : true session.UserID
// RESPONSE : JSON
// Variables : last_message_number  , friend_id
// if last message number value is "" or 0  it will return all messages
func (inms *IndividualMessageHandler) OurEndToEndMessage(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	lang := GetSetLang(inms, response, request)
	session := inms.Session.GetSession(request)
	// context := echo.New().NewContext( )
	res := &struct {
		Success  bool              `json:"success"`
		Message  string            `json:"message"`
		Messages []*entity.Message `json:"messages"`
	}{
		Success: false,
		Message: translation.Translate(lang, "Invalid Input "),
	}
	friendID := request.FormValue("friend_id")
	lastMessageNumber := request.FormValue("last_message_number")
	val := 0
	val, er := strconv.Atoi(lastMessageNumber)
	if er != nil || friendID == "" {
		response.Write(Helper.MarshalThis(res))
		return
	}
	if val < 0 {
		val = 0
	}

	friendExist := inms.UserSer.UserWithIDExist(friendID)
	if !friendExist {
		res.Message = translation.Translate(lang, " Friend With this ID Doesn't Exist !")
		response.Write(Helper.MarshalThis(res))
		return
	}
	messages := inms.MessageSer.AliesMessages(friendID, session.UserID, val)
	if messages == nil || len(messages) == 0 {
		res.Message = translation.Translate(lang, " No Message Found ")
		response.Write(Helper.MarshalThis(res))
		return
	}
	res.Message = fmt.Sprintf(translation.Translate(lang, " Succesfuly Fetched %d Messages "), len(messages))
	res.Success = true
	res.Messages = messages
	response.Write(Helper.MarshalThis(res))
}

// SetTheMessageSeen method for individual message
// method will be called by the receiver person
// ths will set the message seen y specifiying the message number
// and the friendid
// so that the message seen = true and if the person is online  send it to the Sender and
// set the seen_confirmed value = true else if the sender is not online set it to false and save to the database
// this method is related to the Web Socket Service and to be updated later
/*

{
	"status" : 1,
	"body"  : {
		"message_number":"",
		"friend_id":"",
	}
}
METHOD : PUT
INPUT  : JSON
OutPUT  : JSON {
	"success" : false  , // notifying whether the sending of the message succesfuly
	"status" : 1,
	"body"  : {
		"message_number":"",
		"friend_id":"",
	}
}

*/
func (inms *IndividualMessageHandler) SetTheMessageSeen(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	session := inms.Session.GetSession(request)

	in := &struct {
		Status int              `json:"status"`
		Body   *entity.MessageSeen `json:"body"`
	}{}
	jdec := json.NewDecoder(request.Body)
	decError := jdec.Decode(in)


	if (in.Status != entity.MsgSeen) || (in.Body.FriendID == "")  {
		fmt.Println(" Invalid Input Nigga ...  ")
		return
	}


	if decError != nil || in.Body.MessageNumber == 0 || in.Body.FriendID == "" {
		fmt.Println(" Invalid Input Decode Error ...  ")
		return
	}
	friendExist := inms.UserSer.UserWithIDExist(in.Body.FriendID)
	if !friendExist {
		fmt.Println("No Friend By this ID Exist nigga ")
		return
	}

	// setting the message seen
	message, success := inms.MessageSer.SetMessageSeen(in.Body.FriendID, session.UserID, in.Body.MessageNumber)
	if !success {
		fmt.Println(message)
		fmt.Printf(" Setting the Message with Message Number : %d Seen Was Not Succesfuly ", in.Body.MessageNumber)
		return
	}
	fmt.Println(message)
	fmt.Println("The message is set to seen ")

	// notifying the websocket running service to broadcast the seen message is undergoing and
	// if the sender person present send the seen message to him.
	// And, set the seen_confirmed = true

}


// SeenConfirmMessage method for trial saving  the seen_confirmed service 
// this function is later to be implemented by websocket servic 
// --- this is only for trying the Seen Confirmed Service ----
// INPUT : JSON
//   
// OUTPUT : JSON 
//  


