Go templates
============

This tool manages package based templates for the Go language using
"go generate" which requires go 1.4.

[![Build Status](https://travis-ci.org/ncw/gotemplate.png)](https://travis-ci.org/ncw/gotemplate)

Install
-------

Install using go get

    go get github.com/ncw/gotemplate/...

and this will build the `gotemplate` binary in `$GOPATH/bin`.

It will also pull in a set of templates you can start using straight away

  * [set](http://godoc.org/github.com/ncw/gotemplate/set)
  * [list](http://godoc.org/github.com/ncw/gotemplate/list)
  * [sort](http://godoc.org/github.com/ncw/gotemplate/sort)
  * [heap](http://godoc.org/github.com/ncw/gotemplate/heap)
  * [treemap](http://godoc.org/github.com/ncw/gotemplate/treemap)

Using templates
---------------

To use a template, first you must tell `gotemplate` that you want to
use it using a special comment in your code.  For example

    //go:generate gotemplate "github.com/ncw/gotemplate/set" mySet(string)

This tells `go generate` to run `gotemplate` and that you want to use
the set template with the `string` type parameter and with the local
name `mySet`.

Now run `go generate` in your code directory with no arguments.  This
will instantiate the template into a file called `gotemplate_mySet.go`
which will provide a `mySet` type and `newMySet` and `newSizedMySet`
functions to make them. Note that the first letter of your custom name 
is still capitalized when it is not at the beginning of the new name.

    $ go generate
    substituting "github.com/ncw/gotemplate/set" with mySet(string) into package main
    Written 'gotemplate_mySet.go'

If you wish to change what the output file names look like then you
can use the `-outfmt format` flag.  The format must contain a single
instance of the `%v` verb which will be replaced with the template
instance name (default "gotemplate_%v")

Instantiating the templates into your project gives them the ability
to use internal types from your project.

If you use an initial capital when you name your template
instantiation then any external functions will be public.  Eg

    //go:generate gotemplate "github.com/ncw/gotemplate/set" MySet(string)

Would give you `MySet`, `NewMySet` and `NewSizedMySet` instead.

You can use multiple templates and you can use the same template with
different parameters.  In that case you must give it a different name,
eg

    //go:generate gotemplate "github.com/ncw/gotemplate/set" StringSet(string)
    //go:generate gotemplate "github.com/ncw/gotemplate/set" FloatSet(float64)

If the parameters have spaces in then they need to be in quotes, eg

    //go:generate gotemplate "github.com/ncw/gotemplate/sort" "SortGt(string, func(a, b string) bool { return a > b })"

Renaming rules
--------------

All top level identifiers will be substituted when the template is
instantiated.  This is to ensure that they are unique if the template
is instantiated more than once.

Any identifiers with the template name in (eg `Set`) will have the
template name (eg `Set`) part substituted. If the template name does
not begin the identifier, Go's casing style is respected and the 
first letter of your new identifier is capitalized. (eg 'newMySet'
instead of 'newmySet').

Any identifiers without the template name in will just be post-fixed
with the template name.

So if this was run

    //go:generate gotemplate "github.com/ncw/gotemplate/set" MySet(string)

This would substitute these top level identifiers

  * `Set` to `MySet`
  * `NewSet` to `NewMySet`
  * `NewSizedSet` to `NewSizedMySet`
  * `utilityFunc` to `utilityFuncMySet`

Depending on whether the template name is public (initial capital) or
not, all the public external identifiers will have their initial
capitals turned into lower case.  So if this was run

    //go:generate gotemplate "github.com/ncw/gotemplate/set" mySet(string)

This would substitute

  * `Set` to `mySet`
  * `NewSet` to `newMySet`
  * `NewSizedSet` to `newSizedMySet`
  * `utilityFunc` to `utilityFuncMySet`

Installing templates
--------------------

Templates can be installed using `go get` because they are normal Go
packages.  Eg

    go get github.com/someones/template

Will install a template package you can use in your code with

    //go:generate gotemplate "github.com/someones/template" T(Potato)

Then instantiate with

    go generate

Source control for templates
----------------------------

It is expected that the generated files will be checked into version
control, and users of your code will just run `go get` to fetch it.
`go generate` will only be run by developers of the package.

Writing templates
-----------------

Templates are valid go packages.  They should compile and have tests
and be usable as-is.  Because they are packages, if you aren't writing
a public template you should put them in a subdirectory of your
project most likely.

To make a Go package a template it should have one or more
declarations and a special comment signaling to `gotemplate` what the
template is called and what its parameters are. Supported
parameterized declarations are type, const, var and func.

Here is an example from the set package.

    // template type Set(A)
    type A int

This indicates that the base name for the template is `Set` and it has
one type parameter `A`.  When you are writing the template package
make sure you use `A` instead of `int` where you want it to be
substituted with a new type when the template is instantiated.

Similarly, you could write a package with a const parameter.

    // template type Vector(A, N)
    type A int
    const N = 2

    type Vector[N]A

This indicates that the base name for the template is `Vector` and it
has one type parameter `A` and one constant parameter `N`. Again, all
uses of `N` in the template code will be replaced by a literal value
when the template is instantiated.

All the definitions of the template parameters will be removed from
the instantiated template.

All test files are ignored.

Bugs
----

There may be constraints on the types which aren't understood by
`gotemplate`.  For instance the set requires that the types are
comparable.  If you try this you'll get a compile error for example.

    //go:generate gotemplate "github.com/ncw/gotemplate/set" BytesSet([]byte)

Only one .go file is used when reading template definitions at the
moment (programmer laziness - will fix at some point!)

Changelog
---------

  * v0.06 - 2017-05-05
    * Add -outfmt string (thanks Paul Jolly)
  * v0.05 - 2016-02-26
    * Fix docs and examples
    * More set methods - thanks Adam Willis
    * Fix missing error check in code generation
  * v0.04 - 2014-12-23
    * Fixed multi-line type declarations
  * v0.03 - 2014-12-22
    * Allow const and var to be substituted as template parameters
  * v0.02 - 2014-12-15
    * Fixed multi-line const/var declarations
  * v0.01 - 2014-12-10
    * Change renaming rules to make better Go names.  This only affects private exports, eg for `mySet` in the example above,
      * `NewSet` becomes `newMySet` (was `newmySet`)
      * `NewSizedSet` becomes `newSizedMySet` (was `newSizedmySet`)
      * `utilityFunc` becomes `utilityFuncMySet` (`utilityFuncmySet`)
    * This is a backwards incompatible change
  * v0.00 - 2014-10-05
    * First public release

Ideas for the future
--------------------

Make a set type for non comparable things?  Pass in a compare routine?

Make sure that types implement an interface?

Optional parameters?

Philosophy
----------

All code (templates, use of templates and template instantiations)
should be normal Go code - no special types / extensions.

All configuration done with specially formatted comments

Should provide lots practical templates people can use right now.

License
-------

This is free software under the terms of MIT the license (check the
COPYING file included in this package).

Portions of the code have been copied from the Go source.  These are
identified by comments at the head of each file and these are
Copyright (c) The Go Authors.  See the GO-LICENSE file for full details.

Contact and support
-------------------

The project website is at:

  * https://github.com/ncw/gotemplate

There you can file bug reports, ask for help or contribute patches.

Authors
-------

  * Nick Craig-Wood <nick@craig-wood.com>

Contributors
------------

  * Patrick Oyarzun <patrickoyarzun@gmail.com>
  * Adam Willis <akwillis@inbox.com>
  * Paul Jolly <paul@myitcv.org.uk>
  * Igor Mikushkin <igor.mikushkin@gmail.com>
