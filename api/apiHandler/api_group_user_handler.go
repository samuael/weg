package apiHandler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/samuael/Project/Weg/internal/pkg/Group"
	session "github.com/samuael/Project/Weg/internal/pkg/Session"
	"github.com/samuael/Project/Weg/internal/pkg/User"
	entity "github.com/samuael/Project/Weg/internal/pkg/entity"
	"github.com/samuael/Project/Weg/pkg/Helper"
	"github.com/samuael/Project/Weg/pkg/translation"
)

// GroupHandler handling group related Calls
type GroupHandler struct {
	Session  *session.Cookiehandler
	GroupSer Group.GroupService
	UserSer User.UserService
}

// NewGroupHandler function  returning GroupHandler instance pointer
func NewGroupHandler(session *session.Cookiehandler, gser Group.GroupService , userser User.UserService) *GroupHandler {
	return &GroupHandler{
		Session:  session,
		GroupSer: gser,
		UserSer : userser ,
	}
}

// GetSessionHandler to implement the DHandler Interface
func (ghandler *GroupHandler) GetSessionHandler() *session.Cookiehandler {
	return ghandler.Session
}

/*// Group struct representing chat Group
type Group struct {
	ID string `bson:"_id,omitempty"`
	OwnerID string // string representing user ID
	MembersCount int
	Imageurl string
	ActiveCounts int
	GroupName string
	Description string
	LastMessageNumber int
	MembersID []string
	CreatedAt time.Time

}*/

// CreateGroup apimethod to handler group Creation API
// METHOD POST
// INPUT JSON
// VARIABLES group_name , description
func (ghandler *GroupHandler) CreateGroup(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	lang := GetSetLang(ghandler, response, request)
	ret := struct {
		Success bool          `json:"success"`
		Message string        `json:"message"`
		Group   *entity.Group `json:"group"`
	}{
		Success: false,
		Message: translation.Translate(lang, "Invalid Input Please Send A valid "),
	}
	in := &struct {
		GroupName   string `json:"group_name"`
		Description string `json:"description"`
	}{}
	jdec := json.NewDecoder(request.Body)
	er := jdec.Decode(in)
	if er != nil {
		response.Write(Helper.MarshalThis(ret))
		return
	}
	if in.GroupName=="" || in.Description =="" {
		ret.Message= translation.Translate(lang, " Group Name and Description has to be fulfilled Correctly ")
		response.Write(Helper.MarshalThis(ret))
		return 
	}
	session := ghandler.Session.GetSession(request)
	group := &entity.Group{
		OwnerID:           session.UserID,
		CreatedAt:         time.Now(),
		MembersID:         []string{session.UserID},
		MembersCount:      1,
		GroupName:         in.GroupName,
		Description:       in.Description,
		LastMessageNumber: 0,
	}
	group = ghandler.GroupSer.CreateGroup(group)
	if group == nil {
		ret.Message = translation.Translate(lang, "Internal Server ERROR ")
		response.Write(Helper.MarshalThis(ret))
		return
	}
	ret.Success = true
	ret.Message = translation.Translate(lang, " Group Succesfuly Created ")
	ret.Group = group
	response.Write(Helper.MarshalThis(ret))
}

