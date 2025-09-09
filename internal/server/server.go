package server

import (
	"fmt"
	"net"
	"sync/atomic"

	"github.com/PavelVaavra/http-from-tcp/internal/request"
	"github.com/PavelVaavra/http-from-tcp/internal/response"
)

type Handler func(w *response.Writer, req *request.Request)

type Server struct {
	State    atomic.Bool //0..closed, 1..open
	Listener net.Listener
	Handler  Handler
}

func Serve(port int, handler Handler) (*Server, error) {
	l, err := net.Listen("tcp", ":"+fmt.Sprintf("%v", port))
	if err != nil {
		return nil, err
	}
	server := Server{
		Listener: l,
		Handler:  handler,
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
	req, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Printf("could not parse HTTP request: error:%v\n", err.Error())
		return
	}

	w := response.Writer{
		Conn: conn,
	}
	s.Handler(&w, req)

	conn.Close()
	fmt.Println("A connection has been closed...")
}
