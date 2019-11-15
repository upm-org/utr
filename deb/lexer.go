package deb

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

const (
	fieldTypeUndefined  = iota
	fieldTypeSingleLine
	fieldTypeFolded
	fieldTypeMultiline
)

var (
	errIncorrectMatchesCount = errors.New("lexer: incorrect matches count")
	errPermittedSymbol = errors.New("lexer: permitted symbol found")
)

type lexeme struct {
	Field     string
	FieldType byte
	Value     []string
}

func fieldValidator(s string) error {
	for _, sc := range s {
		if (sc < 'A' || sc > 'Z') && (sc < 'a' || sc > 'z') && sc != '-' {
			return errPermittedSymbol
		}
	}
	return nil
}

func trimSpacesSlice(sl []string) []string {
	for i := range sl {
		sl[i] = strings.TrimSpace(sl[i])
	}
	return sl
}

func Parse(reader io.Reader) ([]*lexeme, error) {
	r := bufio.NewReader(reader)
	var res []*lexeme
	// reading until we met error
	for line, err := r.ReadString('\n'); err == nil; line, err = r.ReadString('\n') {
		// empty line check
		if strings.Trim(line, " \n") == "" {
			continue
		}
		// splitting to field and value tokens
		tokens := strings.SplitN(strings.Trim(line, " \n"), ":", 2)

		field := tokens[0]
		if err := fieldValidator(field); err != nil {
			return nil, err
		}

		line := strings.TrimSpace(tokens[1])
		// this is used to accumulate tokens
		value := line
		bytes, err := r.Peek(1)
		if err != nil {
			return nil, err
		}
		if bytes[0] == ' ' {
			for err != io.EOF && bytes[0] == ' ' {
				_, err = r.Discard(1)
				if err != nil {
					return nil, err
				}
				line, err = r.ReadString('\n')
				if err != nil {
					return res, nil
				}
				value += line
				bytes, err = r.Peek(1)
			}
			res = append(res, &lexeme{field, fieldTypeMultiline,
				strings.Split(strings.Trim(value, " \n"), "\n")})
		} else if strings.ContainsRune(value, ',') {
			for value[len(value) - 1] == ',' {
				// parsing as foldedField
				line, err = r.ReadString('\n')
				if err != nil {
					return nil, err
				}

				line = strings.Trim(line, " \n")
				value += line
			}
			res = append(res, &lexeme{field, fieldTypeFolded,
				trimSpacesSlice(strings.Split(value, ","))})
		} else {
			res = append(res, &lexeme{field, fieldTypeSingleLine, []string{value}})
		}
	}
	return res, nil
}
/*func tokenize(reader bufio.Reader) {
	line, _, err := reader.ReadLine()
	if strings.ContainsRune(string(line), ',') {

	}
}*/
