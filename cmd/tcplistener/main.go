package main

import (
	"fmt"
	"io"
	"log"
	"strings"
	"net"
)

const ipPort = "127.0.0.1:42069"

func main() {
	l, err := net.Listen("tcp", ipPort)
	if err != nil {
		log.Fatalf("could not open listener on %s: error: %s\n", ipPort, err)
	}
	defer l.Close()

	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("A connection has been accepted...")
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go func(c net.Conn) {
			ch := getLinesChannel(c)
			
			for line := range ch {
				fmt.Println(line)
			}
			fmt.Println("The connection has been closed...")
		}(conn)
	}	
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func(){
		defer f.Close()
		defer close(ch)

		eightBytes := make([]byte, 8, 8)
		line := ""
		for {
			n, err := f.Read(eightBytes)
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}
			if n == 0 && err == io.EOF {
				ch <- line
				break
			}

			sEightBytes := string(eightBytes[:n])
			parts := strings.Split(sEightBytes, "\n")
			partsLength := len(parts)
			// pseudocode:
			// 	switch partsLength {
			// 	case 1:
			// 		line += parts[0]
			// 	case 2:
			// 		line += parts[0]
			// 		print(line)
					
			// 		line = parts[partsLength - 1]
			// 	case 3:
			// 		line += parts[0]
			// 		print(line)
					
			// 		loop over middle of the parts (without the first and the last element) and print every part

			// 		line = parts[partsLength - 1]
			// 	}
			line += parts[0]
			if partsLength == 1 {
				continue
			} else {
				ch <- line
				if partsLength >= 3 {
					for _, part := range parts[1:partsLength - 2] {
						ch <- part
					}
				}
			}
			line = parts[partsLength - 1]
		}
	}()

	return ch
}