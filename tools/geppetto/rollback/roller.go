package rollback

import (
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/geppetto/config"
	"ponglehub.co.uk/geppetto/rollback/commands"
)

// Roller a manager for rolling back version numbers
type Roller struct {
	commands []commands.Command
}

// FromConfig create a Roller from the config
func FromConfig(cfg *config.Config) Roller {
	roller := Roller{commands: []commands.Command{}}

	for _, repo := range cfg.Repos {
		switch repo.RepoType {
		case config.Node:
			roller.commands = append(
				roller.commands,
				commands.MakeNpmCommand(cfg.BasePath, repo),
			)
		case config.Go:
			roller.commands = append(
				roller.commands,
				commands.MakeGoCommand(cfg.BasePath, repo),
			)
		}
	}

	return roller
}

func (r *Roller) Rollback() {
	for _, command := range r.commands {
		logrus.Infof("Rolling back %s", command.Name())
		err := command.Run()
		if err != nil {
			logrus.Errorf("Failed to rollback %s: %+v", command.Name(), err)
		} else {
			logrus.Infof("Finished %s", command.Name())
		}
	}
}
