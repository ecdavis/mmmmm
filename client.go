package main

import (
	"bufio"
	"log"
	"net"
)

// TODO There are currently two separate ways of knowing a client is closed.
//      Either the read channel closes or we send an int on the quit channel.
//      This is kind of annoying. There should be one way for the manager to
//      know that the client is gone.

type Client struct {
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
	write  chan string
	quit   chan int
}

func NewClient(conn net.Conn) *Client {
	client := new(Client)
	client.conn = conn
	client.reader = bufio.NewReader(conn)
	client.writer = bufio.NewWriter(conn)
	client.write = make(chan string)

	go client.WriteLines()

	return client
}

func (client *Client) ReadLines() <-chan string {
	ch := make(chan string)
	go func() {
		for {
			line, err := client.reader.ReadString('\n')
			if err != nil {
				// TODO Handle the error properly.
				log.Print(err)
				break
			}
			ch <- line
		}
		close(ch)
	}()
	return ch
}

func (client *Client) WriteLines() {
	go func() {
		for line := range client.write {
			_, err := client.writer.WriteString(line)
			if err != nil {
				// TODO Handle the error properly.
				log.Print(err)
				break
			}
			client.writer.Flush()
		}
		client.quit <- 1
	}()
}
