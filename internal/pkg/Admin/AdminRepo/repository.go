package AdminRepo

import (
	"github.com/samuael/Project/Weg/internal/pkg/Admin"
	"go.mongodb.org/mongo-driver/mongo"
)

// AdminRepo  struct representing mongodb Admin Repository Class
type AdminRepo struct {
	DB *mongo.Database
}

// NewAdminRepo function returning AdminRepo 
func NewAdminRepo(  db *mongo.Database ) Admin.AdminRepo {
	return &AdminRepo{
		DB : db  , 
	}
}
