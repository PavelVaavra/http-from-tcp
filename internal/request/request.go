package request

import (
	"io"
	"log"
	"strings"
	"errors"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	b, err := io.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}
	parts := strings.Split(string(b), "\r\n")
	
	rl, err := parseRequestLine(parts[0])
	if err != nil {
		return nil, err
	}

	return &Request{ RequestLine: *rl }, nil
}

func parseRequestLine(s string) (*RequestLine, error) {
	parts := strings.Split(s, " ")
	if len(parts) != 3 {
		return nil, errors.New("Request line doesn't consist of three parts.")
	}
	method := parts[0]
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, errors.New("Method contains other than capital letters.")
		}
	}
	if method != "GET" && method != "POST" && method != "PUT" && method != "DELETE" {
		return nil, errors.New("Not a valid method.")
	}

	requestTarget := parts[1]

	httpVersion := strings.Split(parts[2], "/")[1]
	if httpVersion != "1.1" {
		return nil, errors.New("HTTP version not equal to 1.1")
	}

	return &RequestLine{
		HttpVersion: httpVersion,
		RequestTarget: requestTarget,
		Method: method,
	}, nil
}