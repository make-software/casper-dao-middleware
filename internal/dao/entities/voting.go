package entities

import (
	"encoding/json"
	"time"

	"github.com/make-software/casper-go-sdk/casper"
)

type Voting struct {
	Creator                                  casper.Hash     `json:"creator" db:"creator"`
	DeployHash                               casper.Hash     `json:"deploy_hash" db:"deploy_hash"`
	VotingID                                 uint32          `json:"voting_id" db:"voting_id"`
	VotingTypeID                             VotingTypeID    `json:"voting_type_id" db:"voting_type_id"`
	InformalVotingQuorum                     uint32          `json:"informal_voting_quorum" db:"informal_voting_quorum"`
	InformalVotingStartsAt                   time.Time       `json:"informal_voting_starts_at" db:"informal_voting_starts_at"`
	InformalVotingEndsAt                     time.Time       `json:"informal_voting_ends_at" db:"informal_voting_ends_at"`
	FormalVotingQuorum                       uint32          `json:"formal_voting_quorum" db:"formal_voting_quorum"`
	FormalVotingTime                         uint64          `json:"formal_voting_time" db:"formal_voting_time"`
	FormalVotingStartsAt                     *time.Time      `json:"formal_voting_starts_at" db:"formal_voting_starts_at"`
	FormalVotingEndsAt                       *time.Time      `json:"formal_voting_ends_at" db:"formal_voting_ends_at"`
	Metadata                                 json.RawMessage `json:"metadata" db:"metadata"`
	IsCanceled                               bool            `json:"is_canceled" db:"is_canceled"`
	InformalVotingResult                     *uint8          `json:"informal_voting_result" db:"informal_voting_result"`
	FormalVotingResult                       *uint8          `json:"formal_voting_result" db:"formal_voting_result"`
	ConfigTotalOnboarded                     uint64          `json:"config_total_onboarded" db:"config_total_onboarded"`
	ConfigVotingClearnessDelta               uint64          `json:"config_voting_clearness_delta" db:"config_voting_clearness_delta"`
	ConfigTimeBetweenInformalAndFormalVoting uint64          `json:"config_time_between_informal_and_formal_voting" db:"config_time_between_informal_and_formal_voting"`
}

func NewVoting(
	creator, deployHash casper.Hash,
	votingID uint32,
	votingTypeID VotingTypeID,
	metadata json.RawMessage,
	informalVotingQuorum uint32,
	informalVotingStartsAt, informalVotingEndsAt time.Time,
	formalVotingQuorum uint32,
	formalVotingTime uint64,
	formalVotingStartsAt, formalVotingEndsAt *time.Time,
	configTotalOnboarded, configVotingClearnessDelta, configTimeBetweenInformalAndFormalVoting uint64,

) Voting {
	return Voting{
		Creator:                                  creator,
		DeployHash:                               deployHash,
		VotingID:                                 votingID,
		VotingTypeID:                             votingTypeID,
		InformalVotingQuorum:                     informalVotingQuorum,
		FormalVotingQuorum:                       formalVotingQuorum,
		FormalVotingTime:                         formalVotingTime,
		Metadata:                                 metadata,
		ConfigTotalOnboarded:                     configTotalOnboarded,
		InformalVotingStartsAt:                   informalVotingStartsAt,
		InformalVotingEndsAt:                     informalVotingEndsAt,
		FormalVotingStartsAt:                     formalVotingStartsAt,
		FormalVotingEndsAt:                       formalVotingEndsAt,
		ConfigVotingClearnessDelta:               configVotingClearnessDelta,
		ConfigTimeBetweenInformalAndFormalVoting: configTimeBetweenInformalAndFormalVoting,
	}
}
