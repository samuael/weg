package apiHandler

import (
	"fmt"
	"net/http"

	"github.com/samuael/Project/Weg/internal/pkg/Alie"
	session "github.com/samuael/Project/Weg/internal/pkg/Session"
	"github.com/samuael/Project/Weg/internal/pkg/User"
	"github.com/samuael/Project/Weg/internal/pkg/entity"
	"github.com/samuael/Project/Weg/pkg/Helper"
	"github.com/samuael/Project/Weg/pkg/translation"
)

// AliesHandler struct
type AliesHandler struct {
	Session *session.Cookiehandler
	AlieSer Alie.AlieService
	UserSer User.UserService
}

// NewAliesHandler function returning Alie Handler 
func NewAliesHandler( session *session.Cookiehandler  , alieser Alie.AlieService  , userser User.UserService , ) *AliesHandler{
	return &AliesHandler{
		Session: session,
		AlieSer: alieser,
		UserSer: userser,
	}
}

// GetSessionHandler function 
func (aliehandler *AliesHandler) GetSessionHandler() *session.Cookiehandler {
	return aliehandler.Session
}

// GetListOfAlies method returning your list of alies 
// Method : GET 
// Response : JSON
// AUTHENTICATION : user session 
func (aliehandler *AliesHandler)  GetListOfAlies(response http.ResponseWriter  , request *http.Request ){
	response.Header().Set("Content-Type"  , "application/json")
	lang := GetSetLang(aliehandler , response  , request)
	session := aliehandler.Session.GetSession(request)

	res := &struct{
		Success bool `json:"success"`
		Message string `json:"message"`
		Alies []*entity.User ` json:"alies" `
	}{
		Success: false ,
		Message : translation.Translate(lang  , " No Record Found "),
	}
	
	alies := aliehandler.AlieSer.GetMyAlies( session.UserID )
	if alies == nil {
		response.Write(Helper.MarshalThis(res))
		return 
	}
	res.Message= fmt.Sprintf(translation.Translate(lang  , "Succesfuly Found %d Alies ") , len(alies))
	res.Alies = alies
	res.Success = true
	response.Write(Helper.MarshalThis(res))
}

// DeleteAlie method to delete a firend and all the messages we talked 
// METHOD : Delete
// Variable : friend_id 
// Input : param value 
// Output  : JSON 
// Authorization : valid session and the Alie Should be your friend 
func (aliehandler *AliesHandler) DeleteAlie(response http.ResponseWriter   ,  request *http.Request ){
	response.Header().Set("Content-Type"  ,  "application/json" )
	lang  := GetSetLang(aliehandler  , response  , request)
	session := aliehandler.Session.GetSession(request)

	res := &struct{
		Success bool `json:"success"`
		Message string `json:"message"`
		FriendID  string `json:"friend_id"`
	}{
		Success : false ,
		Message:  translation.Translate(lang , " Invalid Input "),
	}
	if session == nil {
		res.Message= translation.Translate(lang , "UnAuthorized User sign In First ")
		response.Write( Helper.MarshalThis(res))
		return 
	}
	friendid := request.FormValue("friend_id")
	res.FriendID= friendid
	if friendid=="" {
		response.Write(Helper.MarshalThis(res))
		return 
	}
	userExists := aliehandler.UserSer.UserWithIDExist(friendid)
	if !userExists{
		res.Message= translation.Translate(lang , " Unknown Friend id \n User  doesn't exist ")
		response.Write(Helper.MarshalThis(res))
		return 
	}

	aliesExist := aliehandler.AlieSer.AreTheyAlies(friendid  , session.UserID)
	if !aliesExist {
		res.Message = translation.Translate(lang , " You Guys Are Not Alies  ")
		response.Write(Helper.MarshalThis(res))
		return 
	}

	alies := aliehandler.AlieSer.GetAlies(friendid  , session.UserID)
	if alies == nil {
		response.WriteHeader(http.StatusInternalServerError)
		res.Message = translation.Translate(lang  , "Internal Server ERROR ")
		response.Write(Helper.MarshalThis(res))
		return
	}
	success := aliehandler.AlieSer.DeleteAlieByID(alies.ID)
	if !success {
		fmt.Println("...")
		response.WriteHeader(http.StatusInternalServerError)
		res.Message = translation.Translate(lang  , "Internal Server ERROR ")
		response.Write(Helper.MarshalThis(res))
		return
	}
	res.Success = true 
	res.Message = translation.Translate(lang , " Deletion Was Succesful ")
	response.Write(Helper.MarshalThis(res))
} 