package response

import (
	"io"
	"fmt"
	"github.com/PavelVaavra/http-from-tcp/internal/headers"
)

type StatusCode int

const (
	StatusCodeOK StatusCode = 200
	StatusCodeBadRequest StatusCode = 400
	StatusCodeInternalServerError StatusCode = 500 
)

func WriteStatusLine(w io.Writer, statusCode StatusCode, message string) error {
	statusLine := ""
	switch statusCode {
	case StatusCodeOK:
		statusLine = "HTTP/1.1 200 OK\r\n"
	case StatusCodeBadRequest:
		statusLine = "HTTP/1.1 400 " + message + "\r\n"
	case StatusCodeInternalServerError:
		statusLine = "HTTP/1.1 500 " + message + "\r\n"
	}
	_, err := w.Write([]byte(statusLine))
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.Headers{}
	headers.Set("Content-Length", fmt.Sprintf("%v", contentLen))
	headers.Set("Connection", "close")
	headers.Set("Content-Type", "text/plain")

	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		header := k + ": " + v + "\r\n"
		_, err := w.Write([]byte(header))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	if err != nil {
		return err
	}
	return nil
}
