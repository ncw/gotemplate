Go templates
============

This tool manages package based templates for the Go language.

Install
-------

Install using go get

    go get github.com/ncw/gotemplate/...

and this will build the `gotemplate` binary in `$GOPATH/bin`.

It will also pull in a set of templates you can start using straight away

FIXME link to templates

Using templates
---------------

To use a template, first you must tell `gotemplate` that you want to
use it using a special comment in your code.  For example

    // template "github.com/ncw/gotemplate/set" mySet(string)

This tells `gotemplate` that you want to use the set template with the
`string` type parameter and with the local name `mySet`.

Now run `gotemplate` in your code directory with no arguments.  This
will instantiate the template into a file called `gotemplate_mySet.go`
which will provide a `mySet` type and `newmySet` and `newSizesmySet`
functions to make them.

FIXME outptut of the command goes here

If you use an initial capital when you name your template
instantiation then any external functions will be public.  Eg

    // template "github.com/ncw/gotemplate/set" MySet(string)

Would give you `MySet`, `NewMySet` and `NewSizedMySet` instead.

You can use multiple templates and you can use the same template with
different paramters.  In that case you must give it a different name,
eg

    // template "github.com/ncw/gotemplate/set" StringSet(string)
    // template "github.com/ncw/gotemplate/set" FloatSet(float64)

Instantiating the templates into your project gives them the ability
to use internal types from your project.

Installing templates
--------------------

Templates can be installed using `go get` because they are normal Go packages.  Eg

    go get github.com/someones/template

Will install a template package you can use with

    // template "github.com/someones/template" T(Potato)


Writing templates
-----------------

Templates are valid go packages.  They should compile and have tests
and be usable as-is.  Because they are packages, if you aren't writing
a public template you should put them in a subdirectory of your
project most likely.

To make a Go package a template it should have a type definition and a
special comment signaling to `gotemplate` what the type is called and
what its type parameters are.

Here is an example from the set package.

    // template type Set(A)
    type A int

This indicates that the type is called `Set` and it has 1 type
parameter `A`.  When you are writing the template package make sure
you use `A` instead of `int` where you want it to be substituted with
a new type when the template is instantiated.

All non test go files will be used as templates

Bugs
----

There may be constraints on the types which aren't understood by
`gotemplate`.  For instance the set requires that the types are
comparable.  If you try this you'll get a compile error for example.

    // template "github.com/ncw/gotemplate/set" BytesSet([]byte)

FIXME make a set type for non comparable things?

FIXME make sure that types implement an interface? Or pass in a compare routine?

Optional parametrs?

FIXME only 1 .go file is used when reading templates at the moment

Philosophy
----------

All code (templates, use of templates and template instantiations)
should be normal Go code - no special types / extensions.

All configuration done with specially formatted comments

Should provide lots practical templates people can use right now.


Similar Projects
----------------

https://github.com/droundy/gotgo/ - uses .got files

License
-------

This is free software under the terms of MIT the license (check the
COPYING file included in this package).

Contact and support
-------------------

The project website is at:

- https://github.com/ncw/gotemplate

There you can file bug reports, ask for help or contribute patches.

Authors
-------

- Nick Craig-Wood <nick@craig-wood.com>

Contributors
------------

- Your name goes here!
