// Package Idea representing the idea class
package Idea

import "github.com/samuael/Project/Weg/internal/pkg/entity"

// IdeaService representing idea service methods
type IdeaService interface {
	CreteIdea(*entity.Idea ) *entity.Idea
}