package utils

func FindFirst[T any](a []T, f func(T) bool) (_ T, ok bool) {
	for _, i := range a {
		if f(i) {
			return i, true
		}
	}
	return
}
