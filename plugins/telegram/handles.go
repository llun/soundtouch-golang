package telegram

import (
	"fmt"
	"strings"

	"github.com/theovassiliou/soundtouch-golang"
	tb "gopkg.in/tucnak/telebot.v2"
)

// assertSender returns false in case user is not authorized
func (d *TelegramLogger) assertSender(sender *tb.User) bool {
	return isIn(string(sender.ID), d.Config.AuthorizedSender)
}

// /status [speakerName]
func (d *TelegramLogger) status(m *tb.Message) {
	if !d.assertSender(m.Sender) {
		d.bot.Send(m.Sender, "Not Authorized. Use /authorize (authKey)")
		return
	}

	text := m.Text
	speakers := strings.Split(text, " ")
	speakers = speakers[1:]

	gkd := soundtouch.GetKnownDevices()
	var b strings.Builder
	if len(speakers) == 0 {
		b.WriteString("List of Soundtouch Devices\n")
		for sName := range gkd {
			speaker := soundtouch.GetSpeakerByDeviceId(sName)
			fmt.Fprintf(&b, "Device %s-%s with IP %s\n", speaker.Name(), speaker.DeviceID(), speaker.IP)

		}
	} else {
		for _, sName := range speakers {
			speaker := soundtouch.GetSpeakerByName(sName)
			fmt.Fprintf(&b, "Device %s-%s with IP %s\n", speaker.Name(), speaker.DeviceID(), speaker.IP)
			fmt.Fprintf(&b, " isPoweredOn(): %v\n", speaker.IsPoweredOn())
			fmt.Fprintf(&b, " isAlive(): %v\n", speaker.IsAlive())
			zone, _ := speaker.GetZone()

			fmt.Fprintf(&b, " isMaster(): %v\n", speaker.IsMaster())
			fmt.Fprintln(&b, "  zone.Master: ", zone.Master)
			fmt.Fprintln(&b, "  zone.SenderIPAddress: ", zone.SenderIPAddress)
			fmt.Fprintln(&b, "  zone.SenderIsMaster: ", zone.SenderIsMaster)
			fmt.Fprintln(&b, "  zone.Members: ", zone.Members)
		}
	}
	d.bot.Send(m.Sender, b.String())
}

// /authorize [authkey]
func (d *TelegramLogger) authorize(m *tb.Message) {
	authKey := d.Config.AuthKey

	if authKey == "" {
		d.bot.Send(m.Sender, "Authorization temporary disabled")
		return
	}

	text := m.Text
	authParam := strings.Split(text, " ")
	if authKey == authParam[1] {
		d.Config.AuthorizedSender = append(d.Config.AuthorizedSender, string(m.Sender.ID))
		d.bot.Send(m.Sender, "Authorization granted")
		return
	}

	d.bot.Send(m.Sender, fmt.Sprintf("Could not authorize with key %v", authParam[1]))

}
