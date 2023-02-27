package entities

type ReputationChangeReason byte

const (
	ReputationChangeReasonMinted         ReputationChangeReason = 1
	ReputationChangeReasonBurned         ReputationChangeReason = 2
	ReputationChangeReasonStaked         ReputationChangeReason = 3
	ReputationChangeReasonVotingGained   ReputationChangeReason = 4
	ReputationChangeReasonVotingLost     ReputationChangeReason = 5
	ReputationChangeReasonVotingUnstaked ReputationChangeReason = 6
)

type DeployExecutionType struct {
	ID   ReputationChangeReason `json:"id" db:"id"`
	Name string                 `json:"name" db:"name"`
}
