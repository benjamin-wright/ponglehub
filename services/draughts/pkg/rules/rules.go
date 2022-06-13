package rules

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"ponglehub.co.uk/games/draughts/pkg/database"
)

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

func IsYourTurn(player string, game database.Game) bool {
	player1 := player == game.Player1.String() && game.Turn == 0
	player2 := player == game.Player2.String() && game.Turn == 1

	return player1 || player2
}

type Result struct {
	Piece    uuid.UUID
	NewX     int16
	NewY     int16
	ToRemove []uuid.UUID
	King     bool
}

func Process(moves []Move, pieces []database.Piece) (Result, error) {
	result := Result{}

	if len(moves) < 1 {
		return Result{}, errors.New("no moves")
	}

	piece, found := getTargetPiece(moves[0], pieces)
	if !found {
		return Result{}, fmt.Errorf("invalid move, couldn't find piece: %s", moves[0].Piece)
	}

	capturing := false

	for idx, move := range moves {
		if move.Piece != piece.ID {
			return Result{}, fmt.Errorf("more than one piece moved (%s and %s)", move.Piece, piece.ID)
		}

		if isTraversal(move, piece) {
			if idx > 0 {
				return Result{}, errors.New("can't move piece more than once")
			}

			if isBlocked(move, pieces) {
				return Result{}, fmt.Errorf("can't move piece, space already taken: %+v", move)
			}

			result.Piece = move.Piece
			result.NewX = move.X
			result.NewY = move.Y

			if isKing(move, piece) {
				result.King = true
				piece.King = true
			}

			return result, nil
		} else if isCapture(move, piece) {
			capturing = true

			if idx > 0 && !capturing {
				return Result{}, fmt.Errorf("can't capture after a traversal")
			}

			if isBlocked(move, pieces) {
				return Result{}, fmt.Errorf("can't move piece, space already taken: %+v", move)
			}

			captured, found := findCaptured(move, piece, pieces)
			if !found {
				return Result{}, fmt.Errorf("can't capture piece, no piece to capture: %+v", move)
			}

			if captured.Player == piece.Player {
				return Result{}, fmt.Errorf("can't capture piece, it belongs to the current player: %+v", move)
			}

			result.Piece = move.Piece
			result.NewX = move.X
			result.NewY = move.Y
			result.ToRemove = append(result.ToRemove, captured.ID)

			if isKing(move, piece) {
				result.King = true
				piece.King = true
			}
		} else {
			return Result{}, fmt.Errorf("move %+v not recognised", move)
		}
	}

	return result, nil
}

func getTargetPiece(move Move, pieces []database.Piece) (database.Piece, bool) {
	for _, piece := range pieces {
		if piece.ID == move.Piece {
			return piece, true
		}
	}

	return database.Piece{}, false
}

func isTraversal(move Move, piece database.Piece) bool {
	dy := move.Y - piece.Y
	dx := move.X - piece.X

	if piece.King {
		return abs(dy) == 1 && abs(dx) == 1
	}

	switch piece.Player {
	case 0:
		return dy == 1 && abs(dx) == 1
	case 1:
		return dy == -1 && abs(dx) == 1
	default:
		panic(fmt.Sprintf("piece player should be 0 or 1, got %d", piece.Player))
	}
}

func isCapture(move Move, piece database.Piece) bool {
	dy := move.Y - piece.Y
	dx := move.X - piece.X

	if piece.King {
		return abs(dy) == 2 && abs(dx) == 2
	}

	switch piece.Player {
	case 0:
		return dy == 2 && abs(dx) == 2
	case 1:
		return dy == -2 && abs(dx) == 2
	default:
		panic(fmt.Sprintf("piece player should be 0 or 1, got %d", piece.Player))
	}
}

func isBlocked(move Move, pieces []database.Piece) bool {
	for _, piece := range pieces {
		if piece.X == move.X && piece.Y == move.Y {
			return true
		}
	}

	return false
}

func isKing(move Move, piece database.Piece) bool {
	return piece.Player == 0 && move.Y == 7 || piece.Player == 1 && move.Y == 0
}

func findCaptured(move Move, piece database.Piece, pieces []database.Piece) (database.Piece, bool) {
	x := (move.X + piece.X) / 2
	y := (move.Y + piece.Y) / 2

	for _, p := range pieces {
		if p.X == x && p.Y == y {
			return p, true
		}
	}

	return database.Piece{}, false
}
