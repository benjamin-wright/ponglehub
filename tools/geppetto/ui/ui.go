package ui

import (
	"fmt"
	"time"

	tm "github.com/buger/goterm"
	"github.com/gdamore/tcell/v2"
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/geppetto/types"
)

// Display a collection of methods for drawing fullscreen ascii UI outputs
type UI struct {
	screen   tcell.Screen
	display  display
	keyboard keyboard
}

// Stop stop all watchers and close the console writer
func (ui *UI) Stop() {
	ui.display.stop()
	ui.screen.Fini()
}

// Watch UI for file watching
func (ui *UI) Watch(progress <-chan []types.RepoState, errors <-chan error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		logrus.Fatalf("Error making screen: %+v", err)
	}

	err = screen.Init()
	if err != nil {
		logrus.Fatalf("Error initing screen: %+v", err)
	}

	defer screen.Fini()

	ui.keyboard = keyboard{screen: screen}
	ui.display = display{screen: screen}

	commands := ui.keyboard.start()
	ui.display.start(progress, commands)

	// go func() {
	// 	for {
	// 		event := screen.PollEvent()
	// 		switch e := event.(type) {
	// 		case *tcell.EventKey:
	// 			switch e.Key() {
	// 			case tcell.KeyESC:
	// 				screen.Fini()
	// 				logrus.Fatalf("Killed it!")
	// 			case tcell.KeyRune:
	// 				inputs <- e.Rune()
	// 			default:
	// 				inputs <- 'y'
	// 			}
	// 		}
	// 	}
	// }()

	// for {
	// 	// var repo types.Repo
	// 	// var err error
	// 	var input rune
	// 	select {
	// 	// case repo = <-triggers:
	// 	// case err = <-errors:
	// 	case input = <-inputs:
	// 	}

	// 	screen.Clear()
	// 	screen.SetContent(1, 1, input, nil, tcell.StyleDefault)
	// 	screen.Show()
	// }
}

// Start begin drawing updates of build progress
func (d *UI) Start(progress <-chan []types.RepoState, finished chan<- bool) {
	for p := range progress {
		tm.Clear()
		tm.MoveCursor(1, 1)
		tm.Println(time.Now().Format(time.RFC1123))
		tm.MoveCursor(1, 3)
		for _, r := range p {
			icon := " ?"
			switch r.Repo().RepoType {
			case types.Node:
				icon = "🟢"
			case types.Golang:
				icon = "🐹"
			case types.Helm:
				icon = " ⎈"
			}

			if r.Built() {
				tm.Println(fmt.Sprintf("%s %s: ✅", icon, r.Repo().Name))
			} else if r.Skipped() {
				tm.Println(fmt.Sprintf("%s %s: 🔄", icon, r.Repo().Name))
			} else if r.Blocked() {
				tm.Println(fmt.Sprintf("%s %s: ❌", icon, r.Repo().Name))
			} else if r.Errored() != nil {
				tm.Println(fmt.Sprintf("%s %s: 🔥", icon, r.Repo().Name))
				tm.Println(r.Errored())
			} else if r.Building() {
				if r.Phase() == "check" {
					tm.Println(fmt.Sprintf("%s %s: 💡", icon, r.Repo().Name))
				} else {
					tm.Println(fmt.Sprintf("%s %s: 🏗️ (%s)", icon, r.Repo().Name, r.Phase()))
				}
			} else {
				tm.Println(fmt.Sprintf("%s %s: ⏳", icon, r.Repo().Name))
			}
		}

		tm.Println()
		tm.Println("Press...    q: quit")
		tm.Flush()
	}

	finished <- true
}
