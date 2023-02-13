package entities

import (
	"time"

	"casper-dao-middleware/pkg/casper/types"
)

type Voting struct {
	Creator                                  types.Hash `json:"creator" db:"creator"`
	DeployHash                               types.Hash `json:"deploy_hash" db:"deploy_hash"`
	VotingID                                 uint32     `json:"voting_id" db:"voting_id"`
	IsFormal                                 bool       `json:"is_formal" db:"is_formal"`
	HasEnded                                 bool       `json:"has_ended" db:"has_ended"`
	VotingQuorum                             uint32     `json:"voting_quorum" db:"voting_quorum"`
	VotingTime                               uint64     `json:"voting_time" db:"voting_time"`
	DocumentHash                             string     `json:"document_hash" db:"timestamp"`
	ConfigTotalOnboarded                     uint64     `json:"config_total_onboarded" db:"timestamp"`
	ConfigDoubleTimeBetweenVotings           bool       `json:"config_double_time_between_votings" db:"timestamp"`
	ConfigVotingClearnessDelta               uint64     `json:"config_voting_clearness_delta" db:"timestamp"`
	ConfigTimeBetweenInformalAndFormalVoting uint64     `json:"config_time_between_informal_and_formal_voting" db:"config_time_between_informal_and_formal_voting"`
	Timestamp                                time.Time  `json:"timestamp" db:"timestamp"`
}

func NewVoting(
	creator, deployHash types.Hash,
	votingID, votingQuorum uint32,
	votingTime uint64,
	isFormal, configDoubleTimeBetweenVotings bool,
	documentHash string,
	configTotalOnboarded, configVotingClearnessDelta, configTimeBetweenInformalAndFormalVoting uint64,
	timestamp time.Time,
) Voting {
	return Voting{
		Creator:                                  creator,
		DeployHash:                               deployHash,
		VotingID:                                 votingID,
		IsFormal:                                 isFormal,
		VotingQuorum:                             votingQuorum,
		VotingTime:                               votingTime,
		Timestamp:                                timestamp,
		HasEnded:                                 false,
		DocumentHash:                             documentHash,
		ConfigTotalOnboarded:                     configTotalOnboarded,
		ConfigVotingClearnessDelta:               configVotingClearnessDelta,
		ConfigTimeBetweenInformalAndFormalVoting: configTimeBetweenInformalAndFormalVoting,
	}
}
