package ui

import (
	"github.com/gdamore/tcell/v2"
	"ponglehub.co.uk/geppetto/types"
)

type watchUI struct {
	screen tcell.Screen
}

func (w *watchUI) start(progress <-chan []types.RepoState, commands <-chan watchCommand) {}
func (w *watchUI) stop()                                                                 {}
