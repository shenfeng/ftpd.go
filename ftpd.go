package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

var addr = flag.String("addr", ":2121", "The add to listen on (default ':2121')")

func handle(command string, msgs []string, c net.Conn) {
	switch command {
	case "user":
		fmt.Fprintf(c, "331 Username ok, send password.\n\n")
	case "pass":
		fmt.Fprintf(c, "230 Login successful.\r\n")
	case "syst":
		fmt.Fprintf(c, "215 UNIX Type: L8\r\n")
	case "port":
		nums := strings.Split(msgs[0], ",")
		h, _ := strconv.ParseInt(nums[4], 10, 32)
		s, _ := strconv.ParseInt(nums[5], 10, 32)
		port := h*256 + s
		ip := strings.Join(nums[0:4], ".")
		ip = ip + ":" + strconv.Itoa(int(port))
		log.Print("Opened passive connection at: ", ip)

	}
}

func handleConnection(c net.Conn) {
	fmt.Fprintf(c, "220 ftpd.go.\r\n")
	reader := bufio.NewReader(c)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Print(err)
			return
		}
		log.Print("Reqeust: ", line)
		command := strings.TrimSpace(strings.ToLower(line[0:4]))
		msgs := strings.Split(strings.Trim(line, "\r\n "), " ")[1:]
		handle(command, msgs, c)
	}
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
