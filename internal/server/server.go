package server

import (
	"net"
	"sync/atomic"
	"fmt"
	"github.com/PavelVaavra/http-from-tcp/internal/response"
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

// Update the handle method in your server package to use these new functions and methods to return our "default" response:
func (s *Server) handle(conn net.Conn) {
	err := response.WriteStatusLine(conn, response.OK)
	if err != nil {
		return
	}

	headers := response.GetDefaultHeaders(0)

	err = response.WriteHeaders(conn, headers)
	if err != nil {
		return
	}

	conn.Close()
	fmt.Println("A connection has been closed...")
}