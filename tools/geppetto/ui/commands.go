package ui

type command uint8

const (
	redrawCommand     command = 0
	quitCommand       command = 1
	upCommand         command = 2
	downCommand       command = 3
	selectCommand     command = 4
	rebuildCommand    command = 5
	rebuildAllCommand command = 6
	reinstallCommand  command = 7
)
