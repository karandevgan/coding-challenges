package sort

import "cmp"

func MergeSort[T cmp.Ordered](values []T) []T {
	return MergeSortWith(values, func(a T, b T) int {
		if a < b {
			return -1
		}
		if a == b {
			return 0
		}
		return 1
	})
}

func MergeSortWith[T any](values []T, cmp func(T, T) int) []T {
	mid := len(values) / 2
	if len(values) > 1 {
		l1 := MergeSortWith[T](values[:mid], cmp)
		l2 := MergeSortWith[T](values[mid:], cmp)
		return merge(l1, l2, cmp)
	} else {
		return values
	}
}

func merge[T any](l1 []T, l2 []T, cmp func(T, T) int) []T {
	rList := make([]T, len(l1)+len(l2))
	ptr1 := 0
	ptr2 := 0
	ptr3 := 0
	for ptr1 < len(l1) && ptr2 < len(l2) {
		if cmp(l1[ptr1], l2[ptr2]) <= 0 {
			rList[ptr3] = l1[ptr1]
			ptr1 += 1
		} else {
			rList[ptr3] = l2[ptr2]
			ptr2 += 1
		}
		ptr3 += 1
	}

	for ptr1 < len(l1) {
		rList[ptr3] = l1[ptr1]
		ptr1 += 1
		ptr3 += 1
	}

	for ptr2 < len(l2) {
		rList[ptr3] = l2[ptr2]
		ptr2 += 1
		ptr3 += 1
	}

	return rList
}
