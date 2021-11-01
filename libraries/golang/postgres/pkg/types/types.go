package types

import (
	"ponglehub.co.uk/lib/postgres/pkg/connect"
)

type Migration struct {
	Query string
}

type MigrationConfig struct {
	AdminConfig  connect.ConnectConfig
	TargetConfig connect.ConnectConfig
	Migrations   []Migration
}
