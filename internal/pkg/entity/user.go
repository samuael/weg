// Package entity representing 
// holding list of structs and models
package entity

import (
	"time"
)

// User  representing chatting user
type User struct {
	ID string    `bson:"_id,omitempty"  json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Password string   `json:"password,omitempty"`
	Imageurl string  `json:"imageurl,omitempty"`
	Bio string   `json:"bio,omitempty"`
	Email string    `json:"email,omitempty"`
	LastSeen time.Time `bson:"last_seen"  json:"last_seen,omitempty"`
	LastUpdated time.Time `bson:"last_updated"  json:"last_updated,omitempty" `
	MyGroups []string  `json:"my_groups,omitempty"`
	MyAlies []string `json:"my_alies,omitempty"`
}
// Alie struct representing two Clients ID and 
// saving the Begginer of the Caht as StarterID and ObserverID reprsenting the other User 
// message number represents the message number the two clients are in .
type Alie struct {
	ID string `bson:"_id,omitempty"  json:"id,omitempty"`
	MessageNumber int ` bson:"message_number"  json:"message_number,omitempty"`
	A string  `json:"a,omitempty"`// representing the ID of one of Sender client 
	B string `json:"b,omitempty"` // representing the ID of One Of Receiver CLient 
	Messages []*Message `json:"messages,omitempty"`
}