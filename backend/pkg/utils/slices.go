package utils

func NonNil[T any](s []T) []T {
	if s == nil {
		return []T{}
	}
	return s
}
