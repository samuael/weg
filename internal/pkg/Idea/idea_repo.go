package Idea

import "github.com/samuael/Project/Weg/internal/pkg/entity"

// IdeaRepo interface representing the idea methods
type IdeaRepo interface {
	CreteIdea(*entity.Idea ) (*entity.Idea, error)
}