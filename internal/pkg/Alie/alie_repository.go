package Alie

import "github.com/samuael/Project/Weg/internal/pkg/entity"

// AlieRepo interface
type AlieRepo interface {
	UpdateAlies( alies *entity.Alie ) (  *entity.Alie  , error )
	CreateAlieDocument( sender  , receiver string ) (*entity.Alie , error)
	GetAlies(senderid  , receiverid string ) (*entity.Alie , error)
	HaveTheyAliesTable(  SenderID  , ReceiverID string  ) error
	AreTheyAlies( SenderID  , ReceiverID string ) error
	GetMyAlies(  id string  ) ([]*entity.User , error)
	DeleteAlieByID(friendid string ) error
}