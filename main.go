package main

import (
	"fmt"
	. "github.com/ump-org/utr/deb"
	"strings"
)

func main() {
	sll := &SingleLineLexeme{}
	slFSM := SingleLineFSM(sll)
	slRes := slFSM.Match(strings.NewReader("Foo: bar baz guz gaz"))
	fmt.Println(slRes, *sll)

	mll := &MultiLineLexeme{}
	mlFSM := MultiLineFSM(mll)
	mlRes := mlFSM.Match(strings.NewReader(`Bar: lorem ipsum
 popsum
 mopsum gaz gazo
 foobaz bazoo foo`))
	fmt.Println(mlRes, *mll)

	fll := &FoldedLexeme{}
	flFSM := FoldedLineFSM(fll)
	flRes := flFSM.Match(strings.NewReader(`Baz: check, this,
out,
here`))
	fmt.Println(flRes, *fll)
}
