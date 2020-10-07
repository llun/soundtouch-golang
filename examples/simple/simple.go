package main

import (
	"log"
	"time"

	"github.com/theovassiliou/soundtouch-golang"
)

// Storing the speakers
var soundtouchNetwork = make(map[string]*soundtouch.Speaker)

const mySpeaker = "Office"

func main() {

	// Create a speaker object for a known IP
	speaker := soundtouch.NewIPSpeaker("192.168.178.37")

	// --- Get Speaker Info
	info, err := speaker.Info()

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Info %v\n", info)

	// --- Retreive Volume
	vol, err := speaker.Volume()

	log.Printf("The volume is: %d", vol.TargetVolume)

	// --- Vol has also the raw data of communication
	log.Printf("This is the raw data: %s\n", vol.Raw)

	log.Printf("Reduce the volume")
	speaker.SetVolume(vol.TargetVolume / 2)

	// --- Let's wait a little bit
	time.Sleep(5 * time.Second)

	// --- Access the presets
	presets, _ := speaker.Presets()
	for i, pr := range presets.Presets {
		log.Printf("Preset #%v: %v\n", i, pr)
	}

	// --- Check whether the speaker is powered
	log.Println("Is poweredOn:", speaker.IsPoweredOn())

	// --- Pressing the POWER button
	speaker.PressKey(soundtouch.POWER)

	// --- Create a listener to receive updates
	webSocketCh, _ := speaker.Listen()

	go func(msgChan chan *soundtouch.Update) {
		for update := range msgChan {
			log.Println(update)
		}
	}(webSocketCh)

	// --- Wait a little bit so that we can receive something
	time.Sleep(5 * time.Second)

	// --- Check whether the speaker is powered
	log.Println("IsAlive:", speaker.IsAlive())

	// --- Pressing the POWER button again
	speaker.PressKey(soundtouch.POWER)
}
