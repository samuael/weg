package IdeaRepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/samuael/Project/Weg/internal/pkg/Idea"
	"github.com/samuael/Project/Weg/internal/pkg/entity"
	"github.com/samuael/Project/Weg/pkg/Helper"

	// "github.com/samuael/Project/Weg/pkg/Helper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/mongo/primitive"
)

// IdeaRepo the Repo
type IdeaRepo struct {
	DB *mongo.Database
}

// NewIdeaRepo function
func NewIdeaRepo(db *mongo.Database) Idea.IdeaRepo {
	return &IdeaRepo{
		DB: db,
	}
}

// CreateIdea methdo to create an idea
func (idearepo *IdeaRepo) CreateIdea(idea *entity.Idea) (*entity.Idea, error) {

	insertOneResult, era := idearepo.DB.Collection(entity.IDEA).InsertOne(context.TODO(), idea)
	if era != nil {
		return idea, era
	}
	result := Helper.ObjectIDFromInsertResult(insertOneResult)
	fmt.Printf("The Idea ID Is %s ... \n", result)
	if result != "" {
		idea.ID = result
	}
	fmt.Println("The Insert Id ID : ", result)
	return idea, nil
}

// GetIdeas (  userid string   , offset , limit  int )( []*entity.Idea   , error )
func (idearepo *IdeaRepo) GetIdeas(userid string, offset, limit int) ([]*entity.Idea, error) {
	options := &options.FindOptions{}
	fmt.Printf("  Limit : %d  \nOffset : %d"  , limit  , offset)
	options.SetLimit(int64(offset))
	options.SetLimit(int64(limit))
	ideas := []*entity.Idea{}
	cursor, era := idearepo.DB.Collection(entity.IDEA).Find(context.TODO(), bson.D{}, options)
	if era != nil {
		return nil, era
	}
	for cursor.Next(context.TODO()) {
		var idea *entity.Idea
		idea = &entity.Idea{}
		era = cursor.Decode(idea)
		if era == nil {
			ideas = append(ideas, idea)
			continue
		}
	}
	if len(ideas) == 0 {
		return ideas, errors.New("No Record Is Found ... ")
	}
	return ideas, nil
}

// GetIdeaByID (  id string  )(*entity.Idea  , error )
func (idearepo *IdeaRepo) GetIdeaByID(id string) (*entity.Idea, error) {
	if oid, era := primitive.ObjectIDFromHex(id); era == nil {
		idea := &entity.Idea{}
		era = idearepo.DB.Collection(entity.IDEA).FindOne(context.TODO(), bson.D{{"_id", oid}}).Decode(idea)
		if era != nil {
			return nil, era
		}
		return idea, nil
	} else {
		return nil, era
	}
}
// GetMyIdeas ( userid string    )([]*entity.Idea  , error )
func (idearepo *IdeaRepo) GetMyIdeas(userid string) ([]*entity.Idea, error) {
	ideas := []*entity.Idea{}
	if cursor, era := idearepo.DB.Collection(entity.IDEA).Find(context.TODO(), bson.D{{"ownerid", userid}}); era == nil {
		// idea := &entity.Idea{}
		// if era = cursor.Decode(idea); era == nil {
		// 	ideas = append(ideas, idea)
		// } else {
		// 	return ideas, era
		// }
		for cursor.Next(context.TODO()) {
			var idea *entity.Idea
			if era = cursor.Decode(idea); era == nil {
				ideas = append(ideas, idea)
				continue
			}
		}
		return ideas, nil
	} else {
		return ideas, era
	}
}

// DeleteIdeaByID (id string ) error 
func (idearepo *IdeaRepo) DeleteIdeaByID(id string ) error {
	var oid primitive.ObjectID
	oid  , era :=   primitive.ObjectIDFromHex(id)
	if  era != nil {
		return era 
	}
	deleteRes  , era := idearepo.DB.Collection(entity.IDEA).DeleteOne(context.TODO(), bson.D{{  "_id" , oid }})
	if deleteRes.DeletedCount==0 || era != nil {
		return errors.New("Nigga I can't Delete The Message ")
	}
	return nil
}

// UpdateIdea ( idea *entity.Idea  ) (*entity.Idea  , error)
func (idearepo *IdeaRepo) UpdateIdea( idea *entity.Idea  ) (*entity.Idea  , error) {
	var oid primitive.ObjectID
	var era error 
	if oid  , era = primitive.ObjectIDFromHex(idea.ID); era != nil {
		return idea  , era 
	}
	updateRes , err := idearepo.DB.Collection(entity.IDEA).UpdateOne(context.TODO(),  bson.D{{"_id" , oid}}  , 
	bson.D{
		{"$set" , bson.D{
			{ "likes" , idea.Likes  }, 
			{  "imageurl" , idea.ImageURL } ,
			{  "title" , idea.Title } ,
			{  "description" , idea.Description } ,
			{  "dislikes" , idea.Dislikes } ,
			{  "likersid" , idea.LikersID } ,
			{  "dislikersid" , idea.Dislikes } ,
			{  "ownerid" , idea.OwnerID } ,
		}} , 
		} )
		if updateRes.ModifiedCount==0 ||  err != nil {
			print("Modified Count .. "  , updateRes.ModifiedCount)
			return nil , err
		}
		return idea  , nil 
}

