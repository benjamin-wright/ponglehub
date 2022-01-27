package rules

import (
	"fmt"

	"ponglehub.co.uk/games/naughts-and-crosses/pkg/database"
)

func ValidateMark(state *database.GameState, player string, position int) error {
	if player != state.CurrentPlayer() {
		return fmt.Errorf("wrong player: expected %s, got %s", player, state.CurrentPlayer())
	}

	if position < 0 || position > 9 {
		return fmt.Errorf("bad position input: %d", position)
	}

	for _, existing := range state.Player1Marks {
		if existing == position {
			return fmt.Errorf("position %d already taken by player 1", position)
		}
	}

	for _, existing := range state.Player2Marks {
		if existing == position {
			return fmt.Errorf("position %d already taken by player 2", position)
		}
	}

	return nil
}
