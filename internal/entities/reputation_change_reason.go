package entities

type ReputationChangeReason byte

const (
	ReputationChangeReasonMint               ReputationChangeReason = 1
	ReputationChangeReasonBurn               ReputationChangeReason = 2
	ReputationChangeReasonVote               ReputationChangeReason = 3
	ReputationChangeReasonVotingDistribution ReputationChangeReason = 4
	ReputationChangeReasonVotingBurn         ReputationChangeReason = 5
)

type DeployExecutionType struct {
	ID   ReputationChangeReason `json:"id" db:"id"`
	Name string                 `json:"name" db:"name"`
}
