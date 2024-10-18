package sliceutil

import "golang.org/x/exp/slices"

func FindElement[T any](slice []T, predicate func(T) bool) (T, bool) {
	// Find index using slices.IndexFunc (can also use slices.ContainsFunc but doesn't return element)
	index := slices.IndexFunc(slice, predicate)
	if index != -1 {
		return slice[index], true
	}
	var zeroValue T
	return zeroValue, false
}
