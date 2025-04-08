package client

type message struct {
	room   string
	nick   string
	msg    []byte
	sender *Client
}
