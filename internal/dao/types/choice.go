package types

import "errors"

type Choice byte

const (
	ChoiceAgainst Choice = 0
	ChoiceInFavor Choice = 1
)

func NewChoiceFromByte(b byte) (Choice, error) {
	if b != byte(ChoiceAgainst) && b != byte(ChoiceInFavor) {
		return 0, errors.New("invalid choice: expected ChoiceAgainst(1) or ChoiceInFavor(2)")
	}
	return Choice(b), nil
}
