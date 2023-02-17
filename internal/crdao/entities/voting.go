package entities

import (
	"encoding/json"
	"time"

	"casper-dao-middleware/pkg/casper/types"
)

type Voting struct {
	Creator                                  types.Hash      `json:"creator" db:"creator"`
	DeployHash                               types.Hash      `json:"deploy_hash" db:"deploy_hash"`
	VotingID                                 uint32          `json:"voting_id" db:"voting_id"`
	VotingTypeID                             VotingTypeID    `json:"voting_type_id" db:"voting_type_id"`
	IsFormal                                 bool            `json:"is_formal" db:"is_formal"`
	HasEnded                                 bool            `json:"has_ended" db:"has_ended"`
	VotingQuorum                             uint32          `json:"voting_quorum" db:"voting_quorum"`
	VotingTime                               uint64          `json:"voting_time" db:"voting_time"`
	Metadata                                 json.RawMessage `json:"metadata" db:"metadata"`
	ConfigTotalOnboarded                     uint64          `json:"config_total_onboarded" db:"config_total_onboarded"`
	ConfigDoubleTimeBetweenVotings           bool            `json:"config_double_time_between_votings" db:"config_double_time_between_votings"`
	ConfigVotingClearnessDelta               uint64          `json:"config_voting_clearness_delta" db:"config_voting_clearness_delta"`
	ConfigTimeBetweenInformalAndFormalVoting uint64          `json:"config_time_between_informal_and_formal_voting" db:"config_time_between_informal_and_formal_voting"`
	Timestamp                                time.Time       `json:"timestamp" db:"timestamp"`
}

func NewVoting(
	creator, deployHash types.Hash,
	votingID, votingQuorum uint32,
	votingTime uint64,
	votingTypeID VotingTypeID,
	metadata json.RawMessage,
	isFormal, configDoubleTimeBetweenVotings bool,
	configTotalOnboarded, configVotingClearnessDelta, configTimeBetweenInformalAndFormalVoting uint64,
	timestamp time.Time,
) Voting {
	return Voting{
		Creator:                                  creator,
		DeployHash:                               deployHash,
		VotingID:                                 votingID,
		VotingTypeID:                             votingTypeID,
		IsFormal:                                 isFormal,
		VotingQuorum:                             votingQuorum,
		VotingTime:                               votingTime,
		Timestamp:                                timestamp,
		HasEnded:                                 false,
		Metadata:                                 metadata,
		ConfigDoubleTimeBetweenVotings:           configDoubleTimeBetweenVotings,
		ConfigTotalOnboarded:                     configTotalOnboarded,
		ConfigVotingClearnessDelta:               configVotingClearnessDelta,
		ConfigTimeBetweenInformalAndFormalVoting: configTimeBetweenInformalAndFormalVoting,
	}
}
