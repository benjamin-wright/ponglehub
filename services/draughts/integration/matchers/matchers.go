package matchers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"ponglehub.co.uk/games/draughts/pkg/database"
	"ponglehub.co.uk/games/draughts/pkg/rules"
)

func AssertEqualGames(t *testing.T, expected database.Game, actual database.Game) {
	assert.Equal(t, expected.Player1, actual.Player1)
	assert.Equal(t, expected.Player2, actual.Player2)
	assert.Equal(t, expected.Turn, actual.Turn)
	assert.Equal(t, expected.Finished, actual.Finished)

	timeDifference := actual.CreatedTime.UTC().Sub(expected.CreatedTime.UTC())
	if timeDifference < 0 {
		timeDifference = -timeDifference
	}

	assert.Less(t, timeDifference, 5*time.Second)
}

func AssertEqualPieces(t *testing.T, expected []string, actual []database.Piece) {
	assert.Equal(t, expected, rules.ToString(actual))
}
