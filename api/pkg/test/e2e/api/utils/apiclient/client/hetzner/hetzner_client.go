// Code generated by go-swagger; DO NOT EDIT.

package hetzner

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// New creates a new hetzner API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) *Client {
	return &Client{transport: transport, formats: formats}
}

/*
Client for hetzner API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

/*
ListHetznerSizes Lists sizes from hetzner
*/
func (a *Client) ListHetznerSizes(params *ListHetznerSizesParams, authInfo runtime.ClientAuthInfoWriter) (*ListHetznerSizesOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListHetznerSizesParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "listHetznerSizes",
		Method:             "GET",
		PathPattern:        "/api/v1/providers/hetzner/sizes",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &ListHetznerSizesReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*ListHetznerSizesOK), nil

}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
