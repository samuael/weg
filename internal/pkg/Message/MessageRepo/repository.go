package MessageRepo

import (
	"context"
	"errors"

	"github.com/samuael/Project/Weg/internal/pkg/Alie"
	"github.com/samuael/Project/Weg/internal/pkg/Message"
	"github.com/samuael/Project/Weg/internal/pkg/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MessageRepo struct
type MessageRepo struct {
	DB *mongo.Database
	AlieRepo Alie.AlieRepo
}
// NewMessageRepo function
func NewMessageRepo( db *mongo.Database , alierepo Alie.AlieRepo ) Message.MessageRepo {
	return &MessageRepo{
		DB: db,
		AlieRepo : alierepo , 
	}
}
// AliesMessages this is a method returning  alies messages which are after the offset message 
// and newly updated messages that are seen by the receiving body 
// messages where Sent = false & receiver_id == my id 
// seen=true & seen_sent = false & sender_id == my id  
// after sending the sending the Messages 
// there will be some chageds that has to be made on the messages the user is fetching them 
func (messrepo *MessageRepo)  AliesMessages( theid  , myid string  , offset int ) ([]*entity.Message , error) {
	messages := []*entity.Message{}
	alie := &entity.Alie{}
	er := messrepo.DB.Collection(entity.ALIE).FindOne(context.TODO()  , bson.D{{ "$or" ,bson.A{ bson.D{{ "a" , theid } , { "b" , myid }}, bson.D{{ "b" , theid } , { "a" , myid }}}}}).Decode(alie)
	if er != nil {
		return messages  , er 
	}
	for i:=0 ; i <len(alie.Messages);i++ {
		ms := alie.Messages[i]
		if !(ms.Sent) && ms.ReceiverID==myid {
			// i am the Receiver and i am getting the message. 
			// so, you can update the message Sent part because i am getting it nigga 
			alie.Messages[i].Sent = true
			messages= append(messages, ms)
			continue
		}
		if ms.Seen && !(ms.SeenConfirmed) && ( ms.SenderID==myid) {
			//  yes i have seen that the receiver have seen it and dont sent it to me again 
			alie.Messages[i].SeenConfirmed=true 
			// telling yes the sender have told that the message is seen
			messages = append(messages, ms)
			continue
		}
		if ms.MessageNumber-1 >= int(offset){
			messages=append(messages, ms)
			continue
		}
	}
	// saving the update 
	alie  , er = messrepo.AlieRepo.UpdateAlies(alie)
	if er != nil {
		return messages , er 
	}
	return messages , nil 
}

// SetMessageSeen (hisid  , myid string  , messageNumber int )  error
func (messrepo *MessageRepo) SetMessageSeen(hisid  , myid string  , messageNumber int )  error {
	alie := &entity.Alie{}
	er := messrepo.DB.Collection(entity.ALIE).FindOne(context.TODO()  , bson.D{{ "$or" ,bson.A{ bson.D{{ "a" , hisid } , { "b" , myid }}, bson.D{{ "b" , hisid } , { "a" , myid }}}}}).Decode(alie)
	if er != nil {
		return errors.New("This Friends Doesn't Exist ")
	}
	if alie.MessageNumber < messageNumber{
		return errors.New("No Message Found By this Number ")
	}
	change := false
	if alie.Messages[messageNumber-1].MessageNumber== messageNumber {
		if alie.Messages[messageNumber-1].SenderID == myid{
			// You Are the sender and you can't set the message as seen because you are not the receiver 
			return errors.New("UnAuthorized Action you are not Receiver ")
		}
		change= true
		alie.Messages[messageNumber-1].Seen=true
	}else {
		for i:=0;  i < len(alie.Messages); i++ {
			if alie.Messages[i].MessageNumber == messageNumber {
				if alie.Messages[i].SenderID == myid{
					// You Are the sender and you can't set the message as seen because you are not the receiver 
					return errors.New("UnAuthorized Action you are not Receiver ")
				}
				change= true
				alie.Messages[i].Seen=true
			}
		}
	}
	if !change {
		return errors.New("No Change has made ")
	}
	alie  ,er = messrepo.AlieRepo.UpdateAlies(alie)
	if er != nil {
		return errors.New("Internal server ERROR ")
	} 
	return nil 
}

// SetSeenConfirmed (hisid  , myid  string  , messageNumber int ) error
func (messrepo *MessageRepo ) SetSeenConfirmed(hisid  , myid  string  , messageNumber int ) error{
	alie := &entity.Alie{}
	er := messrepo.DB.Collection(entity.ALIE).FindOne(context.TODO()  , bson.D{{ "$or" ,bson.A{ bson.D{{ "a" , hisid } , { "b" , myid }}, bson.D{{ "b" , hisid } , { "a" , myid }}}}}).Decode(alie)
	if er != nil {
		return errors.New("This Friends Doesn't Exist ")
	}
	if alie.MessageNumber < messageNumber{
		return errors.New("No Message Found By this Number ")
	}
	change := false
	if alie.Messages[messageNumber-1].MessageNumber== messageNumber {
		if alie.Messages[messageNumber-1].ReceiverID == myid{
			// You Are the sender and you can't set the message as seen because you are not the receiver 
			return errors.New("UnAuthorized Action you are not Sender ")
		}
		change= true
		alie.Messages[messageNumber-1].SeenConfirmed=true
	}else {
		for i:=0;  i < len(alie.Messages); i++ {
			if alie.Messages[i].MessageNumber == messageNumber {
				if alie.Messages[i].ReceiverID == myid{
					// You Are the sender and you can't set the message as seen because you are not the receiver 
					return errors.New(" UnAuthorized Action you are not Sender ")
				}
				change= true
				alie.Messages[i].SeenConfirmed=true
			}
		}
	}
	if !change {
		return errors.New("No Change has made ")
	}
	alie  ,er = messrepo.AlieRepo.UpdateAlies(alie)
	if er != nil {
		return errors.New("Internal server ERROR ")
	} 
	return nil 
} 

// GetGroupMessages ( groupid string , offset int   ) ([]*entity.Message , error)
func (messrepo *MessageRepo) GetGroupMessages( groupid string , offset int   ) ([]*entity.GroupMessage , error) {
	messages := []*entity.GroupMessage{}
	group := &entity.Group{}

	oid  , er := primitive.ObjectIDFromHex(groupid)
	if er != nil {
		return messages  , er
	}
	er = messrepo.DB.Collection(entity.GROUP).FindOne(context.TODO()  , bson.D{{"_id"  , oid }}).Decode(group)
	if er != nil {
		return messages  , nil 
	}
	for i:=0 ; i <len(group.Messages);i++ {
		ms := group.Messages[i]
		if ms.MessageNumber-1 >= int(offset){
			messages=append(messages, ms)
		}
	}

	return messages , er 

}