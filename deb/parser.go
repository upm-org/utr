package deb

import (
	"fmt"
	"os"
	"strings"
)

type person struct {
	name string
	email string
}

func (p *person) Parse(text string) {
	sp := strings.Split(text, " ")
	p.name = strings.Join(sp[:len(sp)-1], " ")
	p.email = sp[len(sp) - 1]
}

var lexemes []*lexemlexe

type SourcePackage struct {
	source string
	section string
	priority string
	maintainer person
	uploaders []person
	standardsVersion string
	buildDepends []string
	rulesRequiresRoot bool
	homepage string
	vcsGit string
	vcsBrowser string
}

func (s *SourcePackage) ParseFile(path string) error {
	/*data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	splitted := regexp.MustCompile(`\n\S`).Split(string(data), -1)
	for _, s := range splitted {
		if foldedRegexp.MatchString(s) {
			lexemes = append(lexemes, FoldedLineLexeme(foldedRegexp.FindAllString(s, -1)))
			break
		}
		if multiLineRegexp.MatchString(s) {
			lexemes = append(lexemes, MultiLineLexeme(foldedRegexp.FindAllString(s, -1)))
		}
	}*/
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	res, err := Parse(file)
	if err != nil {
		return err
	}
	for _, l := range res {
		fmt.Println(*l)
	}
	return nil
}

type debPackage struct {
	name string
	architecture []string
	depends []string
	recommends []string
	suggests []string
	breaks []string
	conflicts []string
	replaces []string
	description string
}

type control struct {

}

func parseControl() {

}
