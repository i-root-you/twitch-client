package obsws

import (
	"errors"
	"time"
)

// This file is automatically generated.
// https://github.com/christopher-dG/go-obs-websocket/blob/master/codegen/protocol.py

// GetVersionRequest : Returns the latest version of the plugin and the API.
//
// Since obs-websocket version: 0.3.
//
// https://github.com/Palakis/obs-websocket/blob/4.3-maintenance/docs/generated/protocol.md#getversion
type GetVersionRequest struct {
	_request `json:",squash"`
	response chan GetVersionResponse
}

// NewGetVersionRequest returns a new GetVersionRequest.
func NewGetVersionRequest() GetVersionRequest {
	return GetVersionRequest{
		_request{
			ID_:   getMessageID(),
			Type_: "GetVersion",
			err:   make(chan error, 1),
		},
		make(chan GetVersionResponse, 1),
	}
}

// Send sends the request.
func (r *GetVersionRequest) Send(c Client) error {
	if r.sent {
		return ErrAlreadySent
	}
	future, err := c.sendRequest(r)
	if err != nil {
		return err
	}
	r.sent = true
	go func() {
		m := <-future
		var resp GetVersionResponse
		if err = mapToStruct(m, &resp); err != nil {
			r.err <- err
		} else if resp.Status() != StatusOK {
			r.err <- errors.New(resp.Error())
		} else {
			r.response <- resp
		}
	}()
	return nil
}

// Receive waits for the response.
func (r GetVersionRequest) Receive() (GetVersionResponse, error) {
	if !r.sent {
		return GetVersionResponse{}, ErrNotSent
	}
	if receiveTimeout == 0 {
		select {
		case resp := <-r.response:
			return resp, nil
		case err := <-r.err:
			return GetVersionResponse{}, err
		}
	} else {
		select {
		case resp := <-r.response:
			return resp, nil
		case err := <-r.err:
			return GetVersionResponse{}, err
		case <-time.After(receiveTimeout):
			return GetVersionResponse{}, ErrReceiveTimeout
		}
	}
}

// SendReceive sends the request then immediately waits for the response.
func (r GetVersionRequest) SendReceive(c Client) (GetVersionResponse, error) {
	if err := r.Send(c); err != nil {
		return GetVersionResponse{}, err
	}
	return r.Receive()
}

// GetVersionResponse : Response for GetVersionRequest.
//
// Since obs-websocket version: 0.3.
//
// https://github.com/Palakis/obs-websocket/blob/4.3-maintenance/docs/generated/protocol.md#getversion
type GetVersionResponse struct {
	// OBSRemote compatible API version.
	// Fixed to 1.1 for retrocompatibility.
	// Required: Yes.
	Version float64 `json:"version"`
	// obs-websocket plugin version.
	// Required: Yes.
	OBSWebsocketVersion string `json:"obs-websocket-version"`
	// OBS Studio program version.
	// Required: Yes.
	OBSStudioVersion string `json:"obs-studio-version"`
	// List of available request types, formatted as a comma-separated list string (e.g. : "Method1,Method2,Method3").
	// Required: Yes.
	AvailableRequests string `json:"available-requests"`
	_response         `json:",squash"`
}

// GetAuthRequiredRequest : Tells the client if authentication is required
// If so, returns authentication parameters `challenge`
// and `salt` (see "Authentication" for more information).
//
// Since obs-websocket version: 0.3.
//
// https://github.com/Palakis/obs-websocket/blob/4.3-maintenance/docs/generated/protocol.md#getauthrequired
type GetAuthRequiredRequest struct {
	_request `json:",squash"`
	response chan GetAuthRequiredResponse
}

// NewGetAuthRequiredRequest returns a new GetAuthRequiredRequest.
func NewGetAuthRequiredRequest() GetAuthRequiredRequest {
	return GetAuthRequiredRequest{
		_request{
			ID_:   getMessageID(),
			Type_: "GetAuthRequired",
			err:   make(chan error, 1),
		},
		make(chan GetAuthRequiredResponse, 1),
	}
}

