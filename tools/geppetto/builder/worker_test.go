package builder

import "ponglehub.co.uk/geppetto/types"

type call struct {
	repo string
	lang string
}

func makeMockWorker() (chan call, *mockWorker) {
	channel := make(chan call, 20)
	return channel, &mockWorker{calls: channel}
}

type mockWorker struct {
	err   error
	calls chan<- call
}

func (m *mockWorker) buildNPM(repo types.Repo, signals chan<- buildSignal) {
	m.calls <- call{
		repo: repo.Name,
		lang: "npm",
	}

	signals <- buildSignal{
		repo: repo.Name,
		err:  m.err,
	}
}
