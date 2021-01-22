package AlieService

import (
	"github.com/samuael/Project/Weg/internal/pkg/Alie"
	"github.com/samuael/Project/Weg/internal/pkg/entity"
)

// AlieService struct representing
type AlieService struct {
	AlieRepo Alie.AlieRepo
}

// NewAlieService function 
func NewAlieService( repo Alie.AlieRepo ) Alie.AlieService {
	return &AlieService{
		AlieRepo : repo  , 
	}
}

// CreateAlieDocument ( sender  , receiver string ) *entity.Alie
func (alieser *AlieService) CreateAlieDocument( sender  , receiver string ) *entity.Alie {
	alies ,er := alieser.AlieRepo.CreateAlieDocument(sender , receiver )
	if er != nil {
		return nil 
	}
	return alies

}
// GetAlies (senderid  , receiverid string ) *entity.Alie
func(alieser *AlieService  ) GetAlies(senderid  , receiverid string ) *entity.Alie{
	alie  , er := alieser.AlieRepo.GetAlies(senderid  , receiverid )
	if er != nil {
		return nil 
	}
	return alie
}

// UpdateAlies ( alies *entity.Alie ) *entity.Alie
func (alieser *AlieService ) UpdateAlies( alies *entity.Alie ) *entity.Alie{
	alie  , er := alieser.AlieRepo.UpdateAlies( alies )
	if er != nil {
		return nil 
	}
	return alie 
}
// HaveTheyAliesTable ( SenderID  , ReceiverID string )bool
func (alieser  *AlieService)  HaveTheyAliesTable( SenderID  , ReceiverID string )bool {
	er := alieser.AlieRepo.HaveTheyAliesTable(  SenderID  , ReceiverID ) 
	if er != nil {
		return false 
	}
	return true
}

// AreTheyAlies ( SenderID  , ReceiverID string )bool
func (alieser  *AlieService)  AreTheyAlies( SenderID  , ReceiverID string )bool {
	er := alieser.AlieRepo.AreTheyAlies(SenderID  , ReceiverID )
	er = alieser.AlieRepo.HaveTheyAliesTable(  SenderID  , ReceiverID ) 
	if er != nil {
		return false 
	}
	return true
}

// GetMyAlies function returning your Alies as a 
func (alieser *AlieService) GetMyAlies(id string) []*entity.User {
	users  ,er  := alieser.AlieRepo.GetMyAlies(id)
	if er != nil {
		return nil
	}
	return users
}



// DeleteAlieByID (friendid string ) bool
func (alieser *AlieService) DeleteAlieByID(friendid string ) bool {
	er := alieser.AlieRepo.DeleteAlieByID(friendid)
	if er != nil {
		return false 
	} 
	return true 
} 