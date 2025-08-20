package request

import (
	"io"
	"strings"
	"errors"
	"fmt"
)

// Update your parseRequestLine to return the number of bytes it consumed. If it can't find an \r\n (this is important!) it should 
// return 0 and no error. This just means that it needs more data before it can parse the request line.
// 
// Add a new internal "enum" (I just used an int) to your Request struct to track the state of the parser. For now, you just need 2 states:
// "initialized"
// "done".
// 
// Implement a new func (r *Request) parse(data []byte) (int, error) method.
	// It accepts the next slice of bytes that needs to be parsed into the Request struct
	// It updates the "state" of the parser, and the parsed RequestLine field.
	// It returns the number of bytes it consumed (meaning successfully parsed) and an error if it encountered one.
//
// Update the RequestFromReader function.
	// Instead of reading all the bytes, and then parsing the request line, it should use a loop to continually read from the reader and parse new chunks using the parse method.
	// 
	// The loop should continue until the parser is in the "done" state.
	// 
	// You'll need to keep track of:
		// A buffer to read data into ([]byte). I started with a size of 8 and grew it as needed. I also shifted data in and out of it so I don't need to keep storing already-parsed data.
		// How many bytes you've read from the reader
		// How many bytes you've parsed from the buffer
	// The end result is the same (aside from the fact that it properly handles chunks as they arrive) in that it returns a parsed Request struct once the reader is exhausted.

type Request struct {
	RequestLine RequestLine
	// 0..initialized
	// 100..done
	State int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *Request) parse(data []byte) (int, error) {
// If the state of the parser is "initialized", it should call parseRequestLine.
	// If there is an error, it should just return the error.
	// If zero bytes are parsed, but no error is returned, it should return 0 and nil: it needs more data.
	// If bytes are consumed successfully, it should update the .RequestLine field and change the state to "done".
// If the state of the parser is "done", it should return an error that says something like "error: trying to read data in a done state"
// If the state is anything else, it should return an error that says something like "error: unknown state"
	if r.State == 0 {
		rl, n, err := parseRequestLine(string(data))
		if err != nil {
			return 0, err
		}
		if n == 0 && err == nil {
			return 0, nil
		}
		r.RequestLine = *rl
		r.State = 100
		return n, nil
	} else if r.State == 100 {
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
// It shouldn't call io.ReadAll anymore. Instead, it should create a new buffer: buf := make([]byte, bufferSize, bufferSize). 
// Set bufferSize as a constant at the top of the file, and for now, just a size of 8. We want to test with small buffers to make sure 
// our parser can handle it.
// 
// Create a new readToIndex variable and set it to 0. This will keep track of how much data we've read from the io.Reader into the buffer.
// 
// Create a new Request struct and set the state to "initialized".
// 
// While the state of the parser is not "done":
	// If the buffer is full (we've read data into the entire buffer), grow it. Create a new slice that's twice the size and copy the 
	// old data into the new slice.
	// 
	// Read from the io.Reader into the buffer starting at readToIndex.
		// If you hit the end of the reader (io.EOF) set the state to "done" and break out of the loop.
		// 
		// Update readToIndex with the number of bytes you actually read
		// 
		// Call r.parse passing the slice of the buffer that has data that you've actually read so far
		// 
		// Remove the data that was parsed successfully from the buffer (this keeps our buffer small and memory efficient). 
		// I used the copy function and a new slice to do this.
		// 
		// Decrement the readToIndex by the number of bytes that were parsed so that it matches the new length of the buffer.
	
	buff := make([]byte, bufferSize, bufferSize)

	readToIndex := 0

	req := Request{
		RequestLine: RequestLine{},
		State: 0,
	}
	
	for req.State != 100 {
		if (readToIndex + 1) >= cap(buff) {
			buff = append(buff, make([]byte, bufferSize, bufferSize)...)
			// fmt.Printf("len(buff) = %v, cap(buff) = %v\n", len(buff), cap(buff))
		}
		n, err := reader.Read(buff[readToIndex:])
		if n == 0 && err == io.EOF {
			req.State = 100
			continue
		}
		readToIndex += n
		// fmt.Println(readToIndex)
		n, err = req.parse(buff[:readToIndex])
		if err != nil {
			return nil, err
		}
		readToIndex -= n
		// fmt.Println(readToIndex)
	}
	
	return &req, nil
}