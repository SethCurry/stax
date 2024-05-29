package common

func Map[T any, P any](items []T, fn func(T) P) []P {
	ret := make([]P, len(items))

	for k, v := range items {
		ret[k] = fn(v)
	}

	return ret
}

func Filter[T any](items []T, fn func(T) bool) []T {
	var ret []T

	for _, v := range items {
		if fn(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
