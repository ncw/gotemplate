// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This example demonstrates an integer heap built using the heap interface.
package heap_test

import (
	"fmt"

	"github.com/ncw/gotemplate/heap"
)

// This example inserts several ints into an IntHeap, checks the minimum,
// and removes them in order of priority.
func Example_intHeap() {
	h := &heap.Heap{2, 1, 5}
	h.Init()
	h.Push(3)
	fmt.Printf("minimum: %d\n", (*h)[0])
	for len(*h) > 0 {
		fmt.Printf("%d ", h.Pop())
	}
	// Output:
	// minimum: 1
	// 1 2 3 5
}
