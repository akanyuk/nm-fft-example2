package main

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/websocket"

	"github.com/akanyuk/nm-fft-example2/internal/client"
)

var listenPort = flag.Int("port", 9000, "listening port")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func getRoomAndNick(url string) (*string, *string, error) {
	re, err := regexp.Compile(`^/([A-Za-z0-9_]{3,64})/([A-Za-z0-9_]{0,16})$`)
	if err != nil {
		return nil, nil, err
	}
	match := re.FindStringSubmatch(url)
	if len(match) < 2 {
		return nil, nil, errors.New("wrong room name and/or nickname")
	}

	return &match[1], &match[2], nil
}

func main() {
	log.Println("Launching Bonzomatic server")

	referee := client.NewReferee()
	go referee.Run()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		room, nick, err := getRoomAndNick(r.URL.Path)
		if err != nil {
			log.Println(err)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("Client connected to room: '" + *room + "' with nick: '" + *nick + "'")

		cl := client.NewClient(conn, referee, *room, *nick)
		referee.Register <- cl

		go cl.WriteToConnectionPump()
		go cl.ReadFromConnectionPump()
	})

	flag.Parse()

	err := http.ListenAndServe(":"+strconv.Itoa(*listenPort), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
