package apiHandler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/samuael/Project/Weg/cmd/service"
	"github.com/samuael/Project/Weg/internal/pkg/Session"
	"github.com/samuael/Project/Weg/internal/pkg/User"
	entity "github.com/samuael/Project/Weg/internal/pkg/entity"
	"github.com/samuael/Project/Weg/pkg/Helper"
	"github.com/samuael/Project/Weg/pkg/translation"
)

// DHandler general name for handlers Sessions
type DHandler interface {
	GetSessionHandler() *session.Cookiehandler
}

// UserHandler  handling
type UserHandler struct {
	SessionHandler *session.Cookiehandler
	UserSer        User.UserService
	ClientService  *service.ClientService
}

// NewUserHandler function returning client handler single client
func NewUserHandler(
	session *session.Cookiehandler, 
	cs User.UserService, 
	clientService *service.ClientService) *UserHandler {
	return &UserHandler{
		SessionHandler: session,
		UserSer:        cs,
		ClientService:  clientService,
	}
}

// GetSessionHandler () *session.Cookiehandler
func (userh *UserHandler) GetSessionHandler() *session.Cookiehandler {
	return userh.SessionHandler
}

// RegisterClient method
// Input JSON :- username  , password  , confirmpassword  , email (UNIQUE)
func (userh *UserHandler) RegisterClient(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	lang := userh.SessionHandler.GetLang(request)
	if lang == "" {
		lang = "en"
		userh.SessionHandler.SaveLang(response, "en", entity.PROTOCOL+entity.HOST)
	}
	input := struct {
		Username        string `json:"username"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmpassword"`
		Email           string `json:"email"`
	}{}
	retu := struct {
		Success  bool   `json:"success"`
		Message  string `json:"message"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}{
		Success: false,
		Message: "Registration Is NOT SUCCESFUl ",
	}
	jsonDecoder := json.NewDecoder(request.Body)
	era := jsonDecoder.Decode(&input)
	if era != nil {
		retu.Message = "Invalid Input"
		response.Write(Helper.MarshalThis(retu))
		return
	}
	// check Password
	if input.Password != input.ConfirmPassword {
		retu.Message = translation.Translate(lang, "Password Confirmation Failed \nPlease Enter Valid Password ")
		response.Write(Helper.MarshalThis(retu))
		return
	}
	if Helper.ValidateUsername(input.Username, 3) {
		retu.Message = translation.Translate(lang, " Invalid Username Value")
		response.Write(Helper.MarshalThis(retu))
		return
	}
	if len(input.Password) < 4 {
		retu.Message = translation.Translate(lang, " Invalid Password characters length should be greater than 4")
		response.Write(Helper.MarshalThis(retu))
		return
	}
	// finding having this email if so inform that the email is reserved nigga
	if !Helper.MatchesPattern(input.Email, Helper.EmailRX) {
		retu.Message = translation.Translate(lang, "Invalid Email Pattern ")
		response.Write(Helper.MarshalThis(retu))
		return
	}
	exists := userh.UserSer.UserEmailExist(input.Email)
	if exists {
		retu.Message = translation.Translate(lang, "Email Already Registered ... ")
		response.Write(Helper.MarshalThis(retu))
		return
	}
	//Hashing the Password
	pwd, era := Helper.HashPassword(input.Password)
	if era != nil {
		retu.Message = translation.Translate(lang, " Invalid Password Characters ... ")
		response.Write(Helper.MarshalThis(retu))
		return
	}
	user := &entity.User{
		Username: input.Username,
		Password: pwd,
		Email:    input.Email,
		LastSeen: time.Now(),
	}
	user = userh.UserSer.RegisterUser(user)
	if user == nil {
		response.WriteHeader(http.StatusInternalServerError)
		retu.Message = translation.Translate(lang, " INTERNAL SERVER ERROR ")
		response.Write(Helper.MarshalThis(retu))
		return
	}
	response.WriteHeader(http.StatusOK)
	retu.Success = true
	retu.Username = user.Username
	retu.Email = user.Email
	retu.Message = translation.Translate(lang, "Registration Was Succesful ")
	response.Write(Helper.MarshalThis(retu))
}

// UpdateUserProfile function to update the profile of the requesting Body
// METHOD PUT
// AUTHORIZED ONLY the LOGGED IN USER
// USERID from session
// INPUT  : JSON
// OUTPUT : JSON
// EMAIl can not be changed
// METHOD  : PUT
func (userh *UserHandler) UpdateUserProfile(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	lang := GetSetLang(userh, response, request)
	session := userh.SessionHandler.GetSession(request)
	in := &struct {
		Username string `json:"username,omitempty"`
		Bio      string `json:"bio,omitempty"`
	}{}
	ret := &struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		User    *entity.User
	}{
		Success: false,
		Message: translation.Translate(lang, "Invalid Input "),
	}

	jdec := json.NewDecoder(request.Body)
	decodeError := jdec.Decode(in)
	if decodeError != nil || (in.Username == "" && in.Bio == "") {
		response.Write(Helper.MarshalThis(ret))
		return
	}
	user := userh.UserSer.GetUserByEmailAndID(session.Email, session.UserID)
	if user == nil {
		ret.Message = translation.Translate(lang, "UnAuthorized User ")
		response.Write(Helper.MarshalThis(ret))
		return
	}
	if in.Username == user.Username && in.Bio == user.Bio {
		ret.Message = translation.Translate(lang, "No Change has made")
		response.Write(Helper.MarshalThis(ret))
		return
	}
	var chage int
	user.Username = func() string {
		if (in.Username != "") && (in.Username != user.Username) && (len(in.Username) > 3) {
			chage = 1
			return in.Username
		}
		return user.Username
	}()
	user.Bio = func() string {
		if (in.Bio != "") && in.Bio != user.Bio {
			if chage == 1 {
				chage = 4 // meaninng both
			} else {
				chage = 2
			}
			return in.Bio
		}
		return user.Bio
	}()
	if chage == 4 || chage == 1 || chage == 2 {
		user.LastUpdated = time.Now()
		user = userh.UserSer.SaveUser(user)
		if user == nil {
			ret.Message = translation.Translate(lang, "Internal Server ERROR ")
			response.Write(Helper.MarshalThis(ret))
			return
		}

		// The message is succesful and now i am going to broadcast
		// the change to all Alies of   profile owner.
		userh.ClientService.Message <- &entity.AlieProfile{
			Status: entity.MsgAlieProfileChange,
			Body: *user,
			SenderID: user.ID,
		}

		ret.Success = true
		user.Password = ""
		ret.User = user
		ret.Message = translation.Translate(lang, "Profile CHanged Succesfuly ")
		response.Write(Helper.MarshalThis(ret))
		return
	}
	response.Write(Helper.MarshalThis(ret))
}

// ChangeUserPassword functionn only for Users
// METHOD : PUT
// AUTHENTICATION NEEDED
// INPUT : JSON
// OUTPUT : JSON
func (userh *UserHandler) ChangeUserPassword(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	lang := GetSetLang(userh, response, request)
	session := userh.SessionHandler.GetSession(request)

	in := &struct {
		OldPassword     string `json:"password"`
		NewPassword     string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"`
	}{}

	res := &struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{
		Success: false,
		Message: translation.Translate(lang, "Invalid Input "),
	}

	jdec := json.NewDecoder(request.Body)
	er := jdec.Decode(in)
	if er != nil || in.OldPassword == "" || in.NewPassword == "" || in.ConfirmPassword == "" {
		response.Write(Helper.MarshalThis(res))
		return
	} else if in.ConfirmPassword != in.NewPassword {
		res.Message = translation.Translate(lang, "Incompatible Password Input Confirm the Password Again")
		response.Write(Helper.MarshalThis(res))
		return
	}
	user := userh.UserSer.GetUserByEmailAndID(session.Email, session.UserID)
	if user == nil {
		res.Message = translation.Translate(lang, "UnAuthorized user Please Login in again ")
		userh.SessionHandler.DeleteSession(response, request)
		response.Write(Helper.MarshalThis(res))
		return
	}
	if len(in.NewPassword) < 4 {
		res.Message = translation.Translate(lang, "Password Must be greater than 4 ")
		response.Write(Helper.MarshalThis(res))
		return
	}
	newpassword, err := Helper.HashPassword(in.NewPassword)

	if !Helper.ComparePassword(user.Password, in.OldPassword) {
		res.Message = translation.Translate(lang, "INCORRECT Password TRY AGAIN ")
		// userh.SessionHandler.DeleteSession(response , request)
		response.Write(Helper.MarshalThis(res))
		return
	} else if err != nil {
		res.Message = translation.Translate(lang, "Invalid Password TRY AGAIN ")
		// userh.SessionHandler.DeleteSession(response , request)
		response.Write(Helper.MarshalThis(res))
		return
	}

	user.Password = newpassword
	user.LastUpdated = time.Now()
	user = userh.UserSer.SaveUser(user)
	if user == nil {
		response.WriteHeader(http.StatusInternalServerError)
		res.Message = translation.Translate(lang, " INTERNAL SERVER ERROR ")
		response.Write(Helper.MarshalThis(res))
		return
	}
	response.WriteHeader(http.StatusOK)
	res.Success = true
	res.Message = translation.Translate(lang, "Password Changed Succesfuly ")
	response.Write(Helper.MarshalThis(res))
}

// GetSetLang get Language from request header and if it's null set the language
func GetSetLang(userh DHandler, response http.ResponseWriter, request *http.Request) string {
	lang := userh.GetSessionHandler().GetLang(request)
	if lang == "" {
		lang = "en"
		userh.GetSessionHandler().SaveLang(response, "en", entity.PROTOCOL+entity.HOST)
	}
	return lang
}

// Login api Login function
// METHOD POST
// INPUT JSON    {
// email :
// password :
// }
func (userh *UserHandler) Login(response http.ResponseWriter, request *http.Request) {
	lang := GetSetLang(userh, response, request)
	response.Header().Set("Content-Type", "application/json")
	ret := struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{
		Success: false,
		Message: translation.Translate(lang, " INVALID INPUT "),
	}
	reciver := &struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	jdecoder := json.NewDecoder(request.Body)
	derorr := jdecoder.Decode(reciver)
	if derorr != nil {
		response.Write(Helper.MarshalThis(ret))
		return
	}
	// check whether the email is valid or not
	reciver.Email = strings.Trim(reciver.Email, " ")
	if !Helper.MatchesPattern(reciver.Email, Helper.EmailRX) {
		ret.Message = translation.Translate(lang, " Invalid Email Input... ")
		response.Write(Helper.MarshalThis(ret))
		return
	} else if !Helper.ValidatePassword(reciver.Password, 4) {
		ret.Message = translation.Translate(lang, " Invalid Password Input... ")
		response.Write(Helper.MarshalThis(ret))
		return
	}
	exist := userh.UserSer.UserEmailExist(reciver.Email)
	if exist {
		user := userh.UserSer.GetUserByEmail(reciver.Email)
		if user == nil {
			response.WriteHeader(http.StatusInternalServerError)
			ret.Message = translation.Translate(lang, " INTERNAL SERVER ERROR ")
			response.Write(Helper.MarshalThis(ret))
			return
		}
		// check the Password
		if !Helper.ComparePassword(user.Password, reciver.Password) {
		} else {
			if success := userh.SessionHandler.SaveSession(response, &entity.Session{
				UserID:   user.ID,
				Username: user.Username,
				Email:    user.Email,
			}, entity.PROTOCOL+entity.HOST); !success {
				ret.Message = translation.Translate(lang, " Internal Server Error : SESSION ERROR ")
				response.Write(Helper.MarshalThis(ret))
			} else {
				ret.Success = true
				ret.Message = translation.Translate(lang, " Loging in was successful ")
				response.Write(Helper.MarshalThis(ret))
			}
			return
		}
	}
	ret.Message = translation.Translate(lang, " Invalid Username Or Password ")
	response.Write(Helper.MarshalThis(ret))
}

