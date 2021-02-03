package apiHandler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/samuael/Project/Weg/internal/pkg/Idea"
	session "github.com/samuael/Project/Weg/internal/pkg/Session"
	"github.com/samuael/Project/Weg/internal/pkg/User"
	"github.com/samuael/Project/Weg/internal/pkg/entity"
	"github.com/samuael/Project/Weg/pkg/Helper"
	"github.com/samuael/Project/Weg/pkg/translation"
)

// IdeaHandler handler instance representing the idea and it's methods
type IdeaHandler struct {
	Session *session.Cookiehandler
	IdeaSer   Idea.IdeaService
	UserSer User.UserService
}

// NewIdeaHandler returning an IdeaHandler Instance 
func NewIdeaHandler(
	session *session.Cookiehandler , 
	ideaSer Idea.IdeaService , 
	userser User.UserService  , 
) *IdeaHandler{
	return &IdeaHandler{
		Session: session,
		IdeaSer : ideaSer , 
		UserSer:  userser,
	}
}
// GetSessionHandler representing idea handler's DHandler implementation
func (ideah *IdeaHandler) GetSessionHandler() *session.Cookiehandler {
	return ideah.Session
}

// CreateIdea to craete an Idea 
// Input Method POST 
// INPUT Variable FormData 
// OutPut : JSON 
/*
{
	"title" = "this is the title  " , 
	"description" = " this is the description nigga "  , 
	"image" = "this will be 'Form File' "
}
	Owner ID  from session 
*/ 
func (ideah *IdeaHandler)  CreateIdea(  response http.ResponseWriter  , request *http.Request ){
	lang := GetSetLang(ideah , response , request )
	session := ideah.Session.GetSession(request)
	response.Header().Set("Content-Type"  , "application/json")

	res := &struct{
		Success bool `json:"success"`
		Message string `json:"message"`
		Idea *entity.Idea `json:"idea"`
	}{
		Success : false , 
		Message : translation.Translate(lang , "Invalid Input Please Try Again ... "),
	}

	er := request.ParseForm()
	if er != nil {
		response.Write(Helper.MarshalThis(res))
		return 
	}
	title := request.FormValue("title")
	description := request.FormValue("description")
	if title == "" || description == ""{
		res.Message=translation.Translate("Invalid Input \n Tile and Description must be submitted ")
		response.Write(Helper.MarshalThis(res))
		return 
	}
	var hasImage bool 
	var IdeaImagename string 
	var IdeaImageDirectory string 
	hasImage = true 
	var newImage *os.File
	if exists := ideah.UserSer.UserWithIDExist(session.UserID); !exists{
		res.Message= translation.Translate(lang  , "You Are Not allowed To Access this Functionality ")
		ideah.Session.DeleteSession(response , request)
		response.Write(Helper.MarshalThis(res))
		return 
	}
	//  below this in the if statment i will handle copying and handling images 
	if image , header  , err  := request.FormFile("image") ; err == nil {
		defer func(){
			image.Close()
		}()

		randomFilename := fmt.Sprintf("%s.%s" , Helper.GenerateRandomString(10 , Helper.CHARACTERS)  ,  Helper.GetExtension(header.Filename ))
		IdeaImageDirectory := entity.PathToIdeasImageDirectoryFromMain
		if newImage , err  = os.Create(IdeaImageDirectory+randomFilename) ; era == nil {
			// close this image file after the execution of this command 
			defer newImage.Close()
			_ , copyError := io.Copy(newImage  , image)
			if copyError != nil {
				hasImage= false 
				goto copyErrors
			}
			hasImage = true  
			IdeaImagename = fmt.Sprintf("%s%s" , entity.PathToIdeasImageDirectory , randomFilename )
		}else {
			hasImage= false 
		}
	}
	copyErrors : fmt.Println("Copy Error  in the Idea Creation ... ")
	idea := &entity.Idea{
		Title: title,
		Description: description,
		Likes: 0,
		Dislikes: 0,
		LikersID: []string{},
		DislikersID: []string{},
		OwnerID: session.UserID,
		ImageURL: func()string {
			if hasImage {
				return IdeaImagename
			}
			return ""
		}(),
	}
	if idea = ideah.IdeaSer.CreteIdea( idea  ); idea == nil {
		res.Message = translation.Translate(lang  , " Inernal Server Error ")
		log.Println("Internal Server ERROR While Creating an Idea .. ")
		response.Write(Helper.MarshalThis(res))
		return
	}
	res.Message= translation.Translate(lang , "Succesfuly Created an Idea .. ")
	res.Success = true
	res.Idea = idea 
	response.Write(Helper.MarshalThis(res))
}