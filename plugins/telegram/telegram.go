package telegram

import (
	"reflect"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/theovassiliou/soundtouch-golang"
	tb "gopkg.in/tucnak/telebot.v2"
)

var name = "Telegram"

const sampleConfig = `
  ## Enabling logger plugin
  # [telegram]

  ## speakers for which messages should be logged. If empty, all 
  # speakers = ["Office", "Kitchen"]

  ## ignore_messages describes the message types to be ignored
  ## one or more of "ConnectionStateUpdated", "NowPlaying", "Volume"
  ## all if empty
  # ignore_messages = ["ConnectionStateUpdated"] 

  ## Telegram API Key
  # apiKey ="x:y"
  # authorizedSenders = ["999999", "888888"]
`

const description = "Logs all update messages to telegram"

// TelegramLogger describes the plugin. It has a
// Config to store the configuration
// Plugin the plugin function
type TelegramLogger struct {
	Config
	Plugin    soundtouch.PluginFunc
	suspended bool
	bot       *tb.Bot
}

// Config contains the configuration of the plugin
// Speakers list of SpeakerNames the handler is added. All if empty
// Terminate indicates whether this is the last handler to be called
// IgnoreMessages a list of message types to be ignored
// APIKey for the telegram bot
type Config struct {
	Speakers         []string `toml:"speakers"`
	IgnoreMessages   []string `toml:"ignore_messages"`
	APIKey           string   `toml:"apiKey"`
	AuthorizedSender []string `toml:"authorizedSender"`
	AuthKey          string   `toml:"authKey"`
}

// NewTelegramLogger creates a new Logger plugin with the configuration
func NewTelegramLogger(config Config) (d *TelegramLogger) {
	d = &TelegramLogger{}
	d.Config = config

	b, err := tb.NewBot(tb.Settings{
		// You can also set custom API URL.
		// If field is empty it equals to "https://api.telegram.org".
		URL: "https://api.telegram.org",

		Token:  config.APIKey,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}
	d.bot = b

	b.Handle("/status", func(m *tb.Message) {
		d.status(m)
	})

	b.Handle("/authorize", func(m *tb.Message) {
		d.authorize(m)
	})

	mLogger := log.WithFields(log.Fields{
		"Plugin": name,
	})

	b.Handle("/hello", func(m *tb.Message) {
		b.Send(m.Sender, "Hello World!")
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		mLogger.Infof("Recevived telegram message: %#v\n", m.Text)
		mLogger.Infof("  by: %v\n", m.Sender)

	})

	mLogger.Infof("Initialised\n")

	go b.Start()

	return d
}

// Name returns the plugin name
func (d *TelegramLogger) Name() string {
	return name
}

// Description returns a string explaining the purpose of this plugin
func (d *TelegramLogger) Description() string { return description }

// SampleConfig returns text explaining how plugin should be configured
func (d *TelegramLogger) SampleConfig() string { return sampleConfig }

// Terminate indicates that no further plugin will be executed on this speaker
func (d *TelegramLogger) Terminate() bool { return false }

// Disable temporarely the execution of the plugin
func (d *TelegramLogger) Disable() { d.suspended = true }

// Enable temporarely the execution of the plugin
func (d *TelegramLogger) Enable() { d.suspended = false }

// Execute runs the plugin with the given parameter
func (d *TelegramLogger) Execute(pluginName string, update soundtouch.Update, speaker soundtouch.Speaker) {
	if len(d.IgnoreMessages) > 0 && isIn(reflect.TypeOf(update.Value).Name(), d.IgnoreMessages) {
		return
	}
	if len(d.Speakers) > 0 && !isIn(speaker.Name(), d.Speakers) {
		return
	}

	mLogger := log.WithFields(log.Fields{
		"Plugin":        name,
		"Speaker":       speaker.Name(),
		"UpdateMsgType": reflect.TypeOf(update.Value).Name(),
	})
	mLogger.Debugln("Executing", pluginName)
	mLogger.Infof("%v\n", update)
}

func isIn(name string, selected []string) bool {
	for _, s := range selected {
		if name == s {
			return true
		}
	}
	return false
}
