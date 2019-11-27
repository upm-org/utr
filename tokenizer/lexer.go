package tokenizer

import (
	"errors"
	"io"
)

type State struct {
	transitions []func(token rune) *State
}

type transFunc func(match rune) bool

func NewState() *State {
	return &State{}
}

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

func (f *FSM) AddTransition(src *State, dest *State, tf transFunc) {
	src.transitions = append(src.transitions, func(token rune) *State {
		if tf(token) {
			return dest
		}
		return nil
	})
}

func (f *FSM) AddBufferedTransition(src *State, dest *State, tf transFunc) {
	src.transitions = append(src.transitions, func(token rune) *State {
		if tf(token) {
			f.buffer += string(token)
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

func (f *FSM) nextState(t rune) error {
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

func (f *FSM) Match(r io.RuneReader) bool {
	if f.run(r) != nil {
		return false
	}
	return true
}

func (f *FSM) run(r io.RuneReader) error {
	var err error
	for err == nil {
		var t rune
		t, _, err = r.ReadRune()
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

func MatchLetter(token rune) bool {
	return (token >= 'a' && token <= 'z') || (token >= 'A' && token <= 'Z')
}

func MatchNumber(token rune) bool {
	return token >= '0' && token <= '9'
}

func MatchWord(token rune) bool {
	return MatchLetter(token) || token == '\'' || token == '.' || token == ','
}

func MatchValue(tokens ...rune) func(rune) bool {
	return func(m rune) bool {
		for _, t := range tokens {
			if m == t {
				return true
			}
		}
		return false
	}
}
