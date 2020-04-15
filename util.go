package main

import (
	"reflect"

	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Util struct {
}

func (u Util) GetFD(conn *websocket.Conn) uint64 {
	if conn != nil {
		connVal := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn").Elem()
		tcpConn := reflect.Indirect(connVal).FieldByName("conn")
		fdVal := tcpConn.FieldByName("fd")
		pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")
		i64 := pfdVal.FieldByName("Sysfd").Int()

		return uint64(i64)
	}

	return 0
}

func (u Util) FromJson(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	return err
}

func (u Util) ToJson(object interface{}) ([]byte, error) {
	b, err := json.Marshal(&object)
	return b, err
}
