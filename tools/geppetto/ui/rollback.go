package ui

import (
	"github.com/gdamore/tcell/v2"
	"ponglehub.co.uk/geppetto/builder"
	"ponglehub.co.uk/geppetto/scanner"
	"ponglehub.co.uk/geppetto/services"
	"ponglehub.co.uk/geppetto/types"
)

type Rollback struct {
	builder *builder.Builder
	devices *devices
	scanner *scanner.Scanner
}

func NewRollback(chartRepo string) (*Rollback, error) {
	devices, err := makeDevices()
	if err != nil {
		return nil, err
	}

	scanner := scanner.New()
	builder := builder.New(chartRepo)

	return &Rollback{
		builder: builder,
		devices: devices,
		scanner: scanner,
	}, nil
}

func (r *Rollback) Destroy() {
	r.devices.destroy()
}

func (r *Rollback) Start(target string) error {
	repos, err := r.scanner.ScanDir(target)
	if err != nil {
		return err
	}

	rollbackEvents := make(chan string, 5)
	commandEvents := r.devices.listen()

	for index, repo := range repos {
		npm := services.NewNpmService()
		helm := services.NewHelmService()
		go func(repo types.Repo, timeout int) {
			switch repo.RepoType {
			case types.Node:
				npm.SetVersion(repo, "1.0.0")
				npm.Install(repo)
			case types.Helm:
				helm.SetVersion(repo, "1.0.0")
				helm.Install(repo)
			}
			rollbackEvents <- repo.Name
		}(repo, index)
	}

	rolledBack := []string{}
	for {
		width, height := r.devices.getSize()

		r.devices.clear()
		r.devices.drawBorder(width, height)
		r.devices.drawTitle("GEPPETTO", width, len(repos) != len(rolledBack))

		offset := 3

		for line, repo := range repos {
			hasRolled := false
			for _, name := range rolledBack {
				if repo.Name == name {
					hasRolled = true
				}
			}

			style := tcell.StyleDefault

			r.devices.drawIcon(repo.RepoType, 2, line+offset, style)
			r.devices.drawText(repo.Name, 5, line+offset, 50, style)
			if hasRolled {
				r.devices.drawText("✅", 60, line+offset, 5, style)
			} else {
				r.devices.drawText("🏗", 60, line+offset, 5, style)
			}
		}

		r.devices.flush()

		select {
		case repo := <-rollbackEvents:
			if repo != "" {
				rolledBack = append(rolledBack, repo)
			}
		case cmd := <-commandEvents:
			if cmd == quitCommand {
				return nil
			}
		}
	}
}
