package soundtouch

import (
	"fmt"
)

// All network connection types
const (
	WIFI = "NETWORK_WIFI_CONNECTED"
)

// ConnectionStateUpdated defines the message communicated with the soundtouch system
type ConnectionStateUpdated struct {
	// state="NETWORK_WIFI_CONNECTED" up="true" signal="MARGINAL_SIGNAL"
	State  string `xml:"state,attr"`
	Up     string `xml:"up,attr"`
	Signal string `xml:"signal,attr"`
}

// String readable representation of message
func (c ConnectionStateUpdated) String() string {

	return fmt.Sprintf("Connection: {%v, %v = %v}",
		func() string {
			if c.State == WIFI {
				return "WIFI"
			}
			return "NETWORK"

		}(),
		func() string {
			if c.Up == "true" {
				return "UP"
			}
			return "DOWN"

		}(), c.Signal)
}
