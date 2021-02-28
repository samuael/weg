// Package entity representing holding list of structs and models
package entity

// Admin representing the adminstrator of the chat system
type  Admin struct {
	ID  string `bson:"_id,omitempty"  json:"id,omitempty"`
	Email string `json:"email"`
	Username  string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	CreatedBy string `json:"create_by"`
}