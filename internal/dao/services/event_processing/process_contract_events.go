package event_processing

import (
	"errors"
	"fmt"

	"go.uber.org/zap"

	"casper-dao-middleware/internal/dao/events/bid_escrow"
	"casper-dao-middleware/internal/dao/services/bid"
	"casper-dao-middleware/internal/dao/services/job_offer"
	"casper-dao-middleware/internal/dao/services/jobs"
	"casper-dao-middleware/internal/dao/services/settings"
	"casper-dao-middleware/internal/dao/services/votes"

	"casper-dao-middleware/internal/dao/services/voting"

	"casper-dao-middleware/internal/dao/services/account"
	"casper-dao-middleware/internal/dao/services/reputation"

	"casper-dao-middleware/internal/dao/di"
	"casper-dao-middleware/internal/dao/events/admin"
	base_events "casper-dao-middleware/internal/dao/events/base"
	"casper-dao-middleware/internal/dao/events/kyc_nft"
	"casper-dao-middleware/internal/dao/events/kyc_voter"
	"casper-dao-middleware/internal/dao/events/onboarding_request"
	"casper-dao-middleware/internal/dao/events/repo_voter"
	reputation_events "casper-dao-middleware/internal/dao/events/reputation"
	"casper-dao-middleware/internal/dao/events/reputation_voter"
	"casper-dao-middleware/internal/dao/events/simple_voter"
	"casper-dao-middleware/internal/dao/events/slashing_voter"
	"casper-dao-middleware/internal/dao/events/va_nft"
	"casper-dao-middleware/internal/dao/events/variable_repository"
	"casper-dao-middleware/pkg/go-ces-parser"
)

type ProcessContractEvents struct {
	di.EntityManagerAware
	di.CESEventAware
	di.DeployProcessedEventAware
	di.DAOContractsMetadataAware
}

func NewProcessContractEvents() *ProcessContractEvents {
	return &ProcessContractEvents{}
}

func (s *ProcessContractEvents) Execute() error {
	cesEvent := s.GetCESEvent()
	doaContractMetadata := s.GetDAOContractsMetadata()

	switch cesEvent.ContractPackageHash.ToHex() {
	case doaContractMetadata.KycNFTContractPackageHash.ToHex():
		return s.trackKycNFTContract(cesEvent)
	case doaContractMetadata.VANFTContractPackageHash.ToHex():
		return s.trackVANFTContract(cesEvent)
	case doaContractMetadata.ReputationContractPackageHash.ToHex():
		return s.trackReputationContract(cesEvent)
	case doaContractMetadata.RepoVoterContractPackageHash.ToHex():
		return s.trackRepoVoterContract(cesEvent)
	case doaContractMetadata.ReputationVoterContractPackageHash.ToHex():
		return s.trackReputationVoterContract(cesEvent)
	case doaContractMetadata.SimpleVoterContractPackageHash.ToHex():
		return s.trackSimpleVoterContract(cesEvent)
	case doaContractMetadata.SlashingVoterContractPackageHash.ToHex():
		return s.trackSlashingVoterContract(cesEvent)
	case doaContractMetadata.KycVoterContractPackageHash.ToHex():
		return s.trackKycVoterContract(cesEvent)
	case doaContractMetadata.VariableRepositoryContractPackageHash.ToHex():
		return s.trackVariableRepositoryContract(cesEvent)
	case doaContractMetadata.OnboardingRequestContractPackageHash.ToHex():
		return s.trackOnboardingRequestContract(cesEvent)
	case doaContractMetadata.AdminContractPackageHash.ToHex():
		return s.trackAdminContract(cesEvent)
	case doaContractMetadata.BidEscrowContractPackageHash.ToHex():
		return s.trackBidEscrowRepositoryContract(cesEvent)
	default:
		return errors.New("unsupported DAO contract")
	}
}

func (s *ProcessContractEvents) trackKycNFTContract(cesEvent ces.Event) error {
	daoContractMetadata := s.GetDAOContractsMetadata()

	zap.S().With(zap.String("event", cesEvent.Name)).
		With(zap.String("contract", daoContractMetadata.VANFTContractHash.String())).Info("New KYC NFT Contract event")

	switch cesEvent.Name {
	case kyc_nft.TransferEventName:
		trackTransfer := account.NewTrackKycTransfer()
		trackTransfer.SetCESEvent(cesEvent)
		trackTransfer.SetEntityManager(s.GetEntityManager())
		if err := trackTransfer.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.SlashingVoterContractHash.String())).Info("failed to track event")
			return err
		}
	default:
		return fmt.Errorf("unsupported contract event - %s", cesEvent.Name)
	}

	return nil
}

