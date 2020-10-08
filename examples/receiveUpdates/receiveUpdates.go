package main

import (
	"log"
	"reflect"
	"sync"

	"github.com/theovassiliou/soundtouch-golang"
)

func main() {

	// Create a speaker object for a known IP
	// Replace IP address with IP Address of your speaker
	speaker := soundtouch.NewIPSpeaker("192.168.178.52")
	i, _ := speaker.Info()
	log.Println(i.String())

	var wg sync.WaitGroup

	// --- Create a listener to receive updates
	webSocketCh, _ := speaker.Listen()

	wg.Add(1)
	go func(msgChan chan *soundtouch.Update, wg *sync.WaitGroup) {
		defer wg.Done()
		var ConnectionStateUpdatedCounter int
		for update := range msgChan {
			// Handle the messages
			typeName := reflect.TypeOf(update.Value).Name()
			log.Printf("Message type received: %s\n", typeName)
			log.Println("The update message: ", update)

			if typeName == "ConnectionStateUpdated" {
				ConnectionStateUpdatedCounter++
			}

			if ConnectionStateUpdatedCounter > 10 {
				log.Printf("Received %d ConnectionStateUpdated. Terminating. \n", ConnectionStateUpdatedCounter)
				return
			}
		}
	}(webSocketCh, &wg)

	wg.Wait()
}
