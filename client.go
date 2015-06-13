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
	c := new(Client)
	c.conn = conn
	c.reader = bufio.NewReader(conn)
	c.writer = bufio.NewWriter(conn)
	c.write = make(chan string)

	go c.WriteLines()

	return c
}

func (c *Client) ReadLines() <-chan string {
	ch := make(chan string)
	go func() {
		for {
			line, err := c.reader.ReadString('\n')
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

func (c *Client) WriteLines() {
	go func() {
		for line := range c.write {
			_, err := c.writer.WriteString(line)
			if err != nil {
				// TODO Handle the error properly.
				log.Print(err)
				break
			}
			c.writer.Flush()
		}
		c.quit <- 1
	}()
}
