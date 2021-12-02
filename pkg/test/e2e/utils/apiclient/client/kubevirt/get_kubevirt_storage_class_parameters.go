// Code generated by go-swagger; DO NOT EDIT.

package kubevirt

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewGetKubevirtStorageClassParams creates a new GetKubevirtStorageClassParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetKubevirtStorageClassParams() *GetKubevirtStorageClassParams {
	return &GetKubevirtStorageClassParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetKubevirtStorageClassParamsWithTimeout creates a new GetKubevirtStorageClassParams object
// with the ability to set a timeout on a request.
func NewGetKubevirtStorageClassParamsWithTimeout(timeout time.Duration) *GetKubevirtStorageClassParams {
	return &GetKubevirtStorageClassParams{
		timeout: timeout,
	}
}

// NewGetKubevirtStorageClassParamsWithContext creates a new GetKubevirtStorageClassParams object
// with the ability to set a context for a request.
func NewGetKubevirtStorageClassParamsWithContext(ctx context.Context) *GetKubevirtStorageClassParams {
	return &GetKubevirtStorageClassParams{
		Context: ctx,
	}
}

// NewGetKubevirtStorageClassParamsWithHTTPClient creates a new GetKubevirtStorageClassParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetKubevirtStorageClassParamsWithHTTPClient(client *http.Client) *GetKubevirtStorageClassParams {
	return &GetKubevirtStorageClassParams{
		HTTPClient: client,
	}
}

/* GetKubevirtStorageClassParams contains all the parameters to send to the API endpoint
   for the get kubevirt storage class operation.

   Typically these are written to a http.Request.
*/
type GetKubevirtStorageClassParams struct {

	// Credential.
	Credential *string

	// KvKubeconfig.
	KvKubeconfig *string

	// StorageclassID.
	StorageClass string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get kubevirt storage class params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetKubevirtStorageClassParams) WithDefaults() *GetKubevirtStorageClassParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get kubevirt storage class params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetKubevirtStorageClassParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the get kubevirt storage class params
func (o *GetKubevirtStorageClassParams) WithTimeout(timeout time.Duration) *GetKubevirtStorageClassParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get kubevirt storage class params
func (o *GetKubevirtStorageClassParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get kubevirt storage class params
func (o *GetKubevirtStorageClassParams) WithContext(ctx context.Context) *GetKubevirtStorageClassParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get kubevirt storage class params
func (o *GetKubevirtStorageClassParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get kubevirt storage class params
func (o *GetKubevirtStorageClassParams) WithHTTPClient(client *http.Client) *GetKubevirtStorageClassParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get kubevirt storage class params
func (o *GetKubevirtStorageClassParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithCredential adds the credential to the get kubevirt storage class params
func (o *GetKubevirtStorageClassParams) WithCredential(credential *string) *GetKubevirtStorageClassParams {
	o.SetCredential(credential)
	return o
}

// SetCredential adds the credential to the get kubevirt storage class params
func (o *GetKubevirtStorageClassParams) SetCredential(credential *string) {
	o.Credential = credential
}

// WithKvKubeconfig adds the kvKubeconfig to the get kubevirt storage class params
func (o *GetKubevirtStorageClassParams) WithKvKubeconfig(kvKubeconfig *string) *GetKubevirtStorageClassParams {
	o.SetKvKubeconfig(kvKubeconfig)
	return o
}

// SetKvKubeconfig adds the kvKubeconfig to the get kubevirt storage class params
func (o *GetKubevirtStorageClassParams) SetKvKubeconfig(kvKubeconfig *string) {
	o.KvKubeconfig = kvKubeconfig
}

// WithStorageClass adds the storageclassID to the get kubevirt storage class params
func (o *GetKubevirtStorageClassParams) WithStorageClass(storageclassID string) *GetKubevirtStorageClassParams {
	o.SetStorageClass(storageclassID)
	return o
}

// SetStorageClass adds the storageclassId to the get kubevirt storage class params
func (o *GetKubevirtStorageClassParams) SetStorageClass(storageclassID string) {
	o.StorageClass = storageclassID
}

// WriteToRequest writes these params to a swagger request
func (o *GetKubevirtStorageClassParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Credential != nil {

		// header param Credential
		if err := r.SetHeaderParam("Credential", *o.Credential); err != nil {
			return err
		}
	}

	if o.KvKubeconfig != nil {

		// header param KvKubeconfig
		if err := r.SetHeaderParam("KvKubeconfig", *o.KvKubeconfig); err != nil {
			return err
		}
	}

	// path param storageclass_id
	if err := r.SetPathParam("storageclass_id", o.StorageClass); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
