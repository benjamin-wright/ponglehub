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

func getConfig(c *cli.Context) types.Config {
	return types.Config{
		Debug:  c.Bool("debug"),
		Target: c.String("target"),
	}
}

func initLogger(cfg types.Config) {
	if cfg.Debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("Debug logging enabled")
	} else {
		f, err := os.OpenFile("gepetto.log", os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			logrus.Fatalf("Failed to redirect logs to file: %+v", err)
		}
		logrus.SetOutput(f)
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
		},
		Commands: []*cli.Command{
			{
				Name:    "watch",
				Aliases: []string{"w"},
				Usage:   "build everything whenever it changes",
				Action: func(c *cli.Context) error {
					cfg := getConfig(c)
					initLogger(cfg)

					scan := scanner.New()
					disp := display.Display{}

					repos, err := scan.ScanDir(cfg.Target)
					if err != nil {
						return err
					}

					triggers, errors, _ := scan.WatchDir(repos)

					disp.Watch(triggers, errors)
					return nil
				},
			},
			{
				Name:    "build",
				Aliases: []string{"b"},
				Usage:   "build everything",
				Action: func(c *cli.Context) error {
					cfg := getConfig(c)
					initLogger(cfg)

					disp := display.Display{}
					progress := make(chan []types.RepoState, 3)
					finished := make(chan bool)

					if !cfg.Debug {
						go disp.Start(progress, finished)
					} else {
						go func(prg chan []types.RepoState) {
							for range prg {
							}
						}(progress)
					}

					repos, err := scanner.New().ScanDir(cfg.Target)
					if err != nil {
						close(progress)
						return err
					}

					logrus.Infof("Repos: %+v", repos)

					b := builder.New()
					err = b.Build(repos, progress)

					close(progress)

					if !cfg.Debug {
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
					cfg := getConfig(c)
					initLogger(cfg)

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
