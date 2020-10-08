package soundtouch

import (
	"encoding/xml"

	log "github.com/sirupsen/logrus"
)

// Zone defines the zone message to/from soundtouch system
type Zone struct {
	Members         []Member `xml:"member"`
	SenderIPAddress string   `xml:"senderIPAddress,attr"`
	SenderIsMaster  bool     `xml:"senderIsMaster,attr"`
	Master          string   `xml:"master,attr"`
	Raw             string   `xml:"-"`
}

type Member struct {
	IPAddress  string `xml:"ipaddress,attr"`
	MacAddress string `xml:",chardata"`
}

// ZoneSlave defines a zone slave
type ZoneSlave struct {
	XMLName xml.Name `xml:"zone"`
	Members []Member `xml:"member"`
	Master  string   `xml:"master,attr"`
}

// MultiRoomZone defines a zone with members
type MultiRoomZone struct {
	XMLName         xml.Name `xml:"zone"`
	Members         []Member `xml:"member"`
	Master          string   `xml:"master,attr"`
	SenderIPAddress string   `xml:"senderIPAddress,attr"`
}

// AddZoneSlave sends the addZoneSlave command to the soundtouch system
// To be send to the master and should contain only new speakers to add
func (s *Speaker) AddZoneSlave(zi MultiRoomZone) error {

	data, err := xml.Marshal(zi)
	if err != nil {
		return err
	}
	_, err = s.SetData("addZoneSlave", data)
	if err != nil {
		return err
	}
	return nil
}

// SetZone sends the addZoneSlave command to the soundtouch system
// creates a multi-room zone
func (s *Speaker) SetZone(zi MultiRoomZone) error {

	data, err := xml.Marshal(zi)
	if err != nil {
		return err
	}
	_, err = s.SetData("setZone", data)
	if err != nil {
		return err
	}
	return nil
}

// GetZone sends the getZone command to the soundtouch system and returns the zone
func (s *Speaker) GetZone() (Zone, error) {
	body, err := s.GetData("getZone")
	if err != nil {
		return Zone{}, err
	}

	zone := Zone{
		Raw: string(body),
	}
	err = xml.Unmarshal(body, &zone)
	if err != nil {
		return zone, err
	}
	return zone, nil
}

// HasZone returns true if the speaker has a zone
func (s *Speaker) HasZone() bool {
	z, _ := s.GetZone()
	return len(z.Members) > 0
}

// IsMaster returns true if the speaker is the master of a zone
func (s *Speaker) IsMaster() bool {
	z, _ := s.GetZone()

	if z.Master == s.DeviceInfo.DeviceID {
		return true
	}
	return false
}

// GetZoneMembers returns the Members of the zone the speaker is a member
func (s *Speaker) GetZoneMembers() []Member {
	z, _ := s.GetZone()
	return z.Members
}

func (s *Speaker) IsSpeakerMember(members []Member) bool {
	for _, m := range members {
		if m.MacAddress == s.DeviceInfo.DeviceID {
			return true
		}
	}
	return false
}

//AddSlave adds a speaker to an existing zone
func (z *Zone) AddSlave(s Speaker) {
	z.Members = append(z.Members, Member{
		s.IP.String(), s.DeviceInfo.DeviceID,
	})
}

func NewZone(master, slave Speaker) MultiRoomZone {
	z := MultiRoomZone{
		Members: []Member{
			{
				master.IP.String(),
				master.DeviceInfo.DeviceID,
			},
			{
				slave.IP.String(),
				slave.DeviceInfo.DeviceID,
			},
		},
		Master:          master.DeviceInfo.DeviceID,
		SenderIPAddress: master.IP.String(),
	}

	return z

}

func DumpZones(log *log.Entry, s Speaker) {
	if false {
		return
	}
	zone, _ := s.GetZone()
	log.Infoln("  zone.Master", zone.Master)
	log.Infoln("  zone.SenderIPAddress", zone.SenderIPAddress)
	log.Infoln("  zone.SenderIsMaster", zone.SenderIsMaster)
	log.Infoln("  zone.Members", zone.Members)

}
