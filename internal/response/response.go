package response

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/PavelVaavra/http-from-tcp/internal/headers"
)

type Writer struct {
	StatusCode StatusCode
	Message    string
	Conn       net.Conn
}

func (w *Writer) Write() {
	statusCodeStr := strconv.Itoa(int(w.StatusCode))
	statusPhraseStr := statusPhrase[statusCodeStr]

	err := WriteStatusLine(w.Conn, statusCodeStr, statusPhraseStr)
	if err != nil {
		return
	}

	isHtml := strings.Contains(w.Message, "<html>")

	contentType := "plain"
	if isHtml {
		contentType = "html"
		w.Message = fmt.Sprintf(w.Message, statusCodeStr, statusPhraseStr, statusPhraseStr)
	}

	headers := GetDefaultHeaders(len(w.Message), contentType)

	err = WriteHeaders(w.Conn, headers)
	if err != nil {
		return
	}

	_, err = w.Conn.Write([]byte(w.Message))
	if err != nil {
		return
	}
}

type StatusCode int

const (
	StatusCodeOK                  StatusCode = 200
	StatusCodeBadRequest          StatusCode = 400
	StatusCodeInternalServerError StatusCode = 500
)

var statusPhrase = map[string]string{
	"200": "OK",
	"400": "Bad Request",
	"500": "Internal Server Error",
}

func WriteStatusLine(w io.Writer, statusCode string, statusPhrase string) error {
	statusLine := "HTTP/1.1 " + statusCode + " " + statusPhrase + "\r\n"
	_, err := w.Write([]byte(statusLine))
	return err
}

func GetDefaultHeaders(contentLen int, contentType string) headers.Headers {
	headers := headers.Headers{}
	headers.Set("Content-Length", fmt.Sprintf("%v", contentLen))
	headers.Set("Connection", "close")
	headers.Set("Content-Type", fmt.Sprintf("text/%v", contentType))

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
