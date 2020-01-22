package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/hashicorp/mdns"
	"github.com/theovassiliou/soundtouch-golang"
)

const WEBSOCKET_PORT int = 8080
const MESSAGE_BUFFER_SIZE int = 256

type Speaker struct {
	IP           net.IP
	Port         int
	BaseHttpUrl  url.URL
	WebSocketUrl url.URL

	conn *websocket.Conn
}

func main() {

	i, _ := net.InterfaceByName("en0")
	fmt.Printf("Name : %v, supports: %v, HW Address: %v\n", i.Name, i.Flags.String(), i.HardwareAddr)
	speakerCh := lookup(i)
	var wg sync.WaitGroup
	messageCh := make(chan *soundtouch.Update)

	for speaker := range speakerCh {
		log.Printf("Speaker: %v\n", speaker)
		wg.Add(1)
		go func(speaker *Speaker, msgChan chan *soundtouch.Update) {
			defer wg.Done()
			webSocketCh, _ := speaker.Listen()
			for message := range webSocketCh {
				msgChan <- message
			}
		}(speaker, messageCh)

	}
	for m := range messageCh {
		log.Printf("XXX %#v\n", m)
	}
	wg.Wait()
}

// merges multiple channels
// from https://medium.com/justforfunc/two-ways-of-merging-n-channels-in-go-43c0b57cd1de
func merge(cs ...<-chan *soundtouch.Update) <-chan *soundtouch.Update {
	out := make(chan *soundtouch.Update)
	var wg sync.WaitGroup
	wg.Add(len(cs))
	for _, c := range cs {
		go func(c <-chan *soundtouch.Update) {
			for v := range c {
				out <- v
			}
			wg.Done()
		}(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func (s *Speaker) Listen() (chan *soundtouch.Update, error) {
	log.Printf("Dialing %v", s.WebSocketUrl.String())
	conn, _, err := websocket.DefaultDialer.Dial(
		s.WebSocketUrl.String(),
		http.Header{
			"Sec-WebSocket-Protocol": []string{"gabbo"},
		})
	if err != nil {
		return nil, err
	}

	s.conn = conn
	messageCh := make(chan *soundtouch.Update, MESSAGE_BUFFER_SIZE)
	log.Printf("Created channel")
	go func() {
		for {
			_, body, err := conn.ReadMessage()
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Raw Message: %v", string(body))

			update, err := soundtouch.NewUpdate(body)
			log.Printf("Message: %v", update)
			if update != nil {
				messageCh <- update
			}
		}
	}()
	return messageCh, nil

}

func lookup(iface *net.Interface) <-chan *Speaker {
	speakerCh := make(chan *Speaker)
	entriesCh := make(chan *mdns.ServiceEntry, 6)
	defer close(entriesCh)
	go func() {
		defer close(speakerCh)
		for entry := range entriesCh {
			speakerCh <- newSpeaker(entry)
		}
	}()

	params := mdns.DefaultParams("_soundtouch._tcp")
	params.Entries = entriesCh
	if iface != nil {
		params.Interface = iface
	}
	mdns.Query(params)
	return speakerCh
}

func newSpeaker(entry *mdns.ServiceEntry) *Speaker {
	return &Speaker{
		entry.AddrV4,
		entry.Port,
		url.URL{
			Scheme: "http",
			Host:   fmt.Sprintf("%v:%v", entry.AddrV4.String(), entry.Port),
		},
		url.URL{
			Scheme: "ws",
			Host:   fmt.Sprintf("%v:%v", entry.AddrV4.String(), WEBSOCKET_PORT),
		},
		nil,
	}
}
