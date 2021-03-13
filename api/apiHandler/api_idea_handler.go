package apiHandler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

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

Input : 
{
	"title" = "this is the title  " , 
	"description" = " this is the description nigga "  , 
	"image" = "this will be 'Form File' "
}

Out Put : 
{
	"sucess" : false , 
	"message" : "message here " , 
	"idea"  : "idea"  : {
		"title":"",
		"description":"",
		"imageurl":"",
		"like_count":0,
		"dislike_count":0,
		"likers_id":null,
		"dislikers_id":null,
		"owner_id":""
	}
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
	title := request.FormValue("title")
	description := request.FormValue("description")
	fmt.Println(title , description)
	if title == "" || description == ""{
		res.Message=translation.Translate(lang  , "Invalid Input \n Tile and Description must be submitted ")
		response.Write(Helper.MarshalThis(res))
		return 
	}
	var hasImage bool 
	var IdeaImagename string 
	// var IdeaImageDirectory string 
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
		if newImage , err  = os.Create(IdeaImageDirectory+randomFilename) ; err == nil {
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
	if idea = ideah.IdeaSer.CreateIdea( idea  ); idea == nil {
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

// CreateIdeaJSONInput creating an idea usign json input 
// INPUT : JSON 
/*
	{
		"title" : "false"  , 
		"description" : "this is samauel "
	}

	OUTPUT : 
	
Out Put : 
{
	"sucess" : false , 
	"message" : "message here " , 
	"idea"  : "idea"  : {
		"title":"",
		"description":"",
		"imageurl":"",
		"like_count":0,
		"dislike_count":0,
		"likers_id":null,
		"dislikers_id":null,
		"owner_id":""
	}
}
*/
func (ideah *IdeaHandler)  CreateIdeaJSONInput(  response http.ResponseWriter  , request *http.Request ){
	lang := GetSetLang(ideah , response , request )
	session := ideah.Session.GetSession(request)
	response.Header().Set("Content-Type"  , "application/json")

	in := &struct{
		Title string `json:"title"`
		Description string `json:"description"`
	}{}
	res := &struct{
		Success bool `json:"success"`
		Message string `json:"message"`
		Idea *entity.Idea `json:"idea"`
	}{
		Success : false , 
		Message : translation.Translate(lang , "Invalid Input Please Try Again ... "),
	}
	jdec := json.NewDecoder(request.Body)
	era := jdec.Decode(in)
	if era != nil {
		res.Message=translation.Translate(lang  , "Invalid Input \n Tile and Description must be submitted ")
		response.Write(Helper.MarshalThis(res))
		return 
	}
	if in.Title == "" || in.Description == ""{
		res.Message=translation.Translate(lang  , "Invalid Input \n Tile and Description must be submitted ")
		response.Write(Helper.MarshalThis(res))
		return 
	}
	if exists := ideah.UserSer.UserWithIDExist(session.UserID); !exists{
		res.Message= translation.Translate(lang  , "You Are Not allowed To Access this Functionality ")
		ideah.Session.DeleteSession(response , request)
		response.Write(Helper.MarshalThis(res))
		return 
	}
		idea := &entity.Idea{
		Title: in.Title,
		Description: in.Description,
		Likes: 0,
		Dislikes: 0,
		LikersID: []string{},
		DislikersID: []string{},
		OwnerID: session.UserID,
	}
	if idea = ideah.IdeaSer.CreateIdea( idea  ); idea == nil {
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

// GetIdeas function for getting the ideas using the offset and limit 
// offset , limit 
// Method Get 
func (ideah *IdeaHandler)  GetIdeas( response http.ResponseWriter  , request *http.Request ){
	session := ideah.Session.GetSession(request)
	lang := GetSetLang(ideah , response  , request )
	response.Header().Set("Content-Type", "application/json")
	res := &struct{
		Success bool `json:"success"`
		Message string `json:"message"`
		Ideas []*entity.Idea `json:"ideas"`
	}{
		Success: false  , 
		Message : translation.Translate(lang  , " Records Not Found ... ") ,
		Ideas  : nil  , 
	}
	if session == nil {
		response.WriteHeader(http.StatusUnauthorized)
		return 
	}
	offset  , er := strconv.Atoi(request.FormValue("offset"))
	limit  , er := strconv.Atoi(request.FormValue("limit"))
	if er != nil {
		offset = 0
		limit =0	
	}
	ideas := ideah.IdeaSer.GetIdeas( session.UserID  , offset  , limit )
	if ideas== nil {
		res.Message = translation.Translate(lang , "No Records Found ")
		res.Ideas = []*entity.Idea{}
		response.Write(Helper.MarshalThis(res))
		return 
	}
	res.Ideas = ideas 
	res.Success = true
	res.Message = translation.Translate(lang , " Succesfuly Fetched Results ")
	response.Write(Helper.MarshalThis(res))
}

// GetIdeaByID function returning an Idea taking an ID of the Idea as a paramater
// Method : Get 
// Authorizatio : Needed 
func (ideah *IdeaHandler)  GetIdeaByID(response http.ResponseWriter , request *http.Request  ){
	response.Header().Set( "Content-Type"  , "application/json" )
	lang := GetSetLang(ideah  , response , request)
	res := &struct{
		Success bool `json:"success"`
		Message string `json:"message"`
		Idea *entity.Idea `json:"idea"`
	}{
		Success : false  , 
		Message: translation.Translate(lang  ,"Invalid input ... ") , 
		Idea : nil , 
	}
	session := ideah.Session.GetSession(request)
	if session == nil {
		response.WriteHeader(http.StatusUnauthorized)
		return 
	}
	id := request.FormValue("id")
	if id ==""{
		response.Write(Helper.MarshalThis(res))
		return 
	}
	idea := ideah.IdeaSer.GetIdeaByID(id)
	if idea == nil {
		res.Message = translation.Translate(lang  , " Record Not Found ... ")
		response.Write(Helper.MarshalThis(res))
		return 
	}
	res.Idea = idea 
	res.Success = true 
	res.Message = translation.Translate(lang , "Record Found Succesfuly ... ")
	response.Write(Helper.MarshalThis(res))
}

// DeleteIdeaByID method to delete an Idea from REpo this functionality will be visible Only for 
// 1 idea Owner and The Other For Admin 
// The Admin Implementation is not included for a while 
//  but after Admins Handler Completion We Will be adding the admin as A Valid Controller of the ID 
// METHOD : DELETE 
// Authentication : NEDDED 
// REspoonse JSON 
/*
	{
		success : true , 
		message : "Succesfuly Deleted " , 
		"idea_id" : "string id of Idea " , 
	}

*/
// Request Inptut Parameter   : idea_id 
func (ideah	*IdeaHandler) DeleteIdeaByID(response http.ResponseWriter  , request *http.Request  ){
	response.Header().Set("Content-Type", "application/json")
	lang := GetSetLang(ideah , response , request )
	res := &struct{
		Success bool `json:"success"`
		Message string `json:"message"`
		IdeaID string `json:"idea_id"`
	}{
		Success : false  , 
		Message: translation.Translate(lang   , "Un Authorized User "),
	}
	session:= ideah.Session.GetSession(request)
	if session == nil {
		response.Write(Helper.MarshalThis(res))
		return
	}
	var ideaID string 
	if ideaID= request.FormValue("idea_id"); ideaID ==""{
		res.Message= translation.Translate(lang , "Invalid Idea ID ")
		response.Write(Helper.MarshalThis(res))
		return 
	}
	res.IdeaID = ideaID 
	idea := ideah.IdeaSer.GetIdeaByID(ideaID)
	if idea== nil {
		res.Message= translation.Translate(lang , "  Record Not Found  ")
		response.Write(Helper.MarshalThis(res))
		return 
	}
	if !(idea.OwnerID == session.UserID  || session.Role == entity.ADMIN) {
		res.Message= translation.Translate(lang , "You Are Not Authorized To Perform this Action !")
		response.Write(Helper.MarshalThis(res))
		return 
	}
	success := ideah.IdeaSer.DeleteIdeaByID(ideaID)
	if success {
		res.Message = translation.Translate(lang  , "Record Deleted Succesfuly ... ")
		res.Success = true
		response.Write(Helper.MarshalThis(res))
		return 
	}
	res.Message = translation.Translate(lang , " Internal Server ERROR ... ")
	response.Write(Helper.MarshalThis(res))
}

// LikeIdea method to Like an Idea 
func  (ideah *IdeaHandler) LikeIdea(  response http.ResponseWriter  , request *http.Request ){
	response.Header().Set("Content-Type", "application/json")
	lang := GetSetLang( ideah, response , request )
	session := ideah.Session.GetSession(request)
	res := &struct{
		Success bool 
		Message string 
		Idea  *entity.Idea
	}{
		Success: false,
		Message: translation.Translate(lang , " UnAuthorized User ... ") , 
	}
	if session == nil {
		response.Write(Helper.MarshalThis(res))
		return 
	}
	var ideaID string 
	if ideaID = request.FormValue("id"); ideaID == ""{
		res.Message = translation.Translate(lang  , " Invalid Idea ID ")
		response.Write(Helper.MarshalThis(res))
		return
	}
	idea := ideah.IdeaSer.GetIdeaByID(ideaID)
	if idea == nil {
		res.Message = translation.Translate(lang  , " Idea not Found ")
		response.Write(Helper.MarshalThis(res))
		return 
	}
	var liked bool 
	liked = false
	idea.LikersID= func()[]string{
		for _  , id := range idea.LikersID{
			if id == ideaID{
				res.Message = translation.Translate(lang , "You have Already Like it")
				res.Idea = idea 
				response.Write(Helper.MarshalThis(res))
				liked = true
				return idea.LikersID
			}
		}
		ideaso :=  append( idea.LikersID  , ideaID)
		return ideaso
	}()
	if liked{
		return 
	}

	var wasDisliked bool 
	wasDisliked = false 
	idea.DislikersID = func()[]string {
		for ind  , id := range idea.DislikersID{
			if id == ideaID {
				wasDisliked= true 
				// removing the ID from Disliked IDS list 
				return append(idea.DislikersID[:ind] , func()[]string{
					if len( idea.DislikersID ) - 1 > ind {

						return idea.DislikersID[ind+1:]
					} 
					return []string{}
				}()... )
			}
		}
		return []string{}
	}()

	if wasDisliked {
		idea.Dislikes--
	}
	idea.Likes++
	idea = ideah.IdeaSer.UpdateIdea(idea)
	if idea == nil {
		res.Message= translation.Translate(lang  , " Internal Server ERROR ")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write(Helper.MarshalThis(res))
		return 
	}
	res.Message = translation.Translate(lang  , "Liked")
	res.Success = true
	res.Idea= idea
	response.Write(Helper.MarshalThis(res))
}
// DislikeIdea method to DIsLike an Idea 
func  (ideah *IdeaHandler) DislikeIdea(  response http.ResponseWriter  , request *http.Request ){
	response.Header().Set("Content-Type", "application/json")
	lang := GetSetLang( ideah, response , request )
	session := ideah.Session.GetSession(request)
	res := &struct{
		Success bool 
		Message string 
		Idea  *entity.Idea
	}{
		Success: false,
		Message: translation.Translate(lang , " UnAuthorized User ... ") , 
	}
	if session == nil {
		response.Write(Helper.MarshalThis(res))
		return 
	}
	var ideaID string 
	if ideaID = request.FormValue("id"); ideaID == ""{
		res.Message = translation.Translate(lang  , " Invalid Idea ID ")
		response.Write(Helper.MarshalThis(res))
		return
	}
	idea := ideah.IdeaSer.GetIdeaByID(ideaID)
	if idea == nil {
		res.Message = translation.Translate(lang  , " Idea not Found ")
		response.Write(Helper.MarshalThis(res))
		return 
	}
	var disliked bool 
	disliked = false
	idea.DislikersID= func()[]string{
		for _  , id := range idea.DislikersID{
			if id == ideaID{
				res.Message = translation.Translate(lang , "You have Already Disike it")
				res.Idea = idea 
				response.Write(Helper.MarshalThis(res))
				disliked= true
				return idea.DislikersID
			}
		}
		ideaso :=  append( idea.DislikersID  , ideaID)
		return ideaso
	}()
	if disliked{
		return 
	}

	var wasliked bool 
	wasliked = false 
	idea.LikersID = func()[]string {
		for ind  , id := range idea.LikersID{
			if id == ideaID {
				wasliked= true 
				// removing the ID from Disliked IDS list 
				return append(idea.LikersID[:ind] , func()[]string{
					if len( idea.LikersID ) - 1 > ind {

						return idea.LikersID[ind+1:]
					} 
					return []string{}
				}()... )
			}
		}
		return []string{}
	}()

	if wasliked {
		idea.Likes--
	}
	idea.Dislikes++
	idea = ideah.IdeaSer.UpdateIdea(idea)
	if idea == nil {
		res.Message= translation.Translate(lang  , " Internal Server ERROR ")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write(Helper.MarshalThis(res))
		return 
	}
	res.Message = translation.Translate(lang  , "Disliked ! ")
	res.Success = true
	res.Idea= idea
	response.Write(Helper.MarshalThis(res))
}

// Dislike Idea 

// UpdateIdea function to Update an Idea This requires the Idea to By Yours 
// not even the Admin change teh Idea of an Individual 
// For a while it is made only to cahnge the title and the Description of an Idea 
// METHOD : PUT 
/*

{
	"title" : "" , 
	"description"  : "" , 
}

*/
// INPUT : JSON
// OUTPUT : JSON 
/*  
{
	"success" : false  , 
	"message" : "invalid Input Niggoye" , 
	"idea"  : {
		"title":"",
		"description":"",
		"imageurl":"",
		"like_count":0,
		"dislike_count":0,
		"likers_id":null,
		"dislikers_id":null,
		"owner_id":""
	}
}

*/
func (ideah *IdeaHandler ) UpdateIdea(response http.ResponseWriter  , request *http.Request ){
	response.Header().Set( "Content-Type"  , "application/json")
	lang := GetSetLang(ideah , response  , request )
	res := &struct{
		Success bool `json:"success"`
		Message string `json:"message"`
		Idea *entity.Idea `json:"idea"`
	}{
		Success : false  , 
		Message : translation.Translate(lang , "Invalid Input ..") , 
		Idea : &entity.Idea{} , 
	}
	session := ideah.Session.GetSession(request)
	if session == nil {
		res.Message= translation.Translate(lang  , "UnAuthorized User ")
		response.Write(Helper.MarshalThis(res))
		return 
	}
	jsonDecoder := json.NewDecoder(request.Body)
	decErro := jsonDecoder.Decode(res.Idea)
	if (decErro != nil) || (res.Idea == nil) {
		// Error Happened whle Decoding nigga 
		// so  , I am Gonna return Invalid Input Error 
		print(decErro);
		response.Write(Helper.MarshalThis(res))
		return 
	}
	if res.Idea.ID == "" {
		res.Message= translation.Translate(lang , "Invalid Value Idea ID must be specified !")
		response.Write(Helper.MarshalThis(res))
		return 
	}
	idea := ideah.IdeaSer.GetIdeaByID(res.Idea.ID)
	if idea == nil {
		res.Message= translation.Translate(lang , " No Record Found! \nYou Have No Idea By the Specified ID ")
		response.Write( Helper.MarshalThis(res) )
		return		
	}
	if idea.OwnerID != session.UserID {
		res.Message = translation.Translate(lang  , "You Are Not Authorized o change the Idea !")

		response.Write(Helper.MarshalThis(res))
		return 
	}

	idea.Title=func()string{ 
		if res.Idea.Title==""{
			return idea.Title
		}
	return res.Idea.Title}()
	idea.Description = func()string{
		if res.Idea.Description==""{
		return idea.Description
		}
		return res.Idea.Description}()

		idea = ideah.IdeaSer.UpdateIdea(idea)
		if idea == nil {
			res.Message = translation.Translate(lang  ,"Internal Server Error ")
			response.Write(Helper.MarshalThis(res))
			return 
		}
		res.Idea= idea 
		res.Success = true
		res.Message = translation.Translate(lang  ,"Succesfuly Updated !")
		response.Write(Helper.MarshalThis(res))
}

// GetIdeasByUserID to get the Idea of specified User usign the user id as a string 
// METHOD : 	GET 
// INPUT : VARIABLE NAMED 'user_id'
// RESPONSE : JSON 
/* 

{
	"success" : true  , 
	"message" : "mesage String nigga" , 
	"ideas"  : [
		 {
		"title":"",
		"description":"",
		"imageurl":"",
		"like_count":0,
		"dislike_count":0,
		"likers_id":null,
		"dislikers_id":null,
		"owner_id":""
	} , 
	 {
		"title":"",
		"description":"",
		"imageurl":"",
		"like_count":0,
		"dislike_count":0,
		"likers_id":null,
		"dislikers_id":null,
		"owner_id":""
	} , 
	]
}
*/
func (ideah *IdeaHandler)   GetIdeasByUserID(response http.ResponseWriter  , request *http.Request ){
	response.Header().Set("Content-Type", "application/json")
	lang := GetSetLang(ideah , response , request)
	session := ideah.Session.GetSession(request)
	res := &struct{
		Success bool  `json:"success"`
		Message string `json:"message"`
		Ideas []*entity.Idea `json:"ideas"`
	}{
		Success : false , 
		Message : translation.Translate(lang  , "No Record Found ..") ,
		Ideas: []*entity.Idea{} , 
	}
	if session == nil {
		res.Message = translation.Translate(lang , "UnAuthorized User ")
		response.Write(Helper.MarshalThis(res))
		return
	}
	userid := request.FormValue("userid")
	if userid == ""{
		userid = session.UserID
	}
	ideas := ideah.IdeaSer.GetMyIdeas(userid)
	if ideas==nil || len(ideas)==0{
		res.Message = translation.Translate(lang , "No Records Found .. ")
		response.Write(Helper.MarshalThis(res))
		return
	}
	res.Success = true
	res.Message= translation.Translate(lang , " Records Found ")
	res.Ideas = ideas
	response.Write(Helper.MarshalThis(res))
}


// SearchIdeaByTitle  a method to search an ID usign the tile  
// method JSON 
// INPUT  title 
func (ideah *IdeaHandler)   SearchIdeaByTitle(  response http.ResponseWriter  , request *http.Request ){
	response.Header().Set("Content-Type", "application/json")
	res := &struct{
		Success bool `json:"success"`
		Message string `json:"message"`
		Ideas []*entity.Idea `json:"ideas"`
	}{
		Success : false  ,
		Message : "No Record Found" ,  
	}
	title := request.FormValue("title")
	if (title ==""){
		response.Write(Helper.MarshalThis(res))
		return 
	}
	ideas := ideah.IdeaSer.SearchIdeaByTitle(title)
	if ideas == nil  || len(ideas)==0 {
		res.Ideas=[]*entity.Idea{}
		response.Write(Helper.MarshalThis(res))
		return 
	}
	res.Message = fmt.Sprintf("  Succesfuly fetched %d Ideas " , len(ideas)) 
	res.Ideas = ideas
	res.Success = true
	response.Write(Helper.MarshalThis(res))
}

// SearchIdeaByTitle 
