package main

type Game struct {
	sessions []*Session
	add      chan *Session
	remove   chan *Session
	read     chan *SessionInput
	write    chan string
}

func NewGame() *Game {
	game := new(Game)
	game.sessions = make([]*Session, 0)
	game.add = make(chan *Session)
	game.remove = make(chan *Session)
	game.read = make(chan *SessionInput)
	game.write = make(chan string)
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
					game.remove <- session
					return
				} else {
					game.read <- si
				}
			case <-session.quit:
				game.remove <- session
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

func (game *Game) WriteLine(line string) {
	for i := 0; i < len(game.sessions); i++ {
		game.sessions[i].write <- line
	}
}

// TODO Make this return a channel rather than passing one in?
func (game *Game) ProcessCommands(sis chan *SessionInput) {
	for {
		select {
		case session := <-game.add:
			game.AddSession(session)
		case session := <-game.remove:
			game.RemoveSession(session)
		case si := <-game.read:
			sis <- si
		case line := <-game.write:
			game.WriteLine(line)
		}
	}
}
