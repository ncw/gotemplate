package treemap

import (
	"fmt"
)

func ExampleTreeMap_Set() {
	tr := New(less)
	tr.Set(0, "hello")
	v, _ := tr.Get(0)
	fmt.Println(v)
	// Output:
	// hello
}

func ExampleTreeMap_Del() {
	tr := New(less)
	tr.Set(0, "hello")
	tr.Del(0)
	fmt.Println(tr.Contains(0))
	// Output:
	// false
}

func ExampleTreeMap_Get() {
	tr := New(less)
	tr.Set(0, "hello")
	v, _ := tr.Get(0)
	fmt.Println(v)
	// Output:
	// hello
}

func ExampleTreeMap_Contains() {
	tr := New(less)
	tr.Set(0, "hello")
	fmt.Println(tr.Contains(0))
	// Output:
	// true
}

func ExampleTreeMap_Len() {
	tr := New(less)
	tr.Set(0, "hello")
	tr.Set(1, "world")
	fmt.Println(tr.Len())
	// Output:
	// 2
}

func ExampleTreeMap_Clear() {
	tr := New(less)
	tr.Set(0, "hello")
	tr.Set(1, "world")
	tr.Clear()
	fmt.Println(tr.Len())
	// Output:
	// 0
}

func ExampleTreeMap_Iterator() {
	tr := New(less)
	tr.Set(1, "one")
	tr.Set(2, "two")
	tr.Set(3, "three")
	for it := tr.Iterator(); it.Valid(); it.Next() {
		fmt.Println(it.Key(), "-", it.Value())
	}
	// Output:
	// 1 - one
	// 2 - two
	// 3 - three
}

func ExampleTreeMap_Reverse() {
	tr := New(less)
	tr.Set(1, "one")
	tr.Set(2, "two")
	tr.Set(3, "three")
	for it := tr.Reverse(); it.Valid(); it.Next() {
		fmt.Println(it.Key(), "-", it.Value())
	}
	// Output:
	// 3 - three
	// 2 - two
	// 1 - one
}

func ExampleTreeMap_Range() {
	tr := New(less)
	tr.Set(1, "one")
	tr.Set(2, "two")
	tr.Set(3, "three")
	for it, end := tr.Range(1, 2); it != end; it.Next() {
		fmt.Println(it.Key(), "-", it.Value())
	}
	// Output:
	// 1 - one
	// 2 - two
}