func (s *ProcessContractEvents) trackVANFTContract(cesEvent ces.Event) error {
	daoContractMetadata := s.GetDAOContractsMetadata()

	zap.S().With(zap.String("event", cesEvent.Name)).
		With(zap.String("contract", daoContractMetadata.VANFTContractHash.String())).Info("New VA NFT Contract event")

	switch cesEvent.Name {
	case va_nft.TransferEventName:
		trackTransfer := account.NewTrackVATransfer()
		trackTransfer.SetCESEvent(cesEvent)
		trackTransfer.SetEntityManager(s.GetEntityManager())
		if err := trackTransfer.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.VANFTContractHash.String())).Info("failed to track event")
			return err
		}
	default:
		return fmt.Errorf("unsupported contract event - %s", cesEvent.Name)
	}

	return nil
}

func (s *ProcessContractEvents) trackReputationContract(cesEvent ces.Event) error {
	daoContractMetadata := s.GetDAOContractsMetadata()

	zap.S().With(zap.String("event", cesEvent.Name)).
		With(zap.String("contract", daoContractMetadata.VANFTContractHash.String())).Info("New Reputation Contract event")

	switch cesEvent.Name {
	case reputation_events.MintEventName:
		trackMint := reputation.NewTrackMint()
		trackMint.SetCESEvent(cesEvent)
		trackMint.SetEntityManager(s.GetEntityManager())
		trackMint.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		trackMint.SetDAOContractsMetadata(daoContractMetadata)
		if err := trackMint.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.ReputationContractHash.String())).Info("failed to track event")
			return err
		}
	case reputation_events.BurnEventName:
		trackBurn := reputation.NewTrackBurn()
		trackBurn.SetCESEvent(cesEvent)
		trackBurn.SetEntityManager(s.GetEntityManager())
		trackBurn.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		trackBurn.SetDAOContractsMetadata(daoContractMetadata)
		if err := trackBurn.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.ReputationContractHash.String())).Info("failed to track event")
			return err
		}
	default:
		return fmt.Errorf("unsupported contract event - %s", cesEvent.Name)
	}

	return nil
}

func (s *ProcessContractEvents) trackRepoVoterContract(cesEvent ces.Event) error {
	daoContractMetadata := s.GetDAOContractsMetadata()

	zap.S().With(zap.String("event", cesEvent.Name)).
		With(zap.String("contract", daoContractMetadata.VANFTContractHash.String())).Info("New Repo Voter Contract event")

	switch cesEvent.Name {
	case repo_voter.VotingCreatedEventName:
		trackVotingCreated := voting.NewTrackRepoVotingCreated()
		trackVotingCreated.SetCESEvent(cesEvent)
		trackVotingCreated.SetEntityManager(s.GetEntityManager())
		trackVotingCreated.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		if err := trackVotingCreated.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.RepoVoterContractHash.String())).Info("failed to track event")
			return err
		}
	case base_events.VotingEndedEventName:
		trackVotingEnded := voting.NewTrackVotingEnded()
		trackVotingEnded.SetCESEvent(cesEvent)
		trackVotingEnded.SetEntityManager(s.GetEntityManager())
		trackVotingEnded.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		trackVotingEnded.SetDAOContractsMetadata(daoContractMetadata)
		if err := trackVotingEnded.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.RepoVoterContractHash.String())).Info("failed to track event")
			return err
		}
	case base_events.VotingCanceledEventName:
		trackVotingCanceled := voting.NewTrackVotingCanceled()
		trackVotingCanceled.SetCESEvent(cesEvent)
		trackVotingCanceled.SetEntityManager(s.GetEntityManager())
		trackVotingCanceled.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		trackVotingCanceled.SetDAOContractsMetadata(daoContractMetadata)
		if err := trackVotingCanceled.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.RepoVoterContractHash.String())).Info("failed to track event")
			return err
		}
	case base_events.BallotCastEventName:
		trackBallotCast := votes.NewTrackVote()
		trackBallotCast.SetCESEvent(cesEvent)
		trackBallotCast.SetEntityManager(s.GetEntityManager())
		trackBallotCast.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		trackBallotCast.SetDAOContractsMetadata(daoContractMetadata)
		if err := trackBallotCast.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.RepoVoterContractHash.String())).Info("failed to track event")
			return err
		}
	case base_events.BallotCanceledEventName:
		trackBallotCanceled := votes.NewTrackCanceledVote()
		trackBallotCanceled.SetCESEvent(cesEvent)
		trackBallotCanceled.SetEntityManager(s.GetEntityManager())
		if err := trackBallotCanceled.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.RepoVoterContractHash.String())).Info("failed to track event")
			return err
		}
	default:
		return fmt.Errorf("unsupported contract event - %s", cesEvent.Name)
	}

	return nil
}

