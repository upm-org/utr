package tokenizer

import (
	"errors"
	"io"
)

// Error definitions
var (
	errNotEnoughTokens = errors.New("not enough tokens")
)

// Token types
const (
	illegalToken = iota

	SpaceToken
	CommaToken
	NewLineToken
	ColonToken
	NotLetterToken
	LetterToken
)

// Token represents a single Token. Value is the data Token
// holds, TokenType represent it's type.
type Token struct {
	Value     rune
	TokenType int
}

type TokenReader interface {
	Read(t []Token) (n int, err error)
}

type TokenOneReader interface {
	ReadOne() (Token, error)
}

type TokenFullReader interface {
	TokenReader
	TokenOneReader

	ReadAll() (t []Token, err error)
}

type clusterFunc func(r rune) int

type defaultReader struct {
	io.RuneReader

	clusterFunc

	position int
}

func (r *defaultReader) Read(t []Token) (n int, err error) {
	n = 0

	for n < len(t) {
		rn, size, rnErr := r.ReadRune()
		if rnErr != nil {
			err = rnErr
			return
		}
		n += size

		t[r.position] = Token{rn, r.clusterFunc(rn)}
		r.position++
	}

	return
}

func (r *defaultReader) ReadOne() (Token, error) {
	rn, _, rnErr := r.ReadRune()
	if rnErr != nil {
		return Token{}, rnErr
	}
	res := Token{rn, r.clusterFunc(rn)}
	r.position++
	return res, nil
}

func (r *defaultReader) ReadAll() (res []Token, err error) {
	var buff [256]Token
	for err != io.EOF {
		_, err = r.Read(buff[:])
		if err != nil && err != io.EOF {
			return nil, err
		}
		res = append(res, buff[:]...)
	}
	return
}

// Cluster responds for detecting Token type
func defaultClusterFunc(r rune) int {
	switch {
	case r == ' ':
		return SpaceToken
	case r == ',':
		return CommaToken
	case r == '\n':
		return NewLineToken
	case r == ':':
		return ColonToken
	case (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z'):
		return LetterToken
	default:
		return NotLetterToken
	}
}

func NewToken(r rune, f clusterFunc) Token {
	if f != nil {
		return Token{r, f(r)}
	}
	return Token{r, defaultClusterFunc(r)}
}

func NewTokenReader(reader io.RuneReader, f clusterFunc) TokenFullReader {
	if f != nil {
		return &defaultReader{RuneReader: reader, clusterFunc: f}
	}
	return &defaultReader{RuneReader: reader, clusterFunc: defaultClusterFunc}
}
