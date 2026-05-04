package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Println("Server is running on http://localhost:4221")

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		headersChan := handleConnection(conn)
		go func(ch <-chan string) {
			for line := range ch {
				fmt.Println(line)
			}
		}(headersChan)
	}
}

func handleConnection(conn net.Conn) <-chan string {
	out := make(chan string)

	go func() {
		defer conn.Close() // close connection when done
		defer close(out)

		scanner := bufio.NewScanner(conn)
		var requestPath string

		for scanner.Scan() {
			line := scanner.Text()
			// HTTP headers end with a blank line.
			// We MUST break here, otherwise the scanner will wait forever.
			if line == "" {
				break
			}
			out <- line

			if strings.HasPrefix(line, "GET ") {
				parts := strings.Fields(line)
				if len(parts) >= 2 {
					requestPath = parts[1]
				}
			}
		}
		if err := scanner.Err(); err != nil {
			log.Printf("Error reading from connection: %v", err)
		}
		if requestPath == "/" {
			conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		} else {
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		}
	}()
	return out
}
