package game

import (
	"github.com/ecdavis/mmmmm/hooks"
	"github.com/ecdavis/mmmmm/net"
	)

type InputHandler func(*Game, *User, string)

type User struct {
	Session       *net.Session
	InputHandlers []InputHandler
}

func NewUser(session *net.Session) *User {
	user := new(User)
	user.Session = session
	user.InputHandlers = make([]InputHandler, 0)
	return user
}

func (user *User) HandleInput(game *Game, input string) {
	user.InputHandlers[len(user.InputHandlers)-1](game, user, input)
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
	game.users[user.Session] = user
	go func() {
		sessionInputs := user.Session.ReadLines()
		for {
			select {
			case sessionInput, ok := <-sessionInputs:
				if !ok {
					game.RemoveSession(user.Session)
					return
				} else {
					game.HandleSessionInput(sessionInput)
				}
			case <-user.Session.Quit:
				game.RemoveSession(user.Session)
				return
			}
		}
	}()
	hooks.Run("addUser", user)
}

func (game *Game) removeUser(user *User) {
	delete(game.users, user.Session)
	// TODO Move this to a method on Session or User. Also need a way to close the reader.
	close(user.Session.Write)
}

func (game *Game) ProcessTasks() {
	for task := range game.tasks {
		task(game)
	}
}
