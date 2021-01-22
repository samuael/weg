package Message

import "github.com/samuael/Project/Weg/internal/pkg/entity"

// MessageRepo interface
type MessageRepo interface {
	AliesMessages( theid  , myid string  , offset int ) ([]*entity.Message , error)
	SetMessageSeen(hisid  , myid string  , messageNumber int )  error
	SetSeenConfirmed(hisid  , myid  string  , messageNumber int ) error
	GetGroupMessages( groupid string , offset int   ) ([]*entity.GroupMessage , error)
}