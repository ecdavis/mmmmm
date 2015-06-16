package main

import "github.com/ecdavis/mmmmm/net"

type InputHandler func(*Game, *User, string)

type User struct {
	session       *net.Session
	inputHandlers []InputHandler
}

func NewUser(session *net.Session) *User {
	user := new(User)
	user.session = session
	user.inputHandlers = make([]InputHandler, 0)
	user.inputHandlers = append(user.inputHandlers, HandleCommand)
	return user
}

func (user *User) HandleInput(game *Game, input string) {
	user.inputHandlers[len(user.inputHandlers)-1](game, user, input)
}

type Game struct {
	users map[*net.Session]*User
	tasks chan func(*Game)
}

func NewGame() *Game {
	game := new(Game)
	game.users = make(map[*net.Session]*User)
	game.tasks = make(chan func(*Game))
	return game
}

func (game *Game) AddSession(session *net.Session) {
	game.tasks <- func(game *Game) {
		game.addUser(NewUser(session))
	}
}

func (game *Game) RemoveSession(session *net.Session) {
	game.tasks <- func(game *Game) {
		game.removeUser(game.users[session])
	}
}

func (game *Game) HandleSessionInput(input *net.SessionInput) {
	game.tasks <- func(game *Game) {
		game.users[input.Session].HandleInput(game, input.Input)
	}
}

func (game *Game) addUser(user *User) {
	game.users[user.session] = user
	go func() {
		sessionInputs := user.session.ReadLines()
		for {
			select {
			case sessionInput, ok := <-sessionInputs:
				if !ok {
					game.RemoveSession(user.session)
					return
				} else {
					game.HandleSessionInput(sessionInput)
				}
			case <-user.session.Quit:
				game.RemoveSession(user.session)
				return
			}
		}
	}()
}

func (game *Game) removeUser(user *User) {
	delete(game.users, user.session)
	// TODO Move this to a method on Session or User. Also need a way to close the reader.
	close(user.session.Write)
}

func (game *Game) ProcessTasks() {
	for task := range game.tasks {
		task(game)
	}
}
