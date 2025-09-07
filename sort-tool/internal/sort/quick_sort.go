package sort

import (
	"cmp"
)

func QuickSort[T cmp.Ordered](l []T) []T {
	if len(l) < 2 {
		return l
	}
	QuickSortWith(l, func(a T, b T) int {
		if a < b {
			return -1
		}
		if a == b {
			return 0
		}
		return 1
	})
	return l
}

func QuickSortWith[T any](l []T, cmp func(T, T) int) []T {
	if len(l) < 2 {
		return l
	}
	quickSort(l, 0, len(l)-1, cmp)
	return l
}

func quickSort[T any](l []T, first int, last int, cmp func(T, T) int) {
	if first < last {
		ix := partition(l, first, last, cmp)
		quickSort[T](l, first, ix-1, cmp)
		quickSort[T](l, ix+1, last, cmp)
	}
}

func selectPivot[T any](l []T, first int, last int, cmp func(T, T) int) int {
	// select median from first, middle and last elements
	middle := first + (last-first)/2

	if cmp(l[first], l[middle]) > 0 {
		l[first], l[middle] = l[middle], l[first]
	}
	if cmp(l[middle], l[last]) > 0 {
		l[middle], l[last] = l[last], l[middle]
	}
	if cmp(l[first], l[middle]) > 0 {
		l[first], l[middle] = l[middle], l[first]
	}

	// Move median to end for partitioning
	l[middle], l[last] = l[last], l[middle]
	return last

}

func partition[T any](l []T, first int, last int, cmp func(T, T) int) int {
	pv := selectPivot(l, first, last, cmp)
	pvEl := l[pv]
	swapIdx := first - 1
	for i := first; i < last; i++ {
		if cmp(l[i], pvEl) <= 0 {
			swapIdx += 1
			l[i], l[swapIdx] = l[swapIdx], l[i]
		}
	}

	// Move pivot to its final place
	l[swapIdx+1], l[last] = l[last], l[swapIdx+1]
	return swapIdx + 1
}
