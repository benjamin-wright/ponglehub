package integration

import (
	"encoding/json"
	"fmt"
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
	turn     float64
	marks    string
	finished bool
}

type play struct {
	initial  string
	turn     int16
	user     uuid.UUID
	position int
	finished bool
}

func makePlay(initial string, turn int16, user uuid.UUID, position int, finished bool) play {
	return play{
		initial,
		turn,
		user,
		position,
		finished,
	}
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
		play     play
		expected markResponse
	}{
		{
			name: "player 0",
			play: makePlay("---------", 0, userId, 0, false),
			expected: markResponse{
				turn:  1,
				marks: "0--------",
			},
		},
		{
			name: "player 1",
			play: makePlay("---------", 1, opponentId, 0, false),
			expected: markResponse{
				turn:  0,
				marks: "1--------",
			},
		},
		{
			name: "player 0 non-zero",
			play: makePlay("---------", 0, userId, 3, false),
			expected: markResponse{
				turn:  1,
				marks: "---0-----",
			},
		},
		{
			name: "player 1 non-zero",
			play: makePlay("---------", 1, opponentId, 3, false),
			expected: markResponse{
				turn:  0,
				marks: "---1-----",
			},
		},
		{
			name: "messy",
			play: makePlay("0110-0-01", 1, opponentId, 4, false),
			expected: markResponse{
				turn:  0,
				marks: "011010-01",
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
					Turn:    test.play.turn,
				},
				test.play.initial,
			)
			noErr(u, err)

			err = eventClient.Send(
				"naughts-and-crosses.mark",
				map[string]interface{}{
					"game":     gameId.String(),
					"position": test.play.position,
				},
				map[string]interface{}{"userid": test.play.user.String()},
			)
			noErr(u, err)

			event := recorder.WaitForEvent(u, os.Getenv("RECORDER_URL"), "naughts-and-crosses.mark.response")
			var actual map[string]interface{}
			err = json.Unmarshal([]byte(event), &actual)
			noErr(u, err)

			expected := map[string]interface{}{
				"game": map[string]interface{}{
					"ID":       gameId.String(),
					"Player1":  userId.String(),
					"Player2":  opponentId.String(),
					"Created":  created.Format("2006-01-02T15:04:05.999999Z"),
					"Turn":     test.expected.turn,
					"Finished": test.expected.finished,
				},
				"marks": test.expected.marks,
			}

			assert.Equal(u, expected, actual)
		})
	}

	for _, test := range []struct {
		name string
		play play
		err  string
	}{
		{
			name: "wrong player",
			play: makePlay("---------", 0, opponentId, 0, false),
			err:  "not your turn",
		},
		{
			name: "replay",
			play: makePlay("--0------", 0, userId, 2, false),
			err:  "already played",
		},
		{
			name: "overplay",
			play: makePlay("-----0---", 1, opponentId, 5, false),
			err:  "already played",
		},
		{
			name: "too low",
			play: makePlay("-----0---", 1, opponentId, -1, false),
			err:  "illegal position",
		},
		{
			name: "too high",
			play: makePlay("-----0---", 1, opponentId, 9, false),
			err:  "illegal position",
		},
		{
			name: "finished",
			play: makePlay("-----0---", 1, opponentId, 0, true),
			err:  "already finished",
		},
	} {
		t.Run(test.name, func(u *testing.T) {
			recorder.Clear(u, os.Getenv("RECORDER_URL"))
			noErr(u, db.Clear())

			err := db.InsertGame(
				database.Game{
					ID:       gameId,
					Player1:  userId,
					Player2:  opponentId,
					Created:  created,
					Turn:     test.play.turn,
					Finished: test.play.finished,
				},
				test.play.initial,
			)
			noErr(u, err)

			err = eventClient.Send(
				"naughts-and-crosses.mark",
				map[string]interface{}{
					"game":     gameId.String(),
					"position": test.play.position,
				},
				map[string]interface{}{"userid": test.play.user.String()},
			)
			noErr(u, err)

			event := recorder.WaitForEvent(u, os.Getenv("RECORDER_URL"), "naughts-and-crosses.mark.rejection.response")
			var actual map[string]interface{}
			err = json.Unmarshal([]byte(event), &actual)
			noErr(u, err)

			assert.Equal(u, map[string]interface{}{
				"reason": test.err,
			}, actual)
		})
	}

	for index, test := range []struct {
		play     play
		expected string
	}{
		{play: makePlay("0--0-----", 0, userId, 6, false), expected: "0--0--0--"},
		{play: makePlay("0-----0--", 0, userId, 3, false), expected: "0--0--0--"},
		{play: makePlay("---0--0--", 0, userId, 0, false), expected: "0--0--0--"},
		{play: makePlay("0-101----", 0, userId, 6, false), expected: "0-101-0--"},
		{play: makePlay("0-1-1-0--", 0, userId, 3, false), expected: "0-101-0--"},
		{play: makePlay("--101-0--", 0, userId, 0, false), expected: "0-101-0--"},
		{play: makePlay("-0--0----", 0, userId, 7, false), expected: "-0--0--0-"},
		{play: makePlay("-0-----0-", 0, userId, 4, false), expected: "-0--0--0-"},
		{play: makePlay("----0--0-", 0, userId, 1, false), expected: "-0--0--0-"},
		{play: makePlay("--0--0---", 0, userId, 8, false), expected: "--0--0--0"},
		{play: makePlay("--0-----0", 0, userId, 5, false), expected: "--0--0--0"},
		{play: makePlay("-----0--0", 0, userId, 2, false), expected: "--0--0--0"},
		{play: makePlay("-00------", 0, userId, 0, false), expected: "000------"},
		{play: makePlay("0-0------", 0, userId, 1, false), expected: "000------"},
		{play: makePlay("00-------", 0, userId, 2, false), expected: "000------"},
		{play: makePlay("----00---", 0, userId, 3, false), expected: "---000---"},
		{play: makePlay("---0-0---", 0, userId, 4, false), expected: "---000---"},
		{play: makePlay("---00----", 0, userId, 5, false), expected: "---000---"},
		{play: makePlay("-------00", 0, userId, 6, false), expected: "------000"},
		{play: makePlay("------0-0", 0, userId, 7, false), expected: "------000"},
		{play: makePlay("------00-", 0, userId, 8, false), expected: "------000"},
		{play: makePlay("----0---0", 0, userId, 0, false), expected: "0---0---0"},
		{play: makePlay("0-------0", 0, userId, 4, false), expected: "0---0---0"},
		{play: makePlay("0---0----", 0, userId, 8, false), expected: "0---0---0"},
		{play: makePlay("----0-0--", 0, userId, 2, false), expected: "--0-0-0--"},
		{play: makePlay("--0---0--", 0, userId, 4, false), expected: "--0-0-0--"},
		{play: makePlay("--0-0----", 0, userId, 6, false), expected: "--0-0-0--"},
	} {
		t.Run(fmt.Sprintf("win condition %d", index), func(u *testing.T) {
			recorder.Clear(u, os.Getenv("RECORDER_URL"))
			noErr(u, db.Clear())

			err := db.InsertGame(
				database.Game{
					ID:       gameId,
					Player1:  userId,
					Player2:  opponentId,
					Created:  created,
					Turn:     test.play.turn,
					Finished: test.play.finished,
				},
				test.play.initial,
			)
			noErr(u, err)

			err = eventClient.Send(
				"naughts-and-crosses.mark",
				map[string]interface{}{
					"game":     gameId.String(),
					"position": test.play.position,
				},
				map[string]interface{}{"userid": test.play.user.String()},
			)
			noErr(u, err)

			event := recorder.WaitForEvent(u, os.Getenv("RECORDER_URL"), "naughts-and-crosses.mark.response")
			var actual map[string]interface{}
			err = json.Unmarshal([]byte(event), &actual)
			noErr(u, err)

			expected := map[string]interface{}{
				"game": map[string]interface{}{
					"ID":       gameId.String(),
					"Player1":  userId.String(),
					"Player2":  opponentId.String(),
					"Created":  created.Format("2006-01-02T15:04:05.999999Z"),
					"Turn":     float64(test.play.turn),
					"Finished": true,
				},
				"marks": test.expected,
			}

			assert.Equal(u, expected, actual)
		})
	}
}
