package display

import (
	"time"

	tm "github.com/buger/goterm"
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/geppetto/builder"
)

type Display struct{}

func (d *Display) Start(progress <-chan builder.BuildState, stopper <-chan interface{}) {
	tm.Clear()
	running := true
	state := builder.BuildState{}

	for running {
		tm.MoveCursor(1, 1)
		tm.Println(time.Now().Format(time.RFC1123))
		tm.MoveCursor(1, 3)

		for _, r := range state.Repos {
			tm.Println("hi: ", r.)
		}
		tm.Flush()
		time.Sleep(time.Millisecond * 100)

		select {
		case <-stopper:
			running = false
		default:
		}
	}

	logrus.Debug("Stopping display loop")
}
