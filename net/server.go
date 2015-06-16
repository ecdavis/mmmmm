package net

import (
	"log"
	"net"
)

func RunServer() (<-chan *Session, error) {
	ln, err := net.Listen("tcp", ":4040")
	if err != nil {
		return nil, err
	}
	ch := make(chan *Session)
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Print("Accept:", err)
				continue // TODO Perhaps this should break
			}
			ch <- NewSession(conn)
		}
		close(ch)
	}()
	return ch, nil
}
