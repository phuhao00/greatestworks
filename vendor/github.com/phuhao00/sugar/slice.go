package sugar

import "golang.org/x/exp/constraints"

func SliceFiltrate[V any](collection []V, filtrate func(V, int) bool) []V {

	result := []V{}

	for i, v := range collection {
		if filtrate(v, i) {
			result = append(result, v)
		}
	}

	return result
}

func SliceUpdateElement[T any, R any](collection []T, iteratee func(T, int) R) []R {
	result := make([]R, len(collection))

	for i, t := range collection {
		result[i] = iteratee(t, i)
	}

	return result
}

func SliceUniq[T any, U comparable](collection []T, iteratee func(T) U) []T {
	result := make([]T, len(collection))

	seen := make(map[U]struct{}, len(collection))
	for _, item := range collection {
		key := iteratee(item)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
	}

	return result
}

func SliceGroupBy[T any, U comparable](collection []T, iteratee func(T) U) map[U][]T {
	result := map[U][]T{}

	for _, item := range collection {
		key := iteratee(item)

		result[key] = append(result[key], item)
	}
	return result
}

//CheckInSlice  check value in slice
func CheckInSlice[T constraints.Ordered](a T, s []T) bool {
	for _, val := range s {
		if a == val {
			return true
		}
	}
	return false
}

//DelOneInSlice  delete one element of slice  left->right
func DelOneInSlice[T constraints.Ordered](a T, old []T) (new []T) {
	for i, val := range old {
		if a == val {
			new = append(old[:i], old[i+1:]...)
			return
		}
	}
	return old
}
