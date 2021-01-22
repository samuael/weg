package Alie

import "github.com/samuael/Project/Weg/internal/pkg/entity"

// AlieService interface
type AlieService interface {
	UpdateAlies( alies *entity.Alie ) *entity.Alie
	CreateAlieDocument( sender  , receiver string ) *entity.Alie
	GetAlies(senderid  , receiverid string ) *entity.Alie
	HaveTheyAliesTable(  SenderID  , ReceiverID string  ) bool
	AreTheyAlies( SenderID  , ReceiverID string ) bool
	GetMyAlies(  id string  ) []*entity.User
	DeleteAlieByID(friendid string ) bool 
}