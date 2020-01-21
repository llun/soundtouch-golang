package main

import (
	"fmt"
	"log"
	"net"

	"github.com/theovassiliou/soundtouch-golang"
)

func main() {

	i, err := net.InterfaceByName("en0")
	fmt.Printf("Name : %v, supports: %v, HW Address: %v\n", i.Name, i.Flags.String(), i.HardwareAddr)

	// see http://golang.org/pkg/net/#Flags
	// addr, err := i.Addrs()
	// var ipnet *net.IPNet
	// if ipnet, ok := addr[1].(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
	// 	if ipnet.IP.To4() != nil {
	// 		fmt.Println(ipnet.IP.String())
	// 	}

	// }
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
