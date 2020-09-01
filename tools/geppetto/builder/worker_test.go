package builder

import "ponglehub.co.uk/geppetto/types"

func makeMockWorker() (chan string, *mockWorker) {
	channel := make(chan string, 20)
	return channel, &mockWorker{calls: channel}
}

type mockWorker struct {
	err   error
	calls chan<- string
}

func (m *mockWorker) buildNPM(repo types.Repo, signals chan<- buildSignal) {
	m.calls <- repo.Name

	signals <- buildSignal{
		repo: repo.Name,
		err:  m.err,
	}
}
