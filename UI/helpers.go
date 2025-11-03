package UI

func Contains(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}

func Remove(slice []string, val string) []string {
	result := []string{}
	for _, s := range slice {
		if s != val {
			result = append(result, s)
		}
	}
	return result
}
