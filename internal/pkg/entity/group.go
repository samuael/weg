package entity

import (
	"time"

)

// Group struct representing chat Group
type Group struct {
	ID string `bson:"_id,omitempty" json:"id,omitempty"`
	OwnerID string `json:"owner_id"` // string representing user ID 
	MembersCount int `json:"members_count"`
	Imageurl string `json:"imageurl"`
	ActiveCounts int `json:"active_counts"`
	GroupName string `json:"group_name"`
	Description string `json:"description"`
	LastMessageNumber int `json:"last_message_number"`
	MembersID []string `json:"members_id"`
	CreatedAt  time.Time  `json:"created_at"`
	Messages []*GroupMessage `json:"messages"`
}

// GroupMembers representing the group and the members of the group 
type GroupMembers struct {
	ID string `bson:"_id,omitempty"  json:"id,omitempty"`
	GroupID string  `json:"group_id,omitempty"`
	MembersID []int   `json:"members_id,omitempty"`
}

