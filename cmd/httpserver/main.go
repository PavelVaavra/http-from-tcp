package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/PavelVaavra/http-from-tcp/internal/request"
	"github.com/PavelVaavra/http-from-tcp/internal/response"
	"github.com/PavelVaavra/http-from-tcp/internal/server"
)

const port = 42069

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
	handlerError := server.HandlerError{}
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		handlerError.StatusCode = response.StatusCodeBadRequest
		handlerError.Message = "Your problem is not my problem\n"
	case "/myproblem":
		handlerError.StatusCode = response.StatusCodeInternalServerError
		handlerError.Message = "Woopsie, my bad\n"
	default:
		handlerError.StatusCode = response.StatusCodeOK
		handlerError.Message = "All good, frfr\n"
	}
	w.Write([]byte(handlerError.Message))
	return &handlerError
}
