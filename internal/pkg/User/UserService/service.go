package UserService

import (
	entity "github.com/samuael/Project/Weg/internal/pkg/entity"
	"github.com/samuael/Project/Weg/internal/pkg/User"
)

// UserService struct representing client's service
type UserService struct {
	UserRepo User.UserRepo
}

// NewUserService function 
func NewUserService(  repo User.UserRepo  ) User.UserService {
	return &UserService{
		UserRepo : repo  , 
	}
}

// UserEmailExist function checking whether the email exists or not 
// if the Return value is true , it means there is user registered by this email
func (userser *UserService ) UserEmailExist(email string ) bool {
	er := userser.UserRepo.UserEmailExist(email)
	if er != nil {
		return false
	}
	return true
}

// RegisterUser function 
func (userser *UserService) RegisterUser(user *entity.User) *entity.User {
	user  , er := userser.UserRepo.RegisterUser(user)
	if er != nil {
		return nil 
	}
	return user
} 

// GetUserByEmail (Email string ) *entity.User
func (userser *UserService) GetUserByEmail(Email string ) *entity.User {
	user  , era := userser.UserRepo.GetUserByEmail(Email)
	if era != nil {
		return nil 
	}
	return user
}
// GetUserByEmailAndID function 
func  (userser *UserService) GetUserByEmailAndID(Email  , ID string ) *entity.User {
	user  , era := userser.UserRepo.GetUserByEmailAndID(Email , ID)
	if era != nil {
		return nil 
	}
	return user
}

// SaveUser ( *entity.User   ) *entity.User
func (userser *UserService) SaveUser( usr *entity.User ) *entity.User {
	user  , era := userser.UserRepo.SaveUser(usr)
	if era != nil {
		return nil 
	}
	return user
}
// GetUserByID (id string ) *entity.User
func (userser *UserService)  GetUserByID(id string ) *entity.User {
	user  , er :=userser.UserRepo.GetUserByID(id)
	if er != nil {
		return nil 
	}
	return user
}

// UserWithIDExist (friendID string ) error
func (userser *UserService) UserWithIDExist(friendID string ) bool {
	er := userser.UserRepo.UserWithIDExist(friendID )
	if er != nil {
		return false 
	}
	return true 
}

// IsGroupMember returning whether the use is a member or not 
// return value error if the user is not a member of that group
// otherwise it returns nil
func (userser  *UserService) IsGroupMember(userid , groupid string  ) bool {
	er := userser.UserRepo.IsGroupMember(userid  , groupid )
	if er != nil {
		return false 
	}
	return true 
}