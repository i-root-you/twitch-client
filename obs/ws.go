package obs

import (
	"log"
	"time"

	"github.com/i-root-you/obs/ws"
)

func InitOBSClient() (c obsws.Client) {
	c = ws.Client{Host: "localhost", Port: 4444}
	if err := c.Connect(); err != nil {
		log.Fatal(err)
	}
	defer c.Disconnect()
	//' TODO: Is this global? Ugh
	ws.SetReceiveTimeout(time.Second * 2)
	return c
}

// TODO: Not the best name, and not perm but will do for now
// by just mimicing the top func 
// TODO: Return req, then we can have a bit more DRY codde
func GetStatusRequest() {
	// Send and receive a request asynchronously.
	req := ws.NewGetStreamingStatusRequest()
	if err := req.Send(c); err != nil {
		log.Fatal(err)
	}
}

// TODO: This will need to be either passed 'req' or we will need to more
// likely make this and probably the above function a method of the request
// struct.
func Recieve(req ws.Request) (resp ws.Response, err error) {
	// This will block until the response comes (potentially forever).
	// TODO: lol whjy? Seems like we should catch received requests
	// using channels and switch case. 
	return req.Receive()
}

	// Set the amount of time we can wait for a response.
	// TODO: Why are we short polling every 2 seconds? Thats not 
	// very modern. 

	// TODO: We did this twice...
	// Send and receive a request synchronously.
	//req = obsws.NewGetStreamingStatusRequest()


	// Note that we create a new request,
	// because requests have IDs that must be unique.
	// This will block for up to two seconds, since we set a timeout.

// TODO: I really hate this variable name, it makes no sense, you either send or
// receive, or you do neither (as in there is a more descriptive way of
// describing the action of the function). Obviously its redundant because its
// just calling the exact same method name on the 'Request' object and we shou8d
// probably be fixing that function and just deleting this below. 
func SendReceive() {
	resp, err = req.SendReceive(c)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("streaming:", resp.Streaming)


	// TODO:  Does seem to be some event driven programming here at least but
	// seems like this should obviously not really be in the init? 
	// Respond to events by registering handlers.

	c.AddEventHandler("SwitchScenes", func(e obsws.Event) {
		// Make sure to assert the actual event type.
		log.Println("new scene:", e.(obsws.SwitchScenesEvent).SceneName)
	})

	// TODO: What?  Are we just trying to use up processor time for fun?
	//time.Sleep(time.Second * 10)
	// If we are trying to hold it open, a for {} or select {} or something 
	// else would be better than just hoping 10 seconds is right. Thats just
	// begging for race conditions.
}
