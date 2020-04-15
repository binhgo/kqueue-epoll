package main

import (
	"bytes"
	"fmt"
	"log"
	"net"

	"github.com/gobwas/ws/wsutil"
)

type Server struct {
	kq        iPoll
	kqTimeout int64
}

func NewServer(kqTimeout int64) *Server {
	kq := NewKQueue()
	return &Server{kq: kq, kqTimeout: kqTimeout}
}

func (s *Server) Start() {

	for {
		conns, err := s.kq.Wait(s.kqTimeout)
		if err != nil {
			fmt.Printf("Failed to epoll wait %v", err)
			continue
		}

		for _, conn := range conns {
			if conn == nil {
				break
			}

			s.Process(conn)
		}
	}
}

func (s *Server) Process(conn net.Conn) {

	msg, _, err := wsutil.ReadClientData(conn)
	if err != nil {
		if err := s.kq.Remove(conn); err != nil {
			log.Printf("Failed to remove %v", err)
		}
		conn.Close()
	}

	fmt.Println(string(msg))

	req := Request{}
	err = Util{}.FromJson(msg, &req)

	bb := bytes.Buffer{}
	if err != nil {
		bb.WriteString("Error parse request")
		wsutil.WriteServerText(conn, bb.Bytes())
		return
	}

	// query data
	bb.WriteString("OK")
	wsutil.WriteServerText(conn, bb.Bytes())
}
