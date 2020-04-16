package main

import (
	"bytes"
	"fmt"
	"log"
	"net"

	"github.com/gobwas/ws/wsutil"
)

type HandleFunc func(request Request) Response

type Server struct {
	kq        iPoll
	kqTimeout int64
	handlers  map[string]HandleFunc
}

func NewServer(kqTimeout int64) *Server {
	kq := NewKQueue()
	return &Server{
		kq:        kq,
		kqTimeout: kqTimeout,
		handlers:  make(map[string]HandleFunc),
	}
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

	// process data
	handleFunc := s.GetHandler(req.Method)
	response := handleFunc(req)

	jsn, err := Util{}.ToJson(response)
	if err != nil {
		panic(err)
	}

	wsutil.WriteServerText(conn, jsn)
}

func (s *Server) SetHandle(path string, handler HandleFunc) {
	s.handlers[path] = handler
}

func (s *Server) GetHandler(path string) HandleFunc {
	return s.handlers[path]
}