// Logout function api Logging out
// METHOD GET
// VAriables NONE
func (userh *UserHandler) Logout(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	lang := GetSetLang(userh, response, request)
	ret := struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{}
	success := userh.SessionHandler.DeleteSession(response, request)
	if !success {
		ret.Success = false
		ret.Message = translation.Translate(lang, " Logging out was not succesful")
		response.Write(Helper.MarshalThis(ret))
		return
	}
	ret.Success = true
	ret.Message = translation.Translate(lang, " Logging out was succesful")
	response.Write(Helper.MarshalThis(ret))
}

// UploadProfilePic function for user to update their profile Picture
// METHOD POST
// BODY : multipart data
// Variable image :
// INPUT >> Content-Type : application/x-www-form-urlencoded
func (userh *UserHandler) UploadProfilePic(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	lang := GetSetLang(userh, response, request)
	res := struct {
		Success  bool   `json:"success"`
		Message  string `json:"message"`
		Imageurl string `json:"imgurl"`
	}{
		Success: false,
	}
	session := userh.SessionHandler.GetSession(request)
	if session == nil {
		response.WriteHeader(http.StatusUnauthorized)
		res.Message = translation.Translate(lang, " UNAuthorized User")
		response.Write(Helper.MarshalThis(res))
		return
	}
	ear := request.ParseForm()
	if ear != nil {
		res.Message = translation.Translate(lang, " Invalid Input Parse Error ")
		response.Write(Helper.MarshalThis(res))
		return
	}
	file, header, er := request.FormFile("image")
	if er != nil {
		// fmt.Println(er)
		res.Message = translation.Translate(lang, " Invalid Input ")
		response.Write(Helper.MarshalThis(res))
		return
	}
	defer file.Close()
	// validating the extension whether it is an image or not
	if isImage := Helper.IsImage(header.Filename); !isImage {
		res.Message = translation.Translate(lang, " Invalid Input / Only Image Files are Allowed ... ")
		response.Write(Helper.MarshalThis(res))
		return
	}
	randomName := "assets/img/UserImage/" + Helper.GenerateRandomString(5, Helper.CHARACTERS) + "." + Helper.GetExtension(header.Filename)
	user := userh.UserSer.GetUserByEmailAndID(session.Email, session.UserID)
	if user == nil {
		response.WriteHeader(http.StatusInternalServerError)
		res.Message = translation.Translate(lang, "User Doesn't Exist Any More ")
		response.Write(Helper.MarshalThis(res))
		return
	}
	newImage, ere := os.Create(entity.PathToTemplates + randomName)
	defer newImage.Close()
	_, copyError := io.Copy(newImage, file)
	if copyError != nil || ere != nil {
		response.WriteHeader(http.StatusInternalServerError)
		res.Message = translation.Translate(lang, " INTERNAL SERVER ERROR ")
		response.Write(Helper.MarshalThis(res))
		return
	}
	if user.Imageurl != "" {
		os.Remove(entity.PathToTemplates + user.Imageurl)
	}
	user.Imageurl = randomName
	user.LastUpdated = time.Now()
	// saving the User
	user = userh.UserSer.SaveUser(user)
	if user == nil {
		response.WriteHeader(http.StatusInternalServerError)
		res.Message = translation.Translate(lang, " Internal Server Error ...")
		response.Write(Helper.MarshalThis(res))
		return
	}

	// The message is succesful and now i am going to broadcast
	// the change to all Alies of   profile owner.
	userh.ClientService.Message <- &entity.AlieProfile{
		Status: entity.MsgAlieProfileChange,
		Body: *user,
		SenderID: user.ID,
	}
	
	res.Success = true
	response.WriteHeader(http.StatusOK)
	res.Message = translation.Translate(lang, "Succesfully Registered ")
	res.Imageurl = user.Imageurl
	response.Write(Helper.MarshalThis(res))
}

