package service
import "github.com/samuael/Project/Weg/internal/pkg/entity"


type WSGroup struct {
	Group *entity.Group
	ActiveCount int 
	MembersID []string 
}
