// Code generated by go-swagger; DO NOT EDIT.

package kubevirt

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// New creates a new kubevirt API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) ClientService {
	return &Client{transport: transport, formats: formats}
}

/*
Client for kubevirt API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

// ClientOption is the option for Client methods
type ClientOption func(*runtime.ClientOperation)

// ClientService is the interface for Client methods
type ClientService interface {
	GetKubevirtStorageClass(params *GetKubevirtStorageClassParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*GetKubevirtStorageClassOK, error)

	GetKubevirtStorageClassNoCredentials(params *GetKubevirtStorageClassNoCredentialsParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*GetKubevirtStorageClassNoCredentialsOK, error)

	GetKubevirtVmiPreset(params *GetKubevirtVmiPresetParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*GetKubevirtVmiPresetOK, error)

	GetKubevirtVmiPresetNoCredentials(params *GetKubevirtVmiPresetNoCredentialsParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*GetKubevirtVmiPresetNoCredentialsOK, error)

	ListKubevirtStorageClasses(params *ListKubevirtStorageClassesParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListKubevirtStorageClassesOK, error)

	ListKubevirtStorageClassesNoCredentials(params *ListKubevirtStorageClassesNoCredentialsParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListKubevirtStorageClassesNoCredentialsOK, error)

	ListKubevirtVmiPresets(params *ListKubevirtVmiPresetsParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListKubevirtVmiPresetsOK, error)

	ListKubevirtVmiPresetsNoCredentials(params *ListKubevirtVmiPresetsNoCredentialsParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListKubevirtVmiPresetsNoCredentialsOK, error)

	SetTransport(transport runtime.ClientTransport)
}

/*
  GetKubevirtStorageClass gets a k8s storage class in the kubevirt cluster
*/
func (a *Client) GetKubevirtStorageClass(params *GetKubevirtStorageClassParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*GetKubevirtStorageClassOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetKubevirtStorageClassParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "getKubevirtStorageClass",
		Method:             "GET",
		PathPattern:        "/api/v2/providers/kubevirt/storageclasses/{storageclass_id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetKubevirtStorageClassReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetKubevirtStorageClassOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*GetKubevirtStorageClassDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
  GetKubevirtStorageClassNoCredentials Get a Storage Class
*/
func (a *Client) GetKubevirtStorageClassNoCredentials(params *GetKubevirtStorageClassNoCredentialsParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*GetKubevirtStorageClassNoCredentialsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetKubevirtStorageClassNoCredentialsParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "getKubevirtStorageClassNoCredentials",
		Method:             "GET",
		PathPattern:        "/api/v2/projects/{project_id}/clusters/{cluster_id}/providers/kubevirt/storageclasses/{storageclass_id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetKubevirtStorageClassNoCredentialsReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetKubevirtStorageClassNoCredentialsOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*GetKubevirtStorageClassNoCredentialsDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
  GetKubevirtVmiPreset gets a kubevirt virtual machine instance preset
*/
func (a *Client) GetKubevirtVmiPreset(params *GetKubevirtVmiPresetParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*GetKubevirtVmiPresetOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetKubevirtVmiPresetParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "getKubevirtVmiPreset",
		Method:             "GET",
		PathPattern:        "/api/v2/providers/kubevirt/vmipresets/{preset_id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetKubevirtVmiPresetReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetKubevirtVmiPresetOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*GetKubevirtVmiPresetDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
  GetKubevirtVmiPresetNoCredentials Get a VirtualMachineInstancePreset
*/
func (a *Client) GetKubevirtVmiPresetNoCredentials(params *GetKubevirtVmiPresetNoCredentialsParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*GetKubevirtVmiPresetNoCredentialsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetKubevirtVmiPresetNoCredentialsParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "getKubevirtVmiPresetNoCredentials",
		Method:             "GET",
		PathPattern:        "/api/v2/projects/{project_id}/clusters/{cluster_id}/providers/kubevirt/vmipresets/{preset_id}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetKubevirtVmiPresetNoCredentialsReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetKubevirtVmiPresetNoCredentialsOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*GetKubevirtVmiPresetNoCredentialsDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
  ListKubevirtStorageClasses lists available k8s storage classes in the kubevirt cluster
*/
func (a *Client) ListKubevirtStorageClasses(params *ListKubevirtStorageClassesParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListKubevirtStorageClassesOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListKubevirtStorageClassesParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "listKubevirtStorageClasses",
		Method:             "GET",
		PathPattern:        "/api/v2/providers/kubevirt/storageclasses",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &ListKubevirtStorageClassesReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*ListKubevirtStorageClassesOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*ListKubevirtStorageClassesDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
  ListKubevirtStorageClassesNoCredentials List Storage Classes
*/
func (a *Client) ListKubevirtStorageClassesNoCredentials(params *ListKubevirtStorageClassesNoCredentialsParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListKubevirtStorageClassesNoCredentialsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListKubevirtStorageClassesNoCredentialsParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "listKubevirtStorageClassesNoCredentials",
		Method:             "GET",
		PathPattern:        "/api/v2/projects/{project_id}/clusters/{cluster_id}/providers/kubevirt/storageclasses",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &ListKubevirtStorageClassesNoCredentialsReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*ListKubevirtStorageClassesNoCredentialsOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*ListKubevirtStorageClassesNoCredentialsDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
  ListKubevirtVmiPresets lists available kubevirt virtual machine instance preset
*/
func (a *Client) ListKubevirtVmiPresets(params *ListKubevirtVmiPresetsParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListKubevirtVmiPresetsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListKubevirtVmiPresetsParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "listKubevirtVmiPresets",
		Method:             "GET",
		PathPattern:        "/api/v2/providers/kubevirt/vmipresets",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &ListKubevirtVmiPresetsReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*ListKubevirtVmiPresetsOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*ListKubevirtVmiPresetsDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
  ListKubevirtVmiPresetsNoCredentials Lists available VirtualMachineInstancePreset
*/
func (a *Client) ListKubevirtVmiPresetsNoCredentials(params *ListKubevirtVmiPresetsNoCredentialsParams, authInfo runtime.ClientAuthInfoWriter, opts ...ClientOption) (*ListKubevirtVmiPresetsNoCredentialsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewListKubevirtVmiPresetsNoCredentialsParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "listKubevirtVmiPresetsNoCredentials",
		Method:             "GET",
		PathPattern:        "/api/v2/projects/{project_id}/clusters/{cluster_id}/providers/kubevirt/vmipresets",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &ListKubevirtVmiPresetsNoCredentialsReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*ListKubevirtVmiPresetsNoCredentialsOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*ListKubevirtVmiPresetsNoCredentialsDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
