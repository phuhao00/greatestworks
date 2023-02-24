package algorithm

import "math/rand"

// Samples returns N random unique items from collection.
func Samples[T any](collection []T, count int) []T {
	size := len(collection)

	ts := append([]T{}, collection...)

	results := []T{}

	for i := 0; i < size && i < count; i++ {
		copyLength := size - i

		index := rand.Intn(size - i)
		results = append(results, ts[index])

		// Removes element.
		// It is faster to swap with last element and remove it.
		ts[index] = ts[copyLength-1]
		ts = ts[:copyLength-1]
	}

	return results
}
