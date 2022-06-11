package rules

import (
	"fmt"

	"github.com/google/uuid"
	"ponglehub.co.uk/games/draughts/pkg/database"
)

func runeToPiece(char rune) (bool, int16, bool) {
	switch char {
	case 'o':
		return true, 0, false
	case 'O':
		return true, 0, true
	case 'x':
		return true, 1, false
	case 'X':
		return true, 1, true
	default:
		return false, 0, false
	}
}

func ToPieces(gameId uuid.UUID, lines []string) []database.Piece {
	pieces := []database.Piece{}

	for y, line := range lines {
		y = 7 - y

		for x, char := range line {
			if ok, player, king := runeToPiece(char); ok {
				pieces = append(pieces, database.Piece{
					Game:   gameId,
					X:      int16(x),
					Y:      int16(y),
					Player: player,
					King:   king,
				})
			}
		}
	}

	return pieces
}

func pieceToRune(piece database.Piece) rune {
	switch {
	case piece.King && piece.Player == 0:
		return 'O'
	case piece.King && piece.Player == 1:
		return 'X'
	case piece.Player == 0:
		return 'o'
	case piece.Player == 1:
		return 'x'
	default:
		panic(fmt.Sprintf("Illegal piece: %+v", piece))
	}
}

func ToString(pieces []database.Piece) []string {
	lines := []string{}

	for y := int16(0); y < 8; y++ {
		line := []rune{}

		for x := int16(0); x < 8; x++ {
			empty := true
			for _, piece := range pieces {
				if piece.X == x && piece.Y == 7-y {
					empty = false
					line = append(line, pieceToRune(piece))
					break
				}
			}

			if empty {
				line = append(line, ' ')
			}
		}

		lines = append(lines, string(line))
	}

	return lines
}

func NewGame(gameId uuid.UUID) []database.Piece {
	return ToPieces(gameId, []string{
		" x x x x",
		"x x x x ",
		" x x x x",
		"        ",
		"        ",
		"o o o o ",
		" o o o o",
		"o o o o ",
	})
}
