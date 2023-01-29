package sugar

import (
	"golang.org/x/exp/constraints"
)

//Clamp  clamps number within the inclusive lower and upper bounds.
func Clamp[T constraints.Ordered](value, min, max T) T {
	if value < min {
		value = min
	} else if value > max {
		value = max
	}
	return value
}
