package main

import (
	"fmt"
	"utr/tokenizer"
)

func main() {
	//var a deb.SourcePackage

	t, _ := tokenizer.New()
	res, err := t.TokenizeFile("deb/nano/debian/control")
	//err := a.ParseFile("deb/nano/debian/control")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}
