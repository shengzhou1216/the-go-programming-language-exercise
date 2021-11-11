package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Client struct {
	Conn *net.Conn
}

type Server struct {
	Clients []*Client
	Port    string
	net.Listener
}

func NewClient(Conn *net.Conn) *Client {
	return &Client{Conn: Conn}
}

func NewServer(scheme string, port string) (*Server, error) {
	listener, err := net.Listen(scheme, port)
	if err != nil {
		return nil, err
	}
	return &Server{
		Listener: listener,
		Port:     port,
		Clients:  make([]*Client, 0),
	}, nil
}

func main() {
	server, err := NewServer("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ftp server listening on %s\n", server.Addr().String())
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		server.Clients = append(server.Clients, NewClient(&conn))
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		cmd, err := bufio.NewReader(conn).ReadString('\n')
		cmd = strings.Trim(strings.TrimSpace(cmd), "\n")
		if err != nil {
			log.Println(err)
			break
		}
		if cmd == "" {
			log.Println("empty command")
			continue
		}
		log.Printf("read command %s\n", cmd)
		cmds := strings.Split(cmd, " ")
		switch cmds[0] {
		case "close":
			conn.Close()
			return
		case "get":
			// transfer file
			if len(cmds) <= 1 {
				log.Println("Please enter filename")
				continue
			}
			log.Println("received:", cmds[1])
			writeFile(conn, cmds[1])
			// fileInfo, err := os.Stat(cmds[1])
			// if err != nil {
			// 	write(conn, err.Error())
			// 	continue
			// }
			// if !fileInfo.IsDir() {
			// 	// send file
			// 	writeFile(conn, fileInfo)
			// } else {
			// 	// todo: read dir and retrun file to client
			// }
		default:
			c := exec.Command(cmds[0], cmds[1:]...)
			if b, err := c.Output(); err == nil {
				write(conn, string(b))
			} else {

				write(conn, err.Error())
				continue
			}
		}
		time.Sleep(1 * time.Second)
	}
}

// wirte to client
func write(conn net.Conn, res string) {
	if !strings.HasSuffix(res, "\n") {
		res += "\n"
	}
	io.WriteString(conn, res)
}

// wirte file to client
func writeFile(conn net.Conn, fileName string) {
	if conn == nil {
		return
	}
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
	if err != nil {
		write(conn, err.Error())
		return
	}
	defer file.Close()
	bs, err := io.ReadAll(file)
	if err != nil {
		write(conn, err.Error())
		return
	}
	_, err = bufio.NewWriter(conn).Write(bs)
	if err != nil {
		write(conn, err.Error())
		return
	}
}
