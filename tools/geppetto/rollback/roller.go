package rollback

import (
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/geppetto/config"
	"ponglehub.co.uk/geppetto/services"
)

// Roller a manager for rolling back version numbers
type Roller struct {
	npmRepos []npmRoller
}

type npmRoller struct {
	name    string
	service services.NPMRepo
}

// FromConfig create a Roller from the config
func FromConfig(cfg *config.Config) (Roller, error) {
	roller := Roller{npmRepos: []npmRoller{}}

	for _, repo := range cfg.Repos {
		switch repo.RepoType {
		case config.Node:
			npmRepo, err := services.NewNpmRepo(cfg.BasePath + "/" + repo.Path)
			if err != nil {
				return roller, err
			}

			roller.npmRepos = append(roller.npmRepos, npmRoller{name: repo.Name, service: npmRepo})
		case config.Go:
		}
	}

	return roller, nil
}

// Rollback roll back all the versions
func (r *Roller) Rollback() {
	for _, repo := range r.npmRepos {
		err := repo.service.SetVersion("1.0.0")
		if err != nil {
			logrus.Errorf("Failed to roll back %s: %+v", repo.name, err)
		}
	}
}
