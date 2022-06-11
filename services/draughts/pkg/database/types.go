package database

import (
	"time"

	"github.com/google/uuid"
)

type Game struct {
	ID          uuid.UUID `json:"id"`
	Player1     uuid.UUID `json:"player1"`
	Player2     uuid.UUID `json:"player2"`
	Turn        int16     `json:"turn"`
	CreatedTime time.Time `json:"createdTime"`
	Finished    bool      `json:"finished"`
}

type Piece struct {
	ID     uuid.UUID `json:"id"`
	Game   uuid.UUID `json:"game"`
	X      int16     `json:"x"`
	Y      int16     `json:"y"`
	Player int16     `json:"player"`
	King   bool      `json:"king"`
}
