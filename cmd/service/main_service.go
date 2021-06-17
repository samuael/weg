package service

import (
	"encoding/json"
	"time"

	// "time"

	"github.com/samuael/Project/Weg/cmd/service/grpc_service_conn/client"
	"github.com/samuael/Project/Weg/internal/pkg/entity"
)

// MainService struct representing the main service class and
// having continuously running "run()"  method as a handler of messages
type MainService struct {
	ClientMap    map[string]*Client
	GroupMap     map[string]*WSGroup
	EEMBinary    chan entity.EEMBinary
	ADDEEMBinary chan entity.EEMBinary
	GMMBinary    chan entity.GMMBinary
	// GMMessage                     chan entity.GMMessage
	Register                      chan *Client
	UnRegister                    chan *Client
	SeenConfirmIfClientExistCheck chan *entity.SeenConfirmIfClientExist
	// ClientConnExistance to handle deleting of ClientCOnn Object fromuser Conns List
	DeleteClientConn chan *entity.ClientConnExistance
	GrpcClient       *client.GrpcClient
}

// NewMainService  funciton craeting a MainServic instance
func NewMainService() *MainService {
	return &MainService{
		ClientMap:    map[string]*Client{},
		GroupMap:     map[string]*WSGroup{},
		EEMBinary:    make(chan entity.EEMBinary),
		ADDEEMBinary: make(chan entity.EEMBinary),
		GMMBinary:    make(chan entity.GMMBinary),
		// GMMessage:                     make(chan entity.GMMessage),
		Register:                      make(chan *Client),
		UnRegister:                    make(chan *Client),
		SeenConfirmIfClientExistCheck: make(chan *entity.SeenConfirmIfClientExist),
	}
}

// SetGrpcClient this method is to be used for setting the GRPC Client
// which is used to broadcast the Message for the clients
func (mainservice *MainService) SetGrpcHandler(client *client.GrpcClient) {
	mainservice.GrpcClient = client
}

// Run method handling sending messages to clients and Coordinating the accessing of
// Client Map and Group Map ...
func (mainservice *MainService) Run() {
	defer func() {
		close(mainservice.Register)
		close(mainservice.UnRegister)
		close(mainservice.EEMBinary)
		close(mainservice.GMMBinary)
		close(mainservice.ADDEEMBinary)
	}()

	go mainservice.ActiveFriendsNotification()
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
		case seencheck := <-mainservice.SeenConfirmIfClientExistCheck:
			{
				mainservice.SeenConfirmIfClientExist(seencheck)
			}
		case message := <-mainservice.EEMBinary:
			{
				sent_confirm := mainservice.SendEEMessage(message)
				// if !sent_confirm {
				// The Message is not sent so i have to send the Other Peer To send the Message
				success := mainservice.GrpcClient.BroadcastEEMessage(message)
				if success {
					// println("The Message is sent to the Other Server ")
					println(sent_confirm)
				}
				// }

			}
		case message := <-mainservice.ADDEEMBinary:
			{
				mainservice.SendEEMessage(message)

			}
		case gmessage := <-mainservice.GMMBinary:
			{
				mainservice.SendGMessage(gmessage)
			}
		case deleteClientConn := <-mainservice.DeleteClientConn:
			{
				mainservice.DeleteUserClientConn(deleteClientConn)
			}

		}
	}
}

// ActiveFriendsNotification method for notifying users about their active friends
func (mainservice *MainService) ActiveFriendsNotification() {
	ticker := time.NewTicker(time.Duration(time.Second * 12))
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case <-ticker.C:
			{
				mainservice.ActiveFriendsBroadcast()
			}
		}
	}
}

// RegisterClient  function to register the client
func (mainservice *MainService) RegisterClient(client *Client) {
	// check the presense of the client in the client map and if the present
	// create another websocket.Conn object in and Runn REad message on it.
	var ip string
	for ips := range client.Conns {
		ip = ips
	}
	for id, lclient := range mainservice.ClientMap {

		if client.ID == id {
			// check if this device is from another client machine or not
			// by comparing the Ip address
			// if so delete
			for ip, ccon := range client.Conns {
				for ips := range lclient.Conns {
					// meaning the client with same ip present
					if ip == ips {
						delete(lclient.Conns, ip)
						lclient.Conns[ip] = ccon
					}
				}
				//  no ip is found simmilar to mine nigga am gonna
				// join the conns listo
				lclient.Conns[ip] = ccon
				go lclient.ReadMessage(ip)
				go lclient.WriteMessage(ip)
				return
			}
			return
		}
	}
	// The client is not listed in the ClientMap nigga
	// i am gonna create it for the first time

	// Creating an active group since new member is activated.
	for _, group := range client.User.MyGroups {
		wsgroup := mainservice.GroupMap[group]
		if wsgroup != nil {
			wsgroup.ActiveCount++
			wsgroup.MembersID = append(wsgroup.MembersID, group)
			break
		}
		grp := client.ClientService.GroupSer.GetGroupByID(group)
		if grp == nil {
			break
		}
		mainservice.GroupMap[grp.ID] = &WSGroup{
			Group:       grp,
			ActiveCount: 1,
			MembersID:   []string{group},
		}
	}
	mainservice.ClientMap[client.ID] = client
	go client.ReadMessage(ip)
	go client.WriteMessage(ip)

}

// UnregisterClient functio to un register the Client to the clientMap
func (mainservice *MainService) UnregisterClient(client *Client) {
	for id := range mainservice.ClientMap {
		if id == client.ID {
			// deleting the client from the GroupActiveUsers list
			for _, gp := range client.User.MyGroups {
				grp := mainservice.GroupMap[gp]
				if grp != nil {
					grp.ActiveCount--
					grp.MembersID = func() []string {
						membs := []string{}
						for _, id := range grp.MembersID {
							if id != client.ID {
								membs = append(membs, id)
							}
						}
						return membs
					}()
					// if the  length of active members is 0 delete the group from active members list
					if grp.ActiveCount == 0 || len(grp.MembersID) == 0 {
						delete(mainservice.GroupMap, gp)
					}
				}
			}

			// after deleting my id from group.MembersID or after deleting the group
			// from mainservice.GroupMap i am gonna remove the client from the clientMap of mainservice.
			delete(mainservice.ClientMap, id)
		}
	}
}

// SeenConfirmIfClientExist method to ckeck the client existence and if the client exist
// send confirmation message to calling client to update the message
// as seenConfirmed
func (mainservice *MainService) SeenConfirmIfClientExist(seencheck *entity.SeenConfirmIfClientExist) {
	var dclient *Client
	for key, client := range mainservice.ClientMap {
		if key == seencheck.RequesterID {
			dclient = client
		}
		if key == seencheck.WantedID {
			if dclient == nil {
				for id, cl := range mainservice.ClientMap {
					if id == seencheck.RequesterID {
						dclient = cl
						break
					}
				}
			}
			break
		}
	}
	if dclient != nil {
		dclient.SeenConfirmMsg <- entity.SeenConfirmMessage{
			ReceiverID:    seencheck.RequesterID,
			AlieID:        seencheck.WantedID,
			MessageNumber: seencheck.MessageNumber,
		}
	}
}

// SendEEMessage method message specificaly to one user
func (mainservice *MainService) SendEEMessage(message entity.EEMBinary) bool {
	for id, client := range mainservice.ClientMap {
		if id == message.UserID {
			for _, cl := range client.Conns {
				cl.Message <- message.Data
			}
			return true
		}
	}
	return false
}

// SendGMessage method to send group broadcast message to group
func (mainservice *MainService) SendGMessage(message entity.GMMBinary) {
	group := mainservice.GroupMap[message.GroupID]
	if group == nil {
		return
	}
	// for each members of the group send this message
	for _, id := range group.MembersID {
		mainservice.EEMBinary <- entity.EEMBinary{
			UserID: id,
			Data:   message.Data,
		}
	}
}

// DeleteUserClientConn for deleting clientConn instance from client map
func (mainservice *MainService) DeleteUserClientConn(dccex *entity.ClientConnExistance) {
	user := mainservice.ClientMap[dccex.ID]
	if user == nil {
		return
	}
	delete(user.Conns, dccex.IP)
	if len(user.Conns) == 0 {
		mainservice.UnRegister <- user
	}
}

// ActiveFriendsBroadcast method
func (mainservice *MainService) ActiveFriendsBroadcast() {
	for id, cl := range mainservice.ClientMap {
		activeFriens := []string{id}
		for _, conn := range cl.Conns {
			if conn.Conn == nil {
				mainservice.UnRegister <- cl
				break
			}
		}
		for _, alieID := range cl.User.MyAlies {

			if mainservice.ClientMap[alieID] != nil {
				activeFriens = append(activeFriens, alieID)
			}
		}

		if len(activeFriens) > 0 {
			activeInfo := entity.XActiveFriends{
				UserID:        id,
				Status:        int(entity.ActiveFriends),
				ActiveFriends: activeFriens,
			}
			if data, era := json.Marshal(activeInfo); era == nil {
				eemess := entity.EEMBinary{
					Data:   data,
					UserID: activeInfo.UserID,
				}
				mainservice.ADDEEMBinary <- eemess
			}

		}
	}
}
