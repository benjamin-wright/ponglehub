package ui

import (
	tm "github.com/buger/goterm"
	"github.com/gdamore/tcell/v2"
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/geppetto/types"
)

type display struct {
	screen   tcell.Screen
	state    []types.RepoState
	selected int
	building bool
}

func (d *display) draw() {
	d.screen.Clear()
	width, height := d.screen.Size()
	drawBorder(d.screen, width, height)
	drawTitle(d.screen, d.building, width, height)
	drawContent(d.screen, d.state, d.selected, width, height)
	d.screen.Show()
}

func drawBorder(screen tcell.Screen, width int, height int) {
	style := tcell.StyleDefault
	style = style.Foreground(tcell.ColorGreen)

	for x := 1; x < width-1; x++ {
		screen.SetContent(x, 0, '-', nil, style)
		screen.SetContent(x, height-1, '-', nil, style)
	}

	for y := 1; y < height-1; y++ {
		screen.SetContent(0, y, '|', nil, style)
		screen.SetContent(width-1, y, '|', nil, style)
	}

	screen.SetContent(0, 0, '+', nil, style)
	screen.SetContent(0, height-1, '+', nil, style)
	screen.SetContent(width-1, 0, '+', nil, style)
	screen.SetContent(width-1, height-1, '+', nil, style)
}

func drawTitle(screen tcell.Screen, building bool, width int, height int) {
	title := "GEPPETTO"
	titleStart := width/2 - len(title)/2

	style := tcell.StyleDefault
	style = style.Foreground(tcell.ColorGreen)

	drawText(screen, title, titleStart, 1, len(title), style)

	if building {
		drawText(screen, "🏗️", width-4, 1, 10, style)
	} else {
		drawText(screen, "⏳", width-4, 1, 10, style)
	}
}

func drawContent(screen tcell.Screen, state []types.RepoState, selected int, width int, height int) {
	logrus.Infof("Drawing content into %d, %d", width, height)

	for line, repo := range state {
		icon := '?'
		switch repo.Repo().RepoType {
		case types.Node:
			icon = '🟢'
		case types.Golang:
			icon = '🐹'
		case types.Helm:
			icon = '⎈'
		}

		style := tcell.StyleDefault

		if line == selected {
			style = style.Background(tcell.ColorDarkGreen)
			for x := 3; x <= width-3; x++ {
				screen.SetContent(x, line+3, ' ', nil, style)
			}
		}

		screen.SetContent(2, line+3, icon, nil, style)

		drawText(screen, repo.Repo().Name, 5, line+3, 50, style)
		if repo.Built() {
			drawText(screen, "✅", 60, line+3, 5, style)
		} else if repo.Skipped() {
			drawText(screen, "🔄", 60, line+3, 5, style)
		} else if repo.Blocked() {
			drawText(screen, "❌", 60, line+3, 5, style)
		} else if repo.Errored() != nil {
			drawText(screen, "🔥", 60, line+3, 5, style)
			tm.Println(repo.Errored())
		} else if repo.Building() {
			if repo.Phase() == "check" {
				drawText(screen, "💡", 60, line+3, 5, style)
			} else {
				drawText(screen, "🏗️", 60, line+3, 5, style)
			}
		} else {
			drawText(screen, "⏳", 60, line+3, 5, style)
		}
	}
}

func drawText(screen tcell.Screen, content string, x int, y int, maxLength int, style tcell.Style) {
	ellipse := false
	if len(content) > maxLength {
		ellipse = true
	}

	for char, rune := range content {
		if char >= maxLength-3 && ellipse {
			if char < maxLength {
				screen.SetContent(x+char, y, '.', nil, style)
			}

			continue
		}

		screen.SetContent(x+char, y, rune, nil, style)
	}
}
