package builder

type signal struct {
	repo     string
	err      error
	skip     bool
	finished bool
	phase    string
}
