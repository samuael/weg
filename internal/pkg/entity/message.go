// Package entity representing holding list of structs and models
package entity

import (
	"time"
)

// GroupMessage Representing group message
type GroupMessage struct {
	GroupID       string    `json:"group_id"`
	SenderID      string    `json:"sender_id"`
	MessageNumber int       `bson:"message_number" json:"message_number"`
	Text          string    `json:"text"`
	Time          time.Time `json:"time"`
}

// Message representing single end to end message
type Message struct {
	ID            string    `bson:"_id,omitempty"   json:"id"`
	SenderID      string    `json:"sender_id"`
	ReceiverID    string    `json:"receiver_id"`
	Text          string    `json:"text"`
	Time          time.Time `json:"time,omitempty"`
	Seen          bool      `json:"seen"`
	Sent          bool      `json:"sent"`
	SeenConfirmed bool      `json:"seen_confirmed"`
	MessageNumber int       `bson:"message_number" json:"message_number"`
}

// MessageSeen struct
// Status number 1
type MessageSeen struct {
	MessageNumber int    `json:"message_number"`
	FriendID      string `json:"friend_id"`
}
