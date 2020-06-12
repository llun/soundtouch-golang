// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	xj "github.com/basgys/goxml2json"

	"github.com/theovassiliou/soundtouch-golang/server/models"
	"github.com/theovassiliou/soundtouch-golang/server/restapi/operations"
	apiops "github.com/theovassiliou/soundtouch-golang/server/restapi/operations/api"
	"github.com/theovassiliou/soundtouch-golang/server/restapi/operations/device"
	"github.com/theovassiliou/soundtouch-golang/server/restapi/operations/key"

	sndt "github.com/theovassiliou/soundtouch-golang"
)

//go:generate swagger generate server --target ../../server --name SoundtouchRESTfulJSONServer --spec ../swagger/swagger.yml
type speakerMap map[string]bool

// RestSpeaker is the Speaker representation within this server.
type RestSpeaker struct {
	*sndt.Speaker
}

// New converts a Speaker to a restSpeaker
func New(s *sndt.Speaker) *RestSpeaker {
	return &RestSpeaker{s}
}

type speakers map[string]*RestSpeaker

var visibleSpeakers = make(speakers)

type config struct {
	Interface           string   `short:"i" long:"interface" description:"network interface to listen"`
	NoSoundtouchSystems int      `short:"n" long:"noSystems" description:"Number of Soundtouch systems to scan for."`
	Speakers            []string `short:"s" long:"speakers" description:"Speakers to listen for, all if not set"`
	LogLevel            string   `short:"l" long:"log-level" default:"debug" description:"Log level, one of panic, fatal, error, warn or warning, info, debug, trace"`
}

var soundtouchFlags = config{}

func configureFlags(api *operations.SoundtouchRESTfulJSONServerAPI) {
	api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{
		{
			ShortDescription: "Soundtouch Flags",
			LongDescription:  "",
			Options:          &soundtouchFlags,
		},
	}
}

type speakerDevice struct {
	Name      string           `json:"name"`
	Addresses []sndt.IPAddress `json:"addresses"`
}

type speakerDeviceAdvanced map[string]interface{}

