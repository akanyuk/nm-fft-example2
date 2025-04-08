package client

import "log"

type Referee struct {
	Register   chan *Client
	unregister chan *Client
	clients    map[*Client]bool
	broadcast  chan message
}

func NewReferee() *Referee {
	return &Referee{
		Register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		broadcast:  make(chan message),
	}
}

func (r *Referee) Run() {
	for {
		select {
		case client := <-r.Register:
			log.Printf("+++client: %p\n", client)
			r.clients[client] = true

		case client := <-r.unregister:
			if _, ok := r.clients[client]; ok {
				log.Printf("---client: %p\n", client)
				delete(r.clients, client)
				close(client.send)
			}

		case msg := <-r.broadcast:
			for client := range r.clients {
				if client.room == msg.room &&
					client != msg.sender &&
					(client.nick == msg.nick || msg.sender.nick == "") {
					select {
					case client.send <- msg.msg:
					default:
						close(client.send)
						delete(r.clients, client)
					}
				}
			}
		}
	}
}