func (s *ProcessContractEvents) trackReputationVoterContract(cesEvent ces.Event) error {
	daoContractMetadata := s.GetDAOContractsMetadata()

	zap.S().With(zap.String("event", cesEvent.Name)).
		With(zap.String("contract", daoContractMetadata.VANFTContractHash.String())).Info("New Reputation Voter Contract event")

	switch cesEvent.Name {
	case reputation_voter.VotingCreatedEventName:
		trackVotingCreated := voting.NewTrackReputationVotingCreated()
		trackVotingCreated.SetCESEvent(cesEvent)
		trackVotingCreated.SetEntityManager(s.GetEntityManager())
		trackVotingCreated.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		if err := trackVotingCreated.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.ReputationVoterContractHash.String())).Info("failed to track event")
			return err
		}
	case base_events.VotingEndedEventName:
		trackVotingEnded := voting.NewTrackVotingEnded()
		trackVotingEnded.SetCESEvent(cesEvent)
		trackVotingEnded.SetEntityManager(s.GetEntityManager())
		trackVotingEnded.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		trackVotingEnded.SetDAOContractsMetadata(daoContractMetadata)
		if err := trackVotingEnded.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.ReputationVoterContractHash.String())).Info("failed to track event")
			return err
		}
	case base_events.VotingCanceledEventName:
		trackVotingCanceled := voting.NewTrackVotingCanceled()
		trackVotingCanceled.SetCESEvent(cesEvent)
		trackVotingCanceled.SetEntityManager(s.GetEntityManager())
		trackVotingCanceled.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		trackVotingCanceled.SetDAOContractsMetadata(daoContractMetadata)
		if err := trackVotingCanceled.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.ReputationVoterContractHash.String())).Info("failed to track event")
			return err
		}
	case base_events.BallotCastEventName:
		trackBallotCast := votes.NewTrackVote()
		trackBallotCast.SetCESEvent(cesEvent)
		trackBallotCast.SetEntityManager(s.GetEntityManager())
		trackBallotCast.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		trackBallotCast.SetDAOContractsMetadata(daoContractMetadata)
		if err := trackBallotCast.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.ReputationVoterContractHash.String())).Info("failed to track event")
			return err
		}
	case base_events.BallotCanceledEventName:
		trackBallotCanceled := votes.NewTrackCanceledVote()
		trackBallotCanceled.SetCESEvent(cesEvent)
		trackBallotCanceled.SetEntityManager(s.GetEntityManager())
		if err := trackBallotCanceled.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.ReputationVoterContractHash.String())).Info("failed to track event")
			return err
		}
	default:
		return fmt.Errorf("unsupported contract event - %s", cesEvent.Name)
	}

	return nil
}

