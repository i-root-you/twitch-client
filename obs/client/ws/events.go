package ws

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

type rawEvent struct {
	eventType string
	streamTC  time.Duration
	recTC     time.Duration
}

// Event can be returned from the client at any time
type Event interface {
	// UpdateType specifies which kind of event it is
	UpdateType() string
	// StreamTimecode specifies the time since the stream started. It
	// return false if the stream did not yet started.
	StreamTimecode() (time.Duration, bool)
	// RecordTimecode specifies the time since the recording
	// started. It return false of the stream did not yet stated/
	RecordTimecode() (time.Duration, bool)

	copyFromOther(e Event)
}

func (e *rawEvent) UpdateType() string {
	return e.eventType
}

func (e *rawEvent) StreamTimecode() (time.Duration, bool) {
	if e.streamTC < 0 {
		return 0, false
	}
	return e.streamTC, true
}

func (e *rawEvent) RecordTimecode() (time.Duration, bool) {
	if e.recTC < 0 {
		return 0, false
	}
	return e.recTC, true
}

func (e *rawEvent) copyFromOther(other Event) {
	e.eventType = other.UpdateType()
	var ok bool
	if e.streamTC, ok = other.StreamTimecode(); ok == false {
		e.streamTC = -1
	}

	if e.recTC, ok = other.RecordTimecode(); ok == false {
		e.recTC = -1
	}
}

func parseTC(tc string) (time.Duration, error) {
	if len(tc) == 0 {
		return -1, fmt.Errorf("obsws: invalid Timecode format '': empty string")
	}
	var h, m, s, ms int
	_, err := fmt.Sscanf(tc, "%d:%d:%d.%d", &h, &m, &s, &ms)
	if err != nil {
		return -1, fmt.Errorf("obsws: invalid Timecode format '%s': %s", tc, err)
	}

	return time.Duration(h)*time.Hour + time.Duration(m)*time.Minute + time.Duration(s)*time.Second + time.Duration(ms)*time.Millisecond, nil
}

// ErrNotEventMessage is an error returned by UnmarshalEvent when the
// JSON data does not correspond to an Event message.
type ErrNotEventMessage struct{}

func (e ErrNotEventMessage) Error() string {
	return "obsws: message is not an event"
}

type ErrUnknownEventType struct {
	Type string
}

func (e ErrUnknownEventType) Error() string {
	return "obsws: unknown event type '" + e.Type + "'"
}

// UnmarshalEvent unmarshal an Event formatted in JSON. Cannot use the
// standard JSON interface as it would be too complicated to implement
func UnmarshalEvent(data []byte) (Event, error) {
	//first we extract the generic part
	aux := struct {
		StreamTCStr string `json:"stream-timecode"`
		RecTCStr    string `json:"rec-timecode"`
		UpdateType  string `json:"update-type"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return nil, err
	}

	if len(aux.UpdateType) == 0 {
		return nil, ErrNotEventMessage{}
	}
	rawE := &rawEvent{
		eventType: aux.UpdateType,
		recTC:     -1,
		streamTC:  -1,
	}

	if len(aux.StreamTCStr) > 0 {
		var err error
		rawE.streamTC, err = parseTC(aux.StreamTCStr)
		if err != nil {
			return nil, err
		}
	}

	if len(aux.RecTCStr) > 0 {
		var err error
		rawE.recTC, err = parseTC(aux.RecTCStr)
		if err != nil {
			return nil, err
		}
	}

	// now we extract the generic part
	evType, ok := eventFactory[rawE.UpdateType()]
	if ok == false {
		return nil, ErrUnknownEventType{rawE.UpdateType()}
	}

	evInst := reflect.New(evType)
	ev := evInst.Interface().(Event)

	err := json.Unmarshal(data, &ev)
	if err != nil {
		return nil, err
	}
	// copy back initial data
	ev.copyFromOther(rawE)
	return ev, nil
}

var eventFactory map[string]reflect.Type

type EventSwitchScenes struct {
	SceneName string `json:"scene-name"`
	rawEvent
}

type EventScenesChanged struct {
	rawEvent
}

type EventSourceOrderChanged struct {
	SceneName string `json:"scene-name"`
	rawEvent
}

type EventSceneItemAdded struct {
	SceneName string `json:"scene-name"`
	ItemName  string `json:"item-name"`
	rawEvent
}

type EventSceneItemRemoved struct {
	SceneName string `json:"scene-name"`
	ItemName  string `json:"item-name"`
	rawEvent
}

type EventStreamStatus struct {
	Streaming        bool    `json:"streaming"`
	Recording        bool    `json:"recording"`
	PreviewOnly      bool    `json:"preview-only"`
	BytesPerSec      int     `json:"bytes-per-sec"`
	KBitsPerSec      int     `json:"kbits-per-sec"`
	Strain           float64 `json:"strain"`
	TotalStreamTime  int     `json:"total-stream-time"`
	NumTotalFrames   int     `json:"num-total-frames"`
	NumDroppedFrames int     `json:"num-dropped-frames"`
	Fps              float64 `json:"fps"`
	rawEvent
}

func init() {
	eventFactory = map[string]reflect.Type{
		"SwitchScenes":       reflect.TypeOf(EventSwitchScenes{}),
		"ScenesChanged":      reflect.TypeOf(EventScenesChanged{}),
		"SourceOrderChanged": reflect.TypeOf(EventSourceOrderChanged{}),
		"SceneItemAdded":     reflect.TypeOf(EventSceneItemAdded{}),
		"SceneItemRemoved":   reflect.TypeOf(EventSceneItemRemoved{}),
		"StreamStatus":       reflect.TypeOf(EventStreamStatus{}),
	}
}
