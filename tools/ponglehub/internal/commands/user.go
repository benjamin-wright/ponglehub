package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

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
		&cli.StringFlag{
			Name:     "password",
			Aliases:  []string{"p"},
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

		logrus.Infof("User ID: %s", user.ID)

		cli := redis.New("localhost:6379")
		token, err := cli.WaitForKey(fmt.Sprintf("%s.%s", user.ID, "invite"))
		if err != nil {
			return err
		}

		logrus.Infof("Invite token: %s", token)

		json_data, err := json.Marshal(map[string]string{
			"invite":   token,
			"password": c.String("password"),
			"confirm":  c.String("password"),
		})
		if err != nil {
			return err
		}

		res, err := http.Post("http://localhost:4000/auth/set-password", "application/json", bytes.NewBuffer(json_data))
		if err != nil {
			return err
		}
		if res.StatusCode != 200 {
			return fmt.Errorf("failed to set password, response code: %d", res.StatusCode)
		}

		return nil
	},
}
