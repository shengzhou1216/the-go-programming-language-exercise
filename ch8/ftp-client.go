package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"path"
	"strings"
)

func newConn() (net.Conn, error) {
	return net.Dial("tcp", ":8000")
}

func main() {
	conn, err := newConn()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Please enter command...")
	inputReader := bufio.NewReader(os.Stdin)

	// go readConnection(conn, fileChan, outputChan)
	for {
		// read from command line
		input, err := inputReader.ReadBytes('\n')
		if err == nil {
			if len(input) > 0 {
				// send command to server
				cmd := string(input)
				cmd = strings.Trim(strings.TrimSpace(cmd), "\n")
				cmds := strings.Split(cmd, " ")
				switch cmds[0] {
				case "close":
					log.Println("client closed connection")
					conn.Close()
					return
				case "get":
					if len(cmds) <= 1 {
						log.Println("Please enter filename")
						continue
					}
					// todo: 接收文件的流需要与接收输出的流分离开; 即 如何区分文件 与 命令输出结果
					go readFile(cmds[1])
				default:
					go readConnection(conn)
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

func readFile(fileName string) {
	c, err := newConn()
	if err != nil {
		log.Println("Error new connection", err)
		return
	}
	defer c.Close()
	// get base name, default save to current
	fileName = path.Base(fileName)
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
	}
	_, err = io.Copy(file, c)
	if err != nil {
		log.Println(err)
	}
	file.Close()
}
