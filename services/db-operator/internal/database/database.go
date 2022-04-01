package database

type DatabaseClient struct {
	deployments map[string]Client
}

func New() *DatabaseClient {
	return &DatabaseClient{
		deployments: map[string]Client{},
	}
}

func (d *DatabaseClient) CreateClient(client Client) {
	key := client.Deployment + "/" + client.Database + "/" + client.Username
	d.deployments[key] = client
}

func (d *DatabaseClient) DeleteClient(client Client) {
	key := client.Deployment + "/" + client.Database + "/" + client.Username
	delete(d.deployments, key)
}

func (d *DatabaseClient) HasClient(client Client) bool {
	key := client.Deployment + "/" + client.Database + "/" + client.Username
	_, ok := d.deployments[key]

	return ok
}

func DropDeployment(name string) {}
