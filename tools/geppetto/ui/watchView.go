package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/geppetto/types"
)

type watchView struct {
	screen      tcell.Screen
	state       []types.RepoState
	highlighted int
	selected    int
	building    bool
	scroll      int
	maxScroll   int
}

func (d *watchView) draw() {
	d.screen.Clear()
	width, height := d.screen.Size()
	d.drawBorder(width, height)
	d.drawTitle(width, height)
	d.drawContent(width, height)
	d.screen.Show()
}

func (d *watchView) drawBorder(width int, height int) {
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

func (d *watchView) drawTitle(width int, height int) {
	title := "GEPPETTO"
	titleStart := width/2 - len(title)/2

	style := tcell.StyleDefault
	style = style.Foreground(tcell.ColorGreen)

	d.drawText(title, titleStart, 1, len(title), style)

	if d.building {
		d.drawText("🏗️", width-4, 1, 10, style)
	} else {
		d.drawText("⏳", width-4, 1, 10, style)
	}
}

func (d *watchView) drawContent(width int, height int) {
	logrus.Infof("Drawing content into %d, %d", width, height)

	offset := 3
	spareLines := height - 6 - len(d.state)
	d.maxScroll = 0

	for line, repo := range d.state {
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

		if line == d.highlighted {
			style = style.Background(tcell.ColorDarkGreen)
			for x := 3; x <= width-3; x++ {
				d.screen.SetContent(x, line+offset, ' ', nil, style)
			}
		}
		if line == d.selected {
			style = style.Foreground(tcell.ColorLightSlateGray)
			for x := 3; x <= width-3; x++ {
				d.screen.SetContent(x, line+offset, ' ', nil, style)
			}
		}

		d.screen.SetContent(2, line+offset, icon, nil, style)

		d.drawText(repo.Repo().Name, 5, line+offset, 50, style)
		if repo.Built() {
			d.drawText("✅", 60, line+offset, 5, style)
		} else if repo.Skipped() {
			d.drawText("🔄", 60, line+offset, 5, style)
		} else if repo.Blocked() {
			d.drawText("❌", 60, line+offset, 5, style)
		} else if repo.Errored() != nil {
			d.drawText("🔥", 60, line+offset, 5, style)
			if line == d.selected {
				lines := d.drawMultiline(repo.Errored().Error(), 3, line+4, width-6, spareLines, d.scroll, tcell.StyleDefault)
				d.maxScroll = d.getScrollHeight(spareLines, lines)
				if d.scroll > d.maxScroll {
					d.scroll = d.maxScroll
				}
				if lines > spareLines {
					lines = spareLines
				}
				offset = offset + 1 + lines
			}
		} else if repo.Building() {
			if repo.Phase() == "check" {
				d.drawText("💡", 60, line+offset, 5, style)
			} else {
				d.drawText("🏗️", 60, line+offset, 7, style)
				d.drawText(repo.Phase(), 64, line+offset, 20, style)
			}
		} else {
			d.drawText("⏳", 60, line+offset, 5, style)
		}
	}
}

func (d *watchView) drawText(content string, x int, y int, maxLength int, style tcell.Style) {
	ellipse := false
	if len(content) > maxLength {
		logrus.Infof("Content length: %d %s", len(content), content)
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

func (d *watchView) drawMultiline(content string, x int, y int, width int, height int, scroll int, style tcell.Style) int {
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
		}

		if yCoord > scroll && yCoord < height+scroll {
			d.screen.SetContent(x+xCoord, y+yCoord-scroll, rune, nil, style)
		}
	}

	return yCoord + 1
}

func (d *watchView) getScrollHeight(height int, lines int) int {
	if lines > height {
		return lines - height
	}

	return 0
}
