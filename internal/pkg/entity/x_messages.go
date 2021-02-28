package entity

// structs and methods in this class are to be used int the
// exchange messages and this data formats will be simmilar to that of
// messages to be sent in the websocket channel

// XMessage interface representing exchangable messages that have status
// and Method Named(  GetStatus() int  )
type XMessage interface {
	GetStatus() int
	GetBody() interface{}
}

// InMess struct representing in message
// this will be used to read the message from the websocet connection
// In the ReadJSON message
type InMess struct {
	Status int         `json:"status"`
	Body   interface{} `json:"body"`
}

// GetStatus function
func (im *InMess) GetStatus() int {
	return im.Status
}

// GetBody function
func (im *InMess) GetBody() interface{} {
	return im.Body
}

// -------------------------_Seen 1 ---------------------

// SeenBody struct representing the body of seen message
type SeenBody struct {
	MessageNumber int    `json:"message_number"`
	SenderID      string `json:"sender_id"`
	ObserverID    string `json:"observer_id"`
}

// SeenMessage struct
type SeenMessage struct {
	Status   int      `json:"status"`
	Body     SeenBody `json:"body"`
	SenderID string   `json:"-"`
}

// GetStatus implementing XMessage
func (sm *SeenMessage) GetStatus() int {
	return sm.Status
}

// GetBody function
func (sm *SeenMessage) GetBody() interface{} {
	return sm.Body
}

// GetSenderID function
func (sm *SeenMessage) GetSenderID() string {
	return sm.SenderID
}

// ---------------------------------------------------------

// -----------------------------Typing 2  and Stop Typing 3 Messages  -------------

// TypingBody struct
type TypingBody struct {
	TyperID    string `json:"typer_id"`
	ReceiverID string `json:"receiver_id"`
}

// TypingMessage struct
type TypingMessage struct {
	Status   int        `json:"status"`
	Body     TypingBody `json:"body"`
	SenderID string     `json:"-"`
}

// GetStatus function
func (tm *TypingMessage) GetStatus() int {
	return tm.Status
}

// GetBody function
func (tm *TypingMessage) GetBody() interface{} {
	return tm.Body
}

// GetSenderID function
func (tm *TypingMessage) GetSenderID() string {
	return tm.SenderID
}

// -------------------------------------------------

// --------------------EndToEndMessage 4 ---------------------------

// EEMessage struct representing individual message
type EEMessage struct {
	Status   int     `json:"status"`
	Body     Message `json:"body"`
	SenderID string  `json:"-"`
}

// GetStatus implementing the XMessage interface
func (eem *EEMessage) GetStatus() int {
	return eem.Status
}

// GetBody function
func (eem *EEMessage) GetBody() interface{} {
	return eem.Body
}

// GetSenderID function
func (eem *EEMessage) GetSenderID() string {
	return eem.SenderID
}

// ---------------------------------------

// --------------------- Group Message 5 --------------------

// GMMessage struct representting Group Message
type GMMessage struct {
	Status   int          `json:"status"`
	Body     GroupMessage `json:"body"`
	SenderID string       `json:"-"`
}

// GetStatus function
func (gmm *GMMessage) GetStatus() int {
	return gmm.Status
}

// GetBody function
func (gmm *GMMessage) GetBody() interface{} {
	return gmm.Body
}

// GetSenderID function
func (gmm *GMMessage) GetSenderID() string {
	return gmm.SenderID
}

// ----------------------------------------------------

// ------------------- Alie Profile Change 6 and 7 New Alie 8    ----------------

// AlieProfile struct
type AlieProfile struct {
	Status   int    `json:"status"`
	Body     User   `json:"body"`
	SenderID string `json:"-"`
}

// GetStatus representing alie Profile
func (ap *AlieProfile) GetStatus() int {
	return ap.Status
}

// GetBody function
func (ap *AlieProfile) GetBody() interface{} {
	return ap.Body
}

// GetSenderID function
func (ap *AlieProfile) GetSenderID() string {
	return ap.SenderID
}

//----------------- NewAlieBody --------7 ---------------

// NewAlieBody struct
type NewAlieBody struct {
	ReceiverID string `json:"receiver_id"`
	User       *User  `json:"user"`
}

// NewAlie struct representing main body of new alie message
type NewAlie struct {
	Status   int         `json:"status"`
	Body     NewAlieBody `json:"body"`
	SenderID string      `json:"-"`
}

// GetStatus func
func (nal *NewAlie) GetStatus() int {
	return nal.Status
}

// GetBody function
func (nal *NewAlie) GetBody() interface{} {
	return nal.Body
}

// GetSenderID function
func (nal *NewAlie) GetSenderID() string {
	return nal.SenderID
}

// ----------------------------------------------------

// --------------------------Group Profile Change 9 -------------

// GroupProfile struct representing group profile changes
type GroupProfile struct {
	Status   int    `json:"status"`
	Body     Group  `json:"body"`
	SenderID string `json:"-"`
}

// GetStatus func
func (gp *GroupProfile) GetStatus() int {
	return gp.Status
}

// GetBody function
func (gp *GroupProfile) GetBody() interface{} {
	return gp.Body
}

// GetSenderID function
func (gp *GroupProfile) GetSenderID() string {
	return gp.SenderID
}

// -----------------------------------------------------------

// ---------------- Group Join 10 and Group Leave  11  Messages ------------

// JoinLeaveBody struct
// Status Number 9 and 10
type JoinLeaveBody struct {
	User    *User  `json:"user"`
	GroupID string `json:"group_id"`
}

// JoinLeaveMessage struct
type JoinLeaveMessage struct {
	Status   int           `json:"status"`
	Body     JoinLeaveBody `json:"body"`
	SenderID string        `json:"-"`
}

// GetStatus struct
func (jlm *JoinLeaveMessage) GetStatus() int {
	return jlm.Status

}

// GetBody function
func (jlm *JoinLeaveMessage) GetBody() interface{} {
	return jlm.Body
}

// GetSenderID function
func (jlm *JoinLeaveMessage) GetSenderID() string {
	return jlm.SenderID
}

// -----------------------------------------------

// EEMBinary struct representing an end to end message
// that can be passed to the Main Service EEMEssage
// this will serve as a channel
type EEMBinary struct {
	UserID string
	Data   []byte
}

// GMMBinary representing Group message that can be passed to MainService
// this will serve as a channel
type GMMBinary struct {
	GroupID string
	Data    []byte
}

// XActiveFriends struct 
type XActiveFriends  struct {
	Status int   `json:"status"` ;
	UserID string `json:"user_id"` ;
	ActiveFriends []string   `json:"active_friends"`;
  }
  