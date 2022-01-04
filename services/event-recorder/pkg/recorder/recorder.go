package recorder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func noErr(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Not expecting error: %+v", err)
		t.FailNow()
	}
}

func Clear(t *testing.T, url string) {
	resp, err := http.Post(fmt.Sprintf("%s/clear", url), "application/json", nil)
	noErr(t, err)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Failed to clear events: error code %d", resp.StatusCode)
		t.FailNow()
	}
}

func GetEvents(t *testing.T, url string) []string {
	resp, err := http.Get(fmt.Sprintf("%s/events", url))
	noErr(t, err)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Failed to get events: error code %d", resp.StatusCode)
		t.FailNow()
	}

	body, err := ioutil.ReadAll(resp.Body)
	noErr(t, err)

	events := []string{}
	noErr(t, json.Unmarshal(body, &events))

	return events
}

func WaitForEvents(t *testing.T, url string, number int) []string {
	eventChan := make(chan []string, 1)

	go func(eventChan chan<- []string) {
		events := []string{}
		for len(events) < number {
			time.Sleep(100 * time.Millisecond)
			events = GetEvents(t, url)
		}

		eventChan <- events
	}(eventChan)

	select {
	case events := <-eventChan:
		return events
	case <-time.After(5 * time.Second):
		t.Errorf("timed out waiting for %d events", number)
		t.FailNow()
		return nil
	}
}

type Event struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func GetLatest(t *testing.T, url string) (string, string) {
	resp, err := http.Get(fmt.Sprintf("%s/latest", url))
	noErr(t, err)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Failed to get latest event: error code %d", resp.StatusCode)
		t.FailNow()
	}

	body, err := ioutil.ReadAll(resp.Body)
	noErr(t, err)

	event := Event{}
	noErr(t, json.Unmarshal(body, &event))

	return event.Type, event.Data
}

func WaitForEvent(t *testing.T, url string, eventType string) string {
	eventChan := make(chan Event, 1)

	go func(eventChan chan<- Event) {
		var event Event
		for {
			time.Sleep(100 * time.Millisecond)
			resp, err := http.Get(fmt.Sprintf("%s/latest", url))
			if err != nil {
				continue
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				continue
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				continue
			}
			err = json.Unmarshal(body, &event)
			if err != nil {
				continue
			}

			if event.Type != eventType {
				continue
			}

			break
		}

		eventChan <- event
	}(eventChan)

	select {
	case event := <-eventChan:
		return event.Data
	case <-time.After(5 * time.Second):
		t.Errorf("timed out waiting for an %s events", eventType)
		t.FailNow()
		return ""
	}
}
