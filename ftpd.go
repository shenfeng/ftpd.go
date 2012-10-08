package main

import (
	"flag"
	"log"
	"net"
)

var addr = flag.String("addr", ":2121", "The add to listen on (default ':2121')")

func handleConnection(c net.Conn) {
	log.Print(c)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()
	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		handleConnection(conn)
	}
}
