package client

import (
	"encoding/json"
	"log"
)

type Referee struct {
	Register   chan *Client
	PortValues chan PortValues
	unregister chan *Client
	clients    map[*Client]bool
	broadcast  chan message
}

func NewReferee() *Referee {
	return &Referee{
		Register:   make(chan *Client),
		PortValues: make(chan PortValues),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		broadcast:  make(chan message),
	}
}

func (r *Referee) Run() {
	for {
		select {
		case client := <-r.Register:
			log.Printf("+++client: %q/%q\n", client.room, client.nick)
			r.clients[client] = true

		case client := <-r.unregister:
			if _, ok := r.clients[client]; ok {
				log.Printf("---client: %q/%q\n", client.room, client.nick)
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

		case msg := <-r.PortValues:
			m, err := json.Marshal(PortsMessage{Ports: []PortValues{msg}})
			if err != nil {
				log.Printf("Marshaling port message error: %s\n", err)
				continue
			}

			for client := range r.clients {
				select {
				case client.send <- m:
					log.Printf("Sent port values for %q/%q: %s\n", client.room, client.nick, string(m))
				default:
					close(client.send)
					delete(r.clients, client)
				}
			}
		}
	}
}
