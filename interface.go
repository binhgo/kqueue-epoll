package main

import "github.com/gorilla/websocket"

type iPoll interface {
	Add(conn *websocket.Conn) error
	Remove(conn *websocket.Conn) error
	Wait(int64) ([]*websocket.Conn, error)
}
