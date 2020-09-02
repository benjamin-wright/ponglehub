package display

import (
	"fmt"
	"time"

	tm "github.com/buger/goterm"
	"ponglehub.co.uk/geppetto/types"
)

type Display struct{}

func (d *Display) Start(progress <-chan []types.RepoStatus, finished chan<- interface{}) {
	for p := range progress {
		tm.Clear()
		tm.MoveCursor(1, 1)
		tm.Println(time.Now().Format(time.RFC1123))
		tm.MoveCursor(1, 3)
		for _, r := range p {
			if r.Built {
				tm.Println(fmt.Sprintf("%s: ✅", r.Repo.Name))
			} else if r.Skipped {
				tm.Println(fmt.Sprintf("%s: 🔄", r.Repo.Name))
			} else if r.Blocked {
				tm.Println(fmt.Sprintf("%s: ❌", r.Repo.Name))
			} else if r.Error {
				tm.Println(fmt.Sprintf("%s: 🔥", r.Repo.Name))
			} else if r.Building {
				if r.Phase == "check" {
					tm.Println(fmt.Sprintf("%s: 💡", r.Repo.Name))
				} else {
					tm.Println(fmt.Sprintf("%s: 🏗️ (%s)", r.Repo.Name, r.Phase))
				}
			} else {
				tm.Println(fmt.Sprintf("%s: ⏳", r.Repo.Name))
			}
		}
		tm.Flush()
	}

	finished <- "done"
}
