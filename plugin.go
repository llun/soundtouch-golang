package soundtouch

import (
	log "github.com/sirupsen/logrus"
)

func (s *Speaker) AddPlugin(uhc Plugin) {
	if !s.HasPlugin(uhc.Name()) {
		log.Tracef("Adding handler: %s to %s", uhc.Name(), s.Name())
		s.Plugins = append(s.Plugins, uhc)
	}
}

func (s *Speaker) RemovePlugin(name string) {

	for _, suhc := range s.Plugins {
		if suhc.Name() == name {
			suhc.Disable()
		}
	}
}

// HasPlugin returns true if speaker has an UpdateHandler named name. False otherwise
func (s *Speaker) HasPlugin(name string) bool {
	for _, suhc := range s.Plugins {
		if suhc.Name() == name {
			return true
		}
	}
	return false
}

func (s *Speaker) Execute(msgChan chan *Update) {
	for update := range msgChan {
		for _, uh := range s.Plugins {
			uh.Execute(uh.Name(), *update, *s)
			if uh.Terminate() {
				return
			}
		}
	}
}

func (s *Speaker) Disable(pluginName string) {
	for _, suhc := range s.Plugins {
		if suhc.Name() == pluginName {
			suhc.Disable()
		}
	}
}

// PluginFunc turns a function with the right signature into a update handler
type PluginFunc func(pluginName string, update Update, speaker Speaker)

// Execute executing the request and returning a response
func (fn PluginFunc) Execute(pluginName string, update Update, speaker Speaker) {
	fn(pluginName, update, speaker)
}

// Plugin interface for that can handle valid update params
type Plugin interface {
	// SampleConfig returns the default configuration of the Input
	SampleConfig() string

	// Description returns a one-sentence description on the Input
	Description() string

	// Execute operates on one update message
	Execute(pluginName string, update Update, speaker Speaker)

	// Disable temporarely the execution of the plugin
	Disable()

	// Enable temporarely the execution of the plugin
	Enable()

	// Name returns the name of the plugin
	Name() string

	// Terminate indicates that no further plugin will be executed on this speaker
	Terminate() bool
}

// Initializer is an interface that all plugin can optionally implement to initialize the
// plugin.
type Initializer interface {
	// Init performs one time setup of the plugin and returns an error if the
	// configuration is invalid.
	Init() error
}

// PluginConfig describes an Plugin. It has a
// Name to be able to remove again
// Speakers list of SpeakerNames the handler is added. All if empty
// Terminate indicates whether this is the last handler to be called
// Suspended indicates that the plugin is temporarely suspended
// Plugin the plugin function
type PluginConfig struct {
	Name      string
	Speakers  []string
	Terminate bool
	Suspended bool
	Pfunction PluginFunc
}
