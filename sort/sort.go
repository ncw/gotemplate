// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Modified into a gotemplate by Nick Craig-Wood <nick@craig-wood.com>

// Package sort provides primitives for sorting slices of A.
package sort

// An A is the element in the slice []A we are sorting
//
// template type Sort(A, Less)
type A int

// Less is a function to compare two As
func Less(a A, b A) bool {
	return a < b
}

func swap(data []A, i, j int) {
	data[i], data[j] = data[j], data[i]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Insertion sort
func insertionSort(data []A, a, b int) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && Less(data[j], data[j-1]); j-- {
			swap(data, j, j-1)
		}
	}
}

// siftDown implements the heap property on data[lo, hi).
// first is an offset into the array where the root of the heap lies.
func siftDown(data []A, lo, hi, first int) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}
		if child+1 < hi && Less(data[first+child], data[first+child+1]) {
			child++
		}
		if !Less(data[first+root], data[first+child]) {
			return
		}
		swap(data, first+root, first+child)
		root = child
	}
}

func heapSort(data []A, a, b int) {
	first := a
	lo := 0
	hi := b - a

	// Build heap with greatest element at top.
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDown(data, i, hi, first)
	}

	// Pop elements, largest first, into end of data.
	for i := hi - 1; i >= 0; i-- {
		swap(data, first, first+i)
		siftDown(data, lo, i, first)
	}
}

// Quicksort, following Bentley and McIlroy,
// ``Engineering a Sort Function,'' SP&E November 1993.

// medianOfThree moves the median of the three values data[a], data[b], data[c] into data[a].
func medianOfThree(data []A, a, b, c int) {
	m0 := b
	m1 := a
	m2 := c
	// bubble sort on 3 elements
	if Less(data[m1], data[m0]) {
		swap(data, m1, m0)
	}
	if Less(data[m2], data[m1]) {
		swap(data, m2, m1)
	}
	if Less(data[m1], data[m0]) {
		swap(data, m1, m0)
	}
	// now data[m0] <= data[m1] <= data[m2]
}

func swapRange(data []A, a, b, n int) {
	for i := 0; i < n; i++ {
		swap(data, a+i, b+i)
	}
}

func doPivot(data []A, lo, hi int) (midlo, midhi int) {
	m := lo + (hi-lo)/2 // Written like this to avoid integer overflow.
	if hi-lo > 40 {
		// Tukey's ``Ninther,'' median of three medians of three.
		s := (hi - lo) / 8
		medianOfThree(data, lo, lo+s, lo+2*s)
		medianOfThree(data, m, m-s, m+s)
		medianOfThree(data, hi-1, hi-1-s, hi-1-2*s)
	}
	medianOfThree(data, lo, m, hi-1)

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
			if Less(data[b], data[pivot]) { // data[b] < pivot
				b++
			} else if !Less(data[pivot], data[b]) { // data[b] = pivot
				swap(data, a, b)
				a++
				b++
			} else {
				break
			}
		}
		for b < c {
			if Less(data[pivot], data[c-1]) { // data[c-1] > pivot
				c--
			} else if !Less(data[c-1], data[pivot]) { // data[c-1] = pivot
				swap(data, c-1, d-1)
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
		swap(data, b, c-1)
		b++
		c--
	}

	n := min(b-a, a-lo)
	swapRange(data, lo, b-n, n)

	n = min(hi-d, d-c)
	swapRange(data, c, hi-n, n)

	return lo + b - a, hi - (d - c)
}

func quickSort(data []A, a, b, maxDepth int) {
	for b-a > 7 {
		if maxDepth == 0 {
			heapSort(data, a, b)
			return
		}
		maxDepth--
		mlo, mhi := doPivot(data, a, b)
		// Avoiding recursion on the larger subproblem guarantees
		// a stack depth of at most lg(b-a).
		if mlo-a < b-mhi {
			quickSort(data, a, mlo, maxDepth)
			a = mhi // i.e., quickSort(data, mhi, b)
		} else {
			quickSort(data, mhi, b, maxDepth)
			b = mlo // i.e., quickSort(data, a, mlo)
		}
	}
	if b-a > 1 {
		insertionSort(data, a, b)
	}
}

// Sort sorts data.
// It makes one call to data.Len to determine n, and O(n*log(n)) calls to
// data.Less and data.swap. The sort is not guaranteed to be stable.
func Sort(data []A) {
	// Switch to heapsort if depth of 2*ceil(lg(n+1)) is reached.
	n := len(data)
	maxDepth := 0
	for i := n; i > 0; i >>= 1 {
		maxDepth++
	}
	maxDepth *= 2
	quickSort(data, 0, n, maxDepth)
}

// IsSorted reports whether data is sorted.
func IsSorted(data []A) bool {
	n := len(data)
	for i := n - 1; i > 0; i-- {
		if Less(data[i], data[i-1]) {
			return false
		}
	}
	return true
}
