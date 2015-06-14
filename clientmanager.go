package main

type ClientManager struct {
	clients []*Client
	add     chan *Client
	remove  chan *Client
	read    chan string
	write   chan string
}

func NewClientManager() *ClientManager {
	manager := new(ClientManager)
	manager.clients = make([]*Client, 0)
	manager.add = make(chan *Client)
	manager.remove = make(chan *Client)
	manager.read = make(chan string)
	manager.write = make(chan string)
	return manager
}

func (manager *ClientManager) AddClient(client *Client) {
	manager.clients = append(manager.clients, client)
	go func() {
		lines := client.ReadLines()
		for {
			select {
			case line, ok := <-lines:
				if !ok {
					manager.remove <- client
					return
				} else {
					manager.read <- line
				}
			case <-client.quit:
				manager.remove <- client
				return
			}
		}
	}()
}

func (manager *ClientManager) RemoveClient(client *Client) {
	// TODO Super messy. Use a map instead, perhaps?
	found := -1
	for i, c := range manager.clients {
		if c == client {
			found = i
			break
		}
	}
	if found >= 0 {
		// TODO There may be a memory leak here, see: https://github.com/golang/go/wiki/SliceTricks
		manager.clients = append(manager.clients[:found], manager.clients[found+1:]...)
	}
	// TODO Move this to a method on Client. Also need a way to close the reader.
	close(client.write)
}

func (manager *ClientManager) WriteLine(line string) {
	for i := 0; i < len(manager.clients); i++ {
		manager.clients[i].write <- line
	}
}

func (manager *ClientManager) ProcessCommands(commands chan string) {
	for {
		select {
		case client := <-manager.add:
			manager.AddClient(client)
		case client := <-manager.remove:
			manager.RemoveClient(client)
		case line := <-manager.read:
			commands <- line
		case line := <-manager.write:
			manager.WriteLine(line)
		}
	}
}
