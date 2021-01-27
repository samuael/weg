package entity

import "github.com/gorilla/websocket"

// ClientConn  struct representing a
// single client connection
type ClientConn struct {
	IP   string
	Conn *websocket.Conn
}
