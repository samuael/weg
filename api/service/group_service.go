package service

import (
	"fmt"

	"github.com/samuael/Project/Weg/internal/pkg/Message"
	"github.com/samuael/Project/Weg/internal/pkg/entity"
)

// GroupService service
// handling group related functionalities and broadcasting them to the main service
// on which it can access the client map and the Group map
type GroupService struct {
	Message 	   chan entity.XMessage
	MainSer  	   *MainService
	MessageHandler Message.MessageService
}

// Run function
func (gser *GroupService) Run(){
	for {
		select {
		case message := <- gser.Message :
			{
				fmt.Println(message)
			}
		}
	}
}