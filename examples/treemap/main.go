package main

import "fmt"

//go:generate gotemplate "github.com/ncw/gotemplate/treemap" "intStringTreeMap(int, string)"

func less(x, y int) bool { return x < y }

func main() {
	tr := newIntStringTreeMap(less)
	tr.Set(0, "Hello")
	tr.Set(1, "World")

	for it := tr.Iterator(); it.Valid(); it.Next() {
		fmt.Println(it.Key(), it.Value())
	}
}
