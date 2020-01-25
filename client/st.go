package main

import (
	log "github.com/sirupsen/logrus"
	"net"

	"github.com/theovassiliou/soundtouch-golang"
)

func main() {

	i, err := net.InterfaceByName("en0")
	log.Infof("Name : %v, supports: %v, HW Address: %v\n", i.Name, i.Flags.String(), i.HardwareAddr)

	speakerCh := soundtouch.Lookup(i)

	// speakerCh := make(chan *soundtouch.Speaker, 1)
	// soundtouch.Lookup(i, speakerCh)
	speaker := <-speakerCh

	websocketCh, err := speaker.Listen()
	if err != nil {
		log.Fatal(err)
	}

	data, err := speaker.Volume()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%v\n", data)
	log.Printf("%s\n", data.Raw)

	speaker.SetVolume(40)
	log.Printf("Set volume to 40")

	for message := range websocketCh {
		log.Printf("%v\n", message)
	}
}
