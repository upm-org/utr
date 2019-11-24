package tokenizer

import (
	"errors"
	"io"
)

type SingleLineLexeme struct {
	field string
	value string
}

type foldedLexeme struct {
	field string
	value []string
}

type MultiLineLexeme struct {
	field string
	value []string
}

type State struct {
	transitions []func(token Token) *State
}

func NewState() *State {
	return &State{}
}

type transFunc func(match Token) bool

var (
	errNoTransitions = errors.New("fsm did no transitions")
)

func NewFSM(init *State) *FSM {
	return &FSM{
		currentState: init,
	}
}

type FSM struct {
	currentState *State
	buffer       string
	eofHandler   func()
}

func (f *FSM) addTransition(src *State, dest *State, tf transFunc) {
	src.transitions = append(src.transitions, func(token Token) *State {
		if tf(token) {
			return dest
		}
		return nil
	})
}

func (f *FSM) addBufferedTransition(src *State, dest *State, tf transFunc) {
	src.transitions = append(src.transitions, func(token Token) *State {
		if tf(token) {
			f.buffer += string(token.Value)
			return dest
		}
		return nil
	})
}

func (f *FSM) Flush() string {
	s := f.buffer
	f.buffer = ""
	return s
}

func (f *FSM) nextState(t Token) error {
	transitioned := false

	for _, transFn := range f.currentState.transitions {
		if s := transFn(t); s != nil {
			f.currentState = s
			transitioned = true
			break
		}
	}

	if !transitioned {
		return errNoTransitions
	}

	return nil
}

func (f *FSM) EOFHandler(fn func()) {
	f.eofHandler = fn
}

func (f *FSM) Match(tr TokenOneReader) bool {
	if f.run(tr) != nil {
		return false
	}
	return true
}

func (f *FSM) run(tr TokenOneReader) error {
	var err error
	for err == nil {
		var t Token
		t, err = tr.ReadOne()
		if err == io.EOF {
			f.eofHandler()
			return nil
		}
		if err != nil {
			return err
		}

		err = f.nextState(t)
	}
	return err
}

func matchType(tokenType int) func(Token) bool {
	return func(m Token) bool {
		return m.TokenType == tokenType
	}
}

func matchValue(tokenValue rune) func(Token) bool {
	return func(m Token) bool {
		return m.Value == tokenValue
	}
}

// singleLineFSM basically matches next RegExp:
// /([a-zA-Z]+) *: *([a-zA-Z'][ a-zA-Z']*)/gm
func SingleLineFSM(sll *SingleLineLexeme) *FSM {
	firstFieldLetter := NewState()
	fieldLetter := NewState()
	afterFieldSpace := NewState()
	colon := NewState()
	afterColonSpace := NewState()
	firstValueLetter := NewState()
	valueLetter := NewState()
	valueSpace := NewState()

	slFSM := NewFSM(firstFieldLetter)
	slFSM.EOFHandler(func() {
		sll.value = slFSM.Flush()
	})

	slFSM.addBufferedTransition(firstFieldLetter, fieldLetter, matchType(LetterToken))
	slFSM.addBufferedTransition(fieldLetter, fieldLetter, matchType(LetterToken))
	slFSM.addBufferedTransition(fieldLetter, afterFieldSpace, func(match Token) bool {
		if match.TokenType == SpaceToken {
			sll.field = slFSM.Flush()
			return true
		}
		return false
	})
	slFSM.addBufferedTransition(fieldLetter, colon, func(match Token) bool {
		if match.TokenType == ColonToken {
			sll.field = slFSM.Flush()
			return true
		}
		return false
	})

	slFSM.addBufferedTransition(afterFieldSpace, afterFieldSpace, matchType(SpaceToken))
	slFSM.addBufferedTransition(afterFieldSpace, colon, matchType(ColonToken))

	slFSM.addBufferedTransition(colon, afterColonSpace, matchType(SpaceToken))

	slFSM.addBufferedTransition(afterColonSpace, afterColonSpace, matchType(SpaceToken))
	slFSM.addBufferedTransition(afterColonSpace, firstValueLetter, matchType(LetterToken))

	slFSM.addBufferedTransition(firstValueLetter, valueLetter, matchType(LetterToken))
	slFSM.addBufferedTransition(valueLetter, valueLetter, matchType(LetterToken))
	slFSM.addBufferedTransition(valueLetter, valueSpace, matchType(SpaceToken))

	slFSM.addBufferedTransition(valueSpace, valueSpace, matchType(SpaceToken))
	slFSM.addBufferedTransition(valueSpace, valueLetter, matchType(LetterToken))

	return slFSM
}

func MultiLineFSM(mll *MultiLineLexeme) *FSM {
	firstFieldLetter := NewState()
	fieldLetter := NewState()
	afterFieldSpace := NewState()
	colon := NewState()
	afterColonSpace := NewState()
	firstValueLetter := NewState()
	valueLetter := NewState()
	valueSpace := NewState()
	newLine := NewState()
	newLineSpace := NewState()
	dotNewLine := NewState()

	mlFSM := NewFSM(firstFieldLetter)
	var values []string
	mlFSM.EOFHandler(func() {
		mll.value = values
	})

	mlFSM.addBufferedTransition(firstFieldLetter, fieldLetter, matchType(LetterToken))
	mlFSM.addBufferedTransition(fieldLetter, fieldLetter, matchType(LetterToken))
	mlFSM.addBufferedTransition(fieldLetter, afterFieldSpace, func(match Token) bool {
		if match.TokenType == SpaceToken {
			mll.field = mlFSM.Flush()
			return true
		}
		return false
	})
	mlFSM.addTransition(fieldLetter, colon, func(match Token) bool {
		if match.TokenType == ColonToken {
			mll.field = mlFSM.Flush()
			return true
		}
		return false
	})

	mlFSM.addBufferedTransition(afterFieldSpace, afterFieldSpace, matchType(SpaceToken))
	mlFSM.addBufferedTransition(afterFieldSpace, colon, matchType(ColonToken))

	mlFSM.addBufferedTransition(colon, afterColonSpace, matchType(SpaceToken))

	mlFSM.addBufferedTransition(afterColonSpace, afterColonSpace, matchType(SpaceToken))
	mlFSM.addBufferedTransition(afterColonSpace, firstValueLetter, matchType(LetterToken))

	mlFSM.addBufferedTransition(firstValueLetter, valueLetter, matchType(LetterToken))
	mlFSM.addBufferedTransition(valueLetter, valueLetter, matchType(LetterToken))
	mlFSM.addBufferedTransition(valueLetter, valueSpace, matchType(SpaceToken))
	mlFSM.addBufferedTransition(valueLetter, newLine, matchValue('\n'))

	mlFSM.addBufferedTransition(valueSpace, valueSpace, matchType(SpaceToken))
	mlFSM.addBufferedTransition(valueSpace, valueLetter, matchType(LetterToken))
	mlFSM.addBufferedTransition(valueSpace, newLine, matchValue('\n'))

	mlFSM.addBufferedTransition(newLine, newLineSpace, func(match Token) bool {
		if match.TokenType == SpaceToken {
			values = append(values, mlFSM.Flush())
			return true
		}
		return false
	})

	mlFSM.addTransition(newLineSpace, dotNewLine, matchValue('.'))
	mlFSM.addBufferedTransition(newLineSpace, valueSpace, matchType(SpaceToken))
	mlFSM.addBufferedTransition(newLineSpace, valueLetter, matchType(LetterToken))

	return mlFSM
}
