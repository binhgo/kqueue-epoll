package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gobwas/ws"
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

	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		return
	}

	err = sv.kq.Add(conn)
	if err != nil {
		log.Printf("Fail to add connection")
		conn.Close()
	}
}