func (s *ProcessContractEvents) trackSimpleVoterContract(cesEvent ces.Event) error {
	daoContractMetadata := s.GetDAOContractsMetadata()

	zap.S().With(zap.String("event", cesEvent.Name)).
		With(zap.String("contract", daoContractMetadata.VANFTContractHash.String())).Info("New Simple Voter Contract event")

	switch cesEvent.Name {
	case simple_voter.VotingCreatedEventName:
		trackVotingCreated := voting.NewTrackSimpleVotingCreated()
		trackVotingCreated.SetCESEvent(cesEvent)
		trackVotingCreated.SetEntityManager(s.GetEntityManager())
		trackVotingCreated.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		if err := trackVotingCreated.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.SimpleVoterContractHash.String())).Info("failed to track event")
			return err
		}
	case base_events.VotingEndedEventName:
		trackVotingEnded := voting.NewTrackVotingEnded()
		trackVotingEnded.SetCESEvent(cesEvent)
		trackVotingEnded.SetEntityManager(s.GetEntityManager())
		trackVotingEnded.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		trackVotingEnded.SetDAOContractsMetadata(daoContractMetadata)
		if err := trackVotingEnded.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.SimpleVoterContractHash.String())).Info("failed to track event")
			return err
		}
	case base_events.VotingCanceledEventName:
		trackVotingCanceled := voting.NewTrackVotingCanceled()
		trackVotingCanceled.SetCESEvent(cesEvent)
		trackVotingCanceled.SetEntityManager(s.GetEntityManager())
		trackVotingCanceled.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		trackVotingCanceled.SetDAOContractsMetadata(daoContractMetadata)
		if err := trackVotingCanceled.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.SimpleVoterContractHash.String())).Info("failed to track event")
			return err
		}
	case base_events.BallotCastEventName:
		trackBallotCast := votes.NewTrackVote()
		trackBallotCast.SetCESEvent(cesEvent)
		trackBallotCast.SetEntityManager(s.GetEntityManager())
		trackBallotCast.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		trackBallotCast.SetDAOContractsMetadata(daoContractMetadata)
		//if err := trackBallotCast.Execute(); err != nil {
		//	zap.S().With(zap.String("event", cesEvent.Name)).
		//		With(zap.String("contract", daoContractMetadata.SimpleVoterContractHash.String())).Info("failed to track event")
		//	return err
		//}
	case base_events.BallotCanceledEventName:
		trackBallotCanceled := votes.NewTrackCanceledVote()
		trackBallotCanceled.SetCESEvent(cesEvent)
		trackBallotCanceled.SetEntityManager(s.GetEntityManager())
		if err := trackBallotCanceled.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.SimpleVoterContractHash.String())).Info("failed to track event")
			return err
		}
	default:
		return fmt.Errorf("unsupported contract event - %s", cesEvent.Name)
	}

	return nil
}

func (s *ProcessContractEvents) trackSlashingVoterContract(cesEvent ces.Event) error {
	daoContractMetadata := s.GetDAOContractsMetadata()

	zap.S().With(zap.String("event", cesEvent.Name)).
		With(zap.String("contract", daoContractMetadata.VANFTContractHash.String())).Info("New Slashing Voter Contract event")

	switch cesEvent.Name {
	case slashing_voter.VotingCreatedEventName:
		trackVotingCreated := voting.NewTrackSlashingVotingCreated()
		trackVotingCreated.SetCESEvent(cesEvent)
		trackVotingCreated.SetEntityManager(s.GetEntityManager())
		trackVotingCreated.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		if err := trackVotingCreated.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.SlashingVoterContractHash.String())).Info("failed to track event")
			return err
		}
	case base_events.VotingEndedEventName:
		trackVotingEnded := voting.NewTrackVotingEnded()
		trackVotingEnded.SetCESEvent(cesEvent)
		trackVotingEnded.SetEntityManager(s.GetEntityManager())
		trackVotingEnded.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		trackVotingEnded.SetDAOContractsMetadata(daoContractMetadata)
		if err := trackVotingEnded.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.SlashingVoterContractHash.String())).Info("failed to track event")
			return err
		}
	case base_events.VotingCanceledEventName:
		trackVotingCanceled := voting.NewTrackVotingCanceled()
		trackVotingCanceled.SetCESEvent(cesEvent)
		trackVotingCanceled.SetEntityManager(s.GetEntityManager())
		trackVotingCanceled.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		trackVotingCanceled.SetDAOContractsMetadata(daoContractMetadata)
		if err := trackVotingCanceled.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.SlashingVoterContractHash.String())).Info("failed to track event")
			return err
		}
	case base_events.BallotCastEventName:
		trackBallotCast := votes.NewTrackVote()
		trackBallotCast.SetCESEvent(cesEvent)
		trackBallotCast.SetEntityManager(s.GetEntityManager())
		trackBallotCast.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		trackBallotCast.SetDAOContractsMetadata(daoContractMetadata)
		if err := trackBallotCast.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.SlashingVoterContractHash.String())).Info("failed to track event")
			return err
		}
	case base_events.BallotCanceledEventName:
		trackBallotCanceled := votes.NewTrackCanceledVote()
		trackBallotCanceled.SetCESEvent(cesEvent)
		trackBallotCanceled.SetEntityManager(s.GetEntityManager())
		if err := trackBallotCanceled.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.SlashingVoterContractHash.String())).Info("failed to track event")
			return err
		}
	default:
		return fmt.Errorf("unsupported contract event - %s", cesEvent.Name)
	}

	return nil
}

