package IdeaService

import (
	"github.com/samuael/Project/Weg/internal/pkg/Idea"
	"github.com/samuael/Project/Weg/internal/pkg/entity"
)

// IdeaService representing the Idea Service handler
type IdeaService struct {
	IdeaRepo Idea.IdeaRepo 
}

// NewIdeaService function returning an Idea Service Instance 
func NewIdeaService(  repo Idea.IdeaRepo) Idea.IdeaService {
	return &IdeaService{
		IdeaRepo: repo,
	}
}

// CreateIdea function for creating an Idea Instance 
func (ideaser *IdeaService )   CreateIdea(idea  *entity.Idea ) *entity.Idea {
	ideas , era := ideaser.IdeaRepo.CreateIdea(idea )
	if era != nil {
		return nil 
	}
	return ideas 
} 

// GetIdeas returning ideas depending on the offset and the limit of the Ideas count 
func (ideaser *IdeaService)  GetIdeas(  userid string   , offset , limit  int )( []*entity.Idea    ){
	ideas  , err := ideaser.IdeaRepo.GetIdeas(userid  , offset , limit )
	if err != nil{
		return nil 
	}
	return ideas 
}

// GetIdeaByID ( id string  )(*entity.Idea   )
func (ideaser *IdeaService) GetIdeaByID( id string  )*entity.Idea  {
	idea  , era := ideaser.IdeaRepo.GetIdeaByID( id )
	if era != nil {
		return idea 
	}
	return idea 
} 

// GetMyIdeas returning your ideas this one is the service 
func (ideaser *IdeaService )  GetMyIdeas( userid string    )([]*entity.Idea  ){
	ideas  , era := ideaser.IdeaRepo.GetMyIdeas( userid )
	if era != nil {
		return nil 
	}
	return ideas 
}

// DeleteIdeaByID (id string ) error
func (ideaser *IdeaService) DeleteIdeaByID(id string ) bool {
	if erra := ideaser.IdeaRepo.DeleteIdeaByID(id)  ;erra != nil {
		return false
	}
	return true
}

// UpdateIdea functnion returning the Updated Idea or nil value 
func (ideaser *IdeaService )  UpdateIdea( idea *entity.Idea  ) *entity.Idea {
	idea  , ear := ideaser.IdeaRepo.UpdateIdea(idea)
	if ear != nil {
		return nil 
	}
	return idea 
}
// GetIdeasByUserID (userid string ) []*entity.Idea
// func (ideaser *IdeaService)  GetIdeasByUserID(userid string ) []*entity.Idea{
// 	ideas  , era := ideaser.IdeaRepo.GetIdeasByUserID(userid )
// 	if era != nil {
// 		return nil 
// 	} 
// 	return ideas
// }