package Group

import "github.com/samuael/Project/Weg/internal/pkg/entity"

type GroupService interface {
	CreateGroup(group *entity.Group) *entity.Group
	DeleteGroup( ID string  ) bool
	GetGroupByID(ID string )  (*entity.Group)
	UpdateGroup(group *entity.Group) *entity.Group
	IncrementMembersCount( groupID string  ) bool
	DecrementMembersCount( groupID string  ) bool
	SearchGroupsByName(groupname string ) ([]*entity.Group )
	DoesGroupExist(groupID string ) bool 
	IsGroupMember(groupID , memberID string) bool 
}