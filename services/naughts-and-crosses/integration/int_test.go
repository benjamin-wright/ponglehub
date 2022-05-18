package integration

import (
	"encoding/json"
	"io"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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

func getExpected(t *testing.T, games []database.Game) string {
	expected := map[string]interface{}{
		"games": games,
	}

	data, err := json.Marshal(expected)
	if err != nil {
		assert.NoError(t, err)
	}

	return string(data)
}

func TestListGames(t *testing.T) {
	logrus.SetOutput(io.Discard)

	db, eventClient := initClients(t)
	userId := uuid.New()
	otherPlayer := uuid.New()
	randomPlayer := uuid.New()

	for _, test := range []struct {
		name     string
		existing []database.Game
		expected string
	}{
		{
			name:     "empty",
			expected: getExpected(t, []database.Game{}),
		},
		{
			name: "one game",
			existing: []database.Game{
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000000000"), Player1: userId, Player2: otherPlayer, Turn: 0},
			},
			expected: getExpected(t, []database.Game{
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000000000"), Player1: userId, Player2: otherPlayer, Turn: 0},
			}),
		},
		{
			name: "ignore other users games",
			existing: []database.Game{
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000000000"), Player1: userId, Player2: otherPlayer, Turn: 0},
				{ID: uuid.MustParse("10000000-0000-0000-0000-000000000000"), Player1: otherPlayer, Player2: userId, Turn: 0},
				{ID: uuid.MustParse("20000000-0000-0000-0000-000000000000"), Player1: randomPlayer, Player2: otherPlayer, Turn: 0},
			},
			expected: getExpected(t, []database.Game{
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000000000"), Player1: userId, Player2: otherPlayer, Turn: 0},
				{ID: uuid.MustParse("10000000-0000-0000-0000-000000000000"), Player1: otherPlayer, Player2: userId, Turn: 0},
			}),
		},
	} {
		t.Run(test.name, func(u *testing.T) {
			recorder.Clear(u, os.Getenv("RECORDER_URL"))
			noErr(u, db.Clear())

			for _, game := range test.existing {
				err := db.InsertGame(game, "---------")
				noErr(u, err)
			}

			err := eventClient.Send(
				"naughts-and-crosses.list-games",
				nil,
				map[string]interface{}{"userid": userId.String()},
			)
			noErr(u, err)

			event := recorder.WaitForEvent(u, os.Getenv("RECORDER_URL"), "naughts-and-crosses.list-games.response")

			assert.Equal(u, test.expected, event)
		})
	}
}

func TestNewGameEvent(t *testing.T) {
	logrus.SetOutput(io.Discard)

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

type markResponse struct {
	turn  float64
	marks string
}

func TestMarkEvent(t *testing.T) {
	logrus.SetOutput(io.Discard)

	db, eventClient := initClients(t)
	gameId := uuid.New()
	userId := uuid.New()
	opponentId := uuid.New()
	created := time.Now()
	// randomPlayer := uuid.New()

	for _, test := range []struct {
		name     string
		initial  string
		position int
		expected markResponse
	}{
		{
			name:     "success",
			initial:  "---------",
			position: 0,
			expected: markResponse{
				turn:  1,
				marks: "0--------",
			},
		},
	} {
		t.Run(test.name, func(u *testing.T) {
			recorder.Clear(u, os.Getenv("RECORDER_URL"))
			noErr(u, db.Clear())

			err := db.InsertGame(
				database.Game{
					ID:      gameId,
					Player1: userId,
					Player2: opponentId,
					Created: created,
					Turn:    0,
				},
				test.initial,
			)
			noErr(u, err)

			err = eventClient.Send(
				"naughts-and-crosses.mark",
				map[string]interface{}{
					"game":     gameId.String(),
					"position": test.position,
				},
				map[string]interface{}{"userid": userId.String()},
			)
			noErr(u, err)

			event := recorder.WaitForEvent(u, os.Getenv("RECORDER_URL"), "naughts-and-crosses.mark.response")
			var actual map[string]interface{}
			err = json.Unmarshal([]byte(event), &actual)
			noErr(u, err)

			expected := map[string]interface{}{
				"game": map[string]interface{}{
					"ID":      gameId.String(),
					"Player1": userId.String(),
					"Player2": opponentId.String(),
					"Created": created.Format("2006-01-02T15:04:05.999999Z"),
					"Turn":    test.expected.turn,
				},
				"marks": test.expected.marks,
			}

			assert.Equal(u, expected, actual)
		})
	}

}
