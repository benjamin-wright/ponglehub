package main

import "github.com/sirupsen/logrus"

type config struct {
	host string
}

func (c *config) print() {
	logrus.Infof("Config:\n - host: %s", c.host)
}
