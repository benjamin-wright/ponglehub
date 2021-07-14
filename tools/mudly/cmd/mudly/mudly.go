package main

import (
	"flag"
	"os"

	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/tools/mudly/internal/target"
)

func setLogLevel() {
	if logLevel, ok := os.LookupEnv("MUDLY_LOG_LEVEL"); ok {
		parsedLevel, err := logrus.ParseLevel(logLevel)
		if err != nil {
			logrus.Fatalf("Failed to parse log level %s from environment", logLevel)
		}

		logrus.SetLevel(parsedLevel)
	} else if logLevel := flag.String("log-level", "info", "the logging level to use"); logLevel != nil {
		parsedLevel, err := logrus.ParseLevel(*logLevel)
		if err != nil {
			logrus.Fatalf("Failed to parse log level %s from --log-level flag", logLevel)
		}

		logrus.SetLevel(parsedLevel)
	}
}

func main() {
	setLogLevel()

	args := os.Args[1:]
	logrus.Debugf("Running mudly with args: %+v", args)

	if len(args) == 0 {
		logrus.Fatalf("must provide a build target")
	}

	targets := []target.Target{}
	for _, path := range args {
		target, err := target.ParseTarget(path)
		if err != nil {
			logrus.Fatalf("Error parsing target: %+v", err)
		}

		targets = append(targets, *target)
	}

	logrus.Debugf("Targets: %+v", targets)
}
