package tokenizer

import (
	"errors"
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
		initState:    init,
		currentState: init,
		finalState:   &State{},
	}
}

type FSM struct {
	initState    *State
	currentState *State
	finalState   *State
	buffer       string
	finished     bool
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

func (f *FSM) AddFinalTransition(fn func(), srcs ...*State) {
	for _, src := range srcs {
		f.AddTransition(src, f.finalState, func(match rune) bool {
			fn()
			f.finished = true
			return true
		})
	}
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

func (f *FSM) MatchAll(text string) bool {
	f.Match(text)
	return false
}

func (f *FSM) Match(text string) bool {
	for ts := 0; ts < len(text); ts++ {
		f.Flush()
		f.currentState = f.initState
		if f.run(text[ts:]) == nil {
			return true
		}
	}
	return false
}

func (f *FSM) run(text string) error {
	for _, c := range text {
		err := f.nextState(c)
		if f.finished {
			return nil
		}
		if err != nil {
			return err
		}
	}
	return nil
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
