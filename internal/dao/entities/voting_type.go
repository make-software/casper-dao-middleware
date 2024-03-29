package entities

type VotingTypeID byte

const (
	VotingTypeSimple VotingTypeID = iota + 1
	VotingTypeSlashing
	VotingTypeKYC
	VotingTypeRepo
	VotingTypeReputation
	VotingTypeOnboarding
	VotingTypeAdmin
	VotingTypeBidEscrow
)

type VotingType struct {
	ID   VotingTypeID `json:"id" db:"id"`
	Name string       `json:"name" db:"name"`
}
