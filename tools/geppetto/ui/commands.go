package ui

type watchCommand uint8

const (
	redrawCommand watchCommand = 0
	quitCommand   watchCommand = 1
	upCommand     watchCommand = 2
	downCommand   watchCommand = 3
	selectCommand watchCommand = 4
)
