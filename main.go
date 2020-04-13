package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var sv *Server

func main() {

	fmt.Println("Server staring...")

	sv = NewServer(-10)
	go sv.Start()

	fmt.Println("Server started")

	http.HandleFunc("/", WS)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func WS(w http.ResponseWriter, r *http.Request) {

	upgrade := websocket.Upgrader{}
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	err = sv.kq.Add(conn)
	if err != nil {
		log.Printf("Fail to add connection")
		conn.Close()
	}
}
