package ws

import (
	"testing"
	"time"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type EventSuite struct{}

var _ = Suite(&EventSuite{})

func (s *EventSuite) TestTimecodeParse(c *C) {
	tdata := map[string]time.Duration{
		"00:00:00.000": 0,
		"01:00:00.000": 1 * time.Hour,
		"00:01:00.000": 1 * time.Minute,
		"00:00:01.000": 1 * time.Second,
		"00:00:00.042": 42 * time.Millisecond,
	}

	for tc, expected := range tdata {
		res, err := parseTC(tc)
		if ok := c.Check(err, IsNil); ok == false {
			continue
		}
		c.Check(res, Equals, expected)
	}
}

func (s *EventSuite) TestExtractTimecode(c *C) {
	itdata := map[string]string{
		"obsws: invalid Timecode format '.*':.*": "0:0",
	}

	for errorMatch, tc := range itdata {
		res, err := parseTC(tc)
		c.Check(res, Equals, time.Duration(-1))
		c.Check(err, ErrorMatches, errorMatch)
	}

}

func (s *EventSuite) TestEventExtraction(c *C) {
	tdata := map[string]rawEvent{
		`{"update-type":"ScenesChanged"} `:                                  rawEvent{eventType: "ScenesChanged", streamTC: -1, recTC: -1},
		`{"update-type":"ScenesChanged","stream-timecode":"01:00:00.000"} `: rawEvent{eventType: "ScenesChanged", streamTC: 1 * time.Hour, recTC: -1},
		`{"update-type":"ScenesChanged","rec-timecode":"00:01:00.000"} `:    rawEvent{eventType: "ScenesChanged", recTC: 1 * time.Minute, streamTC: -1},
	}

	for jsonData, expected := range tdata {
		res, err := UnmarshalEvent([]byte(jsonData))
		if c.Check(err, IsNil, Commentf("Unexpected error: %s", err)) == false {
			continue
		}
		c.Check(res.UpdateType(), Equals, expected.UpdateType())
		stc, okstc := res.StreamTimecode()
		estc, eokstc := expected.StreamTimecode()
		c.Check(okstc, Equals, eokstc)
		c.Check(stc, Equals, estc)
		rtc, okrtc := res.RecordTimecode()
		ertc, eokrtc := expected.RecordTimecode()
		c.Check(rtc, Equals, ertc)
		c.Check(okrtc, Equals, eokrtc)
	}

	itdata := map[string]string{
		"obsws: message is not an event":           `{"message-id":1234,"status":"ok","error":""}`,
		"json:.*":                                  `{"update-type":42}`,
		"obsws: unknown event type 'foo'":          `{"update-type":"foo"}`,
		"obsws: invalid Timecode format '.*': .*":  `{"update-type":"SwitchScenes","stream-timecode":"a"}`,
		"obsws: invalid Timecode format 'a.*': .*": `{"update-type":"SwitchScenes","rec-timecode":"a"}`,
	}

	for errorMatch, jsonData := range itdata {
		res, err := UnmarshalEvent([]byte(jsonData))
		c.Check(err, ErrorMatches, errorMatch, Commentf(jsonData))
		c.Check(res, IsNil)
	}
}

func (s *EventSuite) TestEventFactory(c *C) {

	tdata := map[Event]string{
		&EventSwitchScenes{
			rawEvent:  rawEvent{"SwitchScenes", -1, -1},
			SceneName: "foo",
		}: `{"update-type":"SwitchScenes","scene-name":"foo"}`,
		&EventStreamStatus{
			rawEvent:         rawEvent{"StreamStatus", -1, -1},
			Fps:              29.97,
			Streaming:        true,
			Recording:        false,
			PreviewOnly:      false,
			BytesPerSec:      1234,
			KBitsPerSec:      1,
			TotalStreamTime:  122,
			Strain:           0.001,
			NumTotalFrames:   200,
			NumDroppedFrames: 1,
		}: `{"update-type":"StreamStatus","fps":29.97,"streaming":true,"bytes-per-sec":1234,"kbits-per-sec":1,"preview-only":false,"strain":0.001,"total-stream-time":122,"num-total-frames":200,"num-dropped-frames":1}`,
	}

	for expected, jsonData := range tdata {
		res, err := UnmarshalEvent([]byte(jsonData))
		if c.Check(err, IsNil, Commentf("Unexpected error: %s", err)) == false {
			continue
		}
		c.Check(res, DeepEquals, expected)
	}

}
