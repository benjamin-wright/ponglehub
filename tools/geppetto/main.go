package main

import (
	"os"

	"ponglehub.co.uk/geppetto/display"
	"ponglehub.co.uk/geppetto/types"

	"ponglehub.co.uk/geppetto/builder"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"ponglehub.co.uk/geppetto/scanner"
)

func initLogger(c *cli.Context) {
	if c.Bool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("Debug logging enabled")
	} else {
		logrus.SetLevel(logrus.FatalLevel)
	}
}

func enableLogging() {
	logrus.SetLevel(logrus.DebugLevel)
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
					disp := display.Display{}
					progress := make(chan []types.RepoStatus, 3)
					finished := make(chan interface{})

					if !c.Bool("debug") {
						go disp.Start(progress, finished)
					} else {
						go func(prg chan []types.RepoStatus) {
							for range prg {
							}
						}(progress)
					}

					repos, err := scanner.New().ScanDir(c.String("target"))
					if err != nil {
						close(progress)
						return err
					}

					logrus.Infof("Repos: %+v", repos)

					b := builder.New()
					err = b.Build(repos, progress)

					close(progress)

					if !c.Bool("debug") {
						<-finished
					}

					return err
				},
			},
			{
				Name:    "rollback",
				Aliases: []string{"r"},
				Usage:   "rollback all versions to 1.0.0",
				Action: func(c *cli.Context) error {
					initLogger(c)
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
