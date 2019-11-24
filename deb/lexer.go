package deb

import . "github.com/ump-org/utr/tokenizer"

type SingleLineLexeme struct {
	field string
	value string
}

type FoldedLexeme struct {
	field string
	value string
}

type MultiLineLexeme struct {
	field string
	value []string
}

func matchLetterOrNumber(token rune) bool {
	return MatchLetter(token) || MatchNumber(token)
}

func matchFoldedValue(token rune) bool {
	return MatchLetter(token) || MatchValue('\'', '.')(token)
}

func matchFieldRune(token rune) bool {
	return MatchLetter(token) || MatchValue('-')(token)
}

// singleLineFSM creates FSM that is bounded to a SingleLineLexeme
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

	slFSM.AddBufferedTransition(firstFieldLetter, fieldLetter, MatchLetter)
	slFSM.AddBufferedTransition(fieldLetter, fieldLetter, matchFieldRune)
	slFSM.AddTransition(fieldLetter, afterFieldSpace, func(Match rune) bool {
		if Match == ' ' {
			sll.field = slFSM.Flush()
			return true
		}
		return false
	})
	slFSM.AddTransition(fieldLetter, colon, func(Match rune) bool {
		if Match == ':' {
			sll.field = slFSM.Flush()
			return true
		}
		return false
	})

	slFSM.AddBufferedTransition(afterFieldSpace, afterFieldSpace, MatchValue(' '))
	slFSM.AddTransition(afterFieldSpace, colon, MatchValue(':'))

	slFSM.AddBufferedTransition(colon, afterColonSpace, MatchValue(' '))

	slFSM.AddBufferedTransition(afterColonSpace, afterColonSpace, MatchValue(' '))
	slFSM.AddBufferedTransition(afterColonSpace, firstValueLetter, MatchWord)

	slFSM.AddBufferedTransition(firstValueLetter, valueLetter, MatchWord)
	slFSM.AddBufferedTransition(valueLetter, valueLetter, MatchWord)
	slFSM.AddBufferedTransition(valueLetter, valueSpace, MatchValue(' '))

	slFSM.AddBufferedTransition(valueSpace, valueSpace, MatchValue(' '))
	slFSM.AddBufferedTransition(valueSpace, valueLetter, MatchWord)

	return slFSM
}

// multiLineFSM creates FSM that is bounded to a MultiLineLexeme
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
		mll.value = append(values, mlFSM.Flush())
	})

	mlFSM.AddBufferedTransition(firstFieldLetter, fieldLetter, MatchLetter)
	mlFSM.AddBufferedTransition(fieldLetter, fieldLetter, matchFieldRune)
	mlFSM.AddTransition(fieldLetter, afterFieldSpace, func(Match rune) bool {
		if Match == ' ' {
			mll.field = mlFSM.Flush()
			return true
		}
		return false
	})
	mlFSM.AddTransition(fieldLetter, colon, func(Match rune) bool {
		if Match == ':' {
			mll.field = mlFSM.Flush()
			return true
		}
		return false
	})

	mlFSM.AddBufferedTransition(afterFieldSpace, afterFieldSpace, MatchValue(' '))
	mlFSM.AddBufferedTransition(afterFieldSpace, colon, MatchValue(':'))

	mlFSM.AddTransition(colon, afterColonSpace, MatchValue(' '))

	mlFSM.AddTransition(afterColonSpace, afterColonSpace, MatchValue(' '))
	mlFSM.AddBufferedTransition(afterColonSpace, firstValueLetter, MatchWord)

	mlFSM.AddBufferedTransition(firstValueLetter, valueLetter, MatchWord)
	mlFSM.AddBufferedTransition(valueLetter, valueLetter, MatchWord)
	mlFSM.AddBufferedTransition(valueLetter, valueSpace, MatchValue(' '))
	mlFSM.AddBufferedTransition(valueLetter, newLine, MatchValue('\n'))

	mlFSM.AddBufferedTransition(valueSpace, valueSpace, MatchValue(' '))
	mlFSM.AddBufferedTransition(valueSpace, valueLetter, MatchWord)
	mlFSM.AddBufferedTransition(valueSpace, newLine, MatchValue('\n'))

	mlFSM.AddBufferedTransition(newLine, newLineSpace, func(Match rune) bool {
		if Match == ' ' {
			values = append(values, mlFSM.Flush())
			return true
		}
		return false
	})

	mlFSM.AddTransition(newLineSpace, dotNewLine, MatchValue('.'))
	mlFSM.AddBufferedTransition(newLineSpace, valueSpace, MatchValue(' '))
	mlFSM.AddBufferedTransition(newLineSpace, valueLetter, MatchLetter)

	mlFSM.AddBufferedTransition(dotNewLine, newLine, MatchValue('\n'))

	return mlFSM
}

// foldedLineFSM creates FSM that is bounded to a FoldedLexeme
func FoldedLineFSM(fll *FoldedLexeme) *FSM {
	firstFieldLetter := NewState()
	fieldLetter := NewState()
	afterFieldSpace := NewState()
	colon := NewState()
	afterColonSpace := NewState()
	firstValueLetter := NewState()
	valueLetter := NewState()
	valueSpace := NewState()
	comma := NewState()
	newLine := NewState()

	flFSM := NewFSM(firstFieldLetter)
	flFSM.EOFHandler(func() {
		fll.value += flFSM.Flush()
	})

	flFSM.AddBufferedTransition(firstFieldLetter, fieldLetter, MatchLetter)
	flFSM.AddBufferedTransition(fieldLetter, fieldLetter, matchFieldRune)
	flFSM.AddTransition(fieldLetter, afterFieldSpace, func(Match rune) bool {
		if Match == ' ' {
			fll.field = flFSM.Flush()
			return true
		}
		return false
	})
	flFSM.AddTransition(fieldLetter, colon, func(Match rune) bool {
		if Match == ':' {
			fll.field = flFSM.Flush()
			return true
		}
		return false
	})

	flFSM.AddTransition(afterFieldSpace, afterFieldSpace, MatchValue(' '))
	flFSM.AddTransition(afterFieldSpace, colon, MatchValue(':'))

	flFSM.AddTransition(colon, afterColonSpace, MatchValue(' '))

	flFSM.AddTransition(afterColonSpace, afterColonSpace, MatchValue(' '))
	flFSM.AddBufferedTransition(afterColonSpace, firstValueLetter, matchFoldedValue)

	flFSM.AddBufferedTransition(firstValueLetter, valueLetter, matchFoldedValue)
	flFSM.AddBufferedTransition(valueLetter, valueLetter, matchFoldedValue)
	flFSM.AddBufferedTransition(valueLetter, valueSpace, MatchValue(' '))
	flFSM.AddTransition(valueLetter, comma, func(match rune) bool {
		if match == ',' {
			fll.value += flFSM.Flush() + " "
			return true
		}
		return false
	})

	flFSM.AddBufferedTransition(valueSpace, valueSpace, MatchValue(' '))
	flFSM.AddBufferedTransition(valueSpace, valueLetter, matchFoldedValue)
	flFSM.AddTransition(valueSpace, comma, func(match rune) bool {
		if match == ',' {
			fll.value += flFSM.Flush() + " "
			return true
		}
		return false
	})

	flFSM.AddTransition(comma, valueSpace, MatchValue(' '))
	flFSM.AddTransition(comma, newLine, MatchValue('\n'))

	flFSM.AddBufferedTransition(newLine, valueSpace, MatchValue(' '))
	flFSM.AddBufferedTransition(newLine, valueLetter, matchFoldedValue)

	return flFSM
}
