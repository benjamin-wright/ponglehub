package builder

type buildSignal struct {
	repo     string
	err      error
	skip     bool
	finished bool
	phase    string
}