// Send sends the request.
func (r *GetAuthRequiredRequest) Send(c Client) error {
	if r.sent {
		return ErrAlreadySent
	}
	future, err := c.sendRequest(r)
	if err != nil {
		return err
	}
	r.sent = true
	go func() {
		m := <-future
		var resp GetAuthRequiredResponse
		if err = mapToStruct(m, &resp); err != nil {
			r.err <- err
		} else if resp.Status() != StatusOK {
			r.err <- errors.New(resp.Error())
		} else {
			r.response <- resp
		}
	}()
	return nil
}

// Receive waits for the response.
func (r GetAuthRequiredRequest) Receive() (GetAuthRequiredResponse, error) {
	if !r.sent {
		return GetAuthRequiredResponse{}, ErrNotSent
	}
	if receiveTimeout == 0 {
		select {
		case resp := <-r.response:
			return resp, nil
		case err := <-r.err:
			return GetAuthRequiredResponse{}, err
		}
	} else {
		select {
		case resp := <-r.response:
			return resp, nil
		case err := <-r.err:
			return GetAuthRequiredResponse{}, err
		case <-time.After(receiveTimeout):
			return GetAuthRequiredResponse{}, ErrReceiveTimeout
		}
	}
}

// SendReceive sends the request then immediately waits for the response.
func (r GetAuthRequiredRequest) SendReceive(c Client) (GetAuthRequiredResponse, error) {
	if err := r.Send(c); err != nil {
		return GetAuthRequiredResponse{}, err
	}
	return r.Receive()
}

// GetAuthRequiredResponse : Response for GetAuthRequiredRequest.
//
// Since obs-websocket version: 0.3.
//
// https://github.com/Palakis/obs-websocket/blob/4.3-maintenance/docs/generated/protocol.md#getauthrequired
type GetAuthRequiredResponse struct {
	// Indicates whether authentication is required.
	// Required: Yes.
	AuthRequired bool `json:"authRequired"`
	// Required: No.
	Challenge string `json:"challenge"`
	// Required: No.
	Salt      string `json:"salt"`
	_response `json:",squash"`
}

// AuthenticateRequest : Attempt to authenticate the client to the server.
//
// Since obs-websocket version: 0.3.
//
// https://github.com/Palakis/obs-websocket/blob/4.3-maintenance/docs/generated/protocol.md#authenticate
type AuthenticateRequest struct {
	// Response to the auth challenge (see "Authentication" for more information).
	// Required: Yes.
	Auth     string `json:"auth"`
	_request `json:",squash"`
	response chan AuthenticateResponse
}

// NewAuthenticateRequest returns a new AuthenticateRequest.
func NewAuthenticateRequest(auth string) AuthenticateRequest {
	return AuthenticateRequest{
		auth,
		_request{
			ID_:   getMessageID(),
			Type_: "Authenticate",
			err:   make(chan error, 1),
		},
		make(chan AuthenticateResponse, 1),
	}
}

// Send sends the request.
func (r *AuthenticateRequest) Send(c Client) error {
	if r.sent {
		return ErrAlreadySent
	}
	future, err := c.sendRequest(r)
	if err != nil {
		return err
	}
	r.sent = true
	go func() {
		m := <-future
		var resp AuthenticateResponse
		if err = mapToStruct(m, &resp); err != nil {
			r.err <- err
		} else if resp.Status() != StatusOK {
			r.err <- errors.New(resp.Error())
		} else {
			r.response <- resp
		}
	}()
	return nil
}

// Receive waits for the response.
func (r AuthenticateRequest) Receive() (AuthenticateResponse, error) {
	if !r.sent {
		return AuthenticateResponse{}, ErrNotSent
	}
	if receiveTimeout == 0 {
		select {
		case resp := <-r.response:
			return resp, nil
		case err := <-r.err:
			return AuthenticateResponse{}, err
		}
	} else {
		select {
		case resp := <-r.response:
			return resp, nil
		case err := <-r.err:
			return AuthenticateResponse{}, err
		case <-time.After(receiveTimeout):
			return AuthenticateResponse{}, ErrReceiveTimeout
		}
	}
}

// SendReceive sends the request then immediately waits for the response.
func (r AuthenticateRequest) SendReceive(c Client) (AuthenticateResponse, error) {
	if err := r.Send(c); err != nil {
		return AuthenticateResponse{}, err
	}
	return r.Receive()
}

// AuthenticateResponse : Response for AuthenticateRequest.
//
// Since obs-websocket version: 0.3.
//
// https://github.com/Palakis/obs-websocket/blob/4.3-maintenance/docs/generated/protocol.md#authenticate
type AuthenticateResponse struct {
	_response `json:",squash"`
}

// SetHeartbeatRequest : Enable/disable sending of the Heartbeat event.
//
// Since obs-websocket version: 4.3.0.
//
// https://github.com/Palakis/obs-websocket/blob/4.3-maintenance/docs/generated/protocol.md#setheartbeat
type SetHeartbeatRequest struct {
	// Starts/Stops emitting heartbeat messages.
	// Required: Yes.
	Enable   bool `json:"enable"`
	_request `json:",squash"`
	response chan SetHeartbeatResponse
}

// NewSetHeartbeatRequest returns a new SetHeartbeatRequest.
func NewSetHeartbeatRequest(enable bool) SetHeartbeatRequest {
	return SetHeartbeatRequest{
		enable,
		_request{
			ID_:   getMessageID(),
			Type_: "SetHeartbeat",
			err:   make(chan error, 1),
		},
		make(chan SetHeartbeatResponse, 1),
	}
}

// Send sends the request.
func (r *SetHeartbeatRequest) Send(c Client) error {
	if r.sent {
		return ErrAlreadySent
	}
	future, err := c.sendRequest(r)
	if err != nil {
		return err
	}
	r.sent = true
	go func() {
		m := <-future
		var resp SetHeartbeatResponse
		if err = mapToStruct(m, &resp); err != nil {
			r.err <- err
		} else if resp.Status() != StatusOK {
			r.err <- errors.New(resp.Error())
		} else {
			r.response <- resp
		}
	}()
	return nil
}

// Receive waits for the response.
func (r SetHeartbeatRequest) Receive() (SetHeartbeatResponse, error) {
	if !r.sent {
		return SetHeartbeatResponse{}, ErrNotSent
	}
	if receiveTimeout == 0 {
		select {
		case resp := <-r.response:
			return resp, nil
		case err := <-r.err:
			return SetHeartbeatResponse{}, err
		}
	} else {
		select {
		case resp := <-r.response:
			return resp, nil
		case err := <-r.err:
			return SetHeartbeatResponse{}, err
		case <-time.After(receiveTimeout):
			return SetHeartbeatResponse{}, ErrReceiveTimeout
		}
	}
}

// SendReceive sends the request then immediately waits for the response.
func (r SetHeartbeatRequest) SendReceive(c Client) (SetHeartbeatResponse, error) {
	if err := r.Send(c); err != nil {
		return SetHeartbeatResponse{}, err
	}
	return r.Receive()
}

// SetHeartbeatResponse : Response for SetHeartbeatRequest.
//
// Since obs-websocket version: 4.3.0.
//
// https://github.com/Palakis/obs-websocket/blob/4.3-maintenance/docs/generated/protocol.md#setheartbeat
type SetHeartbeatResponse struct {
	_response `json:",squash"`
}

// SetFilenameFormattingRequest : Set the filename formatting string.
//
// Since obs-websocket version: 4.3.0.
//
// https://github.com/Palakis/obs-websocket/blob/4.3-maintenance/docs/generated/protocol.md#setfilenameformatting
type SetFilenameFormattingRequest struct {
	// Filename formatting string to set.
	// Required: Yes.
	FilenameFormatting string `json:"filename-formatting"`
	_request           `json:",squash"`
	response           chan SetFilenameFormattingResponse
}

// NewSetFilenameFormattingRequest returns a new SetFilenameFormattingRequest.
func NewSetFilenameFormattingRequest(filenameFormatting string) SetFilenameFormattingRequest {
	return SetFilenameFormattingRequest{
		filenameFormatting,
		_request{
			ID_:   getMessageID(),
			Type_: "SetFilenameFormatting",
			err:   make(chan error, 1),
		},
		make(chan SetFilenameFormattingResponse, 1),
	}
}

// Send sends the request.
func (r *SetFilenameFormattingRequest) Send(c Client) error {
	if r.sent {
		return ErrAlreadySent
	}
	future, err := c.sendRequest(r)
	if err != nil {
		return err
	}
	r.sent = true
	go func() {
		m := <-future
		var resp SetFilenameFormattingResponse
		if err = mapToStruct(m, &resp); err != nil {
			r.err <- err
		} else if resp.Status() != StatusOK {
			r.err <- errors.New(resp.Error())
		} else {
			r.response <- resp
		}
	}()
	return nil
}

// Receive waits for the response.
func (r SetFilenameFormattingRequest) Receive() (SetFilenameFormattingResponse, error) {
	if !r.sent {
		return SetFilenameFormattingResponse{}, ErrNotSent
	}
	if receiveTimeout == 0 {
		select {
		case resp := <-r.response:
			return resp, nil
		case err := <-r.err:
			return SetFilenameFormattingResponse{}, err
		}
	} else {
		select {
		case resp := <-r.response:
			return resp, nil
		case err := <-r.err:
			return SetFilenameFormattingResponse{}, err
		case <-time.After(receiveTimeout):
			return SetFilenameFormattingResponse{}, ErrReceiveTimeout
		}
	}
}

// SendReceive sends the request then immediately waits for the response.
func (r SetFilenameFormattingRequest) SendReceive(c Client) (SetFilenameFormattingResponse, error) {
	if err := r.Send(c); err != nil {
		return SetFilenameFormattingResponse{}, err
	}
	return r.Receive()
}

// SetFilenameFormattingResponse : Response for SetFilenameFormattingRequest.
//
// Since obs-websocket version: 4.3.0.
//
// https://github.com/Palakis/obs-websocket/blob/4.3-maintenance/docs/generated/protocol.md#setfilenameformatting
type SetFilenameFormattingResponse struct {
	_response `json:",squash"`
}

// GetFilenameFormattingRequest : Get the filename formatting string.
//
// Since obs-websocket version: 4.3.0.
//
// https://github.com/Palakis/obs-websocket/blob/4.3-maintenance/docs/generated/protocol.md#getfilenameformatting
type GetFilenameFormattingRequest struct {
	_request `json:",squash"`
	response chan GetFilenameFormattingResponse
}

// NewGetFilenameFormattingRequest returns a new GetFilenameFormattingRequest.
func NewGetFilenameFormattingRequest() GetFilenameFormattingRequest {
	return GetFilenameFormattingRequest{
		_request{
			ID_:   getMessageID(),
			Type_: "GetFilenameFormatting",
			err:   make(chan error, 1),
		},
		make(chan GetFilenameFormattingResponse, 1),
	}
}

// Send sends the request.
func (r *GetFilenameFormattingRequest) Send(c Client) error {
	if r.sent {
		return ErrAlreadySent
	}
	future, err := c.sendRequest(r)
	if err != nil {
		return err
	}
	r.sent = true
	go func() {
		m := <-future
		var resp GetFilenameFormattingResponse
		if err = mapToStruct(m, &resp); err != nil {
			r.err <- err
		} else if resp.Status() != StatusOK {
			r.err <- errors.New(resp.Error())
		} else {
			r.response <- resp
		}
	}()
	return nil
}

// Receive waits for the response.
func (r GetFilenameFormattingRequest) Receive() (GetFilenameFormattingResponse, error) {
	if !r.sent {
		return GetFilenameFormattingResponse{}, ErrNotSent
	}
	if receiveTimeout == 0 {
		select {
		case resp := <-r.response:
			return resp, nil
		case err := <-r.err:
			return GetFilenameFormattingResponse{}, err
		}
	} else {
		select {
		case resp := <-r.response:
			return resp, nil
		case err := <-r.err:
			return GetFilenameFormattingResponse{}, err
		case <-time.After(receiveTimeout):
			return GetFilenameFormattingResponse{}, ErrReceiveTimeout
		}
	}
}

// SendReceive sends the request then immediately waits for the response.
func (r GetFilenameFormattingRequest) SendReceive(c Client) (GetFilenameFormattingResponse, error) {
	if err := r.Send(c); err != nil {
		return GetFilenameFormattingResponse{}, err
	}
	return r.Receive()
}

// GetFilenameFormattingResponse : Response for GetFilenameFormattingRequest.
//
// Since obs-websocket version: 4.3.0.
//
// https://github.com/Palakis/obs-websocket/blob/4.3-maintenance/docs/generated/protocol.md#getfilenameformatting
type GetFilenameFormattingResponse struct {
	// Current filename formatting string.
	// Required: Yes.
	FilenameFormatting string `json:"filename-formatting"`
	_response          `json:",squash"`
}