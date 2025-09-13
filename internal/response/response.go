package response

import (
	"fmt"
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
	BodyVideo    []byte
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

func (w *Writer) WriteBodyVideo() error {
	_, err := w.Conn.Write(w.BodyVideo)
	return err
}

func (w *Writer) WriteChunkedBody(p []byte) error {
	chunk := []byte(fmt.Sprintf("%X", len(p)))
	chunk = append(chunk, []byte("\r\n")...)
	chunk = append(chunk, p...)
	chunk = append(chunk, []byte("\r\n")...)
	_, err := w.Conn.Write(chunk)
	return err
}

func (w *Writer) WriteChunkedBodyDone() error {
	_, err := w.Conn.Write([]byte("0\r\n\r\n"))
	return err
}

func (w *Writer) WriteTrailers() error {
	_, err := w.Conn.Write([]byte("0\r\n"))
	if err != nil {
		return err
	}
	for k, v := range w.Trailers {
		trailer := k + ": " + v + "\r\n"
		_, err := w.Conn.Write([]byte(trailer))
		if err != nil {
			return err
		}
	}
	_, err = w.Conn.Write([]byte("\r\n"))
	if err != nil {
		return err
	}
	return nil
}
