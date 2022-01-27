package main

import (
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/lib/postgres/pkg/connect"
	"ponglehub.co.uk/lib/postgres/pkg/migrate"
	"ponglehub.co.uk/lib/postgres/pkg/types"
)

func main() {
	logrus.Infof("Starting migrations...")

	cfg, err := connect.ConfigFromEnv()
	if err != nil {
		logrus.Fatalf("Failed to get postgres config: %+v", err)
	}

	err = migrate.Migrate(
		cfg,
		[]types.Migration{
			{
				Query: `
					BEGIN;

					SAVEPOINT games_restart;

					CREATE TABLE games (
						id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
						player1 UUID,
						player2 UUID,
						turn INT2
					);

					CREATE TABLE marks (
						id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
						game_id UUID NOT NULL,
						player INT2,
						position INT2,
						CONSTRAINT fk_game
							FOREIGN KEY(game_id) 
							REFERENCES games(id)
							ON DELETE CASCADE
					);

					RELEASE SAVEPOINT games_restart;

					COMMIT;
				`,
			},
		},
	)
	if err != nil {
		logrus.Fatalf("Migrations failed: %+v", err)
	}

	logrus.Infof("Done.")
}
