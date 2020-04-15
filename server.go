package main

import "log"

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
			panic(err)
		}

		for _, conn := range conns {
			if conn != nil {
				_, msg, err := conn.ReadMessage()

				if err != nil {
					s.kq.Remove(conn)
					conn.Close()
					continue
				}

				log.Printf("msg: %s", string(msg))
			}
		}
	}
}
