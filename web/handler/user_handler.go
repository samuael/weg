package handler

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"os"

// 	"github.com/Projects/ScientificNRS/internal/pkg/Admin"
// 	session "github.com/Projects/ScientificNRS/internal/pkg/Session"
// 	"github.com/Projects/ScientificNRS/internal/pkg/entity"
// 	"github.com/Projects/ScientificNRS/pkg/Helper"
// 	// adminse "github.com/Projects/ScientificNRS/internal/pkg/Admin/AdminService"
// )
// // NewAdminHandler  function for creating new Admin Handler
// func NewAdminHandler(
// 	adminservice Admin.AdminService,
// 	sessionHandler *session.Cookiehandler) *AdminHandler {
// 	return &AdminHandler{
// 		AdminService: adminservice,
// 		Session:      sessionHandler,
// 	}
// }
// // AdminHandler struct representing admin handler 
// type AdminHandler struct {
// 	AdminService Admin.AdminService
// 	Session      *session.Cookiehandler
// }

// // RegisterAdmin registering the Admin
// func (adminh *AdminHandler) RegisterAdmin(response http.ResponseWriter, request *http.Request) {

// 	responseJSON := struct {
// 		Success  bool   `json:"success"`
// 		Messaege string `json:"message"`
// 	}{
// 		Success:  false,
// 		Messaege: "Invalid Request Method Not Allowed ",
// 	}

// 	response.Header().Set("Content-Type", "application/json")
// 	if request.Method == http.MethodGet {

// 	} else if request.Method == http.MethodPost {

// 		session := adminh.Session.GetSession(request)

// 		if session == nil || session.UserID == 0 {
// 			responseJSON.Messaege = "UnAuthorized User "
// 			responseJSON.Success = false
// 			jsonReturn, _ := json.Marshal(responseJSON)
// 			response.Write(jsonReturn)
// 			return
// 		}

// 		admin := &entity.Admin{}

// 		recieve := struct {
// 			Username string `json:"username"`
// 			Password string `json:"password"`
// 			Confirm  string `json:"confirm"`
// 		}{}

// 		jsonDecoder := json.NewDecoder(request.Body)
// 		erra := jsonDecoder.Decode(&recieve)
// 		if erra != nil || recieve.Username == "" || recieve.Password == "" {
// 			responseJSON.Messaege = "Invalid Input\n Please Try Again "
// 			jsonReturn, _ := json.Marshal(responseJSON)
// 			response.Write(jsonReturn)
// 			return
// 		}
// 		if recieve.Password != recieve.Confirm {
// 			responseJSON.Messaege = "Invalid Input\n Password Must Match"
// 			jsonReturn, _ := json.Marshal(responseJSON)
// 			response.Write(jsonReturn)
// 			return
// 		}
// 		admin.Username = recieve.Username
// 		admin.Password = recieve.Password
// 		admin = adminh.AdminService.RegisterAdmin(admin)
// 		if admin == nil {
// 			responseJSON.Messaege = "Internal ServerError  "
// 			jsonReturn, _ := json.Marshal(responseJSON)
// 			response.Write(jsonReturn)
// 			return
// 		}
// 		responseJSON.Messaege = " Succesfully Registered User  "
// 		responseJSON.Success = true
// 		jsonReturn, _ := json.Marshal(responseJSON)
// 		response.Write(jsonReturn)
// 		return
// 	}
// }

// // LoginAdmin  login function
// func (adminh *AdminHandler) LoginAdmin(response http.ResponseWriter, request *http.Request) {
// 	response.Header().Add("Content-Type", "application/json")
// 	if request.Method == http.MethodPost {
// 		reciever := struct {
// 			Username string `json:"username"`
// 			Password string `json:"password"`
// 		}{}
// 		responseJSON := struct {
// 			Success bool   `json:"success"`
// 			Message string `json:"message"`
// 		}{
// 			Success: false,
// 			Message: "Invalid Request Body",
// 		}
// 		newDecoder := json.NewDecoder(request.Body)
// 		decodeError := newDecoder.Decode(&reciever)
// 		if decodeError != nil {
// 			fmt.Println(decodeError)
// 			jsonReturn, _ := json.Marshal(responseJSON)
// 			response.Write(jsonReturn)
// 			return
// 		}

// 		admin := adminh.AdminService.GetAdmin(reciever.Username, reciever.Password)

// 		if admin == nil {
// 			responseJSON.Message = " Invalid Username Or Password "
// 			jsonReturn, _ := json.Marshal(responseJSON)
// 			response.Write(jsonReturn)
// 			return
// 		}

// 		session := &entity.Session{
// 			UserID:   admin.ID,
// 			Username: admin.Username,
// 			Password: admin.Password,
// 		}
// 		success := adminh.Session.SaveSession(response, session, entity.PROTOCOL+entity.HOST)

