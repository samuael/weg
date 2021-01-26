package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
	"sync"

	"log"

	// adminre "github.com/samuael/Project/Weg/internal/pkg/Admin/AdminRepo"
	// adminse "github.com/samuael/Project/Weg/internal/pkg/Admin/AdminService"

	"github.com/gorilla/mux"
	"github.com/samuael/Project/Weg/api/apiHandler"
	"github.com/samuael/Project/Weg/internal/pkg/Alie/AlieRepo"
	"github.com/samuael/Project/Weg/internal/pkg/Alie/AlieService"
	"github.com/samuael/Project/Weg/internal/pkg/Group/GroupRepo"
	"github.com/samuael/Project/Weg/internal/pkg/Group/GroupService"

	// "github.com/samuael/Project/Weg/internal/pkg/Message"
	"github.com/samuael/Project/Weg/internal/pkg/Message/MessageRepo"
	"github.com/samuael/Project/Weg/internal/pkg/Message/MessageService"
	session "github.com/samuael/Project/Weg/internal/pkg/Session"
	"github.com/samuael/Project/Weg/internal/pkg/User/UserRepo"
	"github.com/samuael/Project/Weg/internal/pkg/User/UserService"
	DB "github.com/samuael/Project/Weg/internal/pkg/db"
	"go.mongodb.org/mongo-driver/mongo"
)

var once = sync.Once{}

// var th *handler.TemplateHandler
var systemTemplate *template.Template
var db *mongo.Database
var erro error
var sessionHandler *session.Cookiehandler

// var studentHandler *handler.StudentHandler

// For Filtering and Preventing Directory Listening...
func neuter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			return
		}
		next.ServeHTTP(w, r)
	})
}
func startUp() {
	once.Do(
		func() {
			db = DB.ConnectMongodb()
			if db == nil {
				os.Exit(1)
			}
			return
			//Instanciating Templates nigga
			// systemTemplate = (template.Must(template.New("Weg").Funcs(handler.FuncMap).ParseGlob("../../web/templates/*.html")))
		},
	)
}

func init() {
	startUp()
}

func main() {

	if db == nil {
		fmt.Println("DB IS Nill ")
		os.Exit(1)
	}
	defer db.Client().Disconnect(context.TODO())
	sessionHandler = session.NewCookieHandler()

	alierepo := AlieRepo.NewAlieRepo(db)
	alieSer := AlieService.NewAlieService(alierepo)

	userrepo := UserRepo.NewUserRepo(db)
	userser := UserService.NewUserService(userrepo)

	grouprepo := GroupRepo.NewGroupRepo(db)
	gservice := GroupService.NewGroupService(grouprepo)
	inmsRepo := MessageRepo.NewMessageRepo(db, alierepo)
	messageSer := MessageService.NewMessageService(inmsRepo)

	aliehandler := apiHandler.NewAliesHandler(sessionHandler, alieSer, userser)
	userhandler := apiHandler.NewUserHandler(sessionHandler, userser)
	grouphandler := apiHandler.NewGroupHandler(sessionHandler, gservice, userser)
	inmshandler := apiHandler.NewIndvMessageHandler(sessionHandler, messageSer, userser, alieSer)
	gmhandler := apiHandler.NewGroupMessageHandler(gservice, userser, sessionHandler, messageSer)

	mux := mux.NewRouter() //.StrictSlash(true)
	fs := http.FileServer(http.Dir("../../web/templates/assets/"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", neuter(fs)))

	apiroute := mux.PathPrefix("/api/").Subrouter()

	apiroute.HandleFunc("/user/new/", userhandler.RegisterClient).Methods(http.MethodPost)
	apiroute.HandleFunc("/user/login/", userhandler.Login).Methods(http.MethodPost)

	apiroute.HandleFunc("/logout", userhandler.Authenticated(userhandler.Logout)).Methods("GET")
	apiroute.HandleFunc("/user/img/", userhandler.Authenticated(userhandler.UploadProfilePic)).Methods(http.MethodPut)
	apiroute.HandleFunc("/lang/new/", userhandler.Authenticated(userhandler.ChangeLanguage)).Methods(http.MethodGet)
	apiroute.HandleFunc("/user/", userhandler.Authenticated(userhandler.UpdateUserProfile)).Methods(http.MethodPut)
	apiroute.HandleFunc("/user/password/new/", userhandler.Authenticated(userhandler.ChangeUserPassword)).Methods(http.MethodPut)

	// GetGroupMembersList
	apiroute.HandleFunc("/group/new/", userhandler.Authenticated(grouphandler.CreateGroup)).Methods(http.MethodPost)
	apiroute.HandleFunc("/group/", userhandler.Authenticated(grouphandler.DeleteGroup)).Methods(http.MethodDelete)
	apiroute.HandleFunc("/group/img/", userhandler.Authenticated(grouphandler.UpdateGroupProfilePicture)).Methods(http.MethodPut)
	apiroute.HandleFunc("/group/join/", userhandler.Authenticated(grouphandler.JoinGroup)).Methods(http.MethodPut)
	apiroute.HandleFunc("/group/leave/", userhandler.Authenticated(grouphandler.LeaveGroup)).Methods(http.MethodPut)
	apiroute.HandleFunc("/groups/", userhandler.Authenticated(grouphandler.MyGroups)).Methods(http.MethodGet)
	apiroute.HandleFunc("/group/search/", userhandler.Authenticated(grouphandler.SearchGroupByName)).Methods(http.MethodGet)
	apiroute.HandleFunc("/group/members/", userhandler.Authenticated(grouphandler.GetGroupMembersList)).Methods(http.MethodGet)

	// AlieRelated Routed
	apiroute.HandleFunc("/user/friends/", userhandler.Authenticated(aliehandler.GetListOfAlies)).Methods(http.MethodGet)
	apiroute.HandleFunc("/user/friend/", userhandler.Authenticated(aliehandler.DeleteAlie)).Methods(http.MethodDelete)

	// Message Related Routes that are not Websocket Dependant
	apiroute.HandleFunc("/user/friend/messages/", userhandler.Authenticated(inmshandler.OurEndToEndMessage)).Methods(http.MethodGet)
	apiroute.HandleFunc("/group/messages/", userhandler.Authenticated(gmhandler.GetGroupMessage)).Methods(http.MethodGet)

	// Temp Routes    that could be changed to websocket implemmentaation
	apiroute.HandleFunc("/user/message/new/", userhandler.Authenticated(inmshandler.SendAlieMessage)).Methods(http.MethodPost)
	apiroute.HandleFunc("/user/message/seen/", userhandler.Authenticated(inmshandler.SetTheMessageSeen)).Methods(http.MethodPut)
	apiroute.HandleFunc("/group/message/new/", userhandler.Authenticated(gmhandler.SendGroupMessage)).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(":8080", mux))
}

// DirectoryListener  representing
func DirectoryListener(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
