package bot

import "strconv"

func countSelectedOptions(selectedOptions map[int][]string) map[int]int {
	counts := make(map[int]int)

	// Iterate through selected options and count occurrences
	for _, options := range selectedOptions {
		for _, option := range options {
			optionNumber, err := strconv.Atoi(option)
			if err != nil {
				continue
			}
			// Increment the count for the selected option
			counts[optionNumber]++
		}
	}

	return counts
}
func getResultDescription(counts map[int]int) string {
	maxCount := 0
	maxOption := 0

	// Determine which option has the highest count
	for option, count := range counts {
		if count > maxCount {
			maxCount = count
			maxOption = option
		}
	}

	// Return the description based on the option with the highest count
	switch maxOption {
	case 1:
		return Description1
	case 2:
		return Description2
	case 3:
		return Description3
	case 4:
		return Description4
	case 5:
		return Description5
	case 6:
		return Description6
	case 7:
		return Description7
	case 8:
		return Description8
	default:
		return "Нет выбранных вариантов."
	}
}

func remove(slice []string, item string) []string {
	for i, a := range slice {
		if a == item {
			// Удаляем элемент из среза
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}
