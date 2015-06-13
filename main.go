package main

import (
	"bufio"
	"log"
	"net"
	)

var addChannel = make(chan *ChatClient)
var removeChannel = make(chan *ChatClient)
var chatChannel = make(chan string)

type ChatClient struct {
	conn net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
	writeChannel chan string
	closeChannel chan bool
}

func NewChatClient(conn net.Conn) (c *ChatClient) {
	c = new(ChatClient)
	c.conn = conn
	c.reader = bufio.NewReader(conn)
	c.writer = bufio.NewWriter(conn)
	c.writeChannel = make(chan string)
	c.closeChannel = make(chan bool)
	return
}

func readLines(c *ChatClient) {
	for {
		line, err := c.reader.ReadString('\n')
		if err != nil {
			log.Print("readLines:", err)
			close(c.writeChannel)
			<-c.closeChannel // Let the write goroutine exit. Could use a buffered close channel instead?
			return
		}
		select {
		case chatChannel <- line:
			continue
		case <-c.closeChannel: // If a write error occurred.
			close(c.writeChannel)
			return
		}
	}
}

func writeLines(c *ChatClient) {
	for line := range c.writeChannel {
		_, err := c.writer.WriteString(line)
		if err != nil {
			log.Print("writeLines:", err)
			break
		}
		c.writer.Flush()
	}
	c.closeChannel <- true
	removeChannel <- c
	c.conn.Close()
}

func processChatCommands() {
	clients := make([]*ChatClient, 0)
	for {
		select {
		case client := <-addChannel:
			clients = append(clients, client)
		case client := <-removeChannel:
			found := -1
			for i, v := range clients {
				if v == client {
					found = i
					break
				}
			}
			if found >= 0 {
				clients = append(clients[:found], clients[found+1:]...)
			}
		case line := <-chatChannel:
			for i := 0; i < len(clients); i++ {
				clients[i].writeChannel <- line
			}
		}		
	}
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
		c := NewChatClient(conn)
		go readLines(c)
		go writeLines(c)
		addChannel <- c
	}
}

func main() {
	go processChatCommands()
	err := runServer()
	if err != nil {
		log.Fatal("runServer:", err);
	}
}
