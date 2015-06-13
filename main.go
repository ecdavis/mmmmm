package main

import (
	"log"
	"net"
)

var addChannel = make(chan *Client)
var removeChannel = make(chan *Client)
var chatChannel = make(chan string)
var clients = make([]*Client, 0)

func addClient(c *Client) {
	clients = append(clients, c)
	go func() {
		for {
			// TODO Would be nice to use a range here if possible.
			lines := c.ReadLines()
			select {
			case line, ok := <-lines:
				if !ok {
					removeChannel <- c
					return
				} else {
					chatChannel <- line
				}
			case <-c.quit:
				removeChannel <- c
				return
			}
		}
	}()
}

func removeClient(c *Client) {
	// TODO Super messy. Use a map instead, perhaps?
	found := -1
	for i, v := range clients {
		if v == c {
			found = i
			break
		}
	}
	if found >= 0 {
		clients = append(clients[:found], clients[found+1:]...)
	}
	// TODO Move this to a method on Client. Also need a way to close the
	//      reader.
	close(c.write)
}

func sendChatLine(line string) {
	for i := 0; i < len(clients); i++ {
		clients[i].write <- line
	}
}

func processChatCommands() {
	for {
		select {
		case client := <-addChannel:
			addClient(client)
		case client := <-removeChannel:
			removeClient(client)
		case line := <-chatChannel:
			sendChatLine(line)
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
		c := NewClient(conn)
		addChannel <- c
	}
}

func main() {
	go processChatCommands()
	err := runServer()
	if err != nil {
		log.Fatal("runServer:", err)
	}
}
