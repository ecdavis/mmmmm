package main

type Game struct {
	sessions []*Session
	read     chan *SessionInput
	tasks    chan func(*Game)
}

func NewGame() *Game {
	game := new(Game)
	game.sessions = make([]*Session, 0)
	game.read = make(chan *SessionInput)
	game.tasks = make(chan func(*Game))
	return game
}

func (game *Game) AddSession(session *Session) {
	game.sessions = append(game.sessions, session)
	go func() {
		sis := session.ReadLines()
		for {
			select {
			case si, ok := <-sis:
				if !ok {
					game.tasks <- func(game *Game) { game.RemoveSession(session) }
					return
				} else {
					game.read <- si
				}
			case <-session.quit:
				game.tasks <- func(game *Game) { game.RemoveSession(session) }
				return
			}
		}
	}()
}

func (game *Game) RemoveSession(session *Session) {
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

// TODO Make this return a channel rather than passing one in?
func (game *Game) ProcessTasks(sis chan *SessionInput) {
	for {
		select {
		case si := <-game.read:
			sis <- si
		case task := <-game.tasks:
			task(game)
		}
	}
}
