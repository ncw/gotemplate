// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Modified into a gotemplate by Nick Craig-Wood <nick@craig-wood.com>

// Package sort provides primitives for sorting slices of A.
package main

// A is the element in the slice []A we are sorting
// Less is a function to compare two As
//
// template type Sort(A, Less)

func swapSortGt(data []string, i, j int) {
	data[i], data[j] = data[j], data[i]
}

func minSortGt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Insertion sort
func insertionSortGt(data []string, a, b int) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && func(a, b string) bool {
			return a > b
		}(data[j], data[j-1]); j-- {
			swapSortGt(data, j, j-1)
		}
	}
}

// siftDown implements the heap property on data[lo, hi).
// first is an offset into the array where the root of the heap lies.
func siftDownSortGt(data []string, lo, hi, first int) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}
		if child+1 < hi && func(a, b string) bool {
			return a > b
		}(data[first+child], data[first+child+1]) {
			child++
		}
		if !func(a, b string) bool {
			return a > b
		}(data[first+root], data[first+child]) {
			return
		}
		swapSortGt(data, first+root, first+child)
		root = child
	}
}

func heapSortGt(data []string, a, b int) {
	first := a
	lo := 0
	hi := b - a

	// Build heap with greatest element at top.
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDownSortGt(data, i, hi, first)
	}

	// Pop elements, largest first, into end of data.
	for i := hi - 1; i >= 0; i-- {
		swapSortGt(data, first, first+i)
		siftDownSortGt(data, lo, i, first)
	}
}

// Quicksort, following Bentley and McIlroy,
// ``Engineering a Sort Function,'' SP&E November 1993.

// medianOfThree moves the median of the three values data[a], data[b], data[c] into data[a].
func medianOfThreeSortGt(data []string, a, b, c int) {
	m0 := b
	m1 := a
	m2 := c
	// bubble sort on 3 elements
	if func(a, b string) bool {
		return a > b
	}(data[m1], data[m0]) {
		swapSortGt(data, m1, m0)
	}
	if func(a, b string) bool {
		return a > b
	}(data[m2], data[m1]) {
		swapSortGt(data, m2, m1)
	}
	if func(a, b string) bool {
		return a > b
	}(data[m1], data[m0]) {
		swapSortGt(data, m1, m0)
	}
	// now data[m0] <= data[m1] <= data[m2]
}

func swapRangeSortGt(data []string, a, b, n int) {
	for i := 0; i < n; i++ {
		swapSortGt(data, a+i, b+i)
	}
}

func doPivotSortGt(data []string, lo, hi int) (midlo, midhi int) {
	m := lo + (hi-lo)/2 // Written like this to avoid integer overflow.
	if hi-lo > 40 {
		// Tukey's ``Ninther,'' median of three medians of three.
		s := (hi - lo) / 8
		medianOfThreeSortGt(data, lo, lo+s, lo+2*s)
		medianOfThreeSortGt(data, m, m-s, m+s)
		medianOfThreeSortGt(data, hi-1, hi-1-s, hi-1-2*s)
	}
	medianOfThreeSortGt(data, lo, m, hi-1)

	// Invariants are:
	//	data[lo] = pivot (set up by ChoosePivot)
	//	data[lo <= i < a] = pivot
	//	data[a <= i < b] < pivot
	//	data[b <= i < c] is unexamined
	//	data[c <= i < d] > pivot
	//	data[d <= i < hi] = pivot
	//
	// Once b meets c, can swap the "= pivot" sections
	// into the middle of the slice.
	pivot := lo
	a, b, c, d := lo+1, lo+1, hi, hi
	for {
		for b < c {
			if func(a, b string) bool {
				return a > b
			}(data[b], data[pivot]) { // data[b] < pivot
				b++
			} else if !func(a, b string) bool {
				return a > b
			}(data[pivot], data[b]) { // data[b] = pivot
				swapSortGt(data, a, b)
				a++
				b++
			} else {
				break
			}
		}
		for b < c {
			if func(a, b string) bool {
				return a > b
			}(data[pivot], data[c-1]) { // data[c-1] > pivot
				c--
			} else if !func(a, b string) bool {
				return a > b
			}(data[c-1], data[pivot]) { // data[c-1] = pivot
				swapSortGt(data, c-1, d-1)
				c--
				d--
			} else {
				break
			}
		}
		if b >= c {
			break
		}
		// data[b] > pivot; data[c-1] < pivot
		swapSortGt(data, b, c-1)
		b++
		c--
	}

	n := minSortGt(b-a, a-lo)
	swapRangeSortGt(data, lo, b-n, n)

	n = minSortGt(hi-d, d-c)
	swapRangeSortGt(data, c, hi-n, n)

	return lo + b - a, hi - (d - c)
}

func quickSortGt(data []string, a, b, maxDepth int) {
	for b-a > 7 {
		if maxDepth == 0 {
			heapSortGt(data, a, b)
			return
		}
		maxDepth--
		mlo, mhi := doPivotSortGt(data, a, b)
		// Avoiding recursion on the larger subproblem guarantees
		// a stack depth of at most lg(b-a).
		if mlo-a < b-mhi {
			quickSortGt(data, a, mlo, maxDepth)
			a = mhi // i.e., quickSort(data, mhi, b)
		} else {
			quickSortGt(data, mhi, b, maxDepth)
			b = mlo // i.e., quickSort(data, a, mlo)
		}
	}
	if b-a > 1 {
		insertionSortGt(data, a, b)
	}
}

// Sort sorts data.
// It makes one call to data.Len to determine n, and O(n*log(n)) calls to
// data.Less and data.swap. The sort is not guaranteed to be stable.
func SortGt(data []string) {
	// Switch to heapsort if depth of 2*ceil(lg(n+1)) is reached.
	n := len(data)
	maxDepth := 0
	for i := n; i > 0; i >>= 1 {
		maxDepth++
	}
	maxDepth *= 2
	quickSortGt(data, 0, n, maxDepth)
}

// IsSorted reports whether data is sorted.
func IsSortGted(data []string) bool {
	n := len(data)
	for i := n - 1; i > 0; i-- {
		if func(a, b string) bool {
			return a > b
		}(data[i], data[i-1]) {
			return false
		}
	}
	return true
}
