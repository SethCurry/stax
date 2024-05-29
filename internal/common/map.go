package common

func Map[T any, P any](items []T, fn func(T) P) []P {
	ret := make([]P, len(items))

	for k, v := range items {
		ret[k] = fn(v)
	}

	return ret
}
