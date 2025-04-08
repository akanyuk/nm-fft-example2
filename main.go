package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hophiphip/portaudio"

	"github.com/akanyuk/nm-fft-example2/internal/analyzer"
	"github.com/akanyuk/nm-fft-example2/internal/client"
)

var (
	listenPort = flag.Int("port", 9000, "listening port")
	sampleRate = flag.Int("r", 48000, "sample rate")
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func initAudio(buffer []float32, sampleRate float64) (*portaudio.Stream, error) {
	if err := portaudio.Initialize(); err != nil {
		return nil, err
	}

	stream, err := portaudio.OpenDefaultStream(1, 0, sampleRate, len(buffer), &buffer)
	if err != nil {
		return nil, err
	}

	return stream, nil
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

	flag.Parse()

	referee := client.NewReferee()
	go referee.Run()

	buffer := make([]float32, 2048)
	analyzerInstance := analyzer.NewAnalyzer(buffer)

	stream, err := initAudio(buffer, float64(*sampleRate))
	if err != nil {
		log.Fatalf("Error initializing audio: %s\n", err)
		return
	}
	defer func() {
		_ = stream.Close()
	}()

	if err := stream.Start(); err != nil {
		log.Fatalf("Error starting audio stream: %s\n", err)
		return
	}
	defer func() {
		_ = stream.Stop()
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		ticker := time.NewTicker(time.Millisecond * 50)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("Stream read canceled")
				return
			case <-ticker.C:
				if err := stream.Read(); err != nil {
					log.Printf("Audio stream read error: %v\n", err)
				}

				referee.PortValues <- client.PortValues{
					Port:   0x20fb,
					Values: analyzerInstance.Process(),
				}
			}
		}
	}()

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

	if err := http.ListenAndServe(":"+strconv.Itoa(*listenPort), nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
