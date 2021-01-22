package service

import "github.com/samuael/Project/Weg/internal/pkg/entity"

// MainService struct representing
// continuously running service
// This will have || Broadcast Group Message
// looping over all the
type MainService struct {
	ClientMap map[string ] *Client
	GroupMap map[string ] *MGroup
	Register chan *Client 
	UnRegister  chan *Client 

	// End to End Message the string represents the _id of user (client) and []byte the message 
	EEMessage chan entity.EEMBinary

	// Group Members Message  the string holds the group id and the []byte represents the message binary
	GMMessage chan entity.GMMBinary
}


// NewMainService function for instantiating the MainService and it's channels 
func NewMainService() *MainService {
	return &MainService{
		ClientMap : map[string]*Client{} ,
		GroupMap : map[string ] *MGroup{}  , 
		Register : make( chan *Client ) , 
		UnRegister : make( chan *Client ) , 
		EEMessage : make(chan entity.EEMBinary ) ,
		GMMessage :make(chan entity.GMMBinary ) , 
	}
}

// Run service 
func (ms *MainService ) Run(){
	defer func(){
		close(ms.Register)
		close(ms.UnRegister)
		close(ms.EEMessage)
		close(ms.GMMessage)
	}()
	
	for{
		select {
			case client := <- ms.Register:{
				ms.RegisterClient(client)
			}
			case client := <- ms.UnRegister : {
				ms.UnRegisterClient(client)
			}
			case message := <- ms.EEMessage: {
				ms.EndToEndMessage(message)
			}
			case message := <- ms.GMMessage :{
				ms.BroadcastMessage(message)
			}
		}   
	}
}


// RegisterClient  function returning the Client 
func (ms *MainService)  RegisterClient(client *Client) {
	ms.ClientMap[client.User.ID] = client
	if len(client.User.MyGroups) >0 {
		for _ , groupid := range client.User.MyGroups {
			gp := ms.GroupMap[groupid]
			if gp == nil {
				gp = &MGroup{
					ActiveCount: 1,
					ActiveUsers: []string{ 
						client.ID , 
					},
				}
			}else {
				gp.ActiveCount++
				gp.ActiveUsers = append(gp.ActiveUsers, client.ID)
			}
		}
	}
}

// UnRegisterClient for unregistering clinets 
func ( ms *MainService) UnRegisterClient(client *Client) {
	if client == nil {
		return 
	}
	delete(ms.ClientMap  , client.ID)
	client.Conn.Close()
	close(client.Message)

	if len(client.User.MyGroups ) >0 {
		for _ , groupid := range client.User.MyGroups {
			gp := ms.GroupMap[groupid]
			if gp != nil {
				if gp.ActiveCount == 1 && (func() bool {
					for _ , val := range gp.ActiveUsers {
						if val == client.User.ID {
							return true 
						}
						
					}
					return false 
				}()) {
					// if there is only one user and this client is the only
					// group active user then Delete the Group From the stack 
					gp.ActiveCount=0
					delete(ms.GroupMap , groupid)
				}else if gp.ActiveCount >1  {
					for _ , val := range gp.ActiveUsers {
						if val != client.ID {
							gp.ActiveUsers = append(gp.ActiveUsers , client.ID)
						}
					}
				}
			}
		}
	}
}

// EndToEndMessage  function running end to end messages this function takes a message from 
// Client Service and Send to the 
func (ms *MainService ) EndToEndMessage(message entity.EEMBinary) {
	cli := ms.ClientMap[message.UserID]
	if cli != nil {
		return 
	}
	cli.Message <- message.Data
}

// BroadcastMessage func group message 
func (ms *MainService) BroadcastMessage( message entity.GMMBinary) {
	grp := ms.GroupMap[message.GroupID]
	if grp != nil || len(grp.ActiveUsers)==0 {
		return 
	}
	for _ , userid := range grp.ActiveUsers {
		ms.EEMessage <- entity.EEMBinary{
			UserID: userid,
			Data: message.Data,
		}
	}
}