package main

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gobwas/ws"
)

var sv *Server

func main() {

	sv = NewServer(-10)
	go sv.Start()

	sv.SetHandle("GET_ORDER", handleOrder)

	http.HandleFunc("/", UpgradeWebsocket)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handleOrder(request Request) Response {

	md := request.Model
	md = md + "-" + strconv.Itoa(rand.Intn(1000))

	return Response{
		Status:  "OK",
		Message: md,
	}
}

func UpgradeWebsocket(w http.ResponseWriter, r *http.Request) {

	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		return
	}

	err = sv.kq.Add(conn)
	if err != nil {
		log.Printf("FAIL TO ADD CONNECTION")
		conn.Close()
	}
}
