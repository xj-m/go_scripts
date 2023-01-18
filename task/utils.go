package task

func contains[T comparable](elements []T, v T) bool {
	for _, s := range elements {
		if v == s {
			return true
		}
	}
	return false
}

func hasKey[T comparable](elements map[T]string, v T) bool {
	_, ok := elements[v]
	return ok
}
