package handler

import (
	"html/template"
	"net/http"

	"github.com/samuael/Project/Weg/api/apiHandler"
	"github.com/samuael/Project/Weg/internal/pkg/Group"
	session "github.com/samuael/Project/Weg/internal/pkg/Session"
	"github.com/samuael/Project/Weg/internal/pkg/User"
	"github.com/samuael/Project/Weg/internal/pkg/entity"
	// "github.com/samuael/Project/Weg/internal/pkg/User/UserService"
)

// TemplateHandler template handler
type TemplateHandler struct {
	Templs  *template.Template
	Session *session.Cookiehandler
	UserSer  User.UserService
	GroupSer Group.GroupService
}

// NewTemplateHandler function returning template handler pointer 
func NewTemplateHandler(templs *template.Template ,session *session.Cookiehandler , useser User.UserService , groupser Group.GroupService) *TemplateHandler {
	return &TemplateHandler{
		Templs: templs ,
		Session: session,
		UserSer: useser,
		GroupSer: groupser,
	}
}

// GetSessionHandler function returning template handler session instance 
func (th *TemplateHandler) GetSessionHandler() *session.Cookiehandler {
	return th.Session
}

// LoginPage method returning login page 
func (th *TemplateHandler) LoginPage(response http.ResponseWriter  ,request *http.Request) {

	session := th.Session.GetSession(request)
	lang  := apiHandler.GetSetLang(th , response  , request)
	if session != nil {

	}{
		
	}

}