package main

type InputHandler func(*Game, *User, string)

type User struct {
	session       *Session
	inputHandlers []InputHandler
}

func NewUser(session *Session) *User {
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
	users map[*Session]*User
	tasks chan func(*Game)
}

func NewGame() *Game {
	game := new(Game)
	game.users = make(map[*Session]*User)
	game.tasks = make(chan func(*Game))
	return game
}

func (game *Game) AddSession(session *Session) {
	game.tasks <- func(game *Game) {
		game.addUser(NewUser(session))
	}
}

func (game *Game) RemoveSession(session *Session) {
	game.tasks <- func(game *Game) {
		game.removeUser(game.users[session])
	}
}

func (game *Game) HandleSessionInput(input *SessionInput) {
	game.tasks <- func(game *Game) {
		game.users[input.session].HandleInput(game, input.input)
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
			case <-user.session.quit:
				game.RemoveSession(user.session)
				return
			}
		}
	}()
}

func (game *Game) removeUser(user *User) {
	delete(game.users, user.session)
	// TODO Move this to a method on Session or User. Also need a way to close the reader.
	close(user.session.write)
}

func (game *Game) ProcessTasks() {
	for task := range game.tasks {
		task(game)
	}
}
