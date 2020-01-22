package soundtouch

type ConnectionStateUpdated struct {
	// state="NETWORK_WIFI_CONNECTED" up="true" signal="MARGINAL_SIGNAL"
	State  string `xml:"state,attr"`
	Up     string `xml:"up,attr"`
	Signal string `xml:"signal,attr"`
}
