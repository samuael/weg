package service

import "github.com/samuael/Project/Weg/internal/pkg/entity"

// MainService struct representing the main service class and
// having continuously running "run()"  method as a handler of messages
type MainService struct {
	ClientMap map[string]*Client
	GroupMap  map[string]*WSGroup
	XMessage  chan entity.XMessage
	// GMMessage                     chan entity.GMMessage
	Register                      chan *Client
	UnRegister                    chan *Client
	SeenConfirmIfClientExistCheck chan *entity.SeenConfirmIfClientExist
}

// NewMainService  funciton craeting a MainServic instance
func NewMainService() *MainService {
	return &MainService{
		ClientMap: map[string]*Client{},
		GroupMap:  map[string]*WSGroup{},
		XMessage:  make(chan entity.XMessage),
		// GMMessage:                     make(chan entity.GMMessage),
		Register:                      make(chan *Client),
		UnRegister:                    make(chan *Client),
		SeenConfirmIfClientExistCheck: make(chan *entity.SeenConfirmIfClientExist),
	}
}

// Run method handling sending messages to clients and Coordinating the accessing of
// Client Map and Group Map ...
func (mainservice *MainService) Run() {
	defer func() {
		close(mainservice.Register)
		close(mainservice.UnRegister)
		close(mainservice.XMessage)
	}()
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
			// case message <- mainservice.XMessage:
			// 	{

			// 	}
		}
	}
}

// RegisterClient  function to register the client
func (mainservice *MainService) RegisterClient(client *Client) {
	// check the presense of the client in the client map and if the present
	// create another websocket.Conn object in and Runn REad message on it.
	for id, lclient := range mainservice.ClientMap {

		if client.ID == id {
			// check if this device is from another client machine or not
			// by comparing the Ip address
			// if so delete
			for ip, ccon := range client.Conns {
				for ips := range lclient.Conns {
					// meaning the client with same ip present
					if ip == ips {
						return
					}
				}

				//  no ip is found simmilar to mine nigga am gonna
				// join the conns listo
				lclient.Conns[ip] = ccon
				go lclient.ReadMessage(ip)
				go lclient.WriteMessage(ip)
			}
		}
	}
}

// UnregisterClient functio to un register the Client to the clientMap
func (mainservice *MainService) UnregisterClient(client *Client) {
}
