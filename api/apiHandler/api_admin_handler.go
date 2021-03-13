package apiHandler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/samuael/Project/Weg/internal/pkg/Admin"
	session "github.com/samuael/Project/Weg/internal/pkg/Session"
	"github.com/samuael/Project/Weg/internal/pkg/User"
	"github.com/samuael/Project/Weg/internal/pkg/entity"
	"github.com/samuael/Project/Weg/pkg/Helper"
	"github.com/samuael/Project/Weg/pkg/translation"
)

// AdminHandler struct
type AdminHandler struct {
	Session  *session.Cookiehandler
	AdminSer Admin.AdminService
	UserSer  User.UserService
}

// NewAdminHandler function returning a new AdminHandler
func NewAdminHandler(
	session *session.Cookiehandler,
	admiser Admin.AdminService,
	userser User.UserService) *AdminHandler {
	return &AdminHandler{
		Session:  session,
		AdminSer: admiser,
		UserSer:  userser,
	}
}

// GetSessionHandler () *session.Cookiehandler
func (adminh *AdminHandler) GetSessionHandler() *session.Cookiehandler {
	return adminh.Session
}

// CreateAdmin creating an admin using a json only
// INPUT {
//   username : "username "
//     email : abebe@email.com
//   password : "password "
//   confirmPassword : "confirm_password"
//
// }
func (adminh *AdminHandler) CreateAdmin(response http.ResponseWriter, request *http.Request) {
	session := adminh.Session.GetSession(request)
	response.Header().Set("Content-Type", "application/json")
	res := &struct {
		Success bool          `json:"success"`
		Message string        `json:"message"`
		Admin   *entity.Admin `json:"new_admin"`
	}{
		Success: false,
	}

	input := &struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
		Confirm  string `json:"confirm_password"`
	}{}

	jdecoder := json.NewDecoder(request.Body)

	// if session.Role != entity.ADMIN {
	// 	res.Message = "unauthorized user "
	// 	response.Write(Helper.MarshalThis(res))
	// 	return
	// }

	era := jdecoder.Decode(input)
	if era != nil || input.Username == "" || input.Password == "" || input.Confirm == "" || input.Email == "" {
		res.Message = "Invalid Input Please Try again "
		response.Write(Helper.MarshalThis(res))
		return
	}
	if admin := adminh.AdminSer.GetAdminByEmail(input.Email); admin != nil {
		res.Message = "Already Registered Email\n Please Use Another Email"
		response.Write(Helper.MarshalThis(res))
		return
	}
	if input.Password != input.Confirm {
		res.Message = "Password Must Match Please Enter Valid Password "
		response.Write(Helper.MarshalThis(res))
		return
	}
	admin := &entity.Admin{
		Username:  input.Username,
		Password:  input.Password,
		CreatedBy: session.UserID,
		Email:     input.Email,
	}

	admin = adminh.AdminSer.CreateAdmin(admin)
	if admin == nil {
		res.Message = " Internal Server Error "
		response.Write(Helper.MarshalThis(res))
		return
	}

	res.Message = "Admin Created Succefuly "
	res.Admin = admin
	res.Success = true
	response.Write(Helper.MarshalThis(res))
}

// DeleteAdmin method to delete an admin
func (adminh *AdminHandler) DeleteAdmin(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	session := adminh.Session.GetSession(request)
	if session == nil || session.Role != entity.ADMIN {
		response.Write([]byte("UnAuthorized User "))
	}
	res := &struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		AdminID string `json:"admin_id"`
	}{
		Success: false,
	}
	if adminid := request.FormValue("admin_id"); adminid != "" {
		res.AdminID = adminid
		admin := adminh.AdminSer.GetAdminByID(adminid)
		if admin == nil {
			res.Message = " No Admin By this Id "
			response.Write(Helper.MarshalThis(res))
			return
		}
		if admin.CreatedBy != session.UserID {
			res.Message = "You are not authorized to Delete this User "
			response.Write(Helper.MarshalThis(res))
			return
		}
		success := adminh.AdminSer.DeleteAdminByID(adminid)
		if success {
			res.Message = " Admin Deleted Succesfuly "
			res.Success = true
			response.Write(Helper.MarshalThis(res))
			return
		}
		res.Message = "Internal Server Error .. "
	}
	res.Message = " Invalid Input "
	response.Write(Helper.MarshalThis(res))

}

