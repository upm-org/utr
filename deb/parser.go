package deb

import (
	"fmt"
)

type Data struct {
	singleLexemes []SingleLineLexeme
	multiLexemes  []MultiLineLexeme
	foldedLexemes []FoldedLexeme
}

func Parse(text string) {
	tempSL := &SingleLineLexeme{}
	slfsm := SingleLineFSM(tempSL)

	tempML := &MultiLineLexeme{}
	mlfsm := MultiLineFSM(tempML)

	tempFL := &FoldedLexeme{}
	flfsm := FoldedLineFSM(tempFL)

	slfsm.Match(text)
	mlfsm.Match(text)
	flfsm.Match(text)

	fmt.Println(*tempSL)
	fmt.Println(*tempML)
	fmt.Println(*tempFL)
}