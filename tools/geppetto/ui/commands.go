package ui

type command uint8

const (
	redrawCommand command = 0
	quitCommand   command = 1
	upCommand     command = 2
	downCommand   command = 3
	selectCommand command = 4
	unlockCommand command = 5
)
