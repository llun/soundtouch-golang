// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// Keys keys
//
// swagger:model keys
type Keys []string

var keysItemsEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["PLAY","PAUSE","STOP","PREV_TRACK","NEXT_TRACK","POWER","MUTE","VOLUME_UP","VOLUME_DOWN","PRESET_1","PRESET_2","PRESET_3","PRESET_4","PRESET_5","PRESET_6","AUX_INPUT","SHUFFLE_OFF","SHUFFLE_ON","REPEAT_OFF","REPEAT_ONE","REPEAT_ALL","PLAY_PAUSE","ADD_FAVORITE","REMOVE_FAVORITE","BOOKMARK","THUMBS_UP","THUMBS_DOWN"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		keysItemsEnum = append(keysItemsEnum, v)
	}
}

func (m *Keys) validateKeysItemsEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, keysItemsEnum, true); err != nil {
		return err
	}
	return nil
}

// Validate validates this keys
func (m Keys) Validate(formats strfmt.Registry) error {
	var res []error

	for i := 0; i < len(m); i++ {

		// value enum
		if err := m.validateKeysItemsEnum(strconv.Itoa(i), "body", m[i]); err != nil {
			return err
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}