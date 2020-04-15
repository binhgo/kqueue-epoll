package main

import (
	"net"
)

type iPoll interface {
	Add(conn net.Conn) error
	Remove(conn net.Conn) error
	Wait(int64) ([]net.Conn, error)
}
