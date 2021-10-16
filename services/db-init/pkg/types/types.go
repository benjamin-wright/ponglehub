package types

type Migration struct {
	Query string
}

type MigrationConfig struct {
	Host       string
	Port       int
	Username   string
	Password   string
	Database   string
	CertsDir   string
	Migrations []Migration
}
