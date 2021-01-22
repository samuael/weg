package Message

import "github.com/samuael/Project/Weg/internal/pkg/entity"

// MessageService interface
type MessageService interface{
	AliesMessages( theid  , myid string  , offset int ) ([]*entity.Message)
	SetMessageSeen(hisid  , myid string  , messageNumber int ) (string  ,  bool)
	SetSeenConfirmed(hisid  , myid  string  , messageNumber int ) bool  
	GetGroupMessages( groupid string , offset int   ) []*entity.GroupMessage
}