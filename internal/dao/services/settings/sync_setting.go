package settings

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/blake2b"

	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/types"
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

	settingItemKey, err := toDictionaryItemKey(c.setting)
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

	record, err := types.NewRecordFromBytes(result.StoredValue.CLValue.Bytes)
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

func toDictionaryItemKey(key string) (string, error) {
	res := make([]byte, 0)
	blake, err := blake2b.New256(res)
	if err != nil {
		return "", err
	}
	keyBytes := []byte(key)
	blake.Write(binary.LittleEndian.AppendUint32(nil, uint32(len(keyBytes))))
	blake.Write(keyBytes)
	return hex.EncodeToString(blake.Sum(nil)), nil
}
