package integration

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"ponglehub.co.uk/events/recorder/pkg/recorder"
	"ponglehub.co.uk/lib/events"
)

func noErr(t *testing.T, err error) {
	if err != nil {
		t.Logf("Wasn't expecting an error, got %+v", err)
		t.FailNow()
	}
}

func Test_Events(t *testing.T) {
	recorder.Clear(t, os.Getenv("RECORDER_URL"))
	eventClient, err := events.New(events.EventsArgs{
		BrokerEnv: "NAC_URL",
		Source:    "int-test",
	})
	noErr(t, err)

	player1 := uuid.New().String()
	player2 := uuid.New().String()

	err = eventClient.Send("naughts-and-crosses.new-game", map[string]string{
		"player1": player1,
		"player2": player2,
	})
	noErr(t, err)

	recorder.WaitForEvent(t, os.Getenv("RECORDER_URL"), "naughts-and-crosses.new-game-id")
}
