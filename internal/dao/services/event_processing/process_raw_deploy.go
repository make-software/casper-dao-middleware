package event_processing

import (
	"errors"

	"go.uber.org/zap"

	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/services/event_tracking/kyc_nft"
	"casper-dao-middleware/internal/dao/services/event_tracking/kyc_voter"
	"casper-dao-middleware/internal/dao/services/event_tracking/repo_voter"
	"casper-dao-middleware/internal/dao/services/event_tracking/reputation"
	"casper-dao-middleware/internal/dao/services/event_tracking/reputation_voter"
	"casper-dao-middleware/internal/dao/services/event_tracking/simple_voter"
	"casper-dao-middleware/internal/dao/services/event_tracking/slashing_voter"
	"casper-dao-middleware/internal/dao/services/event_tracking/va_nft"
	"casper-dao-middleware/internal/dao/services/event_tracking/varaible_repository"
	"casper-dao-middleware/pkg/go-ces-parser"
)

type ProcessRawDeploy struct {
	di.EntityManagerAware
	di.DeployProcessedEventAware
	di.DAOContractsMetadataAware

	cesParser *ces.EventParser
}

func NewProcessRawDeploy() ProcessRawDeploy {
	return ProcessRawDeploy{}
}

func (c *ProcessRawDeploy) SetCESEventParser(parser *ces.EventParser) {
	c.cesParser = parser
}

func (c *ProcessRawDeploy) Execute() error {
	deployProcessedEvent := c.GetDeployProcessedEvent()
	daoContractsMetadata := c.GetDAOContractsMetadata()

	results, err := c.cesParser.ParseExecutionResults(deployProcessedEvent.DeployProcessed.ExecutionResult)
	if err != nil {
		return err
	}

	for _, result := range results {
		if result.Error != nil {
			zap.S().With(zap.Error(err)).Error("Failed to parse ces events")
			return err
		}
	}

	for _, result := range results {
		switch result.Event.ContractPackageHash.ToHex() {
		case daoContractsMetadata.KycNFTContractPackageHash.ToHex():
			trackKycNFTContract := kyc_nft.NewTrackContract()
			trackKycNFTContract.SetCESEvent(result.Event)
			trackKycNFTContract.SetEntityManager(c.GetEntityManager())
			err = trackKycNFTContract.Execute()
		case daoContractsMetadata.VANFTContractPackageHash.ToHex():
			trackVANFTContract := va_nft.NewTrackContract()
			trackVANFTContract.SetCESEvent(result.Event)
			trackVANFTContract.SetEntityManager(c.GetEntityManager())
			err = trackVANFTContract.Execute()
		case daoContractsMetadata.ReputationContractPackageHash.ToHex():
			trackReputationContract := reputation.NewTrackContract()
			trackReputationContract.SetCESEvent(result.Event)
			trackReputationContract.SetEntityManager(c.GetEntityManager())
			trackReputationContract.SetDeployProcessedEvent(c.GetDeployProcessedEvent())
			trackReputationContract.SetDAOContractsMetadata(daoContractsMetadata)
			err = trackReputationContract.Execute()
		case daoContractsMetadata.RepoVoterContractPackageHash.ToHex():
			trackRepoVoterContract := repo_voter.NewTrackContract()
			trackRepoVoterContract.SetCESEvent(result.Event)
			trackRepoVoterContract.SetEntityManager(c.GetEntityManager())
			trackRepoVoterContract.SetDAOContractsMetadata(daoContractsMetadata)
			trackRepoVoterContract.SetDeployProcessedEvent(c.GetDeployProcessedEvent())
			err = trackRepoVoterContract.Execute()
		case daoContractsMetadata.ReputationVoterContractPackageHash.ToHex():
			trackReputationVoterContract := reputation_voter.NewTrackContract()
			trackReputationVoterContract.SetCESEvent(result.Event)
			trackReputationVoterContract.SetEntityManager(c.GetEntityManager())
			trackReputationVoterContract.SetDAOContractsMetadata(daoContractsMetadata)
			trackReputationVoterContract.SetDeployProcessedEvent(c.GetDeployProcessedEvent())
			err = trackReputationVoterContract.Execute()
		case daoContractsMetadata.SimpleVoterContractPackageHash.ToHex():
			trackSimpleVoterContract := simple_voter.NewTrackContract()
			trackSimpleVoterContract.SetCESEvent(result.Event)
			trackSimpleVoterContract.SetEntityManager(c.GetEntityManager())
			trackSimpleVoterContract.SetDAOContractsMetadata(daoContractsMetadata)
			trackSimpleVoterContract.SetDeployProcessedEvent(c.GetDeployProcessedEvent())
			err = trackSimpleVoterContract.Execute()
		case daoContractsMetadata.SlashingVoterContractPackageHash.ToHex():
			trackSlashingVoterContract := slashing_voter.NewTrackContract()
			trackSlashingVoterContract.SetCESEvent(result.Event)
			trackSlashingVoterContract.SetEntityManager(c.GetEntityManager())
			trackSlashingVoterContract.SetDAOContractsMetadata(daoContractsMetadata)
			trackSlashingVoterContract.SetDeployProcessedEvent(c.GetDeployProcessedEvent())
			err = trackSlashingVoterContract.Execute()
		case daoContractsMetadata.KycVoterContractPackageHash.ToHex():
			trackKycVoterContract := kyc_voter.NewTrackContract()
			trackKycVoterContract.SetCESEvent(result.Event)
			trackKycVoterContract.SetEntityManager(c.GetEntityManager())
			trackKycVoterContract.SetDAOContractsMetadata(daoContractsMetadata)
			trackKycVoterContract.SetDeployProcessedEvent(c.GetDeployProcessedEvent())
			err = trackKycVoterContract.Execute()
		case daoContractsMetadata.VariableRepositoryContractPackageHash.ToHex():
			trackVariableRepositoryContract := varaible_repository.NewTrackContract()
			trackVariableRepositoryContract.SetCESEvent(result.Event)
			trackVariableRepositoryContract.SetEntityManager(c.GetEntityManager())
			trackVariableRepositoryContract.SetDAOContractsMetadata(daoContractsMetadata)
			trackVariableRepositoryContract.SetDeployProcessedEvent(c.GetDeployProcessedEvent())
			err = trackVariableRepositoryContract.Execute()
		default:
			return errors.New("unsupported DAO contract")
		}

		if err != nil {
			zap.S().With(zap.Error(err)).With(zap.String("event-name", result.Event.Name)).Info("Failed to handle DAO event")
			return err
		}
		zap.S().With("event", result.Event.Name).Info("Successfully tracked event")
	}

	return nil
}
