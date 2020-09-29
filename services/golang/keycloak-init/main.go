package main

import "github.com/sirupsen/logrus"

func main() {
	logrus.Info("Starting...")

	cfg := config{
		host: "ponglehub.co.uk",
	}

	cfg.print()
	logrus.Info("Finished")
}
