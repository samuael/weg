package entity

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
	MessageNumber string `json:"message_number"`
	SenderID      string `json:"sender_id"`
	ObserverID    string `json:"observer_id"`
}

// SeenMessage struct
type SeenMessage struct {
	Status int      `json:"status"`
	Body   SeenBody `json:"body"`
}

// GetStatus implementing XMessage
func (sm *SeenMessage) GetStatus() int {
	return sm.Status
}

// GetBody function
func (sm *SeenMessage) GetBody() interface{} {
	return sm.Body
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
	Status int        `json:"status"`
	Body   TypingBody `json:"body"`
}

// GetStatus function
func (tm *TypingMessage) GetStatus() int {
	return tm.Status
}

// GetBody function
func (tm *TypingMessage) GetBody() interface{} {
	return tm.Body
}

// -------------------------------------------------

// --------------------EndToEndMessage 4 ---------------------------

// EEMessage struct representing individual message
type EEMessage struct {
	Status int     `json:"status"`
	Body   Message `json:"body"`
}

// GetStatus implementing the XMessage interface
func (eem *EEMessage) GetStatus() int {
	return eem.Status
}

// GetBody function
func (eem *EEMessage) GetBody() interface{} {
	return eem.Body
}

// ---------------------------------------

// --------------------- Group Message 5 --------------------

// GMMessage struct representting Group Message
type GMMessage struct {
	Status int          `json:"status"`
	Body   GroupMessage `json:"body"`
}

// GetStatus function
func (gmm *GMMessage) GetStatus() int {
	return gmm.Status
}

// GetBody function
func (gmm *GMMessage) GetBody() interface{} {
	return gmm.Body
}

// ----------------------------------------------------

// ------------------- Alie Profile Change 6 and 7 New Alie 8    ----------------
// AlieProfile struct
type AlieProfile struct {
	Status int  `json:"status"`
	Body   User `json:"body"`
}

// GetStatus representing alie Profile
func (ap *AlieProfile) GetStatus() int {
	return ap.Status
}

// GetBody function
func (ap *AlieProfile) GetBody() interface{} {
	return ap.Body
}

//----------------- NewAlieBody --------7 ---------------
type NewAlieBody struct {
	ReceiverID string `json:"receiver_id"`
	User       *User  `json:"user"`
}

// NewAlie struct representing main body of new alie message
type NewAlie struct {
	Status int         `json:"status"`
	Body   NewAlieBody `json:"body"`
}

// ----------------------------------------------------

// --------------------------Group Profile Change 9 -------------

// GroupProfile struct representing group profile changes
type GroupProfile struct {
	Status int   `json:"status"`
	Body   Group `json:"body"`
}

// GetStatus func
func (gp *GroupProfile) GetStatus() int {
	return gp.Status
}

// GetBody function
func (gp *GroupProfile) GetBody() interface{} {
	return gp.Body
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
	Status int           `json:"status"`
	Body   JoinLeaveBody `json:"body"`
}

// GetStatus struct
func (jlm *JoinLeaveMessage) GetStatus() int {
	return jlm.Status

}

// GetBody function
func (jlm *JoinLeaveMessage) GetBody() interface{} {
	return jlm.Body
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
