package ui

import (
	"fmt"
	"time"

	"ponglehub.co.uk/geppetto/builder"

	tm "github.com/buger/goterm"
	"github.com/gdamore/tcell/v2"
	"ponglehub.co.uk/geppetto/scanner"
	"ponglehub.co.uk/geppetto/types"
)

// UI a collection of methods for drawing fullscreen ascii UI outputs
type UI struct {
	builder *builder.Builder
	screen  tcell.Screen
	scan    *scanner.Scanner
}

// Start drawing updates of build progress
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
