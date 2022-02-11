package database

import "github.com/google/uuid"

type GameState struct {
	ID           uuid.UUID
	Player1      uuid.UUID
	Player2      uuid.UUID
	Turn         int
	Player1Marks []int
	Player2Marks []int
}

func (g *GameState) CurrentPlayer() string {
	switch g.Turn {
	case 0:
		return g.Player1.String()
	case 1:
		return g.Player2.String()
	}

	return ""
}

func (g *GameState) Opponent() string {
	switch g.Turn {
	case 0:
		return g.Player2.String()
	case 1:
		return g.Player1.String()
	}

	return ""
}

func (g *GameState) SetMark(position int) {
	switch g.Turn {
	case 0:
		g.Player1Marks = append(g.Player1Marks, position)
	case 1:
		g.Player2Marks = append(g.Player2Marks, position)
	}
}