// DeleteUser   method to delete an admin
func (adminh *AdminHandler) DeleteUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	session := adminh.Session.GetSession(request)
	if session == nil {
		return
	}
	res := &struct {
		Success bool
		Message string
		UserID  string
	}{
		Success: false,
	}
	if session.Role != entity.ADMIN {
		res.Message = " UnAuthorized User "
		res.Success = false
		response.Write(Helper.MarshalThis(res))
		return
	}
	userid := request.FormValue("user_id")
	res.UserID = userid
	user := adminh.UserSer.GetUserByID(userid)
	if user == nil {
		res.Message = " Record Not Found .."
		response.Write(Helper.MarshalThis(res))
		return
	}
	success := adminh.UserSer.DeleteUserByID(userid)
	if !success {
		res.Message = "Internal Server Error "
		response.Write(Helper.MarshalThis(res))
		return
	}
	res.Message = " User Deleted Succesfuly "
	res.Success = true
	response.Write(Helper.MarshalThis(res))
}

// UpdateAdmin method to update the admin account
func (adminh *AdminHandler) UpdateAdmin(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	session := adminh.Session.GetSession(request)
	if session == nil {
		return
	}
	res := &struct {
		Message string
		Success bool
		Admin   *entity.Admin
	}{
		Success: false,
	}
	input := &struct {
		Username string `json:"username"`
		Password string `json:"password"`
		// Email 	 string `json:"email"`
		Confirm string `json:"confirm_password"`
	}{}

	jdecoder := json.NewDecoder(request.Body)
	era := jdecoder.Decode(input)
	if era != nil || input.Username == "" || input.Password == "" || input.Confirm == "" {
		res.Message = "Invalid Input Please Try again "
		response.Write(Helper.MarshalThis(res))
		return
	}
	admin := adminh.AdminSer.GetAdminByID(session.UserID)
	if admin == nil {
		res.Message = " Internal Server Error  "
		adminh.Session.DeleteSession(response, request)
		response.Write(Helper.MarshalThis(res))
		return
	}
	admin.Username = func() string {
		if input.Username != "" && input.Username != admin.Username {
			return input.Username
		}
		return admin.Username
	}()

	admin.Password = func() string {
		if input.Password != "" && input.Password == input.Confirm {
			return input.Password
		}
		return admin.Password
	}()

	admin = adminh.AdminSer.SaveAdmin(admin)
	if admin == nil {
		res.Message = " Internal Server Error "
		response.Write(Helper.MarshalThis(res))
		return
	}

	res.Message = "Admin Updated Succefuly "
	res.Admin = admin
	res.Success = true
	response.Write(Helper.MarshalThis(res))
}

// AdminLogin method
func (adminh *AdminHandler) AdminLogin(response http.ResponseWriter, request *http.Request) {
	lang := GetSetLang(adminh, response, request)
	response.Header().Set("Content-Type", "application/json")
	ret := struct {
		Success bool          `json:"success"`
		Message string        `json:"message"`
		Admin   *entity.Admin `json:"admin"`
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
	exist := adminh.AdminSer.AdminEmailExist(reciver.Email)
	if exist {
		user := adminh.AdminSer.GetAdminByEmail(reciver.Email)
		if user == nil {
			response.WriteHeader(http.StatusInternalServerError)
			ret.Message = translation.Translate(lang, " INTERNAL SERVER ERROR ")
			response.Write(Helper.MarshalThis(ret))
			return
		}
		// check the Password
		if !(user.Password == reciver.Password) {
		} else {
			if success := adminh.Session.SaveSession(response, &entity.Session{
				UserID:   user.ID,
				Username: user.Username,
				Email:    user.Email,
				Role:     entity.ADMIN,
			}, entity.PROTOCOL+entity.HOST); !success {
				ret.Message = translation.Translate(lang, " Internal Server Error : SESSION ERROR ")
				response.Write(Helper.MarshalThis(ret))
			} else {
				ret.Success = true
				ret.Admin = user
				ret.Message = translation.Translate(lang, " Loging in was successful ")
				response.Write(Helper.MarshalThis(ret))
			}
			return
		}
	}
	ret.Message = translation.Translate(lang, " Invalid Username Or Password ")
	response.Write(Helper.MarshalThis(ret))
}

// DeleteIdea
// func (adminh *admin)

// SearchIdea
