package migrate_test

import (
	"io"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"ponglehub.co.uk/lib/postgres/internal/database"
	"ponglehub.co.uk/lib/postgres/pkg/connect"
	"ponglehub.co.uk/lib/postgres/pkg/migrate"
	"ponglehub.co.uk/lib/postgres/pkg/types"
)

func assertSchemas(u *testing.T, db *database.MigrationConn, expected map[string]map[string]string) {
	schemas := map[string]map[string]string{}

	for table := range expected {
		schema, err := db.GetTableSchema(table)
		if err != nil {
			assert.NoError(u, err)
			assert.FailNow(u, "error fetching schema")
		}

		schemas[table] = schema
	}

	assert.Equal(u, expected, schemas)
}

func assertContents(u *testing.T, db *database.MigrationConn, expected map[string][][]interface{}) {
	contents := map[string][][]interface{}{}

	for table := range expected {
		content, err := db.GetContents(table)
		if err != nil {
			assert.NoError(u, err)
			assert.FailNow(u, "error fetching database contents")
		}

		contents[table] = content

		for rowId, row := range expected[table] {
			for colId, obj := range row {
				if obj == mock.Anything {
					if rowId < len(content) && colId < len(content[rowId]) {
						expected[table][rowId][colId] = content[rowId][colId]
					}
				}
			}
		}
	}

	assert.Equal(u, expected, contents)
}

func TestMigrations(t *testing.T) {
	logrus.SetOutput(io.Discard)

	for _, test := range []struct {
		Name       string
		Migrations []types.Migration
		Schemas    map[string]map[string]string
		Contents   map[string][][]interface{}
	}{
		{
			Name: "simples",
			Migrations: []types.Migration{
				{
					Query: `
						CREATE TABLE test_users (
							id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
							name VARCHAR(100) NOT NULL UNIQUE,
							email VARCHAR(100) NOT NULL UNIQUE,
							password VARCHAR(100),
							verified BOOLEAN NOT NULL
						);
					`,
				},
			},
			Schemas: map[string]map[string]string{
				"test_users": {
					"id":       "uuid",
					"name":     "character varying",
					"email":    "character varying",
					"password": "character varying",
					"verified": "boolean",
				},
				"migrations": {
					"id": "bigint",
				},
			},
			Contents: map[string][][]interface{}{
				"test_users": {},
				"migrations": {
					{int64(0)},
				},
			},
		},
		{
			Name: "multi-step",
			Migrations: []types.Migration{
				{Query: `
					CREATE TABLE test_users (
						id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
						name VARCHAR(100) NOT NULL UNIQUE,
						email VARCHAR(100) NOT NULL UNIQUE,
						password VARCHAR(100),
						verified BOOLEAN NOT NULL
					);
				`},
				{Query: `
					INSERT INTO test_users (name, email, password, verified) VALUES ('fred', 'fred@gmail.com', 'my-pass', true)
				`},
			},
			Schemas: map[string]map[string]string{
				"test_users": {
					"id":       "uuid",
					"name":     "character varying",
					"email":    "character varying",
					"password": "character varying",
					"verified": "boolean",
				},
				"migrations": {
					"id": "bigint",
				},
			},
			Contents: map[string][][]interface{}{
				"test_users": {
					{mock.Anything, "fred", "fred@gmail.com", "my-pass", true},
				},
				"migrations": {{int64(0)}, {int64(1)}},
			},
		},
	} {
		t.Run(test.Name, func(u *testing.T) {
			targetCfg, err := connect.ConfigFromEnv()
			assert.NoError(t, err)
			targetCfg.Database = "test_db"

			adminCfg, err := connect.AdminFromEnv()
			assert.NoError(t, err)
			targetCfg.Database = "test_db"

			config := types.MigrationConfig{
				AdminConfig:  adminCfg,
				TargetConfig: targetCfg,
				Migrations:   test.Migrations,
			}

			err = migrate.Clean(&config)
			assert.NoError(t, err)

			err = migrate.Migrate(&config)
			assert.NoError(t, err)

			db, err := database.NewMigrationConn(targetCfg)
			if err != nil {
				assert.NoError(t, err)
				assert.FailNow(t, "error connecting to db")
			}

			assertSchemas(u, db, test.Schemas)
			assertContents(u, db, test.Contents)
		})
	}
}
