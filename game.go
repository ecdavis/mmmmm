package main

type Game struct {
	sessions []*Session
	tasks    chan func(*Game)
}

func NewGame() *Game {
	game := new(Game)
	game.sessions = make([]*Session, 0)
	game.tasks = make(chan func(*Game))
	return game
}

func (game *Game) AddSession(session *Session) {
	game.tasks <- func(game *Game) { game.addSession(session) }
}

func (game *Game) RemoveSession(session *Session) {
	game.tasks <- func(game *Game) { game.removeSession(session) }
}

func (game *Game) HandleSessionInput(input *SessionInput) {
	game.tasks <- func(game *Game) { inputHandlerStack[len(inputHandlerStack)-1](game, input) }
}

func (game *Game) addSession(session *Session) {
	game.sessions = append(game.sessions, session)
	go func() {
		sis := session.ReadLines()
		for {
			select {
			case si, ok := <-sis:
				if !ok {
					game.RemoveSession(session)
					return
				} else {
					game.HandleSessionInput(si)
				}
			case <-session.quit:
				game.RemoveSession(session)
				return
			}
		}
	}()
}

func (game *Game) removeSession(session *Session) {
	// TODO Super messy. Use a map instead, perhaps?
	found := -1
	for i, c := range game.sessions {
		if c == session {
			found = i
			break
		}
	}
	if found >= 0 {
		// TODO There may be a memory leak here, see: https://github.com/golang/go/wiki/SliceTricks
		game.sessions = append(game.sessions[:found], game.sessions[found+1:]...)
	}
	// TODO Move this to a method on Session. Also need a way to close the reader.
	close(session.write)
}

func (game *Game) ProcessTasks() {
	for task := range game.tasks {
		task(game)
	}
}
