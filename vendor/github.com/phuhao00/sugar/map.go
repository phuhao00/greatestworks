package sugar

func Keys[K comparable, V any](in map[K]V) []K {
	result := make([]K, 0, len(in))

	for k := range in {
		result = append(result, k)
	}

	return result
}

func Values[K comparable, V any](in map[K]V) []V {
	result := make([]V, 0, len(in))

	for _, v := range in {
		result = append(result, v)
	}

	return result
}

func FiltrateBy[K comparable, V any](in map[K]V, filtrate func(K, V) bool) map[K]V {
	result := map[K]V{}

	for k, v := range in {
		if filtrate(k, v) {
			result[k] = v
		}
	}

	return result
}

func FiltrateByKeys[K comparable, V any](in map[K]V, keys []K) map[K]V {
	result := map[K]V{}

	for k, v := range in {
		if Contains(keys, k) {
			result[k] = v
		}
	}

	return result
}

func FiltrateByValues[K comparable, V comparable](in map[K]V, values []V) map[K]V {
	result := map[K]V{}

	for k, v := range in {
		if Contains(values, v) {
			result[k] = v
		}
	}

	return result
}

func MapToEntries[K comparable, V any](in map[K]V) []Entry[K, V] {
	result := make([]Entry[K, V], 0, len(in))

	for k, v := range in {
		result = append(result, Entry[K, V]{k, v})
	}

	return result
}

func EntriesToMap[K comparable, V any](entries []Entry[K, V]) map[K]V {
	result := map[K]V{}

	for _, entry := range entries {
		result[entry.Key] = entry.Value
	}

	return result
}

func Invert[K comparable, V comparable](in map[K]V) map[V]K {
	result := map[V]K{}

	for k, v := range in {
		result[v] = k
	}

	return result
}

// Assign merges multiple maps from left to right.
func Assign[K comparable, V any](maps ...map[K]V) map[K]V {
	result := map[K]V{}

	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}

	return result
}

// MapUpdateKeys manipulates a map keys and transforms it to a map of another type.
func MapUpdateKeys[K comparable, V any, R comparable](in map[K]V, iteratee func(K, V) R) map[R]V {
	result := map[R]V{}

	for k, v := range in {
		result[iteratee(k, v)] = v
	}

	return result
}

// MapUpdateValues manipulates a map values and transforms it to a map of another type.
func MapUpdateValues[K comparable, V any, R any](in map[K]V, iteratee func(K, V) R) map[K]R {
	result := map[K]R{}

	for k, v := range in {
		result[k] = iteratee(k, v)
	}

	return result
}
