package steps

type CommandStep struct {
	Name    string `yaml:"name"`
	Command string `yaml:"cmd"`
}

func (c CommandStep) Run() bool {
	return true
}

func (c CommandStep) String() string {
	return c.Name
}
