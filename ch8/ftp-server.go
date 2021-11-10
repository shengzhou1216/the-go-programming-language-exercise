package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os/exec"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ftp server listening on %s\n", listener.Addr().String())
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		var cmd string
		if bs, _, err := bufio.NewReader(conn).ReadLine(); err != nil {
			log.Println(err)
			break
		} else {
			cmd = string(bs)
		}
		if cmd == "" {
			continue
		}
		log.Printf("read command %s\n", cmd)
		switch cmd {
		case "close":
			conn.Close()
			return
		case "get", "send":
			// transfer file
		default:
			c := exec.Command(cmd)
			if b, e := c.Output(); e == nil {
				io.WriteString(conn, string(b) + "\n")
			} else {
				log.Print(e)
				io.WriteString(conn,e.Error() + "\n")
			}
		}
		time.Sleep(1 * time.Second)
	}
}
