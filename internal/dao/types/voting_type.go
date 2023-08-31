package types

import "errors"

type VotingType byte

const (
	VotingTypeInformal VotingType = 0
	VotingTypeFormal   VotingType = 1
)

func NewVotingTypeFromByte(b byte) (VotingType, error) {
	if b != byte(VotingTypeInformal) && b != byte(VotingTypeFormal) {
		return 0, errors.New("invalid voting_type: expected VotingTypeInformal(0) or VotingTypeFormal(1)")
	}
	return VotingType(b), nil
}
