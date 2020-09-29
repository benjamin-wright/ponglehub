package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

type config struct {
	url      string
	interval int
	timeout  int
}

func newConfig() (config, error) {
	cfg := config{}

	if url, ok := os.LookupEnv("WAIT_FOR_URL"); ok {
		cfg.url = url
		logrus.Debugf("Loaded WAIT_FOR_URL value: %s", url)
	} else {
		return cfg, errors.New("Missing value for WAIT_FOR_URL")
	}

	if interval, ok := os.LookupEnv("WAIT_FOR_INTERVAL"); ok {
		intValue, err := strconv.Atoi(interval)
		if err != nil {
			return cfg, fmt.Errorf("Error converting WAIT_FOR_INTERVAL to integer: %+v", err)
		}

		if intValue < 1 {
			return cfg, fmt.Errorf("WAIT_FOR_INTERVAL should be a whole positive number, received %s", interval)
		}

		cfg.interval = intValue
		logrus.Debugf("Loaded WAIT_FOR_INTERVAL value: %s", interval)
	} else {
		cfg.interval = 2
		logrus.Debugf("Defaulted WAIT_FOR_INTERVAL value")
	}

	if timeout, ok := os.LookupEnv("WAIT_FOR_TIMEOUT"); ok {
		intValue, err := strconv.Atoi(timeout)
		if err != nil {
			return cfg, fmt.Errorf("Error converting WAIT_FOR_TIMEOUT to integer: %+v", err)
		}

		if intValue < 1 {
			return cfg, fmt.Errorf("WAIT_FOR_TIMEOUT should be a whole positive number, received %s", timeout)
		}

		cfg.timeout = intValue
		logrus.Debugf("Loaded WAIT_FOR_TIMEOUT value: %s", timeout)
	} else {
		cfg.timeout = 60
		logrus.Debugf("Defaulted WAIT_FOR_TIMEOUT value")
	}

	return cfg, nil
}

func (c *config) print() {
	logrus.Infof("Config:\n url: %s\n interval: %d\n timeout: %d", c.url, c.interval, c.timeout)
}