// DeleteGroup  function handling deleting group
// METHOD : DELETE 
// INPUT : JSON 
// AUTHORIZATION : the one who craeted the Grou[p]
// RESPONSE : JSON
func (ghandler *GroupHandler) DeleteGroup(response http.ResponseWriter , request *http.Request ){
	response.Header().Set("Content-Type" , "application/json" )
	lang := GetSetLang(ghandler, response, request)
	in := &struct{
		GroupID string `json:"group_id"`
	}{}
	res := &struct{
		Success bool `json:"success"`
		Message string `json:"message"`
		GroupID string `json:"group_id"`
	}{
		Success: false ,
		Message: translation.Translate(lang  , "Invalid Input Please Try Again "),
	}
	session := ghandler.Session.GetSession(request)
	jdec := json.NewDecoder(request.Body)
	er := jdec.Decode(in)
	if er != nil || in.GroupID=="" {
		response.Write(Helper.MarshalThis(res))
		return 
	}
	group := ghandler.GroupSer.GetGroupByID(in.GroupID)
	if group == nil {
		res.Message=translation.Translate(lang  , " No Record Found By This ID ")
		response.Write(Helper.MarshalThis(res))
		return
	}
	if group.OwnerID != session.UserID {
		res.Message=translation.Translate(lang  , "Yout Are Authorized For This Action ")
		response.Write(Helper.MarshalThis(res))
		return
	}

	success := ghandler.GroupSer.DeleteGroup(group.ID)
	if !success {
		res.Message=translation.Translate(lang  , "INTERNAL SERVER ERROR PLEASE TRY AGAIN ")
		response.Write(Helper.MarshalThis(res))
		return
	}
	res.Success=true
	res.GroupID= group.ID
	res.Message=translation.Translate(lang  , "Succesfuly Deleted Group Named \""+ group.GroupName+"\" With ID : \"" + group.ID + "\"")
	response.Write(Helper.MarshalThis(res))
}

// UpdateGroupProfilePicture method to change the Profile Picture of the Group 
// METHOD : PUT
// INPUT : image >> file and group_id >> string 
// RESPONSE : JSON
// AUTHORIZATION NEEDED :: ONLY THE ONE WHO CREATED THE gROUP WILL BE ABLE TO DELETE
func (ghandler *GroupHandler)   UpdateGroupProfilePicture( response http.ResponseWriter  , request *http.Request  ){
	response.Header().Set("Content-Type"  , "application/json")
	lang := GetSetLang(ghandler  ,response  ,  request )
	session := ghandler.Session.GetSession(request)
	res := &struct{
		Success bool `json:"success"`
		Message string `json:"message"`
		Imageurl string `json:"imageurl"`
	}{
		Success: false,
		Message: translation.Translate(lang  , "Invalid Input "),
	}
	era := request.ParseMultipartForm(99999999999)
	if era != nil {
		response.Write(Helper.MarshalThis(res))
		return 
	}
	groupid := request.FormValue("group_id")
	file  , header  , er := request.FormFile("image")
	if er != nil || groupid=="" {
		// fmt.Println(er)
		res.Message = translation.Translate( lang  , " Invalid Input ")
		response.Write(Helper.MarshalThis(res))
		return
	} 
	if isImage := Helper.IsImage(header.Filename); !isImage {
		res.Message =translation.Translate( lang  , " Invalid Input / Only Image Files are Allowed ... ")
		response.Write(Helper.MarshalThis(res))
		return
	}
	randomName := "assets/img/GroupImage/"+Helper.GenerateRandomString(5 , Helper.CHARACTERS)+ "."+Helper.GetExtension(header.Filename)
	group := ghandler.GroupSer.GetGroupByID(groupid)
	if group== nil{
		res.Message =translation.Translate( lang  , " There is no group by this id : " + groupid)
		response.Write(Helper.MarshalThis(res))
		return
	}
	if group.OwnerID != session.UserID {
		res.Message =translation.Translate( lang  , " UnAuthorized User ... \nYou Are not allowed to Modify this Group")
		response.Write(Helper.MarshalThis(res))
		return
	}
	newImage , ere := os.Create(entity.PathToTemplates + randomName) 
	_ , copyError := io.Copy(newImage  , file)
	if copyError != nil || ere != nil {
		response.WriteHeader(http.StatusInternalServerError)
		res.Message =translation.Translate( lang  , "  INTERNAL SERVER ERROR  ")
		response.Write(Helper.MarshalThis(res))
		return
	}
	defer newImage.Close()
	defer file.Close()
	if group.Imageurl != "" {
		os.Remove(entity.PathToTemplates+ group.Imageurl)
	}
	group.Imageurl= randomName
	group = ghandler.GroupSer.UpdateGroup(group)
	if group == nil {
		response.WriteHeader(http.StatusInternalServerError)
		res.Message =translation.Translate( lang  , " INTERNAL SERVER ERROR ")
		response.Write(Helper.MarshalThis(res))

		fmt.Println("Here.................")
		
		os.Remove(entity.PathToTemplates+ randomName)

		return
	}
	res.Imageurl= randomName
	res.Success = true
	res.Message =translation.Translate( lang  , " Update Was Succesful ")
	response.Write(Helper.MarshalThis(res))
}

