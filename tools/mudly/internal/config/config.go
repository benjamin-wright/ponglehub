package config

type DevEnv struct {
	Compose *map[string]interface{}
}

type Pipeline struct {
	Name string
}

type Artefact struct {
	Name     string
	Pipeline Pipeline
}

type Config struct {
	DevEnv    *DevEnv
	Artefacts []Artefact
}
