package compare

import (
	"golang.org/x/exp/constraints"
)

func Max[T constraints.Ordered](x ...T) (max T) {
	if len(x) == 0 {
		return
	}
	max = x[0]
	for _, v := range x {
		if v > max {
			max = v
		}
	}
	return
}

func Min[T constraints.Ordered](x ...T) (min T) {
	if len(x) == 0 {
		return
	}
	min = x[0]
	for _, v := range x {
		if v < min {
			min = v
		}
	}
	return
}

func Contains[T comparable](want T, x ...T) bool {
	for _, v := range x {
		if v == want {
			return true
		}
	}
	return false
}
