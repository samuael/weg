package GroupRepo

import (
	"context"
	"errors"
	"fmt"

	// "github.com/gopinath-langote/1build/cmd"
	"github.com/samuael/Project/Weg/internal/pkg/Group"
	"github.com/samuael/Project/Weg/internal/pkg/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GroupRepo struct representing group repositories
type GroupRepo struct {
	DB *mongo.Database
}

// NewGroupRepo function returning group repo pointer instance 
func NewGroupRepo(db *mongo.Database  )  Group.GroupRepo {
	return &GroupRepo{
		DB : db, 
	}
}
// CreatGroup repo method 
func (grepo *GroupRepo )  CreatGroup(group *entity.Group ) (*entity.Group ,error) {
	_ ,er := grepo.DB.Collection(entity.GROUP).InsertOne(context.TODO() , group )
	er = grepo.DB.Collection(entity.GROUP).FindOne(context.TODO() , bson.D{{"createdat"  , group.CreatedAt}}).Decode(group)
	return  group  , er 
}

// DeleteGroup ( group *entity.Group ) error
func (grepo *GroupRepo)  DeleteGroup( ID  string  ) error {

	objID  , era := primitive.ObjectIDFromHex(ID) 
	if era != nil {
		return errors.New("Invalid ID ")
	}
	filter := bson.D{
		{"_id"  ,  objID  } }// bson.ObjectIdHex(group.ID) }
	delRes  ,er := grepo.DB.Collection(entity.GROUP).DeleteOne(context.TODO()  ,filter )
	if er != nil || delRes.DeletedCount==0 {
		return errors.New("No Record Was Found to be Deleted ")
	}
	return nil 
}
// GetGroupByID function returning pointer to group instance repository 
func (grepo *GroupRepo) GetGroupByID(ID string )  (*entity.Group , error) {
	group := &entity.Group{
		ID : ID ,
	}
	id  , err:= primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil , err
	}
	filter := bson.D{{"_id"  , id}}
	era := grepo.DB.Collection(entity.GROUP).FindOne(context.TODO() , filter).Decode(group)
	if era != nil {
		return nil  , era
	}
	return group  , nil 
}



/*type Group struct {
	ID string `bson:"_id,omitempty"`
	OwnerID string // string representing user ID 
	MembersCount int 
	Imageurl string 
	ActiveCounts int 
	GroupName string 
	Description string 
	LastMessageNumber int 
	MembersID []string 
	CreatedAt  time.Time
}
*/
// UpdateGroup (group *entity.Group) (*entity.Group , error)
func (grepo *GroupRepo) UpdateGroup(group *entity.Group) (*entity.Group , error){
	oid  , er := primitive.ObjectIDFromHex(group.ID)
	if er != nil {
		return nil  , er
	}
	if group.MembersCount==0 {
		group.MembersCount=1
	}
	inres  , er := grepo.DB.Collection(entity.GROUP).UpdateOne(context.TODO(), bson.D{{"_id", oid}},
	bson.D{
		{"$set", bson.D{{ "ownerid" , group.OwnerID  },
		{ "imageurl" , group.Imageurl  }, 
		{"membersid"  , group.MembersID  },
		{ "activecounts" , group.ActiveCounts  }, 
		{ "groupname" , group.GroupName  }, 
		{ "description" , group.Description }, 
		{ "lastmessagenumber" , group.LastMessageNumber },
		{ "memberscount" , group.MembersCount },
		{ "createdat" , group.CreatedAt  },
		{ "messages" , group.Messages  },
	},},})
	if er != nil || inres.MatchedCount == 0  {
		return nil  , errors.New("Update Error")
	}
	return group  , nil 
}

// IncrementMembersCount function For Incrementing the count ot 
func (grepo *GroupRepo)  IncrementMembersCount( groupID string  ) error{
	oid  , er := primitive.ObjectIDFromHex(groupID)
	if er != nil {
		return er
	}
	inres , er := grepo.DB.Collection(entity.GROUP).UpdateOne(context.TODO() , bson.D{{ "_id"  , oid}}  , bson.D{{"$inc" , bson.D{{ "memberscount" ,1 }}}})

	if er != nil || inres.UpsertedCount==0 {
		return errors.New("Error ")
	}
	return nil 
}
// DecrementMembersCount function to decrement members count for the Group ID Specified 
func (grepo *GroupRepo) DecrementMembersCount(groupID string ) error {
	oid  , er := primitive.ObjectIDFromHex(groupID)
	if er != nil {
		return er
	}
	inres , er := grepo.DB.Collection(entity.GROUP).UpdateOne(context.TODO() , bson.D{{ "_id"  , oid}}  , bson.D{{"$dec" , bson.D{{ "memberscount" ,1 }}}})

	if er != nil || inres.UpsertedCount==0 {
		return errors.New("Error ")
	}
	return nil  
}

// SearchGroupsByName (groupname string ) ([]*entity.Group , error)
func (grepo *GroupRepo) SearchGroupsByName(groupname string ) ([]*entity.Group , error){
	fmt.Println(groupname)
	groups := []*entity.Group{}
	cursor , er := grepo.DB.Collection(entity.GROUP).Find(context.TODO() , bson.D{{ "groupname" , groupname }} )
	
	if er != nil {
		return groups , nil
	}
	groupo := &entity.Group{}
	err := cursor.Decode(groupo)
	if err != nil {
		return groups , err
	}
	groups = append(groups, groupo)

	for cursor.Next(context.TODO()) {
		grp := &entity.Group{}
		cursor.Decode(grp)
	}
	return groups , er
}

// DoesGroupExist (groupID string ) error
func (grepo *GroupRepo) DoesGroupExist(groupID string ) error {
	id  , err:= primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return err
	}
	filter := bson.D{{"_id"  , id}}
	era := grepo.DB.Collection(entity.GROUP).FindOne(context.TODO() , filter).Err()
	return era
}

// IsGroupMember (groupID string) error
func (grepo *GroupRepo ) IsGroupMember(groupID , memberID string) error {
	grp := &entity.Group{}
	oid  , er := primitive.ObjectIDFromHex(groupID)
	if er != nil {
		return er
	}
	er = grepo.DB.Collection(entity.GROUP).FindOne(context.TODO()  , bson.D{{"_id"  , oid}}).Decode(grp)
	if er != nil {
		return er
	}
	for _  ,val := range grp.MembersID{
		if val == memberID {
			return nil 
		}
	}
	return errors.New(" Is Not Group Member ... ")
} 