// 		if !success {
// 			responseJSON.Message = "Internal Server Error "
// 			jsonReturn, _ := json.Marshal(responseJSON)
// 			response.Write(jsonReturn)
// 			return
// 		}
// 		responseJSON.Success = true
// 		responseJSON.Message = "Succesfuly Logged In "
// 		jsonReturn, _ := json.Marshal(responseJSON)
// 		response.Write(jsonReturn)
// 		return
// 	}
// }

// // ChangeLanguage function to register language in the cookie header of the User
// // METHOD GET 
// // Variable lang 
// // Response application/json
// // Function Changing the Language and inform the Users about it.
// func (adminh *AdminHandler) ChangeLanguage(response http.ResponseWriter, request *http.Request) {
// 	response.Header().Set("Content-Type", "application/json")
// 	lang := request.FormValue("lang")

// 	data := struct {
// 		Success bool   `json:"success"`
// 		Message string `json:"message"`
// 		Lang    string `json:"lang"`
// 	}{
// 		Success: false,
// 		Message: "Invalid Input ",
// 		Lang:    "",
// 	}

// 	noInput := false
// 	if lang == "" {
// 		lang = "amh"
// 		noInput = true
// 	}
// 	success := adminh.Session.SaveLang(response, entity.Langs[lang], entity.PROTOCOL+entity.HOST)

// 	if success == false {
// 		if noInput == true {
// 			data.Message = "Invalid Input"
// 			data.Lang = "amh"
// 		} else {
// 			data.Message = " Internal Server ERROR "
// 			data.Lang = lang
// 		}
// 		responseJSON, _ := json.Marshal(data)
// 		response.Write(responseJSON)
// 		return
// 	}
// 	data.Success = true
// 	data.Message = "Language Changed \"" + lang + "\""
// 	data.Lang = lang
// 	responseJSON, _ := json.Marshal(data)
// 	response.Write(responseJSON)
// 	return
// }

// // GetLanguages function representing the Languages Supported ny the System
// func (adminh *AdminHandler) GetLanguages(response http.ResponseWriter, request *http.Request) {
// 	response.Header().Set("Content-Type", "application/json")
// 	languages := entity.Languages
// 	responseJSON, _ := json.Marshal(languages)
// 	response.Write(responseJSON)
// 	return
// }

// // ChangeProfile function  to change the P{rofile Picture of the User
// func (adminh *AdminHandler) ChangeProfile(response http.ResponseWriter, request *http.Request) {
// 	// INPUT JSON username  , password  , confirm

// 	response.Header().Set("Content-Type", "appllication/json")
// 	resp := struct {
// 		Success     bool   `json:"success"`
// 		Message     string `json:"message"`
// 		NewUsername string `json:"username"`
// 	}{
// 		Success:     false,
// 		Message:     "Invalid Request MEthod ",
// 		NewUsername: "",
// 	}
// 	if request.Method == http.MethodPost {
// 		session := adminh.Session.GetSession(request)
// 		if session == nil {
// 			resp.Message = "You Are Not Authorized Please Try Again "
// 			responseJSON, _ := json.Marshal(resp)
// 			response.Write(responseJSON)
// 			return
// 		}

// 		object := struct {
// 			Username string `json:"username"`
// 			Password string `json:"password"`
// 			Confirm  string `json:"confirm"`
// 		}{}

// 		jdecode := json.NewDecoder(request.Body)
// 		era := jdecode.Decode(&object)

// 		if era != nil || (object.Confirm != object.Password) || (object.Username == "") {
// 			fmt.Println(era)
// 			resp.Message = "Input ERROR ...... "
// 			responseJSON, _ := json.Marshal(resp)
// 			response.Write(responseJSON)
// 			return
// 		}

// 		if (object.Confirm == "") && (object.Username == session.Username) {
// 			resp.Message = " No CHANGE IS MADE "
// 			responseJSON, _ := json.Marshal(resp)
// 			response.Write(responseJSON)
// 			return
// 		}
// 		admin := &entity.Admin{}
// 		admin.ID = session.UserID
// 		if object.Username == "" || (object.Username == session.Username && session.Password != object.Password) {
// 			admin.Username = session.Username
// 			admin.Password = object.Password
// 		} else if object.Username != session.Username && object.Password != session.Username {
// 			admin.Username = object.Username
// 			admin.Password = object.Password
// 		} else {
// 			admin.Username = object.Username
// 			admin.Password = session.Password
// 		}

