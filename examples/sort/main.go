// Example for sort function
package main

import "fmt"

// Instantiate 3 different types of sorting function
//
// Regenerate the templates with "go generate"

// Sort strings using the less function
//go:generate gotemplate "github.com/ncw/gotemplate/sort" "Sort(string, less)"

// Sort floats using the lt function
//go:generate gotemplate "github.com/ncw/gotemplate/sort" "SortF(float64, lt)"

// Sort strings strings using the function passed in
//go:generate gotemplate "github.com/ncw/gotemplate/sort" "SortGt(string, func(a, b string) bool { return a > b })"

func less(a, b string) bool {
	return a < b
}

func lt(a, b float64) bool {
	return a < b
}

func main() {
	s := []string{"b", "c", "a", "e"}
	fmt.Println(s)
	Sort(s)
	fmt.Println(s)
	SortGt(s)
	fmt.Println(s)

	f := []float64{0.1, -1.6, 1.5, +9, -8, 0.01, 1E99}
	fmt.Println(f)
	SortF(f)
	fmt.Println(f)
}
