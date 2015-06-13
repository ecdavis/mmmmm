package main

import (
	"bufio"
	"log"
	"net"
	)

type EchoClient struct {
	conn net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
	writeChannel chan string
	closeChannel chan bool
}

func NewEchoClient(conn net.Conn) (e *EchoClient) {
	e = new(EchoClient)
	e.conn = conn
	e.reader = bufio.NewReader(conn)
	e.writer = bufio.NewWriter(conn)
	e.writeChannel = make(chan string)
	e.closeChannel = make(chan bool)
	return
}

func readLines(e *EchoClient) {
	for {
		line, err := e.reader.ReadString('\n')
		if err != nil {
			log.Print("readLines:", err)
			close(e.writeChannel)
			return
		}
		select {
		case e.writeChannel <- line:
			continue
		case <- e.closeChannel:
			close(e.writeChannel)
			return
		}
	}
}

func writeLines(e *EchoClient) {
	for line := range e.writeChannel {
		_, err := e.writer.WriteString(line)
		if err != nil {
			log.Print("writeLines:", err)
			break
		}
		e.writer.Flush()
	}
	e.closeChannel <- true
	e.conn.Close()
}

func runServer() error {
	ln, err := net.Listen("tcp", ":4040")
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print("Accept:", err)
		}
		e := NewEchoClient(conn)
		go readLines(e)
		go writeLines(e)
	}
}

func main() {
	err := runServer()
	if err != nil {
		log.Fatal("runServer:", err);
	}
}
