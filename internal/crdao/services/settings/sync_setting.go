package settings

import (
	"bytes"
	"encoding/hex"
	"errors"
	"time"

	"go.uber.org/zap"

	"casper-dao-middleware/internal/crdao/di"
	"casper-dao-middleware/internal/crdao/entities"
	"casper-dao-middleware/internal/crdao/types"
	"casper-dao-middleware/internal/crdao/utils"
)

const (
	PostJobDosFee                      = "PostJobDOSFee"
	InternalAuctionTime                = "InternalAuctionTime"
	PublicAuctionTime                  = "PublicAuctionTime"
	DefaultPolicingRate                = "DefaultPolicingRate"
	ReputationConversionRate           = "ReputationConversionRate"
	FiatConversionRateAddress          = "FiatConversionRateAddress"
	ForumKycRequired                   = "ForumKycRequired"
	BidEscrowInformalQuorumRatio       = "BidEscrowInformalQuorumRatio"
	BidEscrowFormalQuorumRatio         = "BidEscrowFormalQuorumRatio"
	BidEscrowFormalVotingTime          = "BidEscrowFormalVotingTime"
	BidEscrowInformalVotingTime        = "BidEscrowInformalVotingTime"
	FormalVotingTime                   = "FormalVotingTime"
	InformalVotingTime                 = "InformalVotingTime"
	FormalQuorumRatio                  = "FormalQuorumRatio"
	InformalQuorumRatio                = "InformalQuorumRatio"
	InformalStakeReputation            = "InformalStakeReputation"
	DistributePaymentToNonVoters       = "DistributePaymentToNonVoters"
	TimeBetweenInformalAndFormalVoting = "TimeBetweenInformalAndFormalVoting"
	VABidAcceptanceTimeout             = "VABidAcceptanceTimeout"
	VACanBidOnPublicAuction            = "VACanBidOnPublicAuction"
	BidEscrowWalletAddress             = "BidEscrowWalletAddress"
	BidEscrowPaymentRatio              = "BidEscrowPaymentRatio"
	VotingClearnessDelta               = "VotingClearnessDelta"
	VotingStartAfterJobSubmission      = "VotingStartAfterJobSubmission"
	DefaultReputationSlash             = "DefaultReputationSlash"
	VotingIdsAddress                   = "VotingIdsAddress"
)

var VariableRepoSettings = []string{
	PostJobDosFee,
	InternalAuctionTime,
	PublicAuctionTime,
	DefaultPolicingRate,
	ReputationConversionRate,
	FiatConversionRateAddress,
	ForumKycRequired,
	BidEscrowInformalQuorumRatio,
	BidEscrowFormalQuorumRatio,
	BidEscrowFormalVotingTime,
	BidEscrowInformalVotingTime,
	FormalVotingTime,
	InformalVotingTime,
	FormalQuorumRatio,
	InformalQuorumRatio,
	InformalStakeReputation,
	DistributePaymentToNonVoters,
	TimeBetweenInformalAndFormalVoting,
	VABidAcceptanceTimeout,
	VACanBidOnPublicAuction,
	BidEscrowWalletAddress,
	BidEscrowPaymentRatio,
	VotingClearnessDelta,
	VotingStartAfterJobSubmission,
	DefaultReputationSlash,
	VotingIdsAddress,
}

type SyncDAOSetting struct {
	di.EntityManagerAware
	di.CasperClientAware

	variableRepositoryContractStorageUref string
	setting                               string
}

func NewSyncDAOSetting() SyncDAOSetting {
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

	record, err := types.NewRecordFromBytes(decoded)
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