// 		admin = adminh.AdminService.SaveAdmin(admin)
// 		if admin == nil {
// 			resp.Message = " Internal Server ERROR "
// 			responseJSON, _ := json.Marshal(resp)
// 			response.Write(responseJSON)
// 			return
// 		}
// 		resp.Message = " Data Edited Successfuly!"
// 		resp.Success = true
// 		resp.NewUsername = admin.Username

// 		session.Username = resp.NewUsername
// 		session.Password = admin.Password
// 		adminh.Session.SaveSession(response, session, entity.PROTOCOL+entity.HOST)
// 		responseJSON, _ := json.Marshal(resp)
// 		response.Write(responseJSON)
// 		return
// 	} else {
// 		responseJSON, _ := json.Marshal(resp)
// 		response.Write(responseJSON)
// 		return
// 	}

// }


// // Logout method 
// // METHOD GET 
// // Variable NONE 
// // RESPONSE JSON
// func (adminh *AdminHandler)  Logout(response http.ResponseWriter  , request *http.Request ){
// 	response.Header().Set("Content-Type" , "application/json")
// 	ret := struct{
// 		Success bool `json:"success"`
// 		Message string `json:"message"`

// 	}{
// 		Success: false,
// 		Message: "Log-Out Was Not Succesful",
// 	}
// 	session := adminh.Session.GetSession(request)
// 	if session == nil {
// 	}else {

// 		success := adminh.Session.DeleteSession(response , request)
// 		if !success {
// 			ret.Message="Can't Log You Out" 
// 		}else {
// 			ret.Message="You Have Logged Out Succesfuly"
// 			ret.Success= true
// 		}
// 	}
// 	jsonReturn  , _ := json.Marshal(ret)
// 	response.Write(jsonReturn )
// }


// // ChangeProfilePicture for uploading and Changing profile Picture  
// // MATHOD POST
// // INPUT FORMDATA 
// // content-Type  : multipart/form-data
// func (adminh  *AdminHandler)  ChangeProfilePicture( response http.ResponseWriter  , request *http.Request ){
// 	//setting  up the headers 
// 	response.Header().Set("Content-Type"  , "application/json")

// 	res:= struct {
// 		Success bool `json:"success"`
// 		Message string `json:"message"`
// 		Imageurl string `json:"imgurl"`
// 	}{
// 		Success : false ,
// 		Message :  "Invalid Request Body " , 
// 		Imageurl : "" , 
// 	}

// 	if request.Method == http.MethodPost {
// 		session := adminh.Session.GetSession(request)
// 		if session == nil {
// 			res.Message="UnAuthorized User "
// 			jsonReturn  , _ := json.Marshal(res)
// 			response.Write(jsonReturn)
// 			return 
// 		}
// 		file  , header  , er := request.FormFile("image")
// 		if er != nil {
// 			res.Message =" Invalid Input / INPUT ERROR "
// 			jsonReturn  , _ := json.Marshal(res)
// 			response.Write(jsonReturn)
// 			return 
// 		}
// 		// validating the extension whether it is an image or not 
// 		if isImage := Helper.IsImage(header.Filename); !isImage {
// 			res.Message =" Invalid Input / Only Image Files are Allowed "
// 			jsonReturn  , _ := json.Marshal(res)
// 			response.Write(jsonReturn)
// 			return 
// 		}

// 		randomName := "img/UserImage/"+Helper.GenerateRandomString(5 , Helper.CHARACTERS)+ header.Filename
// 		newImage , ere := os.Create(entity.PathToTemplates + randomName) 
// 		_ , copyError := io.Copy(newImage  , file)
// 		if copyError != nil || ere != nil {
// 			res.Message =" INTERNAL SERVER ERROR "
// 			response.WriteHeader(http.StatusInternalServerError)
// 			jsonReturn  , _ := json.Marshal(res)
// 			response.Write(jsonReturn)
// 			return
// 		}

// 		admin := adminh.AdminService.GetAdminByID(uint(session.UserID) , session.Username)

// 		admin.Imgurl= randomName
		
// 		admin = adminh.AdminService.SaveAdmin(admin)
// 		if admin == nil {
// 			res.Message="Internal Server Error "
// 			response.WriteHeader(http.StatusInternalServerError)
// 			jsonReturn  , _ := json.Marshal(res)
// 			response.Write(jsonReturn)
// 			return
// 		}

// 		res.Imageurl= entity.PROTOCOL+entity.HOST + admin.Imgurl
// 		res.Message= "Succesfully Registered "
// 		res.Success= true
// 		response.WriteHeader(http.StatusOK)
// 		jsonReturn  , _ := json.Marshal(res)
// 		response.Write(jsonReturn)
// 		return			
// 		// Getting the Admin and chenge the Profile Imgurl variable 
// 	}
// 	jsonReturn  , _ := json.Marshal(res)
// 	response.Write(jsonReturn)
// 	return	
// } 