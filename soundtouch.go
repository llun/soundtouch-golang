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

const websocketPort int = 8080
const messageBufferSize int = 256

// Speaker defines a soundtouch speaker
type Speaker struct {
	IP             net.IP
	Port           int
	BaseHTTPURL    url.URL
	WebSocketURL   url.URL
	DeviceInfo     Info
	conn           *websocket.Conn
	webSocketCh    chan *Update
	UpdateHandlers []UpdateHandlerConfig
}

// Lookup listens via mdns for soundtouch speakers and returns Speaker channel
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

// NewSpeaker returns a new Speaker entity based on a mdns service entry
func NewSpeaker(entry *mdns.ServiceEntry) *Speaker {
	if entry == nil {
		return &Speaker{}
	}

	return &Speaker{
		entry.AddrV4,
		entry.Port,
		url.URL{
			Scheme: "http",
			Host:   fmt.Sprintf("%v:%v", entry.AddrV4.String(), entry.Port),
		},
		url.URL{
			Scheme: "ws",
			Host:   fmt.Sprintf("%v:%v", entry.AddrV4.String(), websocketPort),
		},
		Info{},
		nil,
		nil,
		[]UpdateHandlerConfig{
			{
				Name: "NotConfigured",
				UpdateHandler: UpdateHandlerFunc(func(hndlName string, update Update, speaker Speaker) {
					log.Infof("UpdateHandler not configured.")
				}),
				Terminate: false,
			},
		},
	}
}

// Listen creates a listenes that distributes Update messages via channel
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
	messageCh := make(chan *Update, messageBufferSize)
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
				mLogger.Tracef("Message: unkown")
				mLogger.Tracef(err.Error())
			} else {
				mLogger.Tracef("Message: %v", update)
			}
			if update != nil {
				messageCh <- update
			}
		}
	}()
	return messageCh, nil

}

// Close closes the socket to the soundtouch speaker
func (s *Speaker) Close() error {
	log.Debugf("Closing socket")
	return s.conn.Close()
}

// GetData returns received raw data retrieved a GET for a given soundtouch action
func (s *Speaker) GetData(action string) ([]byte, error) {
	actionURL := s.BaseHTTPURL
	actionURL.Path = action

	mLogger := log.WithFields(log.Fields{
		"Speaker": s.DeviceInfo.Name,
		"ID":      s.DeviceInfo.DeviceID,
	})

	mLogger.Tracef("GET: %s\n", actionURL.String())

	resp, err := http.Get(actionURL.String())
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

// SetData sets raw data via  POST for a given soundtouch action
func (s *Speaker) SetData(action string, input []byte) ([]byte, error) {
	actionURL := s.BaseHTTPURL
	actionURL.Path = action
	buffer := bytes.NewBuffer(input)

	mLogger := log.WithFields(log.Fields{
		"Speaker": s.DeviceInfo.Name,
		"ID":      s.DeviceInfo.DeviceID,
	})

	mLogger.Tracef("POST: %s, %v\n", actionURL.String(), buffer)

	resp, err := http.Post(actionURL.String(), "application/xml", buffer)
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

// Name returns the speakers name as indicated in the info message, or "" if name unknwon
func (s *Speaker) Name() (name string) {
	return s.DeviceInfo.Name
}

// DeviceID returns the speakers DeviceID as indicated in the info message, or "" if name unknwon
func (s *Speaker) DeviceID() (name string) {
	return s.DeviceInfo.DeviceID
}
