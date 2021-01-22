package MessageService

import (
	"github.com/samuael/Project/Weg/internal/pkg/Message"
	"github.com/samuael/Project/Weg/internal/pkg/entity"
)

// MessageService struct
type MessageService struct {
	MessageRepo Message.MessageRepo
}

// NewMessageService function 
func NewMessageService(  repo Message.MessageRepo  ) Message.MessageService {
	return &MessageService{
		MessageRepo : repo  , 
	}
}

// AliesMessages list of messages that are exchanged between sender and receiver and starting from offset number 
func (messer *MessageService ) AliesMessages( theid  , myid string  , offset int ) ([]*entity.Message ){
	messages  , er:= messer.MessageRepo.AliesMessages(theid  , myid , offset)
	if er != nil {
		return nil 
	}
	return messages
}
// SetMessageSeen (hisid  , myid string  , messageNumber int )  error
func (messer *MessageService) SetMessageSeen(hisid  , myid string  , messageNumber int ) (string , bool) {
	er := messer.MessageRepo.SetMessageSeen(hisid , myid  , messageNumber )
	if er != nil {
		return er.Error() , false 
	}
	return "change is made succesfuly" ,true 
}
// SetSeenConfirmed (hisid  , myid  string  , messageNumber int ) bool
func (messer *MessageService) SetSeenConfirmed(hisid  , myid  string  , messageNumber int ) bool {
	er := messer.MessageRepo.SetSeenConfirmed(hisid , myid  , messageNumber )
	if er != nil {
		return false
	}
	return true
} 

// GetGroupMessages ( groupid string , offset int   ) []*entity.Message
func (messer *MessageService) GetGroupMessages( groupid string , offset int   ) []*entity.GroupMessage {
	messages  , er := messer.MessageRepo.GetGroupMessages( groupid, offset )
	if er != nil {
		return nil 
	}
	return messages
}