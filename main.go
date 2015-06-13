package main

import (
	"bufio"
	"io"
	"log"
	"net"
	)

func echoConnection(conn net.Conn) {
	rd := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			log.Print("ReadString:", err)
			break
		}
		_, err = w.WriteString(line)
		if err != nil {
			log.Print("WriteString:", err)
			break
		}
		w.Flush()
	}

	conn.Close()
}

func echoServer() error {
	ln, err := net.Listen("tcp", ":4040")
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print("Accept:", err)
		}
		go echoConnection(conn)
	}
}

func main() {
	err := echoServer()
	if err != nil {
		log.Fatal("echoServer:", err);
	}
}
