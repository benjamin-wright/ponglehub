package rules

import "github.com/google/uuid"

type Move struct {
	Piece uuid.UUID
	X     int16
	Y     int16
}

func (m Move) isLegal() bool {
	return m.X >= 0 &&
		m.X <= 7 &&
		m.Y >= 0 &&
		m.Y <= 7
}
