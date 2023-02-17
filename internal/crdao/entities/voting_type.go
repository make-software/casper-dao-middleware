package entities

type VotingTypeID byte

const (
	VotingTypeSimple VotingTypeID = iota + 1
	VotingTypeSlashing
	VotingTypeKYC
	VotingTypeRepo
	VotingTypeReputation
)

type VotingType struct {
	ID   VotingTypeID `json:"id" db:"id"`
	Name string       `json:"name" db:"name"`
}