func (s *ProcessContractEvents) trackKycVoterContract(cesEvent ces.Event) error {
	daoContractMetadata := s.GetDAOContractsMetadata()

	zap.S().With(zap.String("event", cesEvent.Name)).
		With(zap.String("contract", daoContractMetadata.VANFTContractHash.String())).Info("New KYC Voter Contract event")

	switch cesEvent.Name {
	case kyc_voter.VotingCreatedEventName:
		trackVotingCreated := voting.NewTrackKycVotingCreated()
		trackVotingCreated.SetCESEvent(cesEvent)
		trackVotingCreated.SetEntityManager(s.GetEntityManager())
		trackVotingCreated.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		if err := trackVotingCreated.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.KycVoterContractHash.String())).Info("failed to track event")
			return err
		}
	case base_events.VotingEndedEventName:
		trackVotingEnded := voting.NewTrackVotingEnded()
		trackVotingEnded.SetCESEvent(cesEvent)
		trackVotingEnded.SetEntityManager(s.GetEntityManager())
		trackVotingEnded.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		trackVotingEnded.SetDAOContractsMetadata(daoContractMetadata)
		if err := trackVotingEnded.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.KycVoterContractHash.String())).Info("failed to track event")
			return err
		}
	case base_events.VotingCanceledEventName:
		trackVotingCanceled := voting.NewTrackVotingCanceled()
		trackVotingCanceled.SetCESEvent(cesEvent)
		trackVotingCanceled.SetEntityManager(s.GetEntityManager())
		trackVotingCanceled.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		trackVotingCanceled.SetDAOContractsMetadata(daoContractMetadata)
		if err := trackVotingCanceled.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.KycVoterContractHash.String())).Info("failed to track event")
			return err
		}
	case base_events.BallotCastEventName:
		trackBallotCast := votes.NewTrackVote()
		trackBallotCast.SetCESEvent(cesEvent)
		trackBallotCast.SetEntityManager(s.GetEntityManager())
		trackBallotCast.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		trackBallotCast.SetDAOContractsMetadata(daoContractMetadata)
		if err := trackBallotCast.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.KycVoterContractHash.String())).Info("failed to track event")
			return err
		}
	case base_events.BallotCanceledEventName:
		trackBallotCanceled := votes.NewTrackCanceledVote()
		trackBallotCanceled.SetCESEvent(cesEvent)
		trackBallotCanceled.SetEntityManager(s.GetEntityManager())
		if err := trackBallotCanceled.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.KycVoterContractHash.String())).Info("failed to track event")
			return err
		}
	default:
		return fmt.Errorf("unsupported contract event - %s", cesEvent.Name)
	}

	return nil
}

