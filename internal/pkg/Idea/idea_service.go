// Package Idea representing the idea class
package Idea

import "github.com/samuael/Project/Weg/internal/pkg/entity"

// IdeaService representing idea service methods
type IdeaService interface {
	CreateIdea(*entity.Idea ) *entity.Idea
	GetIdeas(  userid string   , offset , limit  int )( []*entity.Idea    )
	GetIdeaByID( id string  )(*entity.Idea   )
	GetMyIdeas( userid string    )([]*entity.Idea  )
	DeleteIdeaByID(id string ) bool
	UpdateIdea( idea *entity.Idea  ) *entity.Idea 
	SearchIdeaByTitle(title string )  ([]*entity.Idea )
	// GetIdeasByUserID(userid string ) []*entity.Idea
}