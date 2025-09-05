package main

import (
	"log"
	"io"
	"os"
	"os/signal"
	"syscall"
	"github.com/PavelVaavra/http-from-tcp/internal/server"
	"github.com/PavelVaavra/http-from-tcp/internal/request"
	"github.com/PavelVaavra/http-from-tcp/internal/response"
)

const port = 42069

// Back in main, create a handler function. Let's be sure to test our error handling:
// If the request target (path) is /yourproblem return a 400 and the message "Your problem is not my problem\n"
// If the request target (path) is /myproblem return a 500 and the message "Woopsie, my bad\n"
// Otherwise, it should just write the string "All good, frfr\n" to the response body.
func main() {
	server, err := server.Serve(port, myHandler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func myHandler(w io.Writer, req *request.Request) *server.HandlerError {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		return &server.HandlerError{
			StatusCode: response.StatusCodeBadRequest,
			Message: "Your problem is not my problem",
		}
	case "/myproblem":
		return &server.HandlerError{
			StatusCode: response.StatusCodeInternalServerError,
			Message: "Woopsie, my bad",
		}
	default:
		w.Write([]byte("All good, frfr\n"))
		return &server.HandlerError{
			StatusCode: response.StatusCodeOK,
			Message: "OK",
		}
	}
}

// printf "GET /ahoj HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/8.5.0\r\nAccept: */*\r\n\r\n" | nc localhost 42069
// printf "GET /ahoj HTTP/1.1\r\nHost: localhost:42069\r\nConnection: close\r\n\r\n" | nc localhost 42069