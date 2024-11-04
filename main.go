package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	index, err := os.ReadFile("index.html") // cache index.html
	if err != nil {
		log.Fatal(err)
	}

	l, err := net.Listen("tcp", ":10499")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening at %s", l.Addr().String())
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handle(conn, index)
	}
}

func handle(conn net.Conn, index []byte) {
	defer conn.Close()

	log.Printf("Handling connection from %s", conn.RemoteAddr())
	if err := conn.SetDeadline(time.Now().Add(1 * time.Minute)); err != nil {
		log.Printf("Error setting deadline: %s", err)
		return
	}

	reader := bufio.NewReader(conn)
	s, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Error reading line: %s", err)
		return
	}

	if !strings.HasPrefix(s, "GET") {
		_, err := conn.Write([]byte("HTTP/1.1 405 your method is not allowed\r\n\r\n"))
		if err != nil {
			log.Printf("Error writing response: %s", err)
		}
		return
	}

	if _, err := conn.Write([]byte("HTTP/1.1 200 OK\r\n")); err != nil {
		return
	}
	if _, err := conn.Write([]byte("X-Custom-Message: It's working\r\n")); err != nil {
		return
	}
	if _, err := conn.Write([]byte("Cache-Control: no-cache\r\n")); err != nil {
		return
	}

	pieces := strings.Fields(s)
	if pieces[1] == "/" {
		if _, err := conn.Write([]byte("Content-Type: text/html\r\n")); err != nil {
			return
		}
		if _, err := conn.Write([]byte("\r\n")); err != nil {
			return
		}
		if _, err := conn.Write(index); err != nil {
			log.Printf("Error writing index page: %s", err)
		}
		return
	}

	if _, err := conn.Write([]byte("Content-Type: text/event-stream\r\n")); err != nil {
		return
	}
	if _, err := conn.Write([]byte("\r\n")); err != nil {
		return
	}

	comment := false
	for {
		var s string
		if comment {
			s = ": this is a comment\n\n"
		} else {
			s = fmt.Sprintf("event: time\ndata: {\"time\": %d}\n\n", time.Now().Unix())
		}
		comment = !comment
		if _, err := conn.Write([]byte(s)); err != nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
}
