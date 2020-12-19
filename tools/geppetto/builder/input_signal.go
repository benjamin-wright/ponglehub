package builder

type InputSignal struct {
	Repo       string
	Invalidate bool
	Reinstall  bool
	Nuke       bool
}
