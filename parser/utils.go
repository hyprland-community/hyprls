package parser

func not[T any](f func(T) bool) func(T) bool {
	return func(x T) bool {
		return !f(x)
	}
}