// JoinGroup handlign group Joining 
// METHOD  : PUT
// INPUT  : formalue
// OUTPUT : JSON 
// VARIABLES : groupid  , and user id from the session 
// AUTHENTICATION : ONLY LOGGED IN USERS 
func (ghandler *GroupHandler)   JoinGroup( response http.ResponseWriter   , request *http.Request ){
	response.Header().Set("Content-Type"  , "application/json")
	lang := GetSetLang(ghandler  , response  , request )
	res := struct{
		Success bool `json:"success"`
		Message string `json:"message"`
		GroupID string `json:"group_id"`
	}{
		Success:  false ,
		Message:  translation.Translate(lang  , "Invalid Input ..."),
	}
	groupid := request.FormValue("group_id")
	if groupid == "" {
		response.Write(Helper.MarshalThis(res))
		return 
	}
	session := ghandler.Session.GetSession(request)
	if session== nil {
		res.Message=translation.Translate(lang  , " UnAuthorized User ")
		response.Write(Helper.MarshalThis(res))
		return
	}
	res.GroupID= groupid
	group := ghandler.GroupSer.GetGroupByID(groupid)
	user := ghandler.UserSer.GetUserByEmailAndID(session.Email , session.UserID)
	if user == nil {
		res.Message=translation.Translate(lang  , " UnAuthorized User ") 
		ghandler.Session.DeleteSession(response , request)
		response.Write(Helper.MarshalThis(res))
		return
	}
	if group==nil {
		res.Message=translation.Translate(lang  , " No Group IS Registered By this ID ") + groupid
		response.Write(Helper.MarshalThis(res))
		return
	}
	members := group.MembersID
	for _ , id := range members {
		if id == session.UserID {
			res.Message=translation.Translate(lang  , " You Are Already A member ")
			response.Write(Helper.MarshalThis(res))
			return
		}
	}
	group.MembersID= append(group.MembersID, session.UserID)
	group.MembersCount++
	
	group = ghandler.GroupSer.UpdateGroup(group)
	user.MyGroups = append(user.MyGroups, group.ID)
	user = ghandler.UserSer.SaveUser(user)
	if group == nil  || user==nil {
		// fmt.Println("While Saving the Group  "   )
		response.WriteHeader(http.StatusInternalServerError)
		res.Message=translation.Translate(lang  , " Internal Server Error   ")
		response.Write(Helper.MarshalThis(res))
		return
	}
	response.WriteHeader(http.StatusOK)
	res.Success = true
	res.Message= translation.Translate(lang  , " You Joined Group : ") +group.GroupName + translation.Translate(lang  , " Succesfully !") 
	response.Write(Helper.MarshalThis(res))
}


