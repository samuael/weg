package AdminService

import (
	"github.com/samuael/Project/Weg/internal/pkg/Admin"
	"github.com/samuael/Project/Weg/internal/pkg/entity"
)

// AdminService struct representing
type AdminService struct {
	AdminRepo Admin.AdminRepo
}

// NewAdminService function 
func NewAdminService( repo Admin.AdminRepo ) Admin.AdminService {
	return &AdminService{
		AdminRepo : repo  , 
	}
}
// CreateAdmin (admin *entity.Admin)  *entity.Admin
func (adminser *AdminService) CreateAdmin(admin *entity.Admin)  *entity.Admin{
	admin  , era := adminser.AdminRepo.CreateAdmin(admin)
	if era != nil {
		return nil
	}
	return admin 
}

// GetAdminByID () method 
func (adminser *AdminService)  GetAdminByID( id string ) *entity.Admin{
	admin  , era  := adminser.AdminRepo.GetAdminByID(id)
	if era != nil {
		return nil 
	}
	return admin
}

// DeleteAdminByID ( id string  ) bool  
func (adminser *AdminService) DeleteAdminByID( id string  ) bool {
	era := adminser.AdminRepo.DeleteAdminByID(id)
	if era != nil {
		return false // delete
	}
	return true //
}

// GetUserByEmail (email string )  *entity.Admin
func (adminser *AdminService) GetAdminByEmail(email string )  *entity.Admin {
	admin  , era := adminser.AdminRepo.GetAdminByEmail(email )
	if era != nil {
		return nil 
	}
	return admin
} 

// SaveAdmin (admin *entity.Admin)  *entity.Admin
func (adminser  *AdminService ) SaveAdmin(admin *entity.Admin)  *entity.Admin  {
	admin , er := adminser.AdminRepo.SaveAdmin(admin)
	if er != nil {
		return nil 
	}
	return admin 
} 