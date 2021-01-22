// GroupService package name representing
package GroupService

import (
	
	"github.com/samuael/Project/Weg/internal/pkg/Group"
	"github.com/samuael/Project/Weg/internal/pkg/entity"
)

// GroupService struct representing the client Groups Service
type GroupService struct {
	GroupRepo Group.GroupRepo
}
// NewGroupService function returning Group Service pointer instance 
func NewGroupService(  repo Group.GroupRepo ) Group.GroupService {
	return &GroupService{
		GroupRepo : repo, 
	}
}
// CreateGroup handling group Service 
func (groupser *GroupService )  CreateGroup(group *entity.Group) *entity.Group{
	group  , era := groupser.GroupRepo.CreatGroup(group)
	if era != nil {
		return nil 
	}
	return group
}

// DeleteGroup function returning true if succed else return false
func (groupser *GroupService) DeleteGroup( ID  string ) bool  {
	era := groupser.GroupRepo.DeleteGroup(ID )
	if era != nil {
		return false
	}
	return true 
}

// GetGroupByID service function 
func (groupser *GroupService) GetGroupByID(ID string )  (*entity.Group){
	group , er := groupser.GroupRepo.GetGroupByID(ID)
	if er != nil {
		return nil 
	}
	return group
}


// UpdateGroup (group *entity.Group) *entity.Group
func (groupser *GroupService) UpdateGroup(group *entity.Group) *entity.Group{
	group  , er := groupser.GroupRepo.UpdateGroup(group )
	if er != nil {
		return nil 
	} 
	return group
}

// IncrementMembersCount service function
func (groupser *GroupService) IncrementMembersCount( groupID string  ) bool {
	er := groupser.GroupRepo.IncrementMembersCount(groupID)
	if er != nil {
		return false
	}
	return true
}

// DecrementMembersCount service function
func (groupser *GroupService) DecrementMembersCount( groupID string  ) bool {
	er := groupser.GroupRepo.IncrementMembersCount(groupID)
	if er != nil {
		return false
	}
	return true
}
// SearchGroupsByName (id string ) ([]*entity.Group )
func (groupser *GroupService)  SearchGroupsByName(id string ) ([]*entity.Group ){
	grps  , er := groupser.GroupRepo.SearchGroupsByName(id)
	if  er != nil {
		return nil 
	}
	return grps
}

// DoesGroupExist (groupID string ) bool 
func (groupser *GroupService) DoesGroupExist(groupID string ) bool {
	er := groupser.GroupRepo.DoesGroupExist(groupID )
	if er != nil {
		return false 
	}  
	return true 
}

// IsGroupMember (groupID string) bool
func (groupser *GroupService) IsGroupMember(groupID , memberID string) bool {
	er := groupser.GroupRepo.IsGroupMember(groupID , memberID )
	if er != nil {
		return false
	} 
	return true
} 