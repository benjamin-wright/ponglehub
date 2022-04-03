package database

type Client struct {
	Username   string
	Deployment string
	Database   string
	Namespace  string
}

func (c Client) Key() string {
	return c.Namespace + "/" + c.Deployment + "/" + c.Database + "/" + c.Username
}
