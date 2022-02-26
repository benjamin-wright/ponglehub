package integration

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"ponglehub.co.uk/events/recorder/pkg/recorder"
	"ponglehub.co.uk/games/naughts-and-crosses/pkg/database"
	"ponglehub.co.uk/lib/events"
)

func noErr(t *testing.T, err error) {
	if err != nil {
		t.Logf("Wasn't expecting an error, got %+v", err)
		t.FailNow()
	}
}

type Game struct {
	Player1 string
	Player2 string
	Turn    int16
}

func assertGames(t *testing.T, expected []Game, actual []database.Game) {
	if len(expected) != len(actual) {
		assert.Fail(t, "Expected %d games, got %d", len(expected), len(actual))
		return
	}

	for index, _ := range expected {
		assert.Equal(t, expected[index].Player1, actual[index].Player1.String())
		assert.Equal(t, expected[index].Player2, actual[index].Player2.String())
		assert.Equal(t, expected[index].Turn, actual[index].Turn)
	}
}

func initClients(t *testing.T) (*database.Database, *events.Events) {
	db, err := database.New()
	noErr(t, err)

	eventClient, err := events.New(events.EventsArgs{
		BrokerEnv: "NAC_URL",
		Source:    "int-test",
	})
	noErr(t, err)

	return db, eventClient
}

func assertJson(t *testing.T, expected map[string]interface{}, actual string) {
	parsed := map[string]interface{}{}
	err := json.Unmarshal([]byte(actual), &parsed)
	noErr(t, err)

	assert.Equal(t, expected, parsed)
}

func TestNewGameEvent(t *testing.T) {
	db, eventClient := initClients(t)

	recorder.Clear(t, os.Getenv("RECORDER_URL"))
	noErr(t, db.Clear())

	userId := uuid.New().String()
	opponentId := uuid.New().String()

	err := eventClient.Send(
		"naughts-and-crosses.new-game",
		map[string]string{"opponent": opponentId},
		map[string]interface{}{"userid": userId},
	)
	noErr(t, err)

	recorder.WaitForEvent(t, os.Getenv("RECORDER_URL"), "naughts-and-crosses.new-game.response")

	games, err := db.ListGames(userId)
	noErr(t, err)

	assertGames(t, []Game{{Player1: opponentId, Player2: userId, Turn: 0}}, games)
}
