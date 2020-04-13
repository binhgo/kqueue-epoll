package main

import (
	"reflect"

	"github.com/gorilla/websocket"
)

type Util struct {
}

func (u Util) GetFD(conn *websocket.Conn) uint64 {
	connVal := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn").Elem()
	tcpConn := reflect.Indirect(connVal).FieldByName("conn")
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")
	i64 := pfdVal.FieldByName("Sysfd").Int()

	return uint64(i64)
}
