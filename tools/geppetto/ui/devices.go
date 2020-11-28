package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/geppetto/types"
)

type devices struct {
	screen tcell.Screen
}

func makeDevices() (*devices, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, fmt.Errorf("Error making screen: %+v", err)
	}

	err = screen.Init()
	if err != nil {
		return nil, fmt.Errorf("Error initing screen screen: %+v", err)
	}

	return &devices{screen: screen}, nil
}

func (d *devices) destroy() {
	d.screen.Fini()
}

func (d *devices) clear() {
	d.screen.Clear()
}

func (d *devices) flush() {
	d.screen.Show()
}

func (d *devices) drawBorder(width int, height int) {
	style := tcell.StyleDefault
	style = style.Foreground(tcell.ColorGreen)

	for x := 1; x < width-1; x++ {
		d.screen.SetContent(x, 0, '-', nil, style)
		d.screen.SetContent(x, height-1, '-', nil, style)
	}

	for y := 1; y < height-1; y++ {
		d.screen.SetContent(0, y, '|', nil, style)
		d.screen.SetContent(width-1, y, '|', nil, style)
	}

	d.screen.SetContent(0, 0, '+', nil, style)
	d.screen.SetContent(0, height-1, '+', nil, style)
	d.screen.SetContent(width-1, 0, '+', nil, style)
	d.screen.SetContent(width-1, height-1, '+', nil, style)
}

func (d *devices) drawTitle(content string, width int, building bool) {
	title := "GEPPETTO"
	titleStart := width/2 - len(title)/2

	style := tcell.StyleDefault
	style = style.Foreground(tcell.ColorGreen)

	d.drawText(title, titleStart, 1, len(title), style)

	if building {
		d.drawText("🏗", width-4, 1, 10, style)
	} else {
		d.drawText("⏳", width-4, 1, 10, style)
	}
}

func (d *devices) highlightLine(y int, width int, style tcell.Style) {
	for x := 3; x <= width-3; x++ {
		d.screen.SetContent(x, y, ' ', nil, style)
	}
}

func (d *devices) drawText(content string, x int, y int, maxLength int, style tcell.Style) {
	ellipse := false
	if len(content) > maxLength {
		ellipse = true
	}

	for char, rune := range content {
		if char >= maxLength-3 && ellipse {
			if char < maxLength {
				d.screen.SetContent(x+char, y, '.', nil, style)
			}

			continue
		}

		d.screen.SetContent(x+char, y, rune, nil, style)
	}
}

func (d *devices) drawIcon(t types.RepoType, x int, y int, style tcell.Style) {
	icon := '?'
	switch t {
	case types.Node:
		icon = '🟢'
	case types.Golang:
		icon = '🐹'
	case types.Helm:
		icon = '⛵'
	case types.Rust:
		icon = '🦀'
	}

	d.screen.SetContent(x, y, icon, nil, style)
}

func (d *devices) getNumLines(content string, width int) int {
	runes := []rune(content)
	xCoord := 0
	yCoord := 1

	for _, rune := range runes {
		xCoord++
		if xCoord > width {
			xCoord = 0
			yCoord++
			continue
		}

		if rune == '\n' {
			xCoord = 0
			yCoord++
			continue
		}
	}

	return yCoord
}

func (d *devices) getSize() (int, int) { return d.screen.Size() }

func (d *devices) drawMultiline(
	content string,
	x int,
	y int,
	width int,
	height int,
	scroll int,
	style tcell.Style,
) {
	runes := []rune(content)

	xCoord := 0
	yCoord := 0

	for _, rune := range runes {
		xCoord++
		if xCoord > width {
			xCoord = 0
			yCoord++
		}

		if rune == '\n' {
			xCoord = 0
			yCoord++
			continue
		}

		if yCoord >= height+scroll {
			return
		}

		if yCoord > scroll && yCoord < height+scroll {
			d.screen.SetContent(x+xCoord, y+yCoord-scroll, rune, nil, style)
		}
	}
}

func (d *devices) listen() <-chan command {
	commands := make(chan command, 5)

	go func() {
		for {
			event := d.screen.PollEvent()
			switch e := event.(type) {
			case *tcell.EventKey:
				switch e.Key() {
				case tcell.KeyESC:
					fallthrough
				case tcell.KeyCtrlC:
					commands <- quitCommand
				case tcell.KeyUp:
					commands <- upCommand
				case tcell.KeyDown:
					commands <- downCommand
				case tcell.KeyEnter:
					commands <- selectCommand
				case tcell.KeyRune:
					if e.Rune() == ' ' {
						commands <- unlockCommand
					}
				}
			case nil:
				logrus.Debug("Keyboard listener loop stopped: screen finalised")
				return
			}
		}
	}()

	return commands
}
