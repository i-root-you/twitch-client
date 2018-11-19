package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"golang.org/x/net/websocket"
)

type responseData struct {
	channel chan response
	rType   response
}

// A Client connects to a obs-studio websocket to get event and
// perform request on OBS instance remotely
type Client struct {
	eventChannelLock sync.RWMutex
	wg               sync.WaitGroup

	ws *websocket.Conn

	events       chan Event
	requests     chan request
	frames       chan []byte
	responsesMap map[string]responseData
}

// NewClient connects to a websocket instance.
func NewClient(address string, port int) (*Client, error) {

	ws, err := websocket.Dial(fmt.Sprintf("ws://%s:%d/", address, port),
		"",
		fmt.Sprintf("http://%s:%d/", address, port))

	if err != nil {
		return nil, err
	}

	res := &Client{
		ws:           ws,
		requests:     make(chan request),
		responsesMap: make(map[string]responseData),
	}

	go res.internalLoop()
	return res, nil
}

func (c *Client) handleResponse(frame []byte) {
	//check if the message is an event
	ev, err := UnmarshalEvent(frame)
	if err == nil {
		//check if use is listening events
		c.eventChannelLock.RLock()
		defer c.eventChannelLock.RUnlock()
		if c.events != nil {
			c.events <- ev
		}
		return
	}
	if _, ok := err.(ErrNotEventMessage); ok == false {
		//handle error
		if _, ok := err.(ErrUnknownEventType); ok == true {
			//we only log unknown eventype
			log.Printf("%s", err)
			return
		} else {
			panic(fmt.Sprintf("obsws: %s", err))
		}
	}

	// handle response
	var respBase responseBase
	err = json.Unmarshal(frame, &respBase)
	if err != nil {
		panic(fmt.Sprintf("obsws: %s\n'%s'", err, frame))
	}

	respData, ok := c.responsesMap[respBase.messageID()]
	if ok == false {
		panic(fmt.Sprintf("obsws: unknown message-id '%s'\n'%s'", respBase.messageID(), frame))
	}
	err = json.Unmarshal(frame, &(respData.rType))
	if err != nil {
		panic(fmt.Sprintf("obsws: %s\n'%s'", err, frame))
	}
	respData.channel <- respData.rType
	delete(c.responsesMap, respBase.messageID())
	close(respData.channel)
}

func (c *Client) internalLoop() {
	c.wg.Add(1)
	defer c.wg.Done()

	c.frames = make(chan []byte)

	requestUID := 0

	for {
		c.wg.Add(1)
		defer c.wg.Done()
		go func() {
			// read a response asynchronously
			frame := make([]byte, 0, 100)
			websocket.Message.Receive(c.ws, &frame)
			// send response to channel
			c.frames <- frame
		}()

		select {
		case f := <-c.frames:
			c.handleResponse(f)
		case r, ok := <-c.requests:
			if ok == false {
				c.requests = nil
				break
			}
			// send the right request, with an UID
			requestUID++
			rUID := fmt.Sprintf("%d", requestUID)
			c.responsesMap[rUID] = responseData{
				channel: r.getResponseChannel(),
				rType:   r.responseType(),
			}
			r.setMessageID(rUID)
			websocket.JSON.Send(c.ws, r)
		}

		// we are closing the for loop
		if c.requests == nil {
			break
		}

	}

	//will discard the next response, either error or anything...
	c.frames = nil
}

// Authentify performs the authenfication to this websocket instance.
func (c *Client) Authentify(psswd string) error {
	return NotYetImplemented()
}

// Close terminates the connection to the instance
func (c *Client) Close() {
	close(c.requests)
	// wait to be done
	c.wg.Wait()

	c.eventChannelLock.Lock()
	defer c.eventChannelLock.Unlock()
	if c.events != nil {
		close(c.events)
	}
}

// EventChannel returns a channel to read Event from
func (c *Client) EventChannel() <-chan Event {
	c.eventChannelLock.RLock()
	if c.events != nil {
		defer c.eventChannelLock.RUnlock()
		return c.events
	}

	res := make(chan chan Event)
	go func() {
		defer close(res)
		c.eventChannelLock.Lock()
		defer c.eventChannelLock.Unlock()
		if c.events != nil {
			res <- c.events
		}
		c.events = make(chan Event)
		res <- c.events
	}()
	c.eventChannelLock.RUnlock()
	events := <-res
	return events
}

func (c *Client) eventSender() chan<- Event {
	c.eventChannelLock.RLock()
	defer c.eventChannelLock.RUnlock()
	return c.events
}
