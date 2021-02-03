package User

import "github.com/samuael/Project/Weg/internal/pkg/entity"

// UserRepo interface
type UserRepo interface {
	UserEmailExist(email string ) error
	RegisterUser(user *entity.User) (*entity.User , error)
	GetUserByEmail(Email string ) (*entity.User , error)
	GetUserByEmailAndID(Email  , ID string ) (*entity.User , error)
	SaveUser( usr *entity.User   ) (*entity.User , error)
	GetUserByID(id string ) (*entity.User , error )
	UserWithIDExist(friendID string ) error
	IsGroupMember(userid , groupid string  ) error
	SearchUsers( username string  ) ([]*entity.User  , error )
}
