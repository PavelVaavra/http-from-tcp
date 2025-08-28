package request

import (
	"io"
	"strings"
	"errors"
	"strconv"
	"github.com/PavelVaavra/http-from-tcp/internal/headers"
	// "fmt"
)

type Request struct {
	RequestLine RequestLine
	Headers headers.Headers
	Body []byte
	State requestState
}

type requestState int

const (
	requestStateInitialized requestState = iota
	requestStateParsingHeaders
	requestStateParsingBody
	requestStateDone
)


type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *Request) parseSingle(data []byte) (int, error) {
	// fmt.Printf("\"%v|\n", string(data))
	if r.State == requestStateInitialized {
		rl, n, err := parseRequestLine(string(data))
		if err != nil {
			return 0, err
		}
		if n == 0 && err == nil {
			return 0, nil
		}
		r.RequestLine = *rl
		r.State = requestStateParsingHeaders
		return n, nil
	} else if r.State == requestStateParsingHeaders {
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if n == 0 && err == nil {
			return 0, nil
		}
		if done && err == nil {
			r.State = requestStateParsingBody
			return n, nil
		}
		return n, nil
	} else if r.State == requestStateParsingBody {
		contentLengthStr, err := r.Headers.Get("content-length")
		if contentLengthStr == "" && err != nil {
			r.State = requestStateDone
			return len(data), nil
		}
		r.Body = data
		contentLengthInt, err := strconv.Atoi(contentLengthStr)
		if err != nil {
			return 0, errors.New("Content-Length invalid number.")
		}
		if len(r.Body) != contentLengthInt {
			return 0, errors.New("Content-Length doesn't equal to Body length.")
		}
		r.State = requestStateDone
		return len(r.Body), nil
	} else if r.State == requestStateDone {
		return 0,  errors.New("error: trying to read data in a done state")
	} else {
		return 0,  errors.New("error: unknown state")
	}
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.State != requestStateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		totalBytesParsed += n
		if n == 0 && err == nil || r.State == requestStateParsingBody {
			return totalBytesParsed, nil
		}
	}
	return totalBytesParsed, nil
}

func parseRequestLine(s string) (*RequestLine, int, error) {
	requestLine := strings.Split(s, "\r\n")
	// \r\n wasn't found
	if len(requestLine) == 1 {
		return nil, 0, nil
	}

	parts := strings.Split(requestLine[0], " ")
	if len(parts) != 3 {
		return nil, 0, errors.New("Request line doesn't consist of three parts.")
	}
	method := parts[0]
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, 0, errors.New("Method contains other than capital letters.")
		}
	}
	if method != "GET" && method != "POST" && method != "PUT" && method != "DELETE" {
		return nil, 0, errors.New("Not a valid method.")
	}

	requestTarget := parts[1]

	httpVersion := strings.Split(parts[2], "/")[1]
	if httpVersion != "1.1" {
		return nil, 0, errors.New("HTTP version not equal to 1.1")
	}

	return &RequestLine{
		HttpVersion: httpVersion,
		RequestTarget: requestTarget,
		Method: method,
	}, len(requestLine[0]) + len("\r\n"), nil
}

const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {	
	buff := make([]byte, bufferSize, bufferSize)

	readToIndex := 0

	req := Request{
		RequestLine: RequestLine{},
		Headers: headers.Headers{},
		State: requestStateInitialized,
	}
	
	for req.State != requestStateDone {
		if (readToIndex + 1) >= len(buff) {
			buff = append(buff, make([]byte, len(buff), cap(buff))...)
		}
		n, err := reader.Read(buff[readToIndex:])
		if n == 0 && err == io.EOF {
			if req.State != requestStateParsingBody {
				return nil, errors.New("No requestStateParsingBody after EOF.")
			}
			n, err = req.parse(buff[:readToIndex])
			if err != nil {
				return nil, err
			}
			req.State = requestStateDone
			continue
		}
		readToIndex += n
		if req.State != requestStateParsingBody {
			n, err = req.parse(buff[:readToIndex])
			if err != nil {
				return nil, err
			}
			if n != 0 {
				copy(buff, buff[n:])
				readToIndex -= n
			}
		}
	}
	
	return &req, nil
}