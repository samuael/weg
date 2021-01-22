package apiHandler

import (
	"fmt"
	"net/http"

	"github.com/samuael/Project/Weg/internal/pkg/entity"
	"github.com/samuael/Project/Weg/pkg/Helper"
	"github.com/samuael/Project/Weg/pkg/translation"
)

// SearchGroupByName method to return list of groups
// METHOD : GET
// INPUT :
// OUTPUT : JSON  list of groups with less detail
// AUTHENTICATION LOGGING IN
// if the groupname is preceeded by @ sign the following letters will represent the id of the group
// else it's a group name only
func (ghandler *GroupHandler) SearchGroupByName(response http.ResponseWriter  , request *http.Request){
	response.Header().Set("Content-Type"  , "application/json")
	lang := GetSetLang(ghandler , response  , request)
	// session := ghandler.Session.GetSession(request)

	groupname := request.FormValue("group_name")

	res := &struct{
		Success bool `json:"success"`
		Message string `json:"message"`
		Groups []*entity.Group `json:"groups"`
		SearchName string `json:"groupname"`
	}{
		Success : false  , 
		Message : translation.Translate(lang , "Invalid Group Name "),
	}
	if len(groupname)>0 {
		res.Groups = ghandler.GroupSer.SearchGroupsByName(groupname )
		if len(res.Groups) ==0 {
			res.Message= translation.Translate(lang  , " No Record Found ")
			response.Write(Helper.MarshalThis(res))		
			return 
		}
		res.Message = fmt.Sprintf( translation.Translate(lang  , "Succesfully Fetched  %d  %s") , len(res.Groups) , func()string {
			if len(res.Groups) >0 {
				return translation.Translate(lang   , " Groups ")
			}
			return translation.Translate(lang   , " Group ")
		}())
		res.Success = true 
		response.Write(Helper.MarshalThis(res))
		return 
	}
	res.Message= translation.Translate(lang  , " Group Name Can't Be empty String ")
	response.Write(Helper.MarshalThis(res))
}

// GetGroupMembersList function returning Users List inputing group _id 
// METHOD  GET
// INPUT Open Query 
// OUTPUT  :  JOSN 
func (ghandler *GroupHandler)  GetGroupMembersList(response http.ResponseWriter  , request *http.Request){
	response.Header().Set("Content-type"  , "application/json")
	lang := GetSetLang(ghandler  , response  , request)
	// session := ghandler.Session.GetSession(request)
	res := &struct{
		Success bool `json:"success"`
		Message string `json:"message"`
		Members []*entity.User `json:"members"`
		GroupID string `json:"group_id"`
		GroupTitle string `json:"group_title"`
	}{
		Success:  false ,
		Message:  translation.Translate(lang  , "Invalid ID "),
	}
	groupid := request.FormValue("group_id")
	if groupid == "" {
		response.Write(Helper.MarshalThis(res))
		return 
	}
	res.GroupID= groupid
	group := ghandler.GroupSer.GetGroupByID(groupid)
	if group== nil || len(group.MembersID)==0 || group.MembersCount==0 {
		res.Message = translation.Translate(lang  , "Record Not Found ")
		response.Write(Helper.MarshalThis(res))
		return 
	}
	res.GroupTitle = group.GroupName
	members := []*entity.User{}
	for _ , el := range group.MembersID {
		mem := ghandler.UserSer.GetUserByID(el)
		if mem != nil {
			members = append(members, mem)
		}
	}
	if len(members) ==0 {
		res.Message = translation.Translate(lang  , " Record Not Found ")
		response.Write(Helper.MarshalThis(res))
		return 
	}
	res.Message = fmt.Sprintf( translation.Translate(lang  , " Succesfuly Found  %d  %s ")  , len(members)   , Helper.SetPlural(lang  , " Member "  , len(members)) )
	res.Success = true 
	res.Members = members
	response.Write(Helper.MarshalThis(res))	
}