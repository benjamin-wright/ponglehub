package main

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("Running...")

	cfg, err := newConfig()
	if err != nil {
		logrus.Fatalf("Error loading config: %+v", err)
	}

	cfg.print()

	start := time.Now()

	for {
		logrus.Info("Testing url...")

		resp, err := http.Get(cfg.url)

		if resp != nil {
			logrus.Infof("Status: %d", resp.StatusCode)

			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				break
			}
		}

		if time.Since(start) >= time.Duration(cfg.timeout)*time.Second {
			logrus.Fatalf("Timed out: %+v", err)
		}

		time.Sleep(time.Second * time.Duration(cfg.interval))
	}

	logrus.Info("Finished...")
}
