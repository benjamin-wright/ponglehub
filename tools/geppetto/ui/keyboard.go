package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/sirupsen/logrus"
)

type keyboard struct {
	screen  tcell.Screen
	stopper chan<- bool
}

func (k *keyboard) start() <-chan command {
	commands := make(chan command, 5)

	go func() {
		for {
			event := k.screen.PollEvent()
			switch e := event.(type) {
			case *tcell.EventKey:
				switch e.Key() {
				case tcell.KeyESC:
					commands <- quitCommand
				case tcell.KeyUp:
					logrus.Info("Up key pressed")
					commands <- upCommand
				case tcell.KeyDown:
					logrus.Info("Down key pressed")
					commands <- downCommand
				case tcell.KeyEnter:
					commands <- selectCommand
				}
			case nil:
				logrus.Info("Keyboard listener loop stopped: screen finalised")
				return
			}
		}
	}()

	return commands
}
