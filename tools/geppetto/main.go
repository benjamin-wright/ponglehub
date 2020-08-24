package main

import (
	"os"

	"ponglehub.co.uk/geppetto/builder"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"ponglehub.co.uk/geppetto/config"
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
					cfg, err := config.FromFile(c.String("config"))
					if err != nil {
						return err
					}

					b, err := builder.FromConfig(cfg)
					if err != nil {
						logrus.Fatal(err)
					}

					err = b.Build()
					if err != nil {
						logrus.Fatal(err)
					}

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
