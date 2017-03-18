package soundtouch

import (
  "bytes"
  "fmt"
  "io/ioutil"
  "net"
  "net/http"
  "net/url"

  "github.com/hashicorp/mdns"
)

type Speaker struct {
  IP           net.IP
  Port         int
  BaseHttpUrl  url.URL
  WebSocketUrl url.URL
}

func Lookup(speakerCh chan<- *Speaker) {
  entriesCh := make(chan *mdns.ServiceEntry, 1)
  go func() {
    defer close(entriesCh)
    for entry := range entriesCh {
      speakerCh <- NewSpeaker(entry)
    }
  }()
  mdns.Lookup("_soundtouch._tcp", entriesCh)
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
      Host:   fmt.Sprintf("%v:%v", entry.AddrV4.String(), entry.Port),
    },
  }
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
