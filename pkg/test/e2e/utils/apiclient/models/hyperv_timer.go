// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// HypervTimer hyperv timer
//
// swagger:model HypervTimer
type HypervTimer struct {

	// Enabled set to false makes sure that the machine type or a preset can't add the timer.
	// Defaults to true.
	// +optional
	Enabled bool `json:"present,omitempty"`
}

// Validate validates this hyperv timer
func (m *HypervTimer) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this hyperv timer based on context it is used
func (m *HypervTimer) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *HypervTimer) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HypervTimer) UnmarshalBinary(b []byte) error {
	var res HypervTimer
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
