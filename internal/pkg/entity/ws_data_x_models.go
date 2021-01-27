package entity

// ActiveUserExistance struct
// whose instance to be sent through
// ActiveUserExistanceCheck channel of MainService
type ActiveUserExistance struct {
	MyID     string
	FriendID string
}

// ActiveUser  struct
type ActiveUser struct {
	ID     string
	Active bool
}

// SeenConfirmMessage this message is to be sent from websocket main service to
// client Who is active through the channel to update the message specified in the
// MessageNumber int ,
type SeenConfirmMessage struct {
	ReceiverID    string
	AlieID        string
	MessageNumber int
}

// SeenConfirmIfClientExist to be used only by main service loop as a channel
type SeenConfirmIfClientExist struct {
	// the one who is asking the existance of the user
	RequesterID string
	// real message sender id
	WantedID string
	// representing the message number
	MessageNumber int
}
