package builder

import "ponglehub.co.uk/geppetto/types"

type mockWorker struct {
	err error
}

func (m *mockWorker) buildNPM(repo types.Repo, signals chan<- buildSignal) {
	signals <- buildSignal{
		repo: repo.Name,
		err:  m.err,
	}
}
