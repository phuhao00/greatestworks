package sugar

func IndexOf[T comparable](collection []T, predicate func(T) bool) int {
	for i, t := range collection {
		if predicate(t) {
			return i
		}
	}

	return -1
}

func LastIndexOf[T comparable](collection []T, predicate func(T) bool) int {
	l := len(collection)

	for i := l - 1; i >= 0; i-- {
		if predicate(collection[i]) {
			return i
		}
	}
	
	return -1
}
