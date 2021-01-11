package utils

func MapHasAllTrue(a map[string]bool) bool {
	if len(a) == 0 {
		return false
	}

	for _, value := range a {
		if !value {
			return false
		}
	}

	return true
}
