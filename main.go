package main

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gobwas/ws"
	"kqueue-epoll/server"
)

var srv *server.Server

func main() {

	srv = server.NewServer(-10)
	go srv.Start()

	srv.SetHandle("GET-ORDER", handleOrder)

	http.HandleFunc("/", UpgradeWebsocket)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handleOrder(request server.Request) server.Response {

	md := request.Action
	md = md + "-" + strconv.Itoa(rand.Intn(1000))

	return server.Response{
		Status:  "OK",
		Message: md,
	}
}

func UpgradeWebsocket(w http.ResponseWriter, r *http.Request) {

	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		return
	}

	err = srv.Poll.Add(conn)
	if err != nil {
		log.Printf("FAIL TO ADD CONNECTION")
		conn.Close()
	}
}
