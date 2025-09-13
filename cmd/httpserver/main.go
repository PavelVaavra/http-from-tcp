package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/PavelVaavra/http-from-tcp/internal/headers"
	"github.com/PavelVaavra/http-from-tcp/internal/request"
	"github.com/PavelVaavra/http-from-tcp/internal/response"
	"github.com/PavelVaavra/http-from-tcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, videoHandler)
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
		w.StatusPhrase = "Bad Request"
		w.BodyText = "Your problem is not my problem\n"
	case "/myproblem":
		w.StatusCode = response.StatusCodeInternalServerError
		w.StatusPhrase = "Internal Server Error"
		w.BodyText = "Woopsie, my bad\n"
	default:
		w.StatusCode = response.StatusCodeOK
		w.StatusPhrase = "OK"
		w.BodyText = "All good, frfr\n"
	}

	w.Headers = headers.Headers{
		"Connection":     "close",
		"Content-Length": fmt.Sprintf("%v", len(w.BodyText)),
		"Content-Type":   "text/plain",
	}

	w.WriteStatusLine()
	w.WriteHeaders()
	w.WriteBody()
}

func htmlHandler(w *response.Writer, req *request.Request) {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		w.StatusCode = response.StatusCodeBadRequest
		w.StatusPhrase = "Bad Request"
		w.BodyText = fmt.Sprintf(`<html>
  <head>
    <title>%s %s</title>
  </head>
  <body>
    <h1>%s</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`, strconv.Itoa(int(w.StatusCode)), w.StatusPhrase, w.StatusPhrase)
	case "/myproblem":
		w.StatusCode = response.StatusCodeInternalServerError
		w.StatusPhrase = "Internal Server Error"
		w.BodyText = fmt.Sprintf(`<html>
  <head>
    <title>%s %s</title>
  </head>
  <body>
    <h1>%s</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`, strconv.Itoa(int(w.StatusCode)), w.StatusPhrase, w.StatusPhrase)
	default:
		w.StatusCode = response.StatusCodeOK
		w.StatusPhrase = "OK"
		w.BodyText = fmt.Sprintf(`<html>
  <head>
    <title>%s %s</title>
  </head>
  <body>
    <h1>%s</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`, strconv.Itoa(int(w.StatusCode)), w.StatusPhrase, w.StatusPhrase)
	}

	w.Headers = headers.Headers{
		"Connection":     "close",
		"Content-Length": fmt.Sprintf("%v", len(w.BodyText)),
		"Content-Type":   "text/html",
	}

	w.WriteStatusLine()
	w.WriteHeaders()
	w.WriteBody()
}

func chunkHandler(w *response.Writer, req *request.Request) {
	url := "https://httpbin.org"
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {
		url += strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin")
	}

	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("http.Get(\"%v\" err - %v\n", url, err.Error())
	}
	defer res.Body.Close()

	w.StatusCode = response.StatusCode(res.StatusCode)
	w.StatusPhrase = strings.Split(res.Status, " ")[1]
	w.BodyChunked = res.Body

	w.Headers = headers.Headers{
		"Connection":        "close",
		"Transfer-Encoding": "chunked",
		"Content-Type":      "text/plain",
	}

	w.WriteStatusLine()
	w.WriteHeaders()

	buff := make([]byte, 1024)
	totalBytes := 0
	for {
		n, err := w.BodyChunked.Read(buff)
		if n == 0 && err == io.EOF {
			break
		}
		err = w.WriteChunkedBody(buff[:n])
		if err != nil {
			fmt.Printf("w.WriteChunkedBody: err - %v\n", err.Error())
		}
		totalBytes += n
	}
	err = w.WriteChunkedBodyDone()
	if err != nil {
		fmt.Printf("w.WriteChunkedBodyDone: err - %v\n", err.Error())
	}
}

func chunkHandlerTrailers(w *response.Writer, req *request.Request) {
	url := "https://httpbin.org"
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {
		url += strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin")
	}

	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("http.Get(\"%v\" err - %v\n", url, err.Error())
	}
	defer res.Body.Close()

	w.StatusCode = response.StatusCode(res.StatusCode)
	w.StatusPhrase = strings.Split(res.Status, " ")[1]
	w.BodyChunked = res.Body

	w.Headers = headers.Headers{
		"Connection":        "close",
		"Transfer-Encoding": "chunked",
		"Content-Type":      "text/html`",
		"Trailer":           "X-Content-SHA256, X-Content-Length",
	}

	w.WriteStatusLine()
	w.WriteHeaders()

	buff := make([]byte, 1024)
	var body []byte
	totalBytes := 0
	for {
		n, err := w.BodyChunked.Read(buff)
		if n == 0 && err == io.EOF {
			break
		}
		body = append(body, buff[:n]...)
		err = w.WriteChunkedBody(buff[:n])
		if err != nil {
			fmt.Printf("w.WriteChunkedBody: err - %v\n", err.Error())
		}
		totalBytes += n
	}

	w.Trailers = headers.Headers{
		"X-Content-SHA256": fmt.Sprintf("%x", sha256.Sum256(body)),
		"X-Content-Length": strconv.Itoa(totalBytes),
	}

	w.WriteTrailers()
}

func videoHandler(w *response.Writer, req *request.Request) {
	if req.RequestLine.RequestTarget == "/video" {
		w.StatusCode = response.StatusCodeOK
		w.StatusPhrase = "OK"

		data, err := os.ReadFile("assets/vim.mp4")
		if err != nil {
			fmt.Printf("os.ReadFile(\"../../assets/vim.mp4\"): %v\n", err.Error())
		}
		w.BodyVideo = data
	}

	w.Headers = headers.Headers{
		"Connection":     "close",
		"Content-Length": fmt.Sprintf("%v", len(w.BodyVideo)),
		"Content-Type":   "video/mp4",
	}

	w.WriteStatusLine()
	w.WriteHeaders()
	w.WriteBodyVideo()
}
