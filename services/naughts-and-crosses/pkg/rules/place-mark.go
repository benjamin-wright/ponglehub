package rules

import "fmt"

func PlaceMark(marks string, position int, turn int16) string {
	markRunes := []rune(marks)

	switch turn {
	case 0:
		markRunes[position] = '0'
	case 1:
		markRunes[position] = '1'
	default:
		panic(fmt.Sprintf("Invalid turn: %d", turn))
	}

	return string(markRunes)
}

func NextTurn(turn int16) int16 {
	switch turn {
	case 0:
		return 1
	case 1:
		return 0
	default:
		panic(fmt.Sprintf("Invalid turn: %d", turn))
	}
}
