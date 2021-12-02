// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// Bootloader Represents the firmware blob used to assist in the domain creation process.
//
// Used for setting the QEMU BIOS file path for the libvirt domain.
//
// swagger:model Bootloader
type Bootloader struct {

	// bios
	Bios *BIOS `json:"bios,omitempty"`

	// efi
	Efi *EFI `json:"efi,omitempty"`
}

// Validate validates this bootloader
func (m *Bootloader) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateBios(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateEfi(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Bootloader) validateBios(formats strfmt.Registry) error {
	if swag.IsZero(m.Bios) { // not required
		return nil
	}

	if m.Bios != nil {
		if err := m.Bios.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("bios")
			}
			return err
		}
	}

	return nil
}

func (m *Bootloader) validateEfi(formats strfmt.Registry) error {
	if swag.IsZero(m.Efi) { // not required
		return nil
	}

	if m.Efi != nil {
		if err := m.Efi.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("efi")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this bootloader based on the context it is used
func (m *Bootloader) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateBios(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateEfi(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Bootloader) contextValidateBios(ctx context.Context, formats strfmt.Registry) error {

	if m.Bios != nil {
		if err := m.Bios.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("bios")
			}
			return err
		}
	}

	return nil
}

func (m *Bootloader) contextValidateEfi(ctx context.Context, formats strfmt.Registry) error {

	if m.Efi != nil {
		if err := m.Efi.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("efi")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Bootloader) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Bootloader) UnmarshalBinary(b []byte) error {
	var res Bootloader
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
