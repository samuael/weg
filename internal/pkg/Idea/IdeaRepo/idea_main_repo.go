package IdeaRepo

import (
	"context"

	"github.com/samuael/Project/Weg/internal/Idea"
	"github.com/samuael/Project/Weg/internal/pkg/entity"
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/primitive"
)

// IdeaRepo the Repo
type IdeaRepo struct {
	DB *mongo.Database
}

// NewIdeaRepo function 
func NewIdeaRepo( db *mongo.Database ) *Idea.IdeaRepo{
	return &IdeaRepo{
		DB : db  , 
	}
}

// CreateIdea methdo to create an idea 
func (idearepo *IdeaRepo) CreateIdea(  idea *entity.Idea ) *entity.Idea {

	insertResult  ,era := idearepo.DB.Collection(entity.IDEA).InsertOne(context.TODO()  , idea) 



	return nil 
}