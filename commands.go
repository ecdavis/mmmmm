package main

import "strings"

func HandleCommand(manager *SessionManager, si *SessionInput) {
	args := strings.Split(si.input, " ")
	c, a := args[0], strings.Join(args[1:], " ")
	si.session.write <- (c + "\r\n")
	println(a)
}