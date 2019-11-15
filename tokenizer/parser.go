package tokenizer

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"strings"
)

// token represents a single token. Value is the data token
// holds, tokenType represent it's type.
type token struct {
	value     string
	tokenType int
}

// tokenFunc responds for tokenizing function.
// Arguments: First parameter is RuneReader to read from,
// Second parameter is amount of bytes to read
//
// Return values are set of string tokens, bytes read, error
type tokenFunc func(io.RuneReader, int) ([]token, int, error)

// clusterFunc is clustering tokens in different token groups.
// Plays part in tokenFunc.
type clusterFunc func(r rune) int

type TokenizerOption func(*Tokenizer) error

// Tokenizer holds tokenFunc
type Tokenizer struct {
	tokenFunc
	clusterFunc
}

// Error definitions
var (
	errNotEnoughTokens = errors.New("not enough tokens")
)

// Token types
const (
	illegalToken = iota

	spaceToken
	commaToken
	newLineToken
	colonToken
	notLetterToken
	letterToken
)

// defaultClusterBytes responds for detecting token type
func defaultClusterFunc(r rune) int {
	switch {
	case r == ' ':
		return spaceToken
	case r == ',':
		return commaToken
	case r == '\n':
		return newLineToken
	case r == ':':
		return colonToken
	case (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z'):
		return letterToken
	default:
		return notLetterToken
	}
}

// Clue clues nearby standing tokens with identical types.
// Returns error if tokens amount provided is less then 2.
func Clue(tokens []token) ([]token, error) {
	var res []token
	temp := tokens[0]

	if len(tokens) < 2 {
		return nil, errNotEnoughTokens
	}

	for i := 1; i < len(tokens); i++ {
		t := tokens[i]

		if temp.tokenType == t.tokenType {
			temp.value += t.value
		} else {
			res = append(res, temp)
			temp = t
		}
	}
	return res, nil
}

// defaultTokenFunc is the default tokenizing function
func defaultTokenFunc(reader io.RuneReader, n int) ([]token, int, error) {
	const buffSize = 1
	var tokens []token
	read := 0

	for read < n {
		var t token

		r, size, err := reader.ReadRune()
		if err != nil {
			return nil, 0, err
		}
		read += size

		t.tokenType = defaultClusterFunc(r)
		t.value = string(r)

		tokens = append(tokens, t)
	}

	res, err := Clue(tokens)
	if err != nil {
		return nil, 0, err
	}

	return res, read, nil
}

// NewTokenizer creates new Tokenizer
func New(opts ...TokenizerOption) (*Tokenizer, error) {
	t := &Tokenizer{
		tokenFunc:   defaultTokenFunc,
		clusterFunc: defaultClusterFunc,
	}
	for _, opt := range opts {
		if err := opt(t); err != nil {
			return nil, err
		}
	}
	return t, nil
}

// TokenFunc sets f as a tokenizing function.
func TokenFunc(f tokenFunc) TokenizerOption {
	return func(t *Tokenizer) error {
		t.tokenFunc = f
		return nil
	}
}

// ClusterFunc sets f as a clusterizing function.
func ClusterFunc(f clusterFunc) TokenizerOption {
	return func(t *Tokenizer) error {
		t.clusterFunc = f
		return nil
	}
}

func (t *Tokenizer) TokenizeString(s string) ([]token, error) {
	res, _, err := t.tokenFunc(strings.NewReader(s), len(s))
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (t *Tokenizer) TokenizeBytes(b []byte) ([]token, error) {
	res, _, err := t.tokenFunc(bytes.NewBuffer(b), len(b))
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (t *Tokenizer) TokenizeFile(path string) ([]token, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	res, err := t.TokenizeBytes(f)
	if err != nil {
		return nil, err
	}

	return res, nil
}
