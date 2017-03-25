package soundtouch

import (
  "bytes"
  "fmt"
  "io/ioutil"
  "log"
  "net"
  "net/http"
  "net/url"

  "github.com/gorilla/websocket"
  "github.com/hashicorp/mdns"
)

const WEBSOCKET_PORT int = 8080
const MESSAGE_BUFFER_SIZE int = 256

type Speaker struct {
  IP           net.IP
  Port         int
  BaseHttpUrl  url.URL
  WebSocketUrl url.URL
}

func Lookup(iface *net.Interface) <-chan *Speaker {
  speakerCh := make(chan *Speaker)
  entriesCh := make(chan *mdns.ServiceEntry, 1)
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
      Host:   fmt.Sprintf("%v:%v", entry.AddrV4.String(), WEBSOCKET_PORT),
    },
  }
}

func (s *Speaker) Listen() (chan *Update, error) {
  conn, _, err := websocket.DefaultDialer.Dial(
    s.WebSocketUrl.String(),
    http.Header{
      "Sec-WebSocket-Protocol": []string{"gabbo"},
    })
  if err != nil {
    return nil, err
  }

  messageCh := make(chan *Update, MESSAGE_BUFFER_SIZE)
  go func() {
    for {
      _, body, err := conn.ReadMessage()
      if err != nil {
        log.Fatal(err)
      }

      update, err := NewUpdate(body)
      if update != nil {
        messageCh <- update
      }
    }
  }()
  return messageCh, nil

}

func (s *Speaker) GetData(action string) ([]byte, error) {
  actionUrl := s.BaseHttpUrl
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
  actionUrl := s.BaseHttpUrl
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
