package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"ponglehub.co.uk/geppetto/config"
	"ponglehub.co.uk/geppetto/scanner"
)

func initLogger(c *cli.Context) {
	if c.Bool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("Debug logging enabled")
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func main() {
	app := &cli.App{
		Name:  "Geppetto",
		Usage: "Make toys",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "debug",
				Value:   false,
				Usage:   "enable debug logging",
				EnvVars: []string{"GEPETTO_DEBUG"},
			},
			&cli.StringFlag{
				Name:    "target",
				Value:   ".",
				Usage:   "target directory",
				EnvVars: []string{"GEPETTO_TARGET"},
			},
			&cli.StringFlag{
				Name:    "config",
				Value:   ".geppetto.json",
				Usage:   "path to the config file",
				EnvVars: []string{"GEPETTO_CONFIG"},
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "build",
				Aliases: []string{"b"},
				Usage:   "build everything",
				Action: func(c *cli.Context) error {
					initLogger(c)
					_, err := config.FromFile(c.String("config"))
					if err != nil {
						return err
					}

					scanner.New().ScanDir(c.String("target"))
					logrus.Warn("Not implemented yet")

					return nil
				},
			},
			{
				Name:    "rollback",
				Aliases: []string{"r"},
				Usage:   "rollback all versions to 1.0.0",
				Action: func(c *cli.Context) error {
					initLogger(c)
					_, err := config.FromFile(c.String("config"))
					if err != nil {
						return err
					}

					logrus.Warn("Not implemented yet")

					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}
