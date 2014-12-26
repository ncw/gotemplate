package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"testing"
)

type TestReplacement struct {
	title    string
	id       string
	val      ast.Expr
	source   string
	expected string
}

func parseExpr(s string) ast.Expr {
	exp, err := parser.ParseExpr(s)
	if err != nil {
		fatalf("Cannot parse expression %s :%s", s, err.Error())
	}
	return exp
}

var replaceTests = []TestReplacement{
	{
		title: "basic test",
		id:    "A",
		val:   parseExpr("int"),
		source: `package tt

func Add(a, b A) A {
	var sum A = a + b
	return sum
}
`,
		expected: `package tt

func Add(a, b int) int {
	var sum int = a + b
	return sum
}
`,
	},
}

func testReplaceIdent(t *testing.T, tr TestReplacement) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "t.go", tr.source, parser.Mode(0))
	if err != nil {
		fatalf("Failed to parse source: %s", err.Error())
	}

	file = rewriteFile(file, tr.id, tr.val)

	buf := new(bytes.Buffer)
	err = format.Node(buf, fset, file)
	if err != nil {
		fatalf("Failed to format post-replace source: %v", err)
	}

	actual := buf.String()
	if actual != tr.expected {
		t.Errorf(`Output is wrong
Got
-------------
%s
-------------
Expected
-------------
%s
-------------
`, actual, tr.expected)
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

		dir = path.Join(dir, "src")
		err = os.MkdirAll(dir, 0700)
		if err != nil {
			t.Fatalf("Failed to create directory %q: %v", dir, err)
		}
		expectedFile := path.Join(dir, "out.go")
		err = ioutil.WriteFile(expectedFile, []byte(tr.expected), 0600)
		if err != nil {
			t.Fatalf("Failed to write %q: %v", expectedFile, err)
		}
		actualFile := expectedFile + ".actual"
		err = ioutil.WriteFile(actualFile, []byte(actual), 0600)
		if err != nil {
			t.Fatalf("Failed to write %q: %v", actualFile, err)
		}
		cmd := exec.Command("diff", "-u", expectedFile, actualFile)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		_ = cmd.Run()
		t.Errorf("Diff\n----\n%s", out.String())
	}
}

func TestReplacements(t *testing.T) {
	fatalf = func(format string, args ...interface{}) {
		t.Fatalf(format, args...)
	}
	for i, tr := range replaceTests {
		t.Logf("Test[%d] %s", i, tr.title)
		testReplaceIdent(t, tr)
	}
}