// ChangeLanguage function for changing the language
// METHOD GET
// VARIABLE lang
// AUTHORIZATION
func (userh *UserHandler) ChangeLanguage(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	lang := request.FormValue("lang")
	ret := struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Lang    string `json:"lang"`
	}{
		Success: false,
		Message: "Invalid Input",
	}
	if lang == "" {
		response.Write(Helper.MarshalThis(ret))
		return
	}
	userh.SessionHandler.SaveLang(response, Helper.GetAppropriateLangLabel(lang), entity.PROTOCOL+entity.HOST)
	ret.Success = true
	ret.Message = translation.Translate(lang, "Language Succesfully Changed ")
	ret.Lang = Helper.GetAppropriateLangLabel(lang)
	response.Write(Helper.MarshalThis(ret))
	return
}

// ----------------------------------------------------------------------------------------------------------------------------------------------

// LoggedIn checks whether the user is Authenticated or not
func (userh *UserHandler) LoggedIn(request *http.Request) bool {
	session := userh.SessionHandler.GetSession(request)
	if session != nil {
		return true
	}
	return false
}

// Authenticated checks if a user has proper authority to access a give route
func (userh *UserHandler) Authenticated(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !userh.LoggedIn(r) {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		session := userh.SessionHandler.GetSession(r)
		if session == nil {

			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}


// MyProfile returning the profile of the user by reading the Sessio 
func (userh *UserHandler )  MyProfile(response http.ResponseWriter  , request *http.Request ) {
	session := userh.SessionHandler.GetSession(request)
	lang := GetSetLang(userh , response , request )
	response.Header().Set("Content-Type", "application/json" )
	res := &struct{
		Success bool `json:"success"`
		Message string `json:"message"`
		User *entity.User  `json:"user"`
	}{
		Success : false  , 
		Message : translation.Translate(lang , "UnAuthorized User ") , 
		User : nil  , 
	}
	if session == nil {
		response.Write(Helper.MarshalThis(res))
		return 
	}
	user := userh.UserSer.GetUserByID(session.UserID)
	if user == nil {
		res.Message= translation.Translate(lang , "Record Not Found ")
		response.Write(Helper.MarshalThis(res))
		return 
	}
	res.Success = true 
	res.Message = fmt.Sprintf(translation.Translate(lang  , " Welcome %s ") , user.Username )
	res.User = user 
	response.Write(Helper.MarshalThis(res))
}

// SearchUsers using name 
func (userh *UserHandler )  SearchUsers(response http.ResponseWriter  , request *http.Request ) {
	// getting the username and validation of the search query string 
	lang := GetSetLang(userh , response  , request )
	response.Header().Set("Content-Type", "application/json")
	res := &struct{
		Success bool `json:"success"`
		Message string `json:"message"`
		Users []*entity.User `json:"users"`
	}{
		Success : false , 
		Message : translation.Translate(lang , "No Record Found "),
	}
	username := request.FormValue("username")
	if username == "" {
		response.Write(Helper.MarshalThis(res))
		return 
	}
	users := userh.UserSer.SearchUsers(username)
	if users == nil {
		res.Message= translation.Translate(lang   , " No Search result found ")
		response.Write(Helper.MarshalThis(res))
		return 
	}
	res.Success = true
	res.Message = fmt.Sprintf(translation.Translate( lang , "Succesfult %d record" )  , len(users))
	res.Users = users
	response.Write(Helper.MarshalThis(res))
}