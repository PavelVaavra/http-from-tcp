package response

import (
	"io"
	"net"
	"strconv"

	"github.com/PavelVaavra/http-from-tcp/internal/headers"
)

type StatusCode int

const (
	StatusCodeOK                  StatusCode = 200
	StatusCodeBadRequest          StatusCode = 400
	StatusCodeInternalServerError StatusCode = 500
)

type Writer struct {
	StatusCode   StatusCode
	StatusPhrase string
	Headers      headers.Headers
	BodyText     string
	BodyChunked  io.ReadCloser
	Trailers     headers.Headers
	Conn         net.Conn
}

func (w *Writer) WriteStatusLine() error {
	statusLine := "HTTP/1.1 " + strconv.Itoa(int(w.StatusCode)) + " " + w.StatusPhrase + "\r\n"
	_, err := w.Conn.Write([]byte(statusLine))
	return err
}

func (w *Writer) WriteHeaders() error {
	for k, v := range w.Headers {
		header := k + ": " + v + "\r\n"
		_, err := w.Conn.Write([]byte(header))
		if err != nil {
			return err
		}
	}
	_, err := w.Conn.Write([]byte("\r\n"))
	if err != nil {
		return err
	}
	return nil
}

func (w *Writer) WriteBody() error {
	_, err := w.Conn.Write([]byte(w.BodyText))
	return err
}
