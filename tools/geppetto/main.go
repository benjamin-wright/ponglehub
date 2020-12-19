package main

import (
	"os"

	"ponglehub.co.uk/geppetto/types"
	"ponglehub.co.uk/geppetto/ui"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func getConfig(c *cli.Context) types.Config {
	return types.Config{
		Debug:     c.Bool("debug"),
		Target:    c.String("target"),
		ChartRepo: c.String("chartrepo"),
	}
}

func initLogger(cfg types.Config) {
	if cfg.Debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("Debug logging enabled")
	} else {
		logrus.SetLevel(logrus.DebugLevel)
		f, err := os.Create("gepetto.log")
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
			&cli.StringFlag{
				Name:    "chartrepo",
				Value:   "local",
				Usage:   "chart museum repo for local charts",
				EnvVars: []string{"GEPETTO_CHART_REPO"},
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "manual",
				Aliases: []string{"m"},
				Usage:   "build things on command",
				Action: func(c *cli.Context) error {
					cfg := getConfig(c)
					initLogger(cfg)

					manual, err := ui.NewManual(cfg.ChartRepo)
					if err != nil {
						logrus.Fatalf("Failed to create manual instance: %+v", err)
					}

					err = manual.Start(cfg.Target)
					defer manual.Destroy()
					if err != nil {
						logrus.Errorf("Manual exited with unexpected error: %+v", err)
					}

					return nil
				},
			},
			{
				Name:    "rollback",
				Aliases: []string{"r"},
				Usage:   "rollback all versions to 1.0.0",
				Action: func(c *cli.Context) error {
					cfg := getConfig(c)
					initLogger(cfg)

					rollback, err := ui.NewRollback(cfg.ChartRepo)
					if err != nil {
						logrus.Fatalf("Failed to create rollback instance: %+v", err)
					}

					rollback.Start(cfg.Target)
					defer rollback.Destroy()
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
