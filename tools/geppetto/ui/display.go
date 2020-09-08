package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/geppetto/types"
)

type display struct {
	screen tcell.Screen
	state  []types.RepoState
}

func (d *display) draw() {
	d.screen.Clear()
	width, height := d.screen.Size()
	drawBorder(d.screen, width, height)
	drawContent(d.screen, d.state, width, height)
	d.screen.Show()
}

func drawBorder(screen tcell.Screen, width int, height int) {
	for x := 1; x < width-1; x++ {
		screen.SetContent(x, 0, '-', nil, tcell.StyleDefault)
		screen.SetContent(x, height-1, '-', nil, tcell.StyleDefault)
	}

	for y := 1; y < height-1; y++ {
		screen.SetContent(0, y, '|', nil, tcell.StyleDefault)
		screen.SetContent(width-1, y, '|', nil, tcell.StyleDefault)
	}

	screen.SetContent(0, 0, '+', nil, tcell.StyleDefault)
	screen.SetContent(0, height-1, '+', nil, tcell.StyleDefault)
	screen.SetContent(width-1, 0, '+', nil, tcell.StyleDefault)
	screen.SetContent(width-1, height-1, '+', nil, tcell.StyleDefault)
}

func drawContent(screen tcell.Screen, state []types.RepoState, width int, height int) {
	logrus.Infof("Drawing content into %d, %d", width, height)

	for line, repo := range state {
		for char, rune := range repo.Repo().Name {
			screen.SetContent(char+2+char, line+2, rune, nil, tcell.StyleDefault)
		}
	}
}
