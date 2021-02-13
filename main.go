package main

import (
	// "context"
	"fmt"

	"github.com/samuael/Project/Weg/internal/pkg/entity"
	"github.com/samuael/Project/Weg/pkg/Helper"
)

func main(){
	group := entity.Group{}
	user := entity.User{}
	message := entity.Message{}
	gmessage := entity.GroupMessage{}
	joinm := entity.JoinLeaveMessage{}
	alie := entity.AlieProfile{}
	messeen := entity.MessageSeen{}
	typingmes := entity.TypingMessage{}
	idea := entity.Idea{}
	fmt.Println(string(Helper.MarshalThis(group)))
	fmt.Println(string(Helper.MarshalThis(user)))
	fmt.Println(string(Helper.MarshalThis(message)))
	fmt.Println(string(Helper.MarshalThis(gmessage)))
	fmt.Println(string(Helper.MarshalThis(joinm)))
	fmt.Println(string(Helper.MarshalThis(alie)))
	fmt.Println(string(Helper.MarshalThis(messeen)))
	fmt.Println(string(Helper.MarshalThis(typingmes)))
	fmt.Println(string(Helper.MarshalThis(idea)))
	fmt.Println(Helper.RemoveObjectIDPrefix("ObjectID(\"5fe1b21d88b1deda65a9a507\")"))
}