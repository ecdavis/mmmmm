package main

type SessionManager struct {
	sessions []*Session
	add     chan *Session
	remove  chan *Session
	read    chan string
	write   chan string
}

func NewSessionManager() *SessionManager {
	manager := new(SessionManager)
	manager.sessions = make([]*Session, 0)
	manager.add = make(chan *Session)
	manager.remove = make(chan *Session)
	manager.read = make(chan string)
	manager.write = make(chan string)
	return manager
}

func (manager *SessionManager) AddSession(session *Session) {
	manager.sessions = append(manager.sessions, session)
	go func() {
		lines := session.ReadLines()
		for {
			select {
			case line, ok := <-lines:
				if !ok {
					manager.remove <- session
					return
				} else {
					manager.read <- line
				}
			case <-session.quit:
				manager.remove <- session
				return
			}
		}
	}()
}

func (manager *SessionManager) RemoveSession(session *Session) {
	// TODO Super messy. Use a map instead, perhaps?
	found := -1
	for i, c := range manager.sessions {
		if c == session {
			found = i
			break
		}
	}
	if found >= 0 {
		// TODO There may be a memory leak here, see: https://github.com/golang/go/wiki/SliceTricks
		manager.sessions = append(manager.sessions[:found], manager.sessions[found+1:]...)
	}
	// TODO Move this to a method on Session. Also need a way to close the reader.
	close(session.write)
}

func (manager *SessionManager) WriteLine(line string) {
	for i := 0; i < len(manager.sessions); i++ {
		manager.sessions[i].write <- line
	}
}

// TODO Make this return a channel rather than passing one in?
func (manager *SessionManager) ProcessCommands(commands chan string) {
	for {
		select {
		case session := <-manager.add:
			manager.AddSession(session)
		case session := <-manager.remove:
			manager.RemoveSession(session)
		case line := <-manager.read:
			commands <- line
		case line := <-manager.write:
			manager.WriteLine(line)
		}
	}
}