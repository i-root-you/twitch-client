package ws

import "fmt"

type request interface {
	setMessageID(uid string)
	getResponseChannel() chan response
	setResponseChannel(rchan chan response)
	responseType() response
}

type requestBase struct {
	MessageID   string `json:"message-id"`
	RequestType string `json:"request-type"`
	response    chan response
	rType       response
}

func (r *requestBase) setMessageID(ID string) {
	r.MessageID = ID
}

func (r *requestBase) getResponseChannel() chan response {
	return r.response
}

func (r *requestBase) setResponseChannel(rchan chan response) {
	r.response = rchan
}

func (r *requestBase) responseType() response {
	return r.rType
}

func forgeRequest(name string) request {
	return &requestBase{
		RequestType: name,
		rType:       &responseBase{},
	}
}

func forgeRequestWithExpectedResponse(name string, resp response) request {
	return &requestBase{
		RequestType: name,
		rType:       resp,
	}
}

func forgeSetCurrentScene(name string) request {
	type setCurrentScene struct {
		requestBase
		SceneName string `json:"scene-name"`
	}
	return &setCurrentScene{
		requestBase: requestBase{
			RequestType: "SetCurrentScene",
			rType:       &responseBase{},
		},
		SceneName: name,
	}
}

func (c *Client) submitRequest(r request) (response, error) {
	rchan := make(chan response)
	r.setResponseChannel(rchan)
	c.requests <- r
	resp, ok := <-rchan
	if ok == false {
		return nil, fmt.Errorf("obsws: internal error, response channel closed without response")
	}

	if err := resp.error(); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) GetSceneList() (*GetSceneListResponse, error) {
	resp, err := c.submitRequest(forgeRequestWithExpectedResponse("GetSceneList", &GetSceneListResponse{}))
	if err != nil {
		return nil, err
	}
	respCorrect, ok := resp.(*GetSceneListResponse)
	if ok == false {
		return nil, fmt.Errorf("obsws: unexpected response from server: %#v", resp)
	}
	return respCorrect, nil
}

func (c *Client) SetCurrentScene(name string) error {
	_, err := c.submitRequest(forgeSetCurrentScene(name))
	return err
}
