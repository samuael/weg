package Admin

import "github.com/samuael/Project/Weg/internal/pkg/entity"

// AdminRepo interface
type AdminRepo interface {
	CreateAdmin(admin *entity.Admin)  (*entity.Admin ,error) 
	GetAdminByID(id string ) (*entity.Admin , error)
	DeleteAdminByID(id string ) error 
	GetAdminByEmail(email string )  (*entity.Admin  , error )
	SaveAdmin( admin *entity.Admin  ) (*entity.Admin , error)
}