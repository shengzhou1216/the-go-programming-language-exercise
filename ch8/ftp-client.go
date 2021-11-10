package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Please enter command...")
	inputReader := bufio.NewReader(os.Stdin)
	go readConnection(conn)
	for {
		// read from command line
		input, err := inputReader.ReadBytes('\n')
		if err == nil {
			if len(input) > 0 {
				// send command to server
				cmd := string(input)
				if strings.Contains(cmd, "close"){
					log.Println("client closed connection")
					conn.Close()
					return
				}
				conn.Write(input)
			}
		} else {
			log.Println("Error reading", err)
			continue
		}
	}
}

func readConnection(c net.Conn) {
	defer c.Close()
	if _, err := io.Copy(os.Stdout, c); err != nil {
		log.Println(err)
	}
}
