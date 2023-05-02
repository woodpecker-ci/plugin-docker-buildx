package utils

func Map[T any](in []T, fn func(T) T) []T {
	out := in
	for i := range in {
		out[i] = fn(out[i])
	}
	return out
}
