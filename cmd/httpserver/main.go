package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/PavelVaavra/http-from-tcp/internal/headers"
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

// func chunkHandler(w *response.Writer, req *request.Request) {
// 	url := "https://httpbin.org"
// 	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {
// 		url += strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin")
// 	}

// 	res, err := http.Get(url)
// 	if err != nil {
// 		fmt.Errorf("http.Get(\"%v\" err - %v\n", url, err.Error())
// 	}
// 	defer res.Body.Close()

// 	w.StatusCode = response.StatusCode(res.StatusCode)

// 	w.WriteChunked(res.Body)
// }

// Add a new proxy handler to your server that maps /httpbin/x to https://httpbin.org/x, supporting both proxying and chunked responsing.

// 	I used the strings.HasPrefix and strings.TrimPrefix functions to handle routing and route parsing.

// 	Be sure to remove the Content-Length header from the response, and add the Transfer-Encoding: chunked header.

// 	I used http.Get to make the request to httpbin.org and resp.Body.Read to read the response body. I used a buffer size of 1024 bytes, and then
// 	printed n on every call to Read so that I could see how much data was being read. Use n as your chunk size and write that chunk data back to the
// 	client as soon as you get it from httpbin.org. It's pretty cool to see how the data can be forwarded in real-time!

// 	I recommend using netcat to test your chunked responses. Curl will abstract away the chunking for you, so you won't see your hex and cr and lf
// 	characters in your terminal if you use curl. I used this command to see my raw chunked response:

// echo -e "GET /httpbin/stream/100 HTTP/1.1\r\nHost: localhost:42069\r\nConnection: close\r\n\r\n" | nc localhost 42069

// GET localhost:42069/httpbin/stream/100 will trigger a handler on our server that sends a request to https://httpbin.org/stream/100
// and then forwards the response back to the client chunk by chunk.

// HTTP/1.1 200 OK
// Content-Type: text/plain
// Transfer-Encoding: chunked

// 1E
// I could go for a cup of coffee
// B
// But not Java
// 12
// Never go full Java
// 0
