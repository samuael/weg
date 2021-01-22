// Package entity representing holding list of structs and models
package entity

import "time"

// Admin representing the adminstrator of the chat system
type  Admin struct {
	ID  string `bson:"_id,omitempty"  json:"id,omitempty"`
	Username  string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Imageurl string `json:"imageurl,omitempty"`
	Role string `json:"role,omitempty"`
	Time time.Time  `json:"time,omitempty"`
}