package AdminService

import "github.com/samuael/Project/Weg/internal/pkg/Admin"

// AdminService struct representing 
type AdminService struct {
	AdminRepo Admin.AdminRepo
}

// NewAdminService function 
func NewAdminService( repo Admin.AdminRepo ) Admin.AdminService {
	return AdminService{
		AdminRepo : repo  , 
	}
}
