package soundtouch

import (
	log "github.com/sirupsen/logrus"
)

func (s *Speaker) AddUpdateHandler(uhc UpdateHandlerConfig) {
	s.RemoveUpdateHandler("NotConfigured")
	log.Debugf("Adding handler: %v to %v", uhc.Name, s.Name())
	s.UpdateHandlers = append(s.UpdateHandlers, uhc)
}

func (s *Speaker) RemoveUpdateHandler(name string) {
	var newHandler []UpdateHandlerConfig

	for _, suhc := range s.UpdateHandlers {
		if suhc.Name != "NotConfigured" {
			newHandler = append(newHandler, suhc)
		}
	}
	s.UpdateHandlers = newHandler
}

// HasUpdateHandler returns true if speaker has an UpdateHandler named name. False otherwise
func (s *Speaker) HasUpdateHandler(name string) bool {
	for _, suhc := range s.UpdateHandlers {
		if suhc.Name == name {
			return true
		}
	}
	return false
}

func (s *Speaker) Handle(msgChan chan *Update) {
	for update := range msgChan {
		for _, uh := range s.UpdateHandlers {
			uh.UpdateHandler.Handle(uh.Name, *update, *s)
			if uh.Terminate {
				return
			}
		}
	}
}

// UpdateHandlerFunc turns a function with the right signature into a update handler
type UpdateHandlerFunc func(hndlName string, update Update, speaker Speaker)

// Handle executing the request and returning a response
func (fn UpdateHandlerFunc) Handle(hndlName string, update Update, speaker Speaker) {
	fn(hndlName, update, speaker)
}

// UpdateHandler interface for that can handle valid update params
type UpdateHandler interface {
	Handle(hndlName string, update Update, speaker Speaker)
}

// UpdateHandlerConfig describes an UpdateHandler. It has a
// Name to be able to remove again
// Speakers list of SpeakerNames the handler is added. All if empty
// UpdateHandler the function
// Terminate indicates whether this is the last handler to be called
type UpdateHandlerConfig struct {
	Name          string
	Speakers      []string
	UpdateHandler UpdateHandler
	Terminate     bool
}
