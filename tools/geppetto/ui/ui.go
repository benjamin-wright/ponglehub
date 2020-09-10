package ui

import (
	"fmt"
	"time"

	"ponglehub.co.uk/geppetto/builder"

	tm "github.com/buger/goterm"
	"github.com/gdamore/tcell/v2"
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/geppetto/scanner"
	"ponglehub.co.uk/geppetto/types"
)

// UI a collection of methods for drawing fullscreen ascii UI outputs
type UI struct {
	builder  *builder.Builder
	screen   tcell.Screen
	display  watchView
	keyboard keyboard
	scan     *scanner.Scanner
}

// Watch UI for file watching
func (ui *UI) Watch(target string) error {
	ui.scan = scanner.New()
	repos, err := ui.scan.ScanDir(target)
	if err != nil {
		return err
	}

	logrus.Infof("Repos: %+v", repos)

	triggers, _, _ := ui.scan.WatchDir(repos)

	progress := make(chan []types.RepoState, 3)

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
	ui.display = watchView{screen: screen}
	ui.builder = builder.New()

	commands := ui.keyboard.start()
	ui.display.building = true
	ui.display.draw()
	go func() {
		logrus.Info("Building...")
		ui.builder.Build(repos, progress)
		ui.display.building = false
	}()

	for {
		select {
		case cmd := <-commands:
			switch cmd {
			case quitCommand:
				return nil
			case upCommand:
				if ui.display.selected > 0 {
					ui.display.selected--
				}
			case downCommand:
				if ui.display.selected < len(ui.display.state)-1 {
					ui.display.selected++
				}
			}
		case repo := <-triggers:
			logrus.Infof("Got trigger for %s", repo.Name)
			if !ui.display.building {
				ui.display.building = true
				go func() {
					logrus.Info("Building...")
					ui.builder.Build(repos, progress)
					ui.display.building = false
				}()
			}
		case state := <-progress:
			ui.display.state = state
		}
		ui.display.draw()
	}
}

// Start begin drawing updates of build progress
func (ui *UI) Start(progress <-chan []types.RepoState, finished chan<- bool) {
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
