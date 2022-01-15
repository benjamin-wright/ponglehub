package commands

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"ponglehub.co.uk/events/gateway/pkg/crds"
	"ponglehub.co.uk/tools/ponglehub/internal/services/redis"
)

var UserCommand = cli.Command{
	Name:        "users",
	Description: "commands to manage users",
	Subcommands: []*cli.Command{
		&AddUserCommand,
	},
}

var AddUserCommand = cli.Command{
	Name:        "add",
	Description: "add a new user",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "resource-name",
			Aliases:  []string{"r"},
			Required: true,
		},
		&cli.StringFlag{
			Name:     "display-name",
			Aliases:  []string{"d"},
			Required: true,
		},
		&cli.StringFlag{
			Name:     "email",
			Aliases:  []string{"e"},
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		client, err := crds.New(&crds.ClientArgs{
			External: true,
		})
		if err != nil {
			return err
		}

		user, err := client.Create(crds.User{
			Name:    c.String("resource-name"),
			Display: c.String("display-name"),
			Email:   c.String("email"),
		})
		if err != nil {
			return err
		}

		logrus.Infof("User: %+v", user)

		cli := redis.New("localhost:6379")
		token, err := cli.WaitForKey(fmt.Sprintf("%s.%s", user.ID, "invite"))
		if err != nil {
			return err
		}

		logrus.Infof("Invite token: %s", token)

		return nil
	},
}
