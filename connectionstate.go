package soundtouch

import "fmt"

const (
	WIFI = "NETWORK_WIFI_CONNECTED"
)

type ConnectionStateUpdated struct {
	// state="NETWORK_WIFI_CONNECTED" up="true" signal="MARGINAL_SIGNAL"
	State  string `xml:"state,attr"`
	Up     string `xml:"up,attr"`
	Signal string `xml:"signal,attr"`
}

func (c ConnectionStateUpdated) String() string {

	return fmt.Sprintf("Connection: {%v, %v = %v}",
		func() string {
			if c.State == WIFI {
				return "WIFI"
			} else {
				return "NETWORK"
			}
		}(),
		func() string {
			if c.Up == "true" {
				return "UP"
			} else {
				return "DOWN"
			}
		}(), c.Signal)
}
