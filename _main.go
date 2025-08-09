package main

import (
	"fmt"
	"os"
	"io"
	"log"
	"strings"
)

const file = "./messages.txt"

func main() {
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("could not open file %s: error: %s\n", file, err)
	}

	ch := getLinesChannel(f)

	for line := range ch {
		fmt.Printf("read: %s\n", line)
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
				log.Fatalf("could not read from file %s: error: %s\n", file, err)
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