package event_processing

import (
	"errors"

	"go.uber.org/zap"

	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/services/event_tracking"
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
			trackKycNFTContract := event_tracking.NewTrackKycNFTContract()
			trackKycNFTContract.SetCESEvent(result.Event)
			trackKycNFTContract.SetEntityManager(c.GetEntityManager())
			err = trackKycNFTContract.Execute()

		case daoContractsMetadata.VANFTContractPackageHash.ToHex():
			trackVANFTContract := event_tracking.NewTrackVANFTContract()
			trackVANFTContract.SetCESEvent(result.Event)
			trackVANFTContract.SetEntityManager(c.GetEntityManager())
			err = trackVANFTContract.Execute()

		case daoContractsMetadata.ReputationContractPackageHash.ToHex():
			trackReputationContract := event_tracking.NewTrackReputationContract(daoContractsMetadata.ReputationContractPackageHash)
			trackReputationContract.SetCESEvent(result.Event)
			trackReputationContract.SetEntityManager(c.GetEntityManager())
			trackReputationContract.SetDeployProcessedEvent(c.GetDeployProcessedEvent())
			err = trackReputationContract.Execute()
		case daoContractsMetadata.RepoVoterContractPackageHash.ToHex():
			trackRepoVoterContract := event_tracking.NewTrackRepoVoterContract()
			trackRepoVoterContract.SetCESEvent(result.Event)
			trackRepoVoterContract.SetEntityManager(c.GetEntityManager())
			trackRepoVoterContract.SetDAOContractsMetadata(daoContractsMetadata)
			trackRepoVoterContract.SetDeployProcessedEvent(c.GetDeployProcessedEvent())
			err = trackRepoVoterContract.Execute()
		case daoContractsMetadata.ReputationVoterContractPackageHash.ToHex():
			trackReputationVoterContract := event_tracking.NewTrackReputationVoterContract()
			trackReputationVoterContract.SetCESEvent(result.Event)
			trackReputationVoterContract.SetEntityManager(c.GetEntityManager())
			trackReputationVoterContract.SetDAOContractsMetadata(daoContractsMetadata)
			trackReputationVoterContract.SetDeployProcessedEvent(c.GetDeployProcessedEvent())
			err = trackReputationVoterContract.Execute()
		case daoContractsMetadata.SimpleVoterContractPackageHash.ToHex():
			trackSimpleVoterContract := event_tracking.NewTrackSimpleVoterContract()
			trackSimpleVoterContract.SetCESEvent(result.Event)
			trackSimpleVoterContract.SetEntityManager(c.GetEntityManager())
			trackSimpleVoterContract.SetDAOContractsMetadata(daoContractsMetadata)
			trackSimpleVoterContract.SetDeployProcessedEvent(c.GetDeployProcessedEvent())
			err = trackSimpleVoterContract.Execute()

		case daoContractsMetadata.SlashingVoterContractPackageHash.ToHex():
			trackSlashingVoterContract := event_tracking.NewTrackSlashingVoterContract()
			trackSlashingVoterContract.SetCESEvent(result.Event)
			trackSlashingVoterContract.SetEntityManager(c.GetEntityManager())
			trackSlashingVoterContract.SetDAOContractsMetadata(daoContractsMetadata)
			trackSlashingVoterContract.SetDeployProcessedEvent(c.GetDeployProcessedEvent())
			err = trackSlashingVoterContract.Execute()
		case daoContractsMetadata.KycVoterContractPackageHash.ToHex():
			trackKycVoterContract := event_tracking.NewTrackKycVoterContract()
			trackKycVoterContract.SetCESEvent(result.Event)
			trackKycVoterContract.SetEntityManager(c.GetEntityManager())
			trackKycVoterContract.SetDAOContractsMetadata(daoContractsMetadata)
			trackKycVoterContract.SetDeployProcessedEvent(c.GetDeployProcessedEvent())
			err = trackKycVoterContract.Execute()
		case daoContractsMetadata.VariableRepositoryContractPackageHash.ToHex():
			trackVariableRepositoryContract := event_tracking.NewTrackVariableRepositoryContract()
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
	}

	return nil
}
