package Group

import "github.com/samuael/Project/Weg/internal/pkg/entity"


// GroupRepo interface representing clients and interfaces
type GroupRepo interface {
	CreatGroup(group *entity.Group ) (*entity.Group ,error)
	DeleteGroup( ID string  ) error
	GetGroupByID(ID string )  (*entity.Group , error)
	UpdateGroup(group *entity.Group) (*entity.Group , error)
	DecrementMembersCount(groupID string ) error
	IncrementMembersCount( groupID string  ) error
	SearchGroupsByName(groupname string ) ([]*entity.Group , error)
	// SearchGroupsByID(id string ) ([]*entity.Group , error)
	DoesGroupExist(groupID string ) error
	IsGroupMember(groupID , memberID string) error
}