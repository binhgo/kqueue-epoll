package main

import (
	"fmt"
	"syscall"

	"github.com/gorilla/websocket"
)

var timespec = &syscall.Timespec{
	Sec:  0,
	Nsec: 100000,
}

type KQueue struct {
	fd          int
	changes     []syscall.Kevent_t
	events      []syscall.Kevent_t
	connections map[uint64]*websocket.Conn
	// lock        *sync.RWMutex
}

func NewKQueue() *KQueue {
	fd, err := syscall.Kqueue()
	if err != nil {
		return nil
	}

	return &KQueue{
		fd:     fd,
		events: make([]syscall.Kevent_t, 64),
		// lock:        &sync.RWMutex{},
		connections: make(map[uint64]*websocket.Conn),
	}
}

func (k *KQueue) Add(wsConn *websocket.Conn) error {

	fmt.Println("add connection")

	fd := Util{}.GetFD(wsConn)

	readChange := syscall.Kevent_t{
		Ident:  fd,
		Flags:  syscall.EV_ADD,
		Filter: syscall.EVFILT_READ,
	}

	k.changes = append(k.changes, readChange)

	// k.lock.Lock()
	// defer k.lock.Unlock()
	k.connections[fd] = wsConn

	fmt.Println("add connection end")

	return nil
}

func (k *KQueue) Remove(wsConn *websocket.Conn) error {

	fmt.Println("remove connection")

	fd := Util{}.GetFD(wsConn)

	deletedChange := syscall.Kevent_t{
		Ident: fd,
		Flags: syscall.EV_DELETE,
	}

	k.changes = append(k.changes, deletedChange)

	// k.lock.Lock()
	// defer k.lock.Unlock()
	delete(k.connections, fd)

	fmt.Println("remove connection end")

	return nil
}

func (k *KQueue) Wait(timeout int64) ([]*websocket.Conn, error) {

	// fmt.Println("start wait")

	var conns []*websocket.Conn

	var nev int
	var err error
	if timeout >= 0 {
		var ts syscall.Timespec
		ts.Nsec = timeout
		nev, err = syscall.Kevent(k.fd, k.changes, k.events, &ts)

	} else {
		nev, err = syscall.Kevent(k.fd, k.changes, k.events, timespec)
	}

	if err != nil && err != syscall.EINTR {
		panic(err)
	}

	// fmt.Println("wait middle")

	// k.changes = k.changes[:0]
	for i := 0; i < nev; i++ {
		// fmt.Println(nev)
		conn := k.connections[k.events[i].Ident]
		conns = append(conns, conn)
	}

	// fmt.Println("wait return")

	return conns, nil
}
