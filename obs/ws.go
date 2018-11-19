package obs

import (
	"log"
	"time"

	"github.com/i-root-you/obsws"
)

func InitHTTPServer() {
	c := obsws.Client{Host: "localhost", Port: 4444}
	if err := c.Connect(); err != nil {
		log.Fatal(err)
	}
	defer c.Disconnect()

	// Send and receive a request asynchronously.
	req := obsws.NewGetStreamingStatusRequest()
	if err := req.Send(c); err != nil {
		log.Fatal(err)
	}
	// This will block until the response comes (potentially forever).
	resp, err := req.Receive()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("streaming:", resp.Streaming)

	// Set the amount of time we can wait for a response.
	// TODO: Why are we short polling every 2 seconds? Thats not 
	// very modern. 
	obsws.SetReceiveTimeout(time.Second * 2)

	// Send and receive a request synchronously.
	req = obsws.NewGetStreamingStatusRequest()
	// Note that we create a new request,
	// because requests have IDs that must be unique.
	// This will block for up to two seconds, since we set a timeout.
	resp, err = req.SendReceive(c)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("streaming:", resp.Streaming)

	// Respond to events by registering handlers.
	c.AddEventHandler("SwitchScenes", func(e obsws.Event) {
		// Make sure to assert the actual event type.
		log.Println("new scene:", e.(obsws.SwitchScenesEvent).SceneName)
	})

	time.Sleep(time.Second * 10)
}
