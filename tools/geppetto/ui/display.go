package ui

import (
	"github.com/gdamore/tcell/v2"
	"ponglehub.co.uk/geppetto/types"
)

type display struct {
	screen tcell.Screen
}

func (d *display) start(progress <-chan []types.RepoState, commands <-chan command) {
	for {
		d.screen.Clear()

		d.screen.
	}
}
func (d *display) stop()                                                            {}
