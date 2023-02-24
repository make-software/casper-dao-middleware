package types

import "errors"

type VotingType byte

const (
	VotingTypeInformal VotingType = 1
	VotingTypeFormal   VotingType = 2
)

func NewVotingTypeFromByte(b byte) (VotingType, error) {
	if b != byte(VotingTypeInformal) && b != byte(VotingTypeFormal) {
		return 0, errors.New("invalid voting_type: expected VotingTypeInformal(1) or VotingTypeFormal(2)")
	}
	return VotingType(b), nil
}
