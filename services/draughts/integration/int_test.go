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
	"ponglehub.co.uk/games/draughts/pkg/database"
	"ponglehub.co.uk/lib/events"
)

func noErr(t *testing.T, err error) {
	if err != nil {
		t.Logf("Wasn't expecting an error, got %+v", err)
		t.FailNow()
	}
}

func initClients(t *testing.T) (*database.Database, *events.Events) {
	db, err := database.New()
	noErr(t, err)

	eventClient, err := events.New(events.EventsArgs{
		BrokerEnv: "DRAUGHTS_URL",
		Source:    "int-test",
	})
	noErr(t, err)

	return db, eventClient
}

func makeExpected(t *testing.T, games []database.Game) string {
	data, err := json.Marshal(map[string]interface{}{
		"games": games,
	})
	if err != nil {
		t.Logf("failed to parse expected JSON: %+v", err)
		t.FailNow()
	}

	return string(data)
}

func TestListGames(t *testing.T) {
	logrus.SetOutput(io.Discard)

	db, eventClient := initClients(t)
	userId := uuid.New()
	opponentId := uuid.New()

	challengerGame := database.Game{
		ID:          uuid.MustParse("00000000-0000-0000-0000-000000000000"),
		Player1:     userId,
		Player2:     opponentId,
		Turn:        0,
		CreatedTime: time.Unix(5000, 0).UTC(),
		Finished:    false,
	}

	challengerGame2 := database.Game{
		ID:          uuid.MustParse("01000000-0000-0000-0000-000000000000"),
		Player1:     userId,
		Player2:     opponentId,
		Turn:        1,
		CreatedTime: time.Unix(10000, 0).UTC(),
		Finished:    false,
	}

	challengedGame := database.Game{
		ID:          uuid.MustParse("10000000-0000-0000-0000-000000000000"),
		Player1:     opponentId,
		Player2:     userId,
		Turn:        1,
		CreatedTime: time.Unix(7000, 0).UTC(),
		Finished:    true,
	}

	challengedGame2 := database.Game{
		ID:          uuid.MustParse("11000000-0000-0000-0000-000000000000"),
		Player1:     opponentId,
		Player2:     userId,
		Turn:        0,
		CreatedTime: time.Unix(11000, 0).UTC(),
		Finished:    true,
	}

	otherGame := database.Game{
		ID:          uuid.MustParse("20000000-0000-0000-0000-000000000000"),
		Player1:     uuid.New(),
		Player2:     uuid.New(),
		Turn:        0,
		CreatedTime: time.Unix(9000, 0).UTC(),
		Finished:    false,
	}

	for _, test := range []struct {
		name     string
		existing []database.Game
		expected string
	}{
		{
			name:     "empty",
			existing: []database.Game{},
			expected: makeExpected(t, []database.Game{}),
		},
		{
			name:     "player 1",
			existing: []database.Game{challengerGame},
			expected: makeExpected(t, []database.Game{challengerGame}),
		},
		{
			name:     "player 2",
			existing: []database.Game{challengedGame},
			expected: makeExpected(t, []database.Game{challengedGame}),
		},
		{
			name:     "not playing",
			existing: []database.Game{otherGame},
			expected: makeExpected(t, []database.Game{}),
		},
		{
			name:     "mixed",
			existing: []database.Game{challengerGame, challengerGame2, challengedGame, challengedGame2, otherGame},
			expected: makeExpected(t, []database.Game{challengerGame, challengerGame2, challengedGame, challengedGame2}),
		},
	} {
		t.Run(test.name, func(u *testing.T) {
			recorder.Clear(u, os.Getenv("RECORDER_URL"))
			noErr(u, db.Clear())

			for _, game := range test.existing {
				noErr(u, db.InsertGame(game))
			}

			err := eventClient.Send(
				"draughts.list-games",
				nil,
				map[string]interface{}{"userid": userId.String()},
			)
			noErr(u, err)

			event := recorder.WaitForEvent(u, os.Getenv("RECORDER_URL"), "draughts.list-games.response")

			assert.Equal(u, test.expected, event)
		})
	}
}
