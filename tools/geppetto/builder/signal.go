package builder

import (
	"context"
)

type signal struct {
	repo      string
	err       error
	cancelled bool
	skip      bool
	finished  bool
	phase     string
}

func makeErrorSignal(ctx context.Context, repo string, err error) signal {
	if ctx.Err() != nil {
		return signal{repo: repo, cancelled: true}
	}

	return signal{repo: repo, err: err}
}
