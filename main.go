package main

import (
	"fmt"
	. "github.com/ump-org/utr/tokenizer"
	"strings"
)

func main() {
	sll := &MultiLineLexeme{}
	slfsm := MultiLineFSM(sll)
	res := slfsm.Match(NewTokenReader(strings.NewReader(`OK: now its multi
 line
 just
 as I said`), nil))
	fmt.Println(res, *sll)

}
