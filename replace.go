// Copyright 2009 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This code is copied from the gofmt source at
// http://golang.org/src/cmd/gofmt/rewrite.go
//
// Changes from the original:
// rewriteFile creates a FileSet instead of referencing a global.

package main

import (
	"go/ast"
	"go/token"
	"reflect"
)

// rewriteFile replaces identifier with replacement in p returning a new *ast.File
func rewriteFile(p *ast.File, identifier string, replacement ast.Expr) *ast.File {
	cmap := ast.NewCommentMap(token.NewFileSet(), p, p.Comments)
	repl := reflect.ValueOf(replacement)

	var rewriteVal func(val reflect.Value) reflect.Value
	rewriteVal = func(val reflect.Value) reflect.Value {
		// don't bother if val is invalid to start with
		if !val.IsValid() {
			return val
		}
		val = apply(rewriteVal, val)
		// If val is a matching identifier, replace it
		if ident, ok := val.Interface().(*ast.Ident); ok && !val.IsNil() && ident.Name == identifier {
			return subst(repl, reflect.ValueOf(ident.Pos()))
		}
		return val
	}

	r := apply(rewriteVal, reflect.ValueOf(p)).Interface().(*ast.File)
	r.Comments = cmap.Filter(r).Comments() // recreate comments list
	return r
}

// set is a wrapper for x.Set(y); it protects the caller from panics if x cannot be changed to y.
func set(x, y reflect.Value) {
	// don't bother if x cannot be set or y is invalid
	if !x.CanSet() || !y.IsValid() {
		return
	}
	defer func() {
		if e := recover(); e != nil {
			fatalf("Failure while setting value %v to %v: %v", x, y, e)
		}
	}()
	x.Set(y)
}

// Values/types for special cases.
var (
	objectPtrNil = reflect.ValueOf((*ast.Object)(nil))
	scopePtrNil  = reflect.ValueOf((*ast.Scope)(nil))

	objectPtrType = reflect.TypeOf((*ast.Object)(nil))
	positionType  = reflect.TypeOf(token.NoPos)
	scopePtrType  = reflect.TypeOf((*ast.Scope)(nil))
)

// apply replaces each AST field x in val with f(x), returning val.
// To avoid extra conversions, f operates on the reflect.Value form.
func apply(f func(reflect.Value) reflect.Value, val reflect.Value) reflect.Value {
	if !val.IsValid() {
		return reflect.Value{}
	}

	// *ast.Objects introduce cycles and are likely incorrect after
	// rewrite; don't follow them but replace with nil instead
	if val.Type() == objectPtrType {
		return objectPtrNil
	}

	// similarly for scopes: they are likely incorrect after a rewrite;
	// replace them with nil
	if val.Type() == scopePtrType {
		return scopePtrNil
	}

	switch v := reflect.Indirect(val); v.Kind() {
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			e := v.Index(i)
			set(e, f(e))
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			e := v.Field(i)
			set(e, f(e))
		}
	case reflect.Interface:
		e := v.Elem()
		set(v, f(e))
	}
	return val
}

// subst returns a copy of pattern with pos used as the position of
// tokens from the pattern.
func subst(pattern reflect.Value, pos reflect.Value) reflect.Value {
	if !pattern.IsValid() {
		return reflect.Value{}
	}

	// *ast.Objects introduce cycles and are likely incorrect after
	// rewrite; don't follow them but replace with nil instead
	if pattern.Type() == objectPtrType {
		return objectPtrNil
	}

	// similarly for scopes: they are likely incorrect after a rewrite;
	// replace them with nil
	if pattern.Type() == scopePtrType {
		return scopePtrNil
	}

	if pos.IsValid() && pattern.Type() == positionType {
		// use new position only if old position was valid in the first place
		if old := pattern.Interface().(token.Pos); !old.IsValid() {
			return pattern
		}
		return pos
	}

	// Otherwise copy.
	switch p := pattern; p.Kind() {
	case reflect.Slice:
		v := reflect.MakeSlice(p.Type(), p.Len(), p.Len())
		for i := 0; i < p.Len(); i++ {
			v.Index(i).Set(subst(p.Index(i), pos))
		}
		return v

	case reflect.Struct:
		v := reflect.New(p.Type()).Elem()
		for i := 0; i < p.NumField(); i++ {
			v.Field(i).Set(subst(p.Field(i), pos))
		}
		return v

	case reflect.Ptr:
		v := reflect.New(p.Type()).Elem()
		if elem := p.Elem(); elem.IsValid() {
			v.Set(subst(elem, pos).Addr())
		}
		return v

	case reflect.Interface:
		v := reflect.New(p.Type()).Elem()
		if elem := p.Elem(); elem.IsValid() {
			v.Set(subst(elem, pos))
		}
		return v
	}

	return pattern
}
