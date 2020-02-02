package soundtouch

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
	"github.com/hashicorp/mdns"
)

const WebsocketPort int = 8080
const MessageBufferSize int = 256

type Speaker struct {
	IP           net.IP
	Port         int
	BaseHTTPURL  url.URL
	WebSocketURL url.URL
	DeviceInfo   Info
	conn         *websocket.Conn
}

func Lookup(iface *net.Interface) <-chan *Speaker {
	speakerCh := make(chan *Speaker)
	entriesCh := make(chan *mdns.ServiceEntry, 7)
	defer close(entriesCh)
	go func() {
		defer close(speakerCh)
		for entry := range entriesCh {
			speakerCh <- NewSpeaker(entry)
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

func NewSpeaker(entry *mdns.ServiceEntry) *Speaker {
	return &Speaker{
		entry.AddrV4,
		entry.Port,
		url.URL{
			Scheme: "http",
			Host:   fmt.Sprintf("%v:%v", entry.AddrV4.String(), entry.Port),
		},
		url.URL{
			Scheme: "ws",
			Host:   fmt.Sprintf("%v:%v", entry.AddrV4.String(), WebsocketPort),
		},
		Info{},
		nil,
	}
}

func (s *Speaker) Listen() (chan *Update, error) {
	spkLogger := log.WithFields(log.Fields{
		"Speaker": s.DeviceInfo.Name,
		"ID":      s.DeviceInfo.DeviceID,
	})
	spkLogger.Tracef("Dialing %v", s.WebSocketURL.String())
	conn, _, err := websocket.DefaultDialer.Dial(
		s.WebSocketURL.String(),
		http.Header{
			"Sec-WebSocket-Protocol": []string{"gabbo"},
		})
	if err != nil {
		return nil, err
	}

	s.conn = conn
	messageCh := make(chan *Update, MessageBufferSize)
	go func() {
		for {
			mLogger := log.WithFields(log.Fields{
				"Speaker": s.DeviceInfo.Name,
				"ID":      s.DeviceInfo.DeviceID,
			})
			_, body, err := conn.ReadMessage()
			if err != nil {
				log.Fatal(err)
			}
			mLogger.Tracef("Raw Message: %v", string(body))

			update, err := NewUpdate(body)
			if err != nil {
				mLogger.Debugf("Message: unkown")
				mLogger.Tracef(err.Error())
			} else {
				mLogger.Debugf("Message: %v", update)
			}
			if update != nil {
				messageCh <- update
			}
		}
	}()
	return messageCh, nil

}

func (s *Speaker) Close() error {
	log.Debugf("Closing socket")
	return s.conn.Close()
}

func (s *Speaker) GetData(action string) ([]byte, error) {
	actionUrl := s.BaseHTTPURL
	actionUrl.Path = action
	resp, err := http.Get(actionUrl.String())
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (s *Speaker) SetData(action string, input []byte) ([]byte, error) {
	actionUrl := s.BaseHTTPURL
	actionUrl.Path = action
	buffer := bytes.NewBuffer(input)
	resp, err := http.Post(actionUrl.String(), "application/xml", buffer)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return body, nil
}
