package implementer

import (
	"github.com/sirupsen/logrus"
)

func Start(actions <-chan interface{}) chan<- interface{} {
	stopper := make(chan interface{})
	implementer := New()

	go func() {
		running := true
		for running {
			select {
			case <-stopper:
				running = false
				logrus.Info("Stopped implementer")
			case action := <-actions:
				logrus.Infof("ACTION (queued): %T", action)
				implementer.AddAction(action)
			case action := <-implementer.NextAction:
				logrus.Infof("ACTION: %T", action)
				implementer.DoAction(action)
			}
		}
	}()

	return stopper
}
