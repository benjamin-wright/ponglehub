package rules

import (
	"fmt"

	"ponglehub.co.uk/games/naughts-and-crosses/pkg/database"
)

func Validate(game *database.Game, marks string, userId string, position int) *RuleFail {
	markRunes := []rune(marks)

	if isFinished(game) {
		return &RuleFail{
			response: "already finished",
			log:      fmt.Sprintf("user %s tried to make a mark in game %s, but it was already finished", userId, game.ID),
		}
	}

	if isBadTurn(game) {
		return &RuleFail{
			response: "illegal turn",
			log:      fmt.Sprintf("user %s tried to make a mark in game %s, but the turn value was illegal (%d)", userId, game.ID, position),
		}
	}

	if isWrongPlayer(game, userId) {
		return &RuleFail{
			response: "not your turn",
			log:      fmt.Sprintf("user %s tried to make a mark in game %s, but it wasn't their turn", userId, game.ID),
		}
	}

	if isIllegalPosition(position) {
		return &RuleFail{
			response: "illegal position",
			log:      fmt.Sprintf("user %s tried to make a mark in game %s, but it was at an illegal position (%d)", userId, game.ID, position),
		}
	}

	if isAlreadyPlayed(markRunes, position) {
		return &RuleFail{
			response: "already played",
			log:      fmt.Sprintf("user %s tried to make a mark in game %s, but it was already marked", userId, game.ID),
		}
	}

	return nil
}

func isFinished(game *database.Game) bool {
	return game.Finished
}

func isBadTurn(game *database.Game) bool {
	return game.Turn < 0 || game.Turn > 1
}

func isWrongPlayer(game *database.Game, userId string) bool {
	if game.Turn == 0 && game.Player1.String() != userId {
		return true
	}

	if game.Turn == 1 && game.Player2.String() != userId {
		return true
	}

	return false
}

func isIllegalPosition(position int) bool {
	return position < 0 || position > 8
}

func isAlreadyPlayed(marks []rune, position int) bool {
	return marks[position] != '-'
}
