package response

import (
	"io"
	"fmt"
	"github.com/PavelVaavra/http-from-tcp/internal/headers"
)

type StatusCode int

const (
	OK StatusCode = 200
	BadRequest StatusCode = 400
	InternalServerError StatusCode = 500 
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	statusLine := ""
	switch statusCode {
	case OK:
		statusLine = "HTTP/1.1 200 OK\r\n"
	case BadRequest:
		statusLine = "HTTP/1.1 400 Bad Request\r\n"
	case InternalServerError:
		statusLine = "HTTP/1.1 500 Internal Server Error\r\n"
	}
	_, err := w.Write([]byte(statusLine))
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.Headers{}
	headers["content-length"] = fmt.Sprintf("%v", contentLen)
	headers["connection"] = "close"
	headers["content-type"] = "text/plain"

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