func configureAPI(api *operations.SoundtouchRESTfulJSONServerAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	l, err := log.ParseLevel(soundtouchFlags.LogLevel)
	if err != nil {
		log.SetLevel(log.DebugLevel)
		log.Debug("Parse error for log level %s. Using debug instead.", soundtouchFlags.LogLevel)
	} else {
		log.SetLevel(l)
	}
	api.Logger = log.Printf
	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()

	nConf := sndt.NetworkConfig{
		InterfaceName: soundtouchFlags.Interface,
		NoOfSystems:   soundtouchFlags.NoSoundtouchSystems,
		UpdateHandlers: []sndt.UpdateHandlerConfig{
			{
				Name:          "BASIC-HANDLER",
				UpdateHandler: sndt.UpdateHandlerFunc(basicHandler),
				Terminate:     false,
			},
			{
				Name:          "CONNECTION-HANDLER",
				Speakers:      []string{"Office"},
				UpdateHandler: sndt.UpdateHandlerFunc(connectionHandler),
				Terminate:     false,
			},
		},
	}

	speakerCh := sndt.GetDevices(nConf)
	for speaker := range speakerCh {
		s := RestSpeaker{speaker}
		visibleSpeakers[speaker.Name()] = &s
	}

	// GET /{speakerName}/nowPlaying
	api.KeyNowPlayingHandler = key.NowPlayingHandlerFunc(func(params key.NowPlayingParams) middleware.Responder {
		ck, err := checkSpeakerName(params.SpeakerName)
		if !ck {
			return key.NewNowPlayingDefault(404).WithPayload(err)
		}

		s := visibleSpeakers[params.SpeakerName]

		np, _ := s.NowPlaying()
		return key.NewNowPlayingOK().WithPayload(np)
	})

	// GET /{speakerName}/trackInfo
	api.KeyTrackInfoHandler = key.TrackInfoHandlerFunc(func(params key.TrackInfoParams) middleware.Responder {
		return middleware.NotImplemented("/speakerName/trackInfo currently not implemented")
	})

	// GET /{speakerName}/Presets
	api.KeyPresetsHandler = key.PresetsHandlerFunc(func(params key.PresetsParams) middleware.Responder {
		ck, err := checkSpeakerName(params.SpeakerName)
		if !ck {
			return key.NewPresetsDefault(404).WithPayload(err)
		}

		s := visibleSpeakers[params.SpeakerName]

		np, _ := s.Presets()
		return key.NewPresetsOK().WithPayload(np)
	})

	// GET /api/keys-list
	api.APIKeysListHandler = apiops.KeysListHandlerFunc(func(params apiops.KeysListParams) middleware.Responder {
		return apiops.NewKeysListOK().WithPayload(sndt.ALLKEYS)
	})

	// GET /device/list
	api.DeviceListHandler = device.ListHandlerFunc(func(params device.ListParams) middleware.Responder {
		var devices []*models.Device

		for _, s := range visibleSpeakers {
			devices = append(devices, &models.Device{Addresses: []string([]string{s.DeviceInfo.IPAddress[0]}), Name: s.Name()})
		}

		return device.NewListOK().WithPayload(devices)
	})

	// GET /device/listAdvanced
	api.DeviceListAdvancedHandler = device.ListAdvancedHandlerFunc(func(params device.ListAdvancedParams) middleware.Responder {
		var devices []speakerDeviceAdvanced

		for _, s := range visibleSpeakers {
			json1, _ := xj.Convert(strings.NewReader(string(s.DeviceInfo.Raw)))
			var j speakerDeviceAdvanced
			json.Unmarshal(json1.Bytes(), &j)
			devices = append(devices, j)
		}

		return device.NewListAdvancedOK().WithPayload(devices)
	})

	// GET /{speakerName}/playPause
	api.KeyPlayPauseHandler = key.PlayPauseHandlerFunc(func(params key.PlayPauseParams) middleware.Responder {
		ck, err := checkSpeakerName(params.SpeakerName)
		if !ck {
			return key.NewPlayPauseDefault(404).WithPayload(err)
		}
		s := visibleSpeakers[params.SpeakerName]

		s.PowerOn()
		s.PressKey(sndt.PLAYPAUSE)

		return key.NewPlayPauseNoContent()
	})

	api.KeyPressKeyHandler = key.PressKeyHandlerFunc(func(params key.PressKeyParams) middleware.Responder {
		ck, err := checkSpeakerName(params.SpeakerName)
		if !ck {
			return key.NewPressKeyDefault(404).WithPayload(err)
		}

		s := visibleSpeakers[params.SpeakerName]

		s.PressKey(sndt.Key(params.KeyID))
		return key.NewPressKeyNoContent()

	})

	api.KeyPlayHandler = key.PlayHandlerFunc(func(params key.PlayParams) middleware.Responder {
		ck, err := checkSpeakerName(params.SpeakerName)
		if !ck {
			return key.NewPressKeyDefault(404).WithPayload(err)
		}

		s := visibleSpeakers[params.SpeakerName]
		s.PowerOn()
		s.PressKey(sndt.Key(sndt.PLAY))

		return key.NewPlayNoContent()
	})

	api.KeyPowerOnHandler = key.PowerOnHandlerFunc(func(params key.PowerOnParams) middleware.Responder {
		ck, err := checkSpeakerName(params.SpeakerName)
		if !ck {
			return key.NewPowerOnDefault(404).WithPayload(err)
		}

		s := visibleSpeakers[params.SpeakerName]
		result := s.PowerOn()

		return key.NewPowerOnOK().WithPayload(&models.BStatus{Status: &result})
	})

	api.KeyPowerOffHandler = key.PowerOffHandlerFunc(func(params key.PowerOffParams) middleware.Responder {
		ck, err := checkSpeakerName(params.SpeakerName)
		if !ck {
			return key.NewPowerOnDefault(404).WithPayload(err)
		}

		s := visibleSpeakers[params.SpeakerName]
		result := s.PowerOff()

		return key.NewPowerOffOK().WithPayload(&models.BStatus{Status: &result})
	})

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}

func checkInMap(deviceID string, list speakers) bool {
	for _, ms := range list {
		if ms.DeviceInfo.DeviceID == deviceID {
			return true
		}
	}
	return false
}

func basicHandler(hndlName string, update sndt.Update, speaker sndt.Speaker) {
	typeName := reflect.TypeOf(update.Value).Name()
	switch typeName {
	case "ConnectionStateUpdated":
		return
	}
	log.Infof("%s/%s: %s\n", speaker.Name(), typeName, update)
}

func connectionHandler(hndlName string, update sndt.Update, speaker sndt.Speaker) {
	typeName := reflect.TypeOf(update.Value).Name()
	if typeName == "ConnectionStateUpdated" {
		log.Debugf("%s/%s: %s\n", speaker.Name(), typeName, update)
	}
}

func contentItem(u sndt.Update) sndt.ContentItem {
	if hasContentItem(u) {
		return u.Value.(sndt.NowPlaying).Content
	}
	return sndt.ContentItem{}

}

func hasContentItem(u sndt.Update) bool {
	switch reflect.TypeOf(u.Value).Name() {
	case "NowPlaying":
		return true
	}
	return false
}

func checkSpeakerName(speakerName string) (contained bool, err *models.Error) {
	if visibleSpeakers[speakerName] != nil {
		return true, nil
	}
	var errorMsg = fmt.Sprintf("Speaker %s not found", speakerName)

	err = &models.Error{
		Code:    20,
		Message: &errorMsg,
	}
	return false, err
}
