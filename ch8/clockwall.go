package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	var NewYorkServer string
	var TokyoServer string
	var LondonServer string

	log.Println(os.Args[:])

	flag.StringVar(&NewYorkServer, "NewYork", "", "NewYork server addr")
	flag.StringVar(&TokyoServer, "Tokyo", "", "NewYork server addr")
	flag.StringVar(&LondonServer, "London", "", "London server addr")

	flag.Parse()
	go server(NewYorkServer)
	go server(TokyoServer)
	go server(LondonServer)
	for {
		
	}
}

func server(server string) {
	conn, err := net.Dial("tcp", server)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	mustCopy(os.Stdout, conn)
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
