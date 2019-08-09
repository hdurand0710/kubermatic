// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// AddonSpec AddonSpec addon specification
// swagger:model AddonSpec
type AddonSpec struct {

	// Variables is free form data to use for parsing the manifest templates
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// Validate validates this addon spec
func (m *AddonSpec) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *AddonSpec) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *AddonSpec) UnmarshalBinary(b []byte) error {
	var res AddonSpec
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
