package User

import "github.com/samuael/Project/Weg/internal/pkg/entity"

// UserService interface
type UserService interface  {
	UserEmailExist(email string ) bool 
	RegisterUser(user *entity.User) *entity.User
	GetUserByEmail(Email string ) *entity.User
	GetUserByEmailAndID(Email  , ID string ) *entity.User
	SaveUser( *entity.User   ) *entity.User
	GetUserByID(id string ) *entity.User 
	UserWithIDExist(friendID string ) bool 
	IsGroupMember(userid , groupid string  ) bool 
	SearchUsers( username string  ) []*entity.User 
	DeleteUserByID(id string ) bool
}