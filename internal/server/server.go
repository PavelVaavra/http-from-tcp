package server

import (
	"net"
	"sync/atomic"
	"fmt"
	"io"
	// "bytes"
	"github.com/PavelVaavra/http-from-tcp/internal/response"
	"github.com/PavelVaavra/http-from-tcp/internal/request"
)

// Create a Handler function type and HandlerError struct type in your server package. My HandlerError simply contains a status code and a message.
type Handler func(w io.Writer, req *request.Request) *HandlerError

type HandlerError struct {
	StatusCode response.StatusCode
	Message string
}

type Server struct {
	State atomic.Bool //0..closed, 1..open
	Listener net.Listener
}

// Update server.Serve to accept a Handler function as an argument.
// Create some logic that writes a HandlerError to an io.Writer. This will make it easy for us to keep our error handling consistent and DRY.
func Serve(port int, f Handler) (*Server, error) {
	l, err := net.Listen("tcp", ":" + fmt.Sprintf("%v", port))
	if err != nil {
		return nil, err
	}
	server := Server {
		Listener: l,
	}
	server.State.Store(true)
	go server.listen(f)

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

func (s *Server) listen(f Handler) {
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
		go s.handle(conn, f)
	}
}

// Update the handle method to:
	// Parse the request from the connection
	// Create a new empty bytes.Buffer for the handler to write to
	// Call the handler function
	// If the handler errors, write the error to the connection
	// If the handler succeeds:
		// Create new default response headers
		// Write the status line
		// Write the headers
		// Write the response body from the handler's buffer
func (s *Server) handle(conn net.Conn, f Handler) {
	req, err := request.RequestFromReader(conn)
	fmt.Println(req)
	// if err != nil {
	// 	fmt.Printf("could not parse HTTP request: error:%v\n", err.Error())
	// 	return
	// }

	// var b bytes.Buffer
	// handleError := f(&b, req)
	
	// err = response.WriteStatusLine(conn, handleError.StatusCode, handleError.Message)
	err = response.WriteStatusLine(conn, response.StatusCodeOK, "OK")
	if err != nil {
		return
	}

	headers := response.GetDefaultHeaders(0)

	err = response.WriteHeaders(conn, headers)
	if err != nil {
		return
	}

	// _, err = conn.Write(b.Bytes())
	// if err != nil {
	// 	return
	// }

	conn.Close()
	fmt.Println("A connection has been closed...")
}

// Remember, our Handler is responsible for reporting errors or writing the body. Our server implementation takes care of the rest for now.

// package main

// import (
// 	"bytes"
// 	"fmt"
// 	"os"
// )

// func main() {
// 	var b bytes.Buffer // A Buffer needs no initialization.
// 	b.Write([]byte("Hello "))
// 	fmt.Fprintf(&b, "world!")
// 	b.WriteTo(os.Stdout)
// }

// Output:

// Hello world!