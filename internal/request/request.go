package request

import (
	"io"
	"strings"
	"errors"
)

type Request struct {
	RequestLine RequestLine
	// 0..initialized
	// 100..done
	State requestState
}

type requestState int

const (
	requestStateInitialized requestState = iota
	requestStateDone
)


type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *Request) parse(data []byte) (int, error) {
	if r.State == requestStateInitialized {
		rl, n, err := parseRequestLine(string(data))
		if err != nil {
			return 0, err
		}
		if n == 0 && err == nil {
			return 0, nil
		}
		r.RequestLine = *rl
		r.State = requestStateDone
		return n, nil
	} else if r.State == requestStateDone {
		return 0,  errors.New("error: trying to read data in a done state")
	} else {
		return 0,  errors.New("error: unknown state")
	}
}

func parseRequestLine(s string) (*RequestLine, int, error) {
	// fmt.Printf("\"%v\"\n", s)
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
		State: requestStateInitialized,
	}
	
	for req.State != requestStateDone {
		if (readToIndex + 1) >= len(buff) {
			buff = append(buff, make([]byte, len(buff), cap(buff))...)
		}
		n, err := reader.Read(buff[readToIndex:])
		if n == 0 && err == io.EOF {
			req.State = requestStateDone
			continue
		}
		readToIndex += n
		n, err = req.parse(buff[:readToIndex])
		if err != nil {
			return nil, err
		}
		readToIndex -= n
	}
	
	return &req, nil
}