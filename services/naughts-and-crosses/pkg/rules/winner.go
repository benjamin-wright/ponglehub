package rules

func IsWinner(marks string, position int) bool {
	markRunes := []rune(marks)

	return isHorizontalWin(markRunes, position) ||
		isVerticalWin(markRunes, position) ||
		isDiagonalWin(markRunes, position) ||
		isBackDiagonalWin(markRunes, position)
}

func isHorizontalWin(marks []rune, position int) bool {
	row := position % 3

	if marks[row] == '-' {
		return false
	}

	for _, column := range []int{1, 2} {
		if marks[column*3+row] != marks[row] {
			return false
		}
	}

	return true
}

func isVerticalWin(marks []rune, position int) bool {
	column := (position - position%3) / 3

	if marks[column*3] == '-' {
		return false
	}

	for _, row := range []int{1, 2} {
		if marks[column*3+row] != marks[column*3] {
			return false
		}
	}

	return true
}

func isDiagonalWin(marks []rune, position int) bool {
	if position%2 != 0 {
		return false
	}

	if marks[0] == '-' {
		return false
	}

	for _, idx := range []int{1, 2} {
		if marks[0] != marks[4*idx] {
			return false
		}
	}

	return true
}

func isBackDiagonalWin(marks []rune, position int) bool {
	if position%2 != 0 {
		return false
	}

	if marks[2] == '-' {
		return false
	}

	for _, idx := range []int{1, 2} {
		if marks[2] != marks[2+2*idx] {
			return false
		}
	}

	return true
}
