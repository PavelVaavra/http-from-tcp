package server

import (
	"net"
	"sync/atomic"
	"fmt"
)

type Server struct {
	State atomic.Bool //0..closed, 1..open
	Listener net.Listener
}

func Serve(port int) (*Server, error) {
	l, err := net.Listen("tcp", ":" + fmt.Sprintf("%v", port))
	if err != nil {
		return nil, err
	}
	server := Server {
		Listener: l,
	}
	server.State.Store(true)
	go server.listen()

	return &server, nil
}

func (s *Server) Close() error {
	s.State.Store(false)
	err := s.Listener.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) listen() {
	for {
		// Wait for a connection.
		conn, err := s.Listener.Accept()
		if err != nil {
			if !s.State.Load() {
				return
			} else {
				fmt.Println(err.Error())
				continue
			}
		}
		fmt.Println("A connection has been accepted...")
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	conn.Write([]byte("HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Length: 13\r\n" +
		"\r\n" +
		"Hello World!\n"))

	conn.Close()
	fmt.Println("A connection has been closed...")
}