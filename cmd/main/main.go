package main

import (
	"context"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"

	"log"

	// adminre "github.com/samuael/Project/Weg/internal/pkg/Admin/AdminRepo"
	// adminse "github.com/samuael/Project/Weg/internal/pkg/Admin/AdminService"

	"github.com/gorilla/mux"
	"github.com/samuael/Project/Weg/api/apiHandler"
	"github.com/samuael/Project/Weg/cmd/service"
	"github.com/samuael/Project/Weg/cmd/service/grpc_service_conn/client"
	"github.com/samuael/Project/Weg/cmd/service/grpc_service_conn/server"
	"github.com/samuael/Project/Weg/internal/pkg/Admin/AdminRepo"
	"github.com/samuael/Project/Weg/internal/pkg/Admin/AdminService"
	"github.com/samuael/Project/Weg/internal/pkg/Alie/AlieRepo"
	"github.com/samuael/Project/Weg/internal/pkg/Alie/AlieService"
	"github.com/samuael/Project/Weg/internal/pkg/Group/GroupRepo"
	"github.com/samuael/Project/Weg/internal/pkg/Group/GroupService"
	"github.com/samuael/Project/Weg/internal/pkg/Idea/IdeaRepo"
	"github.com/samuael/Project/Weg/internal/pkg/Idea/IdeaService"
	"github.com/samuael/Project/Weg/internal/pkg/entity"
	"google.golang.org/grpc"

	// "github.com/samuael/Project/Weg/internal/pkg/Message"
	pb "github.com/samuael/Project/Weg/cmd/service/grpc_service_conn/proto"
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
	print("I got calles ... ")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			return
		}
		print("Serving ... ")
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

	idearepo := IdeaRepo.NewIdeaRepo(db)
	ideaser := IdeaService.NewIdeaService(idearepo)
	ideahand := apiHandler.NewIdeaHandler(sessionHandler, ideaser, userser)

	// adminRelated instances
	adminrepo := AdminRepo.NewAdminRepo(db)
	adminser := AdminService.NewAdminService(adminrepo)
	adminhandler := apiHandler.NewAdminHandler(sessionHandler, adminser, userser)

	// Continuously Running service objects instantiation

	mainservice := service.NewMainService()
	groupservice := service.NewGroupService()
	clientservice := service.NewClientService(
		mainservice,
		messageSer,
		groupservice,
		userser,
		gservice,
		alieSer,
		sessionHandler)

	//  Registering the Grpc handler
	grpcHandler := server.NewGrpcEEServer(mainservice)
	grpcClient := client.NewGrpcClient()
	// server instantiated at the default port of 8070

	mainservice.SetGrpcHandler(grpcClient)

	// Instantiating the server

	// Instantiating the GRPC sserver
	go grpcClient.UpdatePeerServers()
	go InstantiateGrpcServer(grpcHandler)

	go mainservice.Run()
	go clientservice.Run()
	// -----------------end --------------------

	userhandler := apiHandler.NewUserHandler(sessionHandler, userser, clientservice)
	grouphandler := apiHandler.NewGroupHandler(sessionHandler, gservice, userser)
	inmshandler := apiHandler.NewIndvMessageHandler(sessionHandler, messageSer, userser, alieSer)
	gmhandler := apiHandler.NewGroupMessageHandler(gservice, userser, sessionHandler, messageSer)

	mux := mux.NewRouter() //.StrictSlash(true)
	fs := http.FileServer(http.Dir("../../web/templates/assets/"))
	http.Handle("/assets/", http.Handler(http.StripPrefix("/assets/", neuter(fs))))
	http.Handle("/", mux)

	// waiting for chat ws or wss web socket services with a route /chat/
	// this creates a web socket client object and create a continuously running loop for each
	// client connection ? Not User - but , clients of user
	// meaning One User may have multiple clients From Mobile , Browser etc and he/she can use
	// multiple device on one account.
	mux.Handle("/chat/", clientservice)

	apiroute := mux.PathPrefix("/api/").Subrouter()

	adminsRoute := apiroute.PathPrefix("/admin/").Subrouter()
	adminsRoute.HandleFunc("/new/", userhandler.Authenticated(adminhandler.CreateAdmin)).Methods(http.MethodPost)
	adminsRoute.HandleFunc("/", userhandler.Authenticated(adminhandler.DeleteAdmin)).Methods(http.MethodDelete)
	adminsRoute.HandleFunc("/", userhandler.Authenticated(adminhandler.UpdateAdmin)).Methods(http.MethodPut)
	adminsRoute.HandleFunc("/login/", adminhandler.AdminLogin).Methods(http.MethodPost)
	// CreateIdea  ownerid    DeleteIdeaByID DislikeIdea  GetIdeasByUserID CreateIdeaJSONInput
	apiroute.HandleFunc("/idea/new/", userhandler.Authenticated(ideahand.CreateIdea)).Methods(http.MethodPost)
	apiroute.HandleFunc("/idea/new/", userhandler.Authenticated(ideahand.CreateIdeaJSONInput)).Methods(http.MethodPut)
	apiroute.HandleFunc("/idea/", userhandler.Authenticated(ideahand.UpdateIdea)).Methods(http.MethodPut)
	apiroute.HandleFunc("/ideas/", userhandler.Authenticated(ideahand.GetIdeas)).Methods(http.MethodGet)
	apiroute.HandleFunc("/idea/", userhandler.Authenticated(ideahand.GetIdeaByID)).Methods(http.MethodGet)
	apiroute.HandleFunc("/idea/", userhandler.Authenticated(ideahand.DeleteIdeaByID)).Methods(http.MethodDelete)
	apiroute.HandleFunc("/idea/like", userhandler.Authenticated(ideahand.LikeIdea)).Methods(http.MethodGet)
	apiroute.HandleFunc("/idea/dislike", userhandler.Authenticated(ideahand.DislikeIdea)).Methods(http.MethodGet)
	apiroute.HandleFunc("/user/ideas/", userhandler.Authenticated(ideahand.GetIdeasByUserID)).Methods(http.MethodGet) // SearchIdeaByTitle
	apiroute.HandleFunc("/idea/search/", ideahand.SearchIdeaByTitle).Methods(http.MethodGet)                          // SearchIdeaByTitle

	apiroute.HandleFunc("/user/new/", userhandler.RegisterClient).Methods(http.MethodPost)
	apiroute.HandleFunc("/user/login/", userhandler.Login).Methods(http.MethodPost)

	apiroute.HandleFunc("/logout", userhandler.Authenticated(userhandler.Logout)).Methods("GET")
	apiroute.HandleFunc("/user/img/", userhandler.Authenticated(userhandler.UploadProfilePic)).Methods(http.MethodPut)
	apiroute.HandleFunc("/lang/new/", userhandler.Authenticated(userhandler.ChangeLanguage)).Methods(http.MethodGet)
	apiroute.HandleFunc("/user/", userhandler.Authenticated(userhandler.UpdateUserProfile)).Methods(http.MethodPut)
	apiroute.HandleFunc("/user/myprofile/", userhandler.Authenticated(userhandler.MyProfile)).Methods(http.MethodGet)
	apiroute.HandleFunc("/user/password/new/", userhandler.Authenticated(userhandler.ChangeUserPassword)).Methods(http.MethodPut)
	apiroute.HandleFunc("/user/search/", userhandler.Authenticated(userhandler.SearchUsers)).Methods(http.MethodGet) //
	apiroute.HandleFunc("/user/", userhandler.Authenticated(userhandler.DeleteMyAccount)).Methods(http.MethodDelete) //  GetUserByID
	apiroute.HandleFunc("/user/", userhandler.GetUserByID).Methods(http.MethodGet)                                   //

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
	apiroute.HandleFunc("/user/message/", userhandler.Authenticated(inmshandler.DeleteMessage)).Methods(http.MethodDelete)
	apiroute.HandleFunc("/user/message/seen/", userhandler.Authenticated(inmshandler.SetTheMessageSeen)).Methods(http.MethodPut)
	apiroute.HandleFunc("/group/message/new/", userhandler.Authenticated(gmhandler.SendGroupMessage)).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(entity.SERVER_PORT, nil))
}

func InstantiateGrpcServer(serv *server.GrpcEEHandler) {
	lis, err := net.Listen("tcp", entity.GRPC_PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMessageServiceServer(s, serv)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
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
