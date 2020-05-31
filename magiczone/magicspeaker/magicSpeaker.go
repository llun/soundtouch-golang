package magicspeaker

import (
	"reflect"

	log "github.com/sirupsen/logrus"

	soundtouch "github.com/theovassiliou/soundtouch-golang"
)

type MagicSpeaker struct {
	*soundtouch.Speaker
	SpeakerName   string
	WebSocketCh   chan *soundtouch.Update
	KnownSpeakers *MagicSpeakers
}
type MagicSpeakers map[string]*MagicSpeaker

func New(s *soundtouch.Speaker) *MagicSpeaker {
	return &MagicSpeaker{s, "", nil, nil}
}

func (m *MagicSpeaker) MessageLoop() {
	for message := range m.WebSocketCh {
		mu := MagicUpdate{Update: *message}
		m.HandleUpdate(mu, m.WebSocketCh)
	}
}

// handle message per speaker
func (m *MagicSpeaker) HandleUpdate(msg MagicUpdate, webSocketCh chan *soundtouch.Update) {
	typeName := reflect.TypeOf(msg.Value).Name()
	mLogger := log.WithFields(log.Fields{
		"Speaker":     m.SpeakerName,
		"MessageType": typeName,
	})

	if !(msg.Is("NowPlaying")) {
		if !msg.Is("ConnectionStateUpdated") {
			mLogger.Debugf("Ignoring %s\n", typeName)
		}
		return
	}
	np := msg.Value.(soundtouch.NowPlaying)
	if !(np.PlayStatus == soundtouch.PlayState) {
		return
	}

	mLogger.Debugln("PlayStatus == PlayState")

	if !(np.StreamType == soundtouch.RadioStreaming) {
		return
	}
	mLogger.Debugln("StreamType == RadioStreaming")
	compatibleStreamers := make([]MagicSpeaker, 0)
	for _, spk := range *m.KnownSpeakers {
		if m.DeviceInfo.DeviceID == spk.DeviceInfo.DeviceID {
			continue
		}
		snp, _ := spk.Speaker.NowPlaying()
		if np.Content == snp.Content {
			mLogger.Debugln("Found other speaker streaming the same content")
			compatibleStreamers = append(compatibleStreamers, *spk)
		}
	}

	if len(compatibleStreamers) == 0 {
		return // as there are no other speakers streaming the same content
	}

	// 1. Check: Already any zones defined?
	for _, c := range compatibleStreamers {
		if c.HasZone() {
			// search for the one server that is indicated as master "zone.master == c.ownDeviceId"
			zone, _ := c.GetZone()
			if zone.Master == c.Speaker.DeviceInfo.DeviceID {
				if !m.IsSpeakerMember(zone.Members) {
					mLogger.Infof("Adding myself to master %v zone.\n", zone.Master)
					newZone := soundtouch.NewZone(*c.Speaker, *m.Speaker)
					c.AddZoneSlave(newZone)
					soundtouch.DumpZones(mLogger, *c.Speaker)
					return
				}
			}
		}
	}

	choosenAsNewMaster := compatibleStreamers[0]
	if !choosenAsNewMaster.HasZone() {
		newZone := soundtouch.NewZone(*choosenAsNewMaster.Speaker, *m.Speaker)
		mLogger.Infof("Creating new zone with %v as master.\n", newZone.Master)
		choosenAsNewMaster.SetZone(newZone)
		soundtouch.DumpZones(mLogger, *choosenAsNewMaster.Speaker)
		return
	}

}