func (s *ProcessContractEvents) trackOnboardingRequestContract(cesEvent ces.Event) error {
	daoContractMetadata := s.GetDAOContractsMetadata()

	zap.S().With(zap.String("event", cesEvent.Name)).
		With(zap.String("contract", daoContractMetadata.VANFTContractHash.String())).Info("New Onboarding Request Contract event")

	switch cesEvent.Name {
	case onboarding_request.VotingCreatedEventName:
		trackVotingCreated := voting.NewTrackOnboardingVotingCreated()
		trackVotingCreated.SetCESEvent(cesEvent)
		trackVotingCreated.SetEntityManager(s.GetEntityManager())
		trackVotingCreated.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		if err := trackVotingCreated.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.OnboardingRequestContractHash.String())).Info("failed to track event")
			return err
		}
	case base_events.VotingEndedEventName:
		trackVotingEnded := voting.NewTrackVotingEnded()
		trackVotingEnded.SetCESEvent(cesEvent)
		trackVotingEnded.SetEntityManager(s.GetEntityManager())
		trackVotingEnded.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		trackVotingEnded.SetDAOContractsMetadata(daoContractMetadata)
		if err := trackVotingEnded.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.OnboardingRequestContractHash.String())).Info("failed to track event")
			return err
		}
	case base_events.VotingCanceledEventName:
		trackVotingCanceled := voting.NewTrackVotingCanceled()
		trackVotingCanceled.SetCESEvent(cesEvent)
		trackVotingCanceled.SetEntityManager(s.GetEntityManager())
		trackVotingCanceled.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		trackVotingCanceled.SetDAOContractsMetadata(daoContractMetadata)
		if err := trackVotingCanceled.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.OnboardingRequestContractHash.String())).Info("failed to track event")
			return err
		}
	case base_events.BallotCastEventName:
		trackBallotCast := votes.NewTrackVote()
		trackBallotCast.SetCESEvent(cesEvent)
		trackBallotCast.SetEntityManager(s.GetEntityManager())
		trackBallotCast.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		trackBallotCast.SetDAOContractsMetadata(daoContractMetadata)
		if err := trackBallotCast.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.OnboardingRequestContractHash.String())).Info("failed to track event")
			return err
		}
	case base_events.BallotCanceledEventName:
		trackBallotCanceled := votes.NewTrackCanceledVote()
		trackBallotCanceled.SetCESEvent(cesEvent)
		trackBallotCanceled.SetEntityManager(s.GetEntityManager())
		if err := trackBallotCanceled.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.OnboardingRequestContractHash.String())).Info("failed to track event")
			return err
		}

	default:
		return fmt.Errorf("unsupported contract event - %s", cesEvent.Name)
	}

	return nil
}

func (s *ProcessContractEvents) trackAdminContract(cesEvent ces.Event) error {
	daoContractMetadata := s.GetDAOContractsMetadata()

	zap.S().With(zap.String("event", cesEvent.Name)).
		With(zap.String("contract", daoContractMetadata.VANFTContractHash.String())).Info("New Admin Contract event")

	switch cesEvent.Name {
	case admin.VotingCreatedEventName:
		trackVotingCreated := voting.NewTrackAdminVotingCreated()
		trackVotingCreated.SetCESEvent(cesEvent)
		trackVotingCreated.SetEntityManager(s.GetEntityManager())
		trackVotingCreated.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		if err := trackVotingCreated.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.AdminContractHash.String())).Info("failed to track event")
			return err
		}
	case base_events.VotingEndedEventName:
		trackVotingEnded := voting.NewTrackVotingEnded()
		trackVotingEnded.SetCESEvent(cesEvent)
		trackVotingEnded.SetEntityManager(s.GetEntityManager())
		trackVotingEnded.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		trackVotingEnded.SetDAOContractsMetadata(daoContractMetadata)
		if err := trackVotingEnded.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.AdminContractHash.String())).Info("failed to track event")
			return err
		}
	case base_events.VotingCanceledEventName:
		trackVotingCanceled := voting.NewTrackVotingCanceled()
		trackVotingCanceled.SetCESEvent(cesEvent)
		trackVotingCanceled.SetEntityManager(s.GetEntityManager())
		trackVotingCanceled.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		trackVotingCanceled.SetDAOContractsMetadata(daoContractMetadata)
		if err := trackVotingCanceled.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.AdminContractHash.String())).Info("failed to track event")
			return err
		}
	case base_events.BallotCastEventName:
		trackBallotCast := votes.NewTrackVote()
		trackBallotCast.SetCESEvent(cesEvent)
		trackBallotCast.SetEntityManager(s.GetEntityManager())
		trackBallotCast.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		trackBallotCast.SetDAOContractsMetadata(daoContractMetadata)
		if err := trackBallotCast.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.AdminContractHash.String())).Info("failed to track event")
			return err
		}
	case base_events.BallotCanceledEventName:
		trackBallotCanceled := votes.NewTrackCanceledVote()
		trackBallotCanceled.SetCESEvent(cesEvent)
		trackBallotCanceled.SetEntityManager(s.GetEntityManager())
		if err := trackBallotCanceled.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.AdminContractHash.String())).Info("failed to track event")
			return err
		}
	default:
		return fmt.Errorf("unsupported contract event - %s", cesEvent.Name)
	}

	return nil
}