// LeaveGroup  function to leave a group 
// METHOD PUT 
// INPUT  formvalue  group_id 
// RESPONSE : JSON 
// AUTHORIZATION  : ANY USR 
func (ghandler *GroupHandler)   LeaveGroup(response http.ResponseWriter  , request *http.Request) {
	response.Header().Set("Content-Type"  , "application/json")
	lang := GetSetLang(ghandler  , response  ,request)
	session := ghandler.Session.GetSession(request)
	groupID := request.FormValue("group_id")
	res :=struct{
		Success bool `json:"success"`
		Message string `json:"message"`
		GroupID string `json:"group_id"`
	}{
		Success : false ,
		Message: translation.Translate(lang   ,"Invalid Input missing variable ") + "group_id",
	}
	if groupID == ""{
		response.Write(Helper.MarshalThis(res))
		return 
	}
	res.GroupID= groupID
	// user := 
	group := ghandler.GroupSer.GetGroupByID(groupID)
	user := ghandler.UserSer.GetUserByEmailAndID(session.Email , session.UserID)
	if user == nil {
		res.Message=translation.Translate(lang  , " UnAuthorized User ") 
		ghandler.Session.DeleteSession(response , request)
		response.Write(Helper.MarshalThis(res))
		return
	}
	if group == nil {
		res.Message = translation.Translate(lang  , "Group with ID ") + groupID + translation.Translate(lang  ," Not Found ")
		response.Write(Helper.MarshalThis(res))
		return 
	}
	for _  , mem := range group.MembersID {
		if mem== session.UserID {
			group.MembersID = func()[]string{
				val := []string{}
				for _ , v := range group.MembersID {
					if v != session.UserID {
						val = append(val  , v)
					}
				}
				return val 
			}()
			group.MembersCount--
			group = ghandler.GroupSer.UpdateGroup(group)
			
			user.MyGroups = func()[]string {  
				val := []string{}
			for _ , v := range user.MyGroups {
				if v != group.ID {
					val = append(val, v)
				}
			}
			return val 
			}()
			user= ghandler.UserSer.SaveUser(user)
			if group == nil || user == nil {
				res.Message = translation.Translate(lang  , " INTERNAL SERVER ERROR ") 
				response.Write(Helper.MarshalThis(res))
				return 
			}
			res.Success = true 
			res.Message = translation.Translate(lang  , " You have succesfuly left a group ") 
			response.Write(Helper.MarshalThis(res))
			return 
		}
	}
	res.Message =fmt.Sprintf( translation.Translate(lang  , " You are not a member at group  Named : \"%s\" "), group.GroupName )
	response.Write(Helper.MarshalThis(res))
}

// MyGroups returning list of groups uou Are In 
// METHOD GET
// AUTHORIZATION ONLY FOR GROUP MEMBERS 
// RESPONSE : JSON 
func (ghandler *GroupHandler) MyGroups(response http.ResponseWriter  , request *http.Request) {
	response.Header().Set("Content-Type"  , "application/json")
	lang := GetSetLang(ghandler  ,response  , request)
	session := ghandler.Session.GetSession(request)
	res := &struct{
		Success bool `json:"success"`
		Message string `json:"message"`
		Groups []*entity.Group `json:"groups"`
	}{
		Success:  false,
		Message: translation.Translate(lang  , "Invalid Input Body"),
	}
	user := ghandler.UserSer.GetUserByEmailAndID(session.Email  , session.UserID)
	if user == nil {
		res.Message= translation.Translate(lang  , " UnAuthorized User ")
		ghandler.Session.DeleteSession(response  , request)
		response.Write(Helper.MarshalThis(res))
		return 
	}

	// Getting Groups By Ther ID 
	groups := []*entity.Group{}
	for _ , g := range user.MyGroups{
		gr := ghandler.GroupSer.GetGroupByID(g)
		if gr != nil {
			groups = append(groups, gr)
		}
	}
	if len(user.MyGroups)==0 || len(groups)==0 {
		res.Message= translation.Translate(lang  , " You are not a member in any group\nNo Groups Found ")
		response.Write(Helper.MarshalThis(res))
		return
	}
	res.Success=true
	res.Groups = groups
	res.Message = fmt.Sprintf(translation.Translate(lang , "We have succesfuly found ")+" %d  %s" , len(groups)  , func()string { if len(groups)>1 {
		return translation.Translate(lang , " Groups ")
	} 
	return translation.Translate(lang , " Group ") }() )
	response.Write(Helper.MarshalThis(res))
}

// 