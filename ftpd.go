package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

var addr = flag.String("addr", ":2121", "The addr to listen (':2121')")

type FTPConn struct {
	command net.Conn
	data    net.Conn
	pasv    bool
	root    string
	cwd     string
}

func list(c *FTPConn, msgs []string) {
	cwd := c.cwd
	if len(msgs) > 0 {
		cwd = path.Join(cwd, msgs[0])
	}
	cmd := exec.Command("ls", "-l", cwd)
	log.Print("dir: ", cwd)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		c.data.Close()
		reply(c, "550 Folder not found.")
	} else {
		if c.pasv {
			reply(c, "150 File status okay. About to open data connection.")
		} else {
			reply(c, "125 Data connection already open. Transfer starting.")
		}
		d := strings.Join(strings.Split(out.String(), "\n"), "\r\n") + "\r\n"
		fmt.Fprintf(c.data, d)
		c.data.Close()
		reply(c, "226 Transfer complete.")
	}
}

func reply(c *FTPConn, msg string) {
	fmt.Fprintf(c.command, msg+"\r\n")
}

func port(msgs []string, c *FTPConn) {
	nums := strings.Split(msgs[0], ",")
	h, _ := strconv.ParseInt(nums[4], 10, 32)
	s, _ := strconv.ParseInt(nums[5], 10, 32)
	port := h*256 + s
	ip := strings.Join(nums[0:4], ".") + ":" + strconv.Itoa(int(port))
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		c.pasv = true
		reply(c, "501 Can't connect to a foreign address")
	} else {
		log.Print("Opened passive connection at: ", ip)
		c.data = conn
		reply(c, "200 Active data connection established.")
	}
}

func pasv(c *FTPConn) {
	if addr, ok := c.command.LocalAddr().(*net.TCPAddr); ok {
		ip := addr.IP.String()
		for i := 40000; i < 40050; i++ {
			ln, err := net.Listen("tcp", ip+":"+strconv.Itoa(i))
			if err == nil {
				ip = strings.Join(strings.Split(ip, "."), ",")
				h := strconv.Itoa(i >> 8)
				l := strconv.Itoa(i % 256)
				msg := "227 Entering passive mode (" + ip + "," + h + "," + l + ")"
				log.Print(msg + "; port " + strconv.Itoa(i))
				reply(c, msg)
				go (func() {
					data, err := ln.Accept()
					if err == nil {
						c.data = data
						c.pasv = true
					}
				})()
			}
		}
	}
}

func retr(c *FTPConn, msgs []string) {
	f := path.Join(c.cwd, msgs[0])
	file, err := os.Open(f)
	log.Print("get file: " + f)
	if err == nil {
		if c.pasv {
			reply(c, "150 File status okay. About to open data connection.")
		} else {
			reply(c, "125 Data connection already open. Transfer starting.")
		}
		io.Copy(c.data, file)
		file.Close()
		c.data.Close()

		stat, _ := os.Stat(f)
		size := strconv.FormatInt(stat.Size(), 10)
		reply(c, "226 Closing data connection, sent "+size+" bytes")
		// reply(c, "226 Transfer complete.")
	} else {
		reply(c, "550 file not found.")
	}
}

func handle(command string, msgs []string, c *FTPConn) {
	switch command {
	case "user":
		reply(c, "331 Username ok, send password.")
	case "pass":
		reply(c, "230 Login successful.")
	case "syst":
		reply(c, "215 UNIX Type: L8.")
	case "port":
		port(msgs, c)
	case "list":
		list(c, msgs)
	case "pwd":
		reply(c, "257 \"/\" is the current directory.")
	case "type":
		reply(c, "200 Type set to binary.")
	case "cwd":
		dir := msgs[0]
		c.cwd = path.Join(c.root, dir)
		reply(c, "250 Directory changed to "+dir+".")
	case "pasv":
		pasv(c)
	case "retr":
		retr(c, msgs)
	case "quit":
		reply(c, "221 bye")
		c.command.Close()
	}
}

func handleConnection(c *FTPConn) {
	reply(c, "220 ftpd.go.")
	reader := bufio.NewReader(c.command)
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

	wd, _ := os.Getwd()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		c := FTPConn{command: conn, pasv: false, cwd: wd, root: wd}
		go handleConnection(&c)
	}
}
