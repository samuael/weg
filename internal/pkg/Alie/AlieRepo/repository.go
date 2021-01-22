package AlieRepo

import (
	"github.com/samuael/Project/Weg/internal/pkg/Alie"
	"github.com/samuael/Project/Weg/internal/pkg/entity"
	"github.com/samuael/Project/Weg/pkg/Helper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"context"
	"errors"
	"fmt"
)

// AlieRepo  struct representing mongodb Admin Repository Class
type AlieRepo struct {
	DB *mongo.Database
}

// NewAlieRepo function returning AdminRepo 
func NewAlieRepo(  db *mongo.Database ) Alie.AlieRepo {
	return &AlieRepo{
		DB : db  , 
	}
}

// CreateAlieDocument ( sender  , receiver string ) (*entity.Alie , error)
func (alierepo *AlieRepo) CreateAlieDocument( sender  , receiver string ) (*entity.Alie , error){
	alie := &entity.Alie{}
	id:= ""
	// id := strconv.Itoa(int(time.Now().UnixNano()))
	// alie.ID= id
	alie.A=sender
	alie.B=receiver
	sr , er := alierepo.DB.Collection(entity.ALIE).InsertOne(context.TODO() , alie)
	
	if oid  , ok := sr.InsertedID.(primitive.ObjectID);ok {
		id = Helper.RemoveObjectIDPrefix(oid.String())
	}
	alie.ID= id
	
	if er != nil {
		return nil ,er 
	}
	return alie , nil 
}

// GetAlies (senderid  , receiverid string ) (*entity.Alie , error)
func (alierepo *AlieRepo) GetAlies(senderid  , receiverid string ) (*entity.Alie , error){
	alies := &entity.Alie{}
	er := alierepo.DB.Collection(entity.ALIE ).FindOne(context.TODO() , bson.D{{"$or" , bson.A{ bson.D{{"a" ,senderid } ,{"b" , receiverid  }}, bson.D{{"b" ,senderid } ,{"a" , receiverid  }} }}}).Decode(alies)
	if er != nil || alies.ID=="" || alies.A=="" || alies.B =="" {
		return nil , er
	}
	return alies , nil 
}


// UpdateAlies ( alies *entity.Alie ) *entity.Alie 
func (alierepo *AlieRepo ) UpdateAlies( alies *entity.Alie ) (*entity.Alie , error ) {
	// alies = &entity.Alie{}
	oid  , er := primitive.ObjectIDFromHex(alies.ID)
	if er != nil {
		fmt.Println("It's Here Nigga ")
		return nil  , er 
	}
	_ , er = alierepo.DB.Collection(entity.ALIE).UpdateOne( context.TODO()  , bson.D{{"_id" , oid }} , bson.D{{ "$set" , bson.D{
		{ "message_number"  , alies.MessageNumber } , 
		{ "messages"   ,  alies.Messages } , 
		{  "a" , alies.A   } , 
		{ "b" , alies.B } , }    }} )
		fmt.Println(er)
		return alies  , er
}
// AreTheyAlies ( SenderID  , ReceiverID string ) error
func (alierepo *AlieRepo) AreTheyAlies( SenderID  , ReceiverID string ) error{
	// Fetch their alies List and check if there is 
	// senderAlies := []string{}
	// receiverAlies := []string{}
	soid  , er := primitive.ObjectIDFromHex(SenderID)
	roid  , er := primitive.ObjectIDFromHex(ReceiverID)
	if er != nil {
		return er 
	}
	sender := &entity.User{}
	receiver := &entity.User{}
	er = alierepo.DB.Collection(entity.USER).FindOne(context.TODO()  , bson.D{{"_id" ,soid  }}  ).Decode(sender)
	er = alierepo.DB.Collection(entity.USER).FindOne(context.TODO()  , bson.D{{"_id" ,roid  }}  ).Decode(receiver)
	if er != nil {
		return er 
	}
	pass := false 
	for _ , vals := range sender.MyAlies{
		if vals == receiver.ID{
			pass= true
		}
	}
	if !pass{
		return errors.New("New Error")
	}
	for _ , vals := range sender.MyAlies{
		if vals == receiver.ID{
			return nil 
		}
	} 
	return errors.New("They Are Not Alies  ")

}

// HaveTheyAliesTable method for craeting Alies Document in  Alies Collection
func (alierepo *AlieRepo)  HaveTheyAliesTable(  SenderID  , ReceiverID string  ) error {

	alies := &entity.Alie{}
	er := alierepo.DB.Collection(entity.ALIE ).FindOne(context.TODO()  , bson.D{{"$or" , bson.A{ bson.D{{"a" ,SenderID } ,{"b" , ReceiverID  }}, bson.D{{"b" ,SenderID } ,{"a" , ReceiverID  }} }}}).Decode(alies)
	if er != nil || alies.ID=="" || alies.A=="" || alies.B =="" {
		return er
	}
	return nil 
} 


// GetMyAlies (  id string  ) ([]*entity.Alie , error)
func (alierepo *AlieRepo) GetMyAlies(  id string  ) ([]*entity.User , error) {
	users := []*entity.User{}
	alies := []*entity.Alie{}
	cursor  , er := alierepo.DB.Collection(entity.ALIE).Find(context.TODO()  ,
	bson.D{{ "$or" , 
		bson.A{ bson.D{{ "a" , id }}, bson.D{{ "b" , id }}}}})
		if er != nil {
			return users  , er 
		}
		alie := &entity.Alie{}
		if er := cursor.Decode(alie); er ==nil {
			alies = append(alies, alie)
		}
		// Finding the alies table in which the user id exist  
	for cursor.Next(context.TODO()) {
		alie := &entity.Alie{}
		cursor.Decode(alie)
		if alie.ID == ""{
			continue
		}
		alies = append(alies, alie)
	}
	fmt.Println(len(alies) , alies) 
	// from each alies table select for new the ID of user that's not you 
	for _ , al := range alies {
		friendID := ""
		if al.A==id {
			friendID= al.B
		}else {
			friendID= al.A
		}
		friend ,er  := alierepo.GetUserByID(friendID)
		if er != nil {
			continue
		}
		users = append(users, friend)
	}
	if len(users)==0 {
		return users  , errors.New("Alies Repo 164 :  No Record Found ")
	}
	return users  , nil
} 

// GetUserByID (id string ) (*entity.User , error )
func (alierepo *AlieRepo) GetUserByID(id string ) (*entity.User , error ){
	user := &entity.User{}
	oid  , er := primitive.ObjectIDFromHex(id)
	if er != nil {
		return nil  , er 
	}
	er = alierepo.DB.Collection(entity.USER).FindOne( context.TODO() , bson.D{{"_id", oid }} ).Decode(user)
	return user  , er 
}

// DeleteAlieByID (friendid string ) error 
func (alierepo *AlieRepo)  DeleteAlieByID(friendid string ) error{
	oid  , er := primitive.ObjectIDFromHex(friendid)
	if er != nil {
		return er 
	}
	delRes , er := alierepo.DB.Collection(entity.ALIE).DeleteOne(context.TODO()  , bson.D{{"_id"  , oid }})
	if er != nil || delRes.DeletedCount==0 {
		return errors.New(" No Record Deleted ")
	}
	return nil 
} 