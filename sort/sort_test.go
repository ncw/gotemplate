// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Modified into a gotemplate by Nick Craig-Wood <nick@craig-wood.com>

package sort

import (
	"math/rand"
	"testing"
)

var ints = [...]A{74, 59, 238, -784, 9845, 959, 905, 0, 0, 42, 7586, -5467984, 7586}

func TestSortIntSlice(t *testing.T) {
	data := ints
	a := data[0:]
	Sort(a)
	if !IsSorted(a) {
		t.Errorf("sorted %v", ints)
		t.Errorf("   got %v", data)
	}
}

func TestSortLarge_Random(t *testing.T) {
	n := 1000000
	if testing.Short() {
		n /= 100
	}
	data := make([]A, n)
	for i := 0; i < len(data); i++ {
		data[i] = A(rand.Intn(100))
	}
	if IsSorted(data) {
		t.Fatalf("terrible rand.rand")
	}
	Sort(data)
	if !IsSorted(data) {
		t.Errorf("sort didn't sort - 1M ints")
	}
}

func BenchmarkSortInt1K(b *testing.B) {
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		data := make([]A, 1<<10)
		for i := 0; i < len(data); i++ {
			data[i] = A(i ^ 0x2cc)
		}
		b.StartTimer()
		Sort(data)
		b.StopTimer()
	}
}

func BenchmarkSortInt64K(b *testing.B) {
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		data := make([]A, 1<<16)
		for i := 0; i < len(data); i++ {
			data[i] = A(i ^ 0xcccc)
		}
		b.StartTimer()
		Sort(data)
		b.StopTimer()
	}
}
