package rules

func IsTie(marks string) bool {
	for _, mark := range []rune(marks) {
		if mark == '-' {
			return false
		}
	}

	return true
}
