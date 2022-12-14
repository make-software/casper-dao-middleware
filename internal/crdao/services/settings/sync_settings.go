package settings

import (
	"bytes"
	"encoding/hex"
	"errors"
	"go.uber.org/zap"
	"time"

	"casper-dao-middleware/internal/crdao/dao_event_parser/utils"
	"casper-dao-middleware/internal/crdao/di"
	"casper-dao-middleware/internal/crdao/entities"
)

const (
	DefaultPolicingRate         = "default_policing_rate"
	ReputationConversionRate    = "reputation_conversion_rate"
	ForumKycRequired            = "forum_kyc_required"
	FormalVotingQuorum          = "formal_voting_quorum"
	InformalVotingQuorum        = "informal_voting_quorum"
	VotingQuorum                = "voting_quorum"
	FormalVotingTime            = "formal_voting_time"
	InformalVotingTime          = "informal_voting_time"
	VotingTime                  = "voting_time"
	MinimumGovernanceReputation = "minimum_governance_reputation"
	MinimumVotingReputation     = "minimum_voting_reputation"
)

var DaoSettings = []string{
	DefaultPolicingRate,
	ReputationConversionRate,
	ForumKycRequired,
	FormalVotingQuorum,
	InformalVotingQuorum,
	VotingQuorum,
	FormalVotingTime,
	InformalVotingTime,
	VotingTime,
	MinimumGovernanceReputation,
	MinimumVotingReputation,
}

type SyncDAOSetting struct {
	di.EntityManagerAware
	di.CasperClientAware

	variableRepositoryContractStorageUref string
	setting                               string
}

func NewSyncDAOSettings() SyncDAOSetting {
	return SyncDAOSetting{}
}

func (c *SyncDAOSetting) SetVariableRepositoryContractStorageUref(uref string) {
	c.variableRepositoryContractStorageUref = uref
}

func (c *SyncDAOSetting) SetSetting(setting string) {
	c.setting = setting
}

func (c *SyncDAOSetting) Execute() error {
	stateRootHashRes, err := c.GetCasperClient().GetStateRootHashByHash("")
	if err != nil {
		return err
	}

	settingItemKey, err := utils.ToDictionaryItemKey(c.setting)
	if err != nil {
		return err
	}

	result, err := c.GetCasperClient().GetDictionaryItem(stateRootHashRes.StateRootHash, c.variableRepositoryContractStorageUref, settingItemKey)
	if err != nil {
		return err
	}

	if result.StoredValue.CLValue == nil {
		zap.S().With(zap.String("setting", c.setting)).Debug("expected initialized CLValue")
		return errors.New("expected initialized CLValue")
	}

	decoded := make([]byte, len(result.StoredValue.CLValue.Bytes))
	if _, err := hex.Decode(decoded, bytes.Trim(result.StoredValue.CLValue.Bytes, "\"")); err != nil {
		return err
	}

	record, err := utils.NewRecordFromBytes(decoded)
	if err != nil {
		return err
	}

	setting := entities.NewSetting(c.setting, record.Value.String(), nil, nil)

	if futureVal := record.FutureValue; futureVal != nil {
		nextVal := futureVal.Value.String()
		setting.NextValue = &nextVal

		activationTime := time.Unix(int64(futureVal.ActivationTime), 0)
		setting.ActivationTime = &activationTime
	}

	if err := c.GetEntityManager().SettingRepository().Upsert(setting); err != nil {
		return err
	}

	zap.S().With(zap.String("setting", c.setting)).Info("Variable contract setting tracked")
	return nil
}
