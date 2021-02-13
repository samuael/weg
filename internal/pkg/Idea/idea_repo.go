package Idea

import "github.com/samuael/Project/Weg/internal/pkg/entity"

// IdeaRepo interface representing the idea methods
type IdeaRepo interface {
	CreateIdea(idea *entity.Idea ) (*entity.Idea, error)
	GetIdeas(  userid string   , offset , limit  int )( []*entity.Idea   , error )
	GetIdeaByID( id string  )(*entity.Idea  , error )
	GetMyIdeas( userid string    )([]*entity.Idea  , error )
	DeleteIdeaByID(id string ) error  
	UpdateIdea( idea *entity.Idea  ) (*entity.Idea  , error)
	// GetIdeasByUserID(userid string ) ([]*entity.Idea , error)
}