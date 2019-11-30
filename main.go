package main

import (
	. "github.com/ump-org/utr/deb"
)

func main() {
	Parse(`Foo: a b c
Bar: d
 e
 f
Baz: g,
h,
i
`)
}
