package common

func POINTER[T any](val T) *T {
	return &val
}
