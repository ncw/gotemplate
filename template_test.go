// Tests for template

package main

import (
	"bytes"
	"go/build"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"testing"
)

type TestTemplate struct {
	title   string
	args    string
	pkg     string
	in      string
	outName string
	out     string
}

var tests = []TestTemplate{
	{
		title: "Simple test public",
		args:  "Min(int)",
		pkg:   "main",
		in: `package tt

import "fmt"

// template type TT(A)
type A int

func TT(a A) A { return A(0) }
func TTone(a A) A { return A(1) }
func AFunc() {}
func (a A) f0() {}
func (a *A) F1() {}
var AVar int
`,
		outName: "gotemplate_Min.go",
		out: `package main

import "fmt"

// template type TT(A)

func Min(a int) int    { return int(0) }
func Minone(a int) int { return int(1) }
func AFuncMin()        {}
func (a int) f0()      {}
func (a *int) F1()     {}

var AVarMin int
`,
	},
	{
		title: "Simple test private",
		args:  "min(float64)",
		pkg:   "main",
		in: `package tt

// template type TT(A)
type A int

func TT(a A) A { return A(0) }
func TTone(a A) A { return A(1) }
func AFunc() {}
var AVar int
`,
		outName: "gotemplate_min.go",
		out: `package main

// template type TT(A)

func min(a float64) float64    { return float64(0) }
func minone(a float64) float64 { return float64(1) }
func aFuncMin()                {}

var aVarMin int
`,
	},
	{
		title: "Test function",
		args:  "Min(int8, func(a int8, b int8) bool { return a < b })",
		pkg:   "main",
		in: `package tt

// template type TT(A, Less)
type A int
func Less(a, b A) { return a < b }

func TT(a, b A) A { return Less(a, b) }
func TTone(a A) A { return !Less(a, b) }
`,
		outName: "gotemplate_Min.go",
		out: `package main

// template type TT(A, Less)

func Min(a, b int8) int8 {
	return func(a int8, b int8) bool {
		return a < b
	}(a, b)
}
func Minone(a int8) int8 {
	return !func(a int8, b int8) bool {
		return a < b
	}(a, b)
}
`,
	},
}

func testTemplate(t *testing.T, test *TestTemplate) {
	// Disable logging
	log.SetOutput(ioutil.Discard)

	// Make temporary directory
	dir, err := ioutil.TempDir("", "gotemplate_test")
	if err != nil {
		t.Fatalf("Failed to make temp dir: %v", err)
	}
	defer func() {
		err := os.RemoveAll(dir)
		if err != nil {
			t.Logf("Failed to remove temp dir: %v", err)
		}
	}()

	// Make subdirectories
	src := path.Join(dir, "src")
	err = os.Mkdir(src, 0700)
	if err != nil {
		t.Fatalf("Failed to make dir %q: %v", src, err)
	}
	input := path.Join(src, "input")
	err = os.Mkdir(input, 0700)
	if err != nil {
		t.Fatalf("Failed to make dir %q: %v", input, err)
	}
	output := path.Join(src, "output")
	err = os.Mkdir(output, 0700)
	if err != nil {
		t.Fatalf("Failed to make dir %q: %v", output, err)
	}

	// Change directory to output directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to read cwd: %v", err)
	}
	err = os.Chdir(output)
	if err != nil {
		t.Fatalf("Failed to cd %q dir: %v", output, err)
	}
	defer func() {
		err := os.Chdir(cwd)
		if err != nil {
			t.Logf("Failed to change back to cwd: %v", err)
		}
	}()

	// Set GOPATH to directory
	build.Default.GOPATH = dir

	// Write template input
	tmpl := path.Join(input, "main.go")
	err = ioutil.WriteFile(tmpl, []byte(test.in), 0600)
	if err != nil {
		t.Fatalf("Failed to write %q: %v", tmpl, err)
	}

	// Write main.go for output
	main := path.Join(output, "main.go")
	err = ioutil.WriteFile(main, []byte("package main"), 0600)
	if err != nil {
		t.Fatalf("Failed to write %q: %v", main, err)
	}

	// Instantiate template
	template := newTemplate(output, "input", test.args)
	template.instantiate()

	// Check output
	expectedFile := path.Join(output, test.outName)
	actualBytes, err := ioutil.ReadFile(expectedFile)
	if err != nil {
		t.Fatalf("Failed to read %q: %v", expectedFile, err)
	}
	actual := string(actualBytes)
	if actual != test.out {
		t.Errorf(`Output is wrong
Got
-------------
%s
-------------
Expected
-------------
%s
-------------
`, actual, test.out)
		actualFile := expectedFile + ".actual"
		err = ioutil.WriteFile(actualFile, []byte(test.out), 0600)
		if err != nil {
			t.Fatalf("Failed to write %q: %v", actualFile, err)
		}
		cmd := exec.Command("diff", "-u", actualFile, expectedFile)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		_ = cmd.Run()
		t.Errorf("Diff\n----\n%s", out.String())
	}

}

func TestSub(t *testing.T) {
	for i := range tests {
		t.Logf("Test[%d] %q", i, tests[i].title)
		testTemplate(t, &tests[i])
	}
}
