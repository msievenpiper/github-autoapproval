package internal

func Filter[T any](items []T, test func(T) bool) (ret []T) {
	for _, el := range items {
		if test(el) {
			ret = append(ret, el)
		}
	}

	return ret
}

func Map[T any, E any](items []T, mapFunc func(T) E) []E {
	var ret []E
	for _, el := range items {
		ret = append(ret, mapFunc(el))
	}

	return ret
}
