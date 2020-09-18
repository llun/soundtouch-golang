// Code generated by go-swagger; DO NOT EDIT.

package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

var (
	// SwaggerJSON embedded version of the swagger document used at generation time
	SwaggerJSON json.RawMessage
	// FlatSwaggerJSON embedded flattened version of the swagger document used at generation time
	FlatSwaggerJSON json.RawMessage
)

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/berlin.vassiliou-pohl.soundtouch-golang.v1+json"
  ],
  "produces": [
    "application/berlin.vassiliou-pohl.soundtouch-golang.v1+json"
  ],
  "swagger": "2.0",
  "info": {
    "description": "This applications implements a RESTful/json based API for the soundtouch system.\nIt's implemented as an API which is described by this swagger spec document.\n\nThe server discovers automatically all Soundtouch devices present in the LAN.\nIt's intend is to ease the communication with Boses Soundtouch system.\n",
    "title": "Soundtouch RESTful/json server",
    "version": "0.0.1"
  },
  "host": "localhost:5006",
  "paths": {
    "/api/keys-list": {
      "get": {
        "description": "This method will get all possible keys. These keys can be used in the '/{deviceName}/key/{key}' to replace the {key} placeholder.",
        "produces": [
          "application/berlin.vassiliou-pohl.soundtouch-golang.v1+json"
        ],
        "tags": [
          "API"
        ],
        "operationId": "keysList",
        "responses": {
          "200": {
            "description": "a JSON object with all key literals.",
            "schema": {
              "$ref": "#/definitions/keys"
            }
          },
          "default": {
            "description": "generic error response",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    },
    "/device/list": {
      "get": {
        "description": "This url will return a JSON object with the found soundtouch devices on your network. When unplugging one of your soundtouches, the device will not be in the list when making a new requests. Even the soundtouches that are powered off are returned. Per soundtouch, only the basic information is displayed.",
        "produces": [
          "application/berlin.vassiliou-pohl.soundtouch-golang.v1+json"
        ],
        "tags": [
          "device"
        ],
        "operationId": "list",
        "responses": {
          "200": {
            "description": "a JSON object with the found soundtouch devices on your network.",
            "schema": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/device"
              }
            }
          },
          "default": {
            "description": "generic error response",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    },
    "/device/listAdvanced": {
      "get": {
        "description": "This url will return a JSON object with the found soundtouch devices on your network. When unplugging one of your soundtouches, the device will not be in the list when making a new requests. Even the soundtouches that are powered off are returned. This advanced view will display all information known about the soundtouch.",
        "produces": [
          "application/berlin.vassiliou-pohl.soundtouch-golang.v1+json"
        ],
        "tags": [
          "device"
        ],
        "operationId": "listAdvanced",
        "responses": {
          "200": {
            "description": "a JSON object with the found soundtouch devices on your network.",
            "schema": {
              "$ref": "#/definitions/bObject"
            }
          },
          "default": {
            "description": "generic error response",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    },
    "/{speakerName}/key/{keyId}": {
      "get": {
        "description": "Presses and releases a key on selected deviceId",
        "produces": [
          "application/berlin.vassiliou-pohl.soundtouch-golang.v1+json"
        ],
        "tags": [
          "key"
        ],
        "operationId": "pressKey",
        "parameters": [
          {
            "type": "string",
            "description": "The name of device",
            "name": "speakerName",
            "in": "path",
            "required": true
          },
          {
            "enum": [
              "PLAY",
              "PAUSE",
              "PREV_TRACK",
              "NEXT_TRACK",
              "POWER",
              "MUTE",
              "PRESET_1",
              "PRESET_2",
              "PRESET_3",
              "PRESET_4",
              "PRESET_5",
              "PRESET_6",
              "SHUFFLE_OFF",
              "SHUFFLE_ON",
              "REPEAT_OFF",
              "REPEAT_ONE",
              "REPEAT_ALL",
              "PLAY_PAUSE",
              "ADD_FAVORITE",
              "REMOVE_FAVORITE",
              "BOOKMARK",
              "THUMBS_UP",
              "THUMBS_DOWN"
            ],
            "type": "string",
            "name": "keyId",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "204": {
            "description": "empty on success"
          },
          "default": {
            "description": "generic error response",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    },
    "/{speakerName}/nowPlaying": {
      "get": {
        "description": "This method will indicate what's playing at this moment.",
        "produces": [
          "application/berlin.vassiliou-pohl.soundtouch-golang.v1+json"
        ],
        "tags": [
          "key",
          "device"
        ],
        "operationId": "nowPlaying",
        "parameters": [
          {
            "type": "string",
            "description": "The name of device",
            "name": "speakerName",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "a JSON A JSON object indicting what is being played. Returns information both in Standby and PoweredOn mode",
            "schema": {
              "$ref": "#/definitions/bObject"
            }
          },
          "default": {
            "description": "generic error response",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    },
    "/{speakerName}/play": {
      "get": {
        "description": "starts playing",
        "produces": [
          "application/berlin.vassiliou-pohl.soundtouch-golang.v1+json"
        ],
        "tags": [
          "key"
        ],
        "operationId": "play",
        "parameters": [
          {
            "type": "string",
            "name": "speakerName",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "204": {
            "description": "empty on success"
          },
          "default": {
            "description": "generic error response",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    },
    "/{speakerName}/playPause": {
      "get": {
        "description": "This method will play if the device was paused or pause when it was playing.",
        "produces": [
          "application/berlin.vassiliou-pohl.soundtouch-golang.v1+json"
        ],
        "tags": [
          "key"
        ],
        "operationId": "playPause",
        "parameters": [
          {
            "type": "string",
            "name": "speakerName",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "204": {
            "description": "empty on success"
          },
          "default": {
            "description": "generic error response",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    },
    "/{speakerName}/powerOff": {
      "get": {
        "description": "powers off a device",
        "produces": [
          "application/berlin.vassiliou-pohl.soundtouch-golang.v1+json"
        ],
        "tags": [
          "key"
        ],
        "operationId": "powerOff",
        "parameters": [
          {
            "type": "string",
            "name": "speakerName",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "The returned status will be true if the SoundTouch is turned off. The returned status will be false if the SoundTouch was already turned off.",
            "schema": {
              "$ref": "#/definitions/bStatus"
            }
          },
          "default": {
            "description": "generic error response",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    },
    "/{speakerName}/powerOn": {
      "get": {
        "description": "powers on a device",
        "tags": [
          "key"
        ],
        "operationId": "powerOn",
        "parameters": [
          {
            "type": "string",
            "name": "speakerName",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "The returned status will be true if the SoundTouch is turned on. The returned status will be false if the SoundTouch was already turned on.",
            "schema": {
              "$ref": "#/definitions/bStatus"
            }
          },
          "default": {
            "description": "generic error response",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    },
    "/{speakerName}/presets": {
      "get": {
        "description": "This method will return all 6 presets that are configured on the SoundTouch device.",
        "produces": [
          "application/berlin.vassiliou-pohl.soundtouch-golang.v1+json"
        ],
        "tags": [
          "key",
          "device"
        ],
        "operationId": "presets",
        "parameters": [
          {
            "type": "string",
            "description": "The name of device",
            "name": "speakerName",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "a JSON A JSON object describing the presets. Returns information both in Standby and PoweredOn mode",
            "schema": {
              "$ref": "#/definitions/bObject"
            }
          },
          "default": {
            "description": "generic error response",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    },
    "/{speakerName}/trackInfo": {
      "get": {
        "description": "Get more information on what is currently played",
        "produces": [
          "application/berlin.vassiliou-pohl.soundtouch-golang.v1+json"
        ],
        "tags": [
          "key",
          "device"
        ],
        "operationId": "trackInfo",
        "parameters": [
          {
            "type": "string",
            "description": "The name of device",
            "name": "speakerName",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "a JSON A JSON object indicting what is being played. Returns information both in Standby and PoweredOn mode",
            "schema": {
              "$ref": "#/definitions/bObject"
            }
          },
          "default": {
            "description": "generic error response",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "bObject": {
      "type": "object"
    },
    "bStatus": {
      "type": "object",
      "required": [
        "status"
      ],
      "properties": {
        "status": {
          "type": "boolean",
          "format": "int64"
        }
      }
    },
    "device": {
      "type": "object",
      "properties": {
        "addresses": {
          "type": "array",
          "items": {
            "type": "string",
            "example": "10.0.0.7"
          }
        },
        "name": {
          "type": "string",
          "example": "Bathroom"
        }
      }
    },
    "error": {
      "type": "object",
      "required": [
        "message"
      ],
      "properties": {
        "code": {
          "type": "integer",
          "format": "int64"
        },
        "message": {
          "type": "string"
        }
      }
    },
    "keys": {
      "type": "array",
      "items": {
        "type": "string",
        "enum": [
          "PLAY",
          "PAUSE",
          "STOP",
          "PREV_TRACK",
          "NEXT_TRACK",
          "POWER",
          "MUTE",
          "VOLUME_UP",
          "VOLUME_DOWN",
          "PRESET_1",
          "PRESET_2",
          "PRESET_3",
          "PRESET_4",
          "PRESET_5",
          "PRESET_6",
          "AUX_INPUT",
          "SHUFFLE_OFF",
          "SHUFFLE_ON",
          "REPEAT_OFF",
          "REPEAT_ONE",
          "REPEAT_ALL",
          "PLAY_PAUSE",
          "ADD_FAVORITE",
          "REMOVE_FAVORITE",
          "BOOKMARK",
          "THUMBS_UP",
          "THUMBS_DOWN"
        ]
      }
    }
  }
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/berlin.vassiliou-pohl.soundtouch-golang.v1+json"
  ],
  "produces": [
    "application/berlin.vassiliou-pohl.soundtouch-golang.v1+json"
  ],
  "swagger": "2.0",
  "info": {
    "description": "This applications implements a RESTful/json based API for the soundtouch system.\nIt's implemented as an API which is described by this swagger spec document.\n\nThe server discovers automatically all Soundtouch devices present in the LAN.\nIt's intend is to ease the communication with Boses Soundtouch system.\n",
    "title": "Soundtouch RESTful/json server",
    "version": "0.0.1"
  },
  "host": "localhost:5006",
  "paths": {
    "/api/keys-list": {
      "get": {
        "description": "This method will get all possible keys. These keys can be used in the '/{deviceName}/key/{key}' to replace the {key} placeholder.",
        "produces": [
          "application/berlin.vassiliou-pohl.soundtouch-golang.v1+json"
        ],
        "tags": [
          "API"
        ],
        "operationId": "keysList",
        "responses": {
          "200": {
            "description": "a JSON object with all key literals.",
            "schema": {
              "$ref": "#/definitions/keys"
            }
          },
          "default": {
            "description": "generic error response",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    },
    "/device/list": {
      "get": {
        "description": "This url will return a JSON object with the found soundtouch devices on your network. When unplugging one of your soundtouches, the device will not be in the list when making a new requests. Even the soundtouches that are powered off are returned. Per soundtouch, only the basic information is displayed.",
        "produces": [
          "application/berlin.vassiliou-pohl.soundtouch-golang.v1+json"
        ],
        "tags": [
          "device"
        ],
        "operationId": "list",
        "responses": {
          "200": {
            "description": "a JSON object with the found soundtouch devices on your network.",
            "schema": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/device"
              }
            }
          },
          "default": {
            "description": "generic error response",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    },
    "/device/listAdvanced": {
      "get": {
        "description": "This url will return a JSON object with the found soundtouch devices on your network. When unplugging one of your soundtouches, the device will not be in the list when making a new requests. Even the soundtouches that are powered off are returned. This advanced view will display all information known about the soundtouch.",
        "produces": [
          "application/berlin.vassiliou-pohl.soundtouch-golang.v1+json"
        ],
        "tags": [
          "device"
        ],
        "operationId": "listAdvanced",
        "responses": {
          "200": {
            "description": "a JSON object with the found soundtouch devices on your network.",
            "schema": {
              "$ref": "#/definitions/bObject"
            }
          },
          "default": {
            "description": "generic error response",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    },
    "/{speakerName}/key/{keyId}": {
      "get": {
        "description": "Presses and releases a key on selected deviceId",
        "produces": [
          "application/berlin.vassiliou-pohl.soundtouch-golang.v1+json"
        ],
        "tags": [
          "key"
        ],
        "operationId": "pressKey",
        "parameters": [
          {
            "type": "string",
            "description": "The name of device",
            "name": "speakerName",
            "in": "path",
            "required": true
          },
          {
            "enum": [
              "PLAY",
              "PAUSE",
              "PREV_TRACK",
              "NEXT_TRACK",
              "POWER",
              "MUTE",
              "PRESET_1",
              "PRESET_2",
              "PRESET_3",
              "PRESET_4",
              "PRESET_5",
              "PRESET_6",
              "SHUFFLE_OFF",
              "SHUFFLE_ON",
              "REPEAT_OFF",
              "REPEAT_ONE",
              "REPEAT_ALL",
              "PLAY_PAUSE",
              "ADD_FAVORITE",
              "REMOVE_FAVORITE",
              "BOOKMARK",
              "THUMBS_UP",
              "THUMBS_DOWN"
            ],
            "type": "string",
            "name": "keyId",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "204": {
            "description": "empty on success"
          },
          "default": {
            "description": "generic error response",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    },
    "/{speakerName}/nowPlaying": {
      "get": {
        "description": "This method will indicate what's playing at this moment.",
        "produces": [
          "application/berlin.vassiliou-pohl.soundtouch-golang.v1+json"
        ],
        "tags": [
          "key",
          "device"
        ],
        "operationId": "nowPlaying",
        "parameters": [
          {
            "type": "string",
            "description": "The name of device",
            "name": "speakerName",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "a JSON A JSON object indicting what is being played. Returns information both in Standby and PoweredOn mode",
            "schema": {
              "$ref": "#/definitions/bObject"
            }
          },
          "default": {
            "description": "generic error response",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    },
    "/{speakerName}/play": {
      "get": {
        "description": "starts playing",
        "produces": [
          "application/berlin.vassiliou-pohl.soundtouch-golang.v1+json"
        ],
        "tags": [
          "key"
        ],
        "operationId": "play",
        "parameters": [
          {
            "type": "string",
            "name": "speakerName",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "204": {
            "description": "empty on success"
          },
          "default": {
            "description": "generic error response",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    },
    "/{speakerName}/playPause": {
      "get": {
        "description": "This method will play if the device was paused or pause when it was playing.",
        "produces": [
          "application/berlin.vassiliou-pohl.soundtouch-golang.v1+json"
        ],
        "tags": [
          "key"
        ],
        "operationId": "playPause",
        "parameters": [
          {
            "type": "string",
            "name": "speakerName",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "204": {
            "description": "empty on success"
          },
          "default": {
            "description": "generic error response",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    },
    "/{speakerName}/powerOff": {
      "get": {
        "description": "powers off a device",
        "produces": [
          "application/berlin.vassiliou-pohl.soundtouch-golang.v1+json"
        ],
        "tags": [
          "key"
        ],
        "operationId": "powerOff",
        "parameters": [
          {
            "type": "string",
            "name": "speakerName",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "The returned status will be true if the SoundTouch is turned off. The returned status will be false if the SoundTouch was already turned off.",
            "schema": {
              "$ref": "#/definitions/bStatus"
            }
          },
          "default": {
            "description": "generic error response",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    },
    "/{speakerName}/powerOn": {
      "get": {
        "description": "powers on a device",
        "tags": [
          "key"
        ],
        "operationId": "powerOn",
        "parameters": [
          {
            "type": "string",
            "name": "speakerName",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "The returned status will be true if the SoundTouch is turned on. The returned status will be false if the SoundTouch was already turned on.",
            "schema": {
              "$ref": "#/definitions/bStatus"
            }
          },
          "default": {
            "description": "generic error response",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    },
    "/{speakerName}/presets": {
      "get": {
        "description": "This method will return all 6 presets that are configured on the SoundTouch device.",
        "produces": [
          "application/berlin.vassiliou-pohl.soundtouch-golang.v1+json"
        ],
        "tags": [
          "key",
          "device"
        ],
        "operationId": "presets",
        "parameters": [
          {
            "type": "string",
            "description": "The name of device",
            "name": "speakerName",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "a JSON A JSON object describing the presets. Returns information both in Standby and PoweredOn mode",
            "schema": {
              "$ref": "#/definitions/bObject"
            }
          },
          "default": {
            "description": "generic error response",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    },
    "/{speakerName}/trackInfo": {
      "get": {
        "description": "Get more information on what is currently played",
        "produces": [
          "application/berlin.vassiliou-pohl.soundtouch-golang.v1+json"
        ],
        "tags": [
          "key",
          "device"
        ],
        "operationId": "trackInfo",
        "parameters": [
          {
            "type": "string",
            "description": "The name of device",
            "name": "speakerName",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "a JSON A JSON object indicting what is being played. Returns information both in Standby and PoweredOn mode",
            "schema": {
              "$ref": "#/definitions/bObject"
            }
          },
          "default": {
            "description": "generic error response",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "bObject": {
      "type": "object"
    },
    "bStatus": {
      "type": "object",
      "required": [
        "status"
      ],
      "properties": {
        "status": {
          "type": "boolean",
          "format": "int64"
        }
      }
    },
    "device": {
      "type": "object",
      "properties": {
        "addresses": {
          "type": "array",
          "items": {
            "type": "string",
            "example": "10.0.0.7"
          }
        },
        "name": {
          "type": "string",
          "example": "Bathroom"
        }
      }
    },
    "error": {
      "type": "object",
      "required": [
        "message"
      ],
      "properties": {
        "code": {
          "type": "integer",
          "format": "int64"
        },
        "message": {
          "type": "string"
        }
      }
    },
    "keys": {
      "type": "array",
      "items": {
        "type": "string",
        "enum": [
          "PLAY",
          "PAUSE",
          "STOP",
          "PREV_TRACK",
          "NEXT_TRACK",
          "POWER",
          "MUTE",
          "VOLUME_UP",
          "VOLUME_DOWN",
          "PRESET_1",
          "PRESET_2",
          "PRESET_3",
          "PRESET_4",
          "PRESET_5",
          "PRESET_6",
          "AUX_INPUT",
          "SHUFFLE_OFF",
          "SHUFFLE_ON",
          "REPEAT_OFF",
          "REPEAT_ONE",
          "REPEAT_ALL",
          "PLAY_PAUSE",
          "ADD_FAVORITE",
          "REMOVE_FAVORITE",
          "BOOKMARK",
          "THUMBS_UP",
          "THUMBS_DOWN"
        ]
      }
    }
  }
}`))
}