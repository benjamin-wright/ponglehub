package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"k8s.io/client-go/kubernetes/scheme"
	"ponglehub.co.uk/events/gateway/pkg/crds"
	"ponglehub.co.uk/tools/ponglehub/internal/commands"
)

func main() {
	crds.AddToScheme(scheme.Scheme)

	app := &cli.App{
		Name:        "ponglehub",
		Description: "admin cli for the ponglehub app",
		Commands: []*cli.Command{
			&commands.UserCommand,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