func (s *ProcessContractEvents) trackVariableRepositoryContract(cesEvent ces.Event) error {
	daoContractMetadata := s.GetDAOContractsMetadata()

	zap.S().With(zap.String("event", cesEvent.Name)).
		With(zap.String("contract", daoContractMetadata.VANFTContractHash.String())).Info("New VariableRepository Contract event")

	switch cesEvent.Name {
	case variable_repository.ValueUpdatedEventName:
		trackValueUpdated := settings.NewTrackUpdatedSetting()
		trackValueUpdated.SetCESEvent(cesEvent)
		trackValueUpdated.SetEntityManager(s.GetEntityManager())
		if err := trackValueUpdated.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.VariableRepositoryContractHash.String())).Info("failed to track event")
			return err
		}

	default:
		return fmt.Errorf("unsupported contract event - %s", cesEvent.Name)
	}

	return nil
}

func (s *ProcessContractEvents) trackBidEscrowRepositoryContract(cesEvent ces.Event) error {
	daoContractMetadata := s.GetDAOContractsMetadata()

	zap.S().With(zap.String("event", cesEvent.Name)).
		With(zap.String("contract", daoContractMetadata.BidEscrowContractHash.String())).Info("New BidEscrow Contract event")

	switch cesEvent.Name {
	case bid_escrow.JobOfferCreatedEventName:
		trackJobOffer := job_offer.NewTrackJobOfferCreated()
		trackJobOffer.SetCESEvent(cesEvent)
		trackJobOffer.SetEntityManager(s.GetEntityManager())
		trackJobOffer.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		if err := trackJobOffer.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.BidEscrowContractHash.String())).Info("failed to track event")
			return err
		}
	case bid_escrow.BidSubmittedEventName:
		trackSubmittedBid := bid.NewTrackBidSubmitted()
		trackSubmittedBid.SetCESEvent(cesEvent)
		trackSubmittedBid.SetEntityManager(s.GetEntityManager())
		trackSubmittedBid.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		if err := trackSubmittedBid.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.BidEscrowContractHash.String())).Info("failed to track event")
			return err
		}
	case bid_escrow.JobCreatedEventName:
		trackJobCreated := jobs.NewTrackJobCreated()
		trackJobCreated.SetCESEvent(cesEvent)
		trackJobCreated.SetEntityManager(s.GetEntityManager())
		trackJobCreated.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		if err := trackJobCreated.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.BidEscrowContractHash.String())).Info("failed to track event")
			return err
		}
	case bid_escrow.JobSubmittedEventName:
		trackJobSubmitted := jobs.NewTrackJobSubmitted()
		trackJobSubmitted.SetCESEvent(cesEvent)
		trackJobSubmitted.SetEntityManager(s.GetEntityManager())
		trackJobSubmitted.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		if err := trackJobSubmitted.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.BidEscrowContractHash.String())).Info("failed to track event")
			return err
		}
	case bid_escrow.JobCancelledEventName:
		trackJobCancelled := jobs.NewTrackJobCancelled()
		trackJobCancelled.SetCESEvent(cesEvent)
		trackJobCancelled.SetEntityManager(s.GetEntityManager())
		trackJobCancelled.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		if err := trackJobCancelled.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.BidEscrowContractHash.String())).Info("failed to track event")
			return err
		}
	case bid_escrow.JobRejectedEventName:
		trackJobRejected := jobs.NewTrackJobRejected()
		trackJobRejected.SetCESEvent(cesEvent)
		trackJobRejected.SetEntityManager(s.GetEntityManager())
		trackJobRejected.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		if err := trackJobRejected.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.BidEscrowContractHash.String())).Info("failed to track event")
			return err
		}
	case bid_escrow.JobDoneEventName:
		trackJobDone := jobs.NewTrackJobDone()
		trackJobDone.SetCESEvent(cesEvent)
		trackJobDone.SetEntityManager(s.GetEntityManager())
		trackJobDone.SetDeployProcessedEvent(s.GetDeployProcessedEvent())
		if err := trackJobDone.Execute(); err != nil {
			zap.S().With(zap.String("event", cesEvent.Name)).
				With(zap.String("contract", daoContractMetadata.BidEscrowContractHash.String())).Info("failed to track event")
			return err
		}

	default:
		return fmt.Errorf("unsupported contract event - %s", cesEvent.Name)
	}

	return nil
}
