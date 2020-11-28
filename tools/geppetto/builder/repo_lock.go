package builder

type repoLock struct {
	locks []string
}

func newRepoLock() repoLock {
	return repoLock{
		locks: []string{},
	}
}

func (r *repoLock) lock(repo string) {
	for _, lock := range r.locks {
		if lock == repo {
			return
		}
	}

	r.locks = append(r.locks, repo)
}

func (r *repoLock) unlock(repo string) {
	index := -1

	for i, lock := range r.locks {
		if lock == repo {
			index = i
			break
		}
	}

	if index == -1 {
		return
	}

	length := len(r.locks)
	r.locks[index] = r.locks[length-1]
	r.locks[length-1] = ""
	r.locks = r.locks[:length-1]
}

func (r *repoLock) isLocked(repo string) bool {
	for _, lock := range r.locks {
		if lock == repo {
			return true
		}
	}

	return false
}
