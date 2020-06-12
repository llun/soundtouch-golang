package soundtouch

import (
	log "github.com/sirupsen/logrus"
)

func (s *Speaker) AddPlugin(uhc PluginConfig) {
	s.RemovePlugin("NotConfigured")
	log.Debugf("Adding handler: %v to %v", uhc.Name, s.Name())
	s.UpdateHandlers = append(s.UpdateHandlers, uhc)
}

func (s *Speaker) RemovePlugin(name string) {
	var newHandler []PluginConfig

	for _, suhc := range s.UpdateHandlers {
		if suhc.Name != "NotConfigured" {
			newHandler = append(newHandler, suhc)
		}
	}
	s.UpdateHandlers = newHandler
}

// HasPlugin returns true if speaker has an UpdateHandler named name. False otherwise
func (s *Speaker) HasPlugin(name string) bool {
	for _, suhc := range s.UpdateHandlers {
		if suhc.Name == name {
			return true
		}
	}
	return false
}

func (s *Speaker) Execute(msgChan chan *Update) {
	for update := range msgChan {
		for _, uh := range s.UpdateHandlers {
			uh.Plugin.Execute(uh.Name, *update, *s)
			if uh.Terminate {
				return
			}
		}
	}
}

// PluginFunc turns a function with the right signature into a update handler
type PluginFunc func(hndlName string, update Update, speaker Speaker)

// Handle executing the request and returning a response
func (fn PluginFunc) Execute(hndlName string, update Update, speaker Speaker) {
	fn(hndlName, update, speaker)
}

// UpdateHandler interface for that can handle valid update params
type Plugin interface {
	Execute(hndlName string, update Update, speaker Speaker)
}

// PluginConfig describes an UpdateHandler. It has a
// Name to be able to remove again
// Speakers list of SpeakerNames the handler is added. All if empty
// Plugin the plugin function
// Terminate indicates whether this is the last handler to be called
type PluginConfig struct {
	Name      string
	Speakers  []string
	Plugin    Plugin
	Terminate bool
}
