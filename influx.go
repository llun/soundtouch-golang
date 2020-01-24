package soundtouch

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
)

type InfluxDB struct {
	BaseHTTPURL url.URL
	Database    string
}

func (i *InfluxDB) SetData(action string, input []byte) ([]byte, error) {
	actionURL, _ := url.Parse("http://localhost:8086/write?db=soundtouch")
	buffer := bytes.NewBuffer(input)
	// log.Debugf("Going to send action: %v, %v", action, string(input))
	resp, err := http.Post(actionURL.String(), "", buffer)
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
