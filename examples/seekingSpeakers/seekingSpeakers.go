package main

import (
	"log"

	"github.com/theovassiliou/soundtouch-golang"
)

// Storing the speakers
var soundtouchNetwork = make(map[string]*soundtouch.Speaker)

const mySpeaker = "Office"

func main() {

	interfaceName := "en0"
	noOfSoundtouchSystems := -1

	// --- Build the NetworkConfiguration
	nConf := soundtouch.NetworkConfig{
		InterfaceName: interfaceName,
		NoOfSystems:   noOfSoundtouchSystems,
	}

	// --- Create a channel to get the mdns' found speakers
	speakerCh := soundtouch.GetDevices(nConf)
	for speaker := range speakerCh {
		log.Printf("Found device %s with IP: %v", speaker.Name(), speaker.IP)
		soundtouchNetwork[speaker.Name()] = speaker
	}

	// --- Access the speaker via name and retrieve the Info
	data, err := soundtouchNetwork[mySpeaker].Info()

	if err != nil {
		log.Fatal(err)
	}
	// --- Print only the most relevant information
	log.Printf("Info: %s", data.String())

	// --- Print the received data of communication
	log.Printf("This is the raw data: %s\n", data.Raw)

}
