package main

import (
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
	server, err := server.Serve(port, htmlHandler)
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

func textHandler(w *response.Writer, req *request.Request) {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		w.StatusCode = response.StatusCodeBadRequest
		w.Message = "Your problem is not my problem\n"
	case "/myproblem":
		w.StatusCode = response.StatusCodeInternalServerError
		w.Message = "Woopsie, my bad\n"
	default:
		w.StatusCode = response.StatusCodeOK
		w.Message = "All good, frfr\n"
	}
	w.Write()
}

func htmlHandler(w *response.Writer, req *request.Request) {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		w.StatusCode = response.StatusCodeBadRequest
		w.Message = `<html>
  <head>
    <title>%s %s</title>
  </head>
  <body>
    <h1>%s</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`
	case "/myproblem":
		w.StatusCode = response.StatusCodeInternalServerError
		w.Message = `<html>
  <head>
    <title>%s %s</title>
  </head>
  <body>
    <h1>%s</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`
	default:
		w.StatusCode = response.StatusCodeOK
		w.Message = `<html>
  <head>
    <title>%s %s</title>
  </head>
  <body>
    <h1>%s</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`
	}
	w.Write()
}
