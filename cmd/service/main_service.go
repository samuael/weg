package service

import "github.com/samuael/Project/Weg/internal/pkg/entity"

// MainService struct representing the main service class and
// having continuously running "run()"  method as a handler of messages
type MainService struct {
	ClientMap  map[string]*Client
	GroupMap   map[string]*WSGroup
	EEMessage  chan entity.EEMessage
	GMMessage  chan entity.GMMessage
	Register   chan *Client
	UnRegister chan *Client
}

// NewMainService  funciton craeting a MainServic instance
func NewMainService() *MainService {
	return &MainService{
		ClientMap:  map[string]*Client{},
		GroupMap:   map[string]*WSGroup{},
		EEMessage:  make(chan entity.EEMessage),
		GMMessage:  make(chan entity.GMMessage),
		Register:   make(chan *Client),
		UnRegister: make(chan *Client),
	}
}

// Run method handling sending messages to clients and Coordinating the accessing of
// Client Map and Group Map ...
func (mainservice *MainService) Run() {

	for {

		select {
		case client := <-mainservice.Register:
			{
				mainservice.RegisterClient(client)
			}
		case client := <-mainservice.UnRegister:
			{
				mainservice.UnregisterClient(client)
			}

		}
	}
}

// RegisterClient  function to register the client
func (mainservice *MainService) RegisterClient(client *Client) {

}

// UnregisterClient functio to un register the Client to the clientMap
func (mainservice *MainService) UnregisterClient(client *Client) {

}
