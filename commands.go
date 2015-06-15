package main

import "strings"

func HandleCommand(game *Game, user *User, input string) {
	args := strings.Split(input, " ")
	c, a := args[0], strings.Join(args[1:], " ")
	user.session.write <- (c + "\r\n")
	println(a)
}
