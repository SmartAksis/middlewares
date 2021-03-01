package array_utils

func StringArrayContains(array []string, contain string) bool {
	for _, a := range array {
		if a == contain {
			return true
		}
	}
	return false
}
