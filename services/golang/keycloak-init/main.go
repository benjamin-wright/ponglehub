package main

import "github.com/sirupsen/logrus"

func main() {
	logrus.Info("Starting...")

	cfg, err := newConfig()

	if err != nil {
		logrus.Fatalf("Failed to load config: %+v", err)
	}

	cfg.print()
	logrus.Info("Finished")
}
