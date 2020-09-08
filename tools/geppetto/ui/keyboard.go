package ui

import "github.com/gdamore/tcell/v2"

type watchInput struct {
	screen tcell.Screen
}

func (w *watchInput) start() <-chan watchCommand { return nil }
func (w *watchInput) stop()                      {}
