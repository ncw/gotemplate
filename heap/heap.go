// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package heap provides heap operations for a type A and a comparison
// function Less.  A heap is a tree with the property that each node
// is the minimum-valued node in its subtree.
//
// The minimum element in the tree is the root, at index 0.
//
// This provides a min-heap with the following invariants (established
// after Init has been called or if the data is empty or sorted):
//
//	!Less(h[j], h[i]) for 0 <= i < len(h) and j = 2*i+1 or 2*i+2 and j < len(h)
//
// A heap is a common way to implement a priority queue. To build a
// priority queue, use the (negative) priority as the ordering for the
// Less method, so Push adds items while Pop removes the
// highest-priority item from the queue. The Examples include such an
// implementation; the file example_pq_test.go has the complete
// source.
//
package heap

// An A is the element in the slice []A we are keeping as a heap
//
// template type Heap(A, Less)
type A int

// Less is a function to compare two As
func Less(a A, b A) bool {
	return a < b
}

// Heap stored in an slice
type Heap []A

// Init is compulsory before any of the heap operations
// can be used. Init is idempotent with respect to the heap invariants
// and may be called whenever the heap invariants may have been invalidated.
// Its complexity is O(n) where n = len(h).
//
func (h *Heap) Init() {
	// heapify
	n := len(*h)
	for i := n/2 - 1; i >= 0; i-- {
		h.down(i, n)
	}
}

// Push pushes the element x onto the heap. The complexity is
// O(log(n)) where n = len(h).
//
func (h *Heap) Push(x A) {
	*h = append(*h, x)
	h.up(len(*h) - 1)
}

// Pop removes the minimum element (according to Less) from the heap
// and returns it. The complexity is O(log(n)) where n = len(h).
// It is equivalent to h.Remove(0).
//
func (h *Heap) Pop() A {
	hs := *h
	n := len(hs) - 1
	hs[0], hs[n] = hs[n], hs[0]
	h.down(0, n)
	result := hs[n]
	*h = hs[:n]
	return result
}

// Remove removes the element at index i from the heap.
// The complexity is O(log(n)) where n = len(h).
//
func (h *Heap) Remove(i int) A {
	hs := *h
	n := len(hs) - 1
	if n != i {
		hs[i], hs[n] = hs[n], hs[i]
		h.down(i, n)
		h.up(i)
	}
	result := hs[n]
	*h = hs[:n]
	return result
}

// Fix re-establishes the heap ordering after the element at index i has changed its value.
// Changing the value of the element at index i and then calling Fix is equivalent to,
// but less expensive than, calling h.Remove(i) followed by a Push of the new value.
// The complexity is O(log(n)) where n = len(h).
func (h *Heap) Fix(i int) {
	h.down(i, len(*h))
	h.up(i)
}

func (h *Heap) up(j int) {
	hs := *h
	for {
		i := (j - 1) / 2 // parent
		if i == j || !Less(hs[j], hs[i]) {
			break
		}
		hs[i], hs[j] = hs[j], hs[i]
		j = i
	}
}

func (h *Heap) down(i, n int) {
	hs := *h
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && !Less(hs[j1], hs[j2]) {
			j = j2 // = 2*i + 2  // right child
		}
		if !Less(hs[j], hs[i]) {
			break
		}
		hs[i], hs[j] = hs[j], hs[i]
		i = j
	}
}
