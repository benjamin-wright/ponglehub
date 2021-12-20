package types

type Client struct {
	Name       string
	Username   string
	Namespace  string
	Deployment string
	Database   string
	Secret     string
	Ready      bool
}
