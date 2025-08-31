package server

import (
	"net"
	"sync/atomic"
	"fmt"
)

// In your new server package, implement the following methods and types (I won't hold your hand too much here):
// 
// type Server struct - Contains the state of the server
// 
// func Serve(port int) (*Server, error) - Creates a net.Listener and returns a new Server instance. Starts listening for requests inside a goroutine.
// 
// func (s *Server) Close() error - Closes the listener and the server
// 
// func (s *Server) listen() - Uses a loop to .Accept new connections as they come in, and handles each one in a new goroutine. I used an atomic.Bool to 
// track whether the server is closed or not so that I can ignore connection errors after the server is closed.
// 
// func (s *Server) handle(conn net.Conn) - Handles a single connection by writing the following response and then closing the connection:

type Server struct {
	State atomic.Bool
	Server net.Listener
}

func Serve(port int) (*Server, error) {
	l, err := net.Listen("tcp", ":" + fmt.Sprintf("%v", port))
	if err != nil {
		return nil, err
	}
	server := Server {
		Server: l,
	}
	server.State.Store(true)
	go server.listen()

	return &server, nil
}

func (s *Server) Close() error {
	err := s.Server.Close()
	if err != nil {
		return err
	}
	s.State.Store(false)
	return nil
}

func (s *Server) listen() {
	for {
		// Wait for a connection.
		conn, err := s.Server.Accept()
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