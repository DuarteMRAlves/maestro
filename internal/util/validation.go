package util

func ValidateEqualElementsString(expected []string, actual []string) bool {
	if len(expected) != len(actual) {
		return false
	}
	counts := map[string]int{}
	for _, e := range expected {
		counts[e]++
	}
	for _, a := range actual {
		_, exists := counts[a]
		if !exists {
			return false
		}
		counts[a]--
	}
	for _, c := range counts {
		if c != 0 {
			return false
		}
	}
	return true
}
