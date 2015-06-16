package net

import (
	"bufio"
	"log"
	"net"
	"strings"
)

// TODO There are currently two separate ways of knowing a session is closed.
//      Either the read channel closes or we send an int on the quit channel.
//      This is kind of annoying. There should be one way for the manager to
//      know that the session is gone.
//		We also don't actually close the connection on an error.

type Session struct {
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
	Write  chan string
	Quit   chan int
}

type SessionInput struct {
	Session *Session
	Input   string
}

func NewSession(conn net.Conn) *Session {
	session := new(Session)
	session.conn = conn
	session.reader = bufio.NewReader(conn)
	session.writer = bufio.NewWriter(conn)
	session.Write = make(chan string)

	go session.WriteLines()

	return session
}

func (session *Session) ReadLines() <-chan *SessionInput {
	ch := make(chan *SessionInput)
	go func() {
		for {
			line, err := session.reader.ReadString('\n')
			if err != nil {
				// TODO Handle the error properly.
				log.Print(err)
				break
			}
			line = strings.TrimRight(line, "\r\n")
			ch <- &SessionInput{session, line}
		}
		close(ch)
	}()
	return ch
}

func (session *Session) WriteLines() {
	go func() {
		for line := range session.Write {
			_, err := session.writer.WriteString(line)
			if err != nil {
				// TODO Handle the error properly.
				log.Print(err)
				break
			}
			session.writer.Flush()
		}
		session.Quit <- 1
	}()
}
