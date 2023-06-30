package utils

import (
	"context"
	"errors"
	"fmt"

	"github.com/make-software/casper-go-sdk/casper"

	"casper-dao-middleware/pkg/config"
)

type DAOContractsMetadata struct {
	VariableRepositoryContractPackageHash casper.ContractPackageHash
	VariableRepositoryContractHash        casper.Hash

	ReputationContractPackageHash casper.ContractPackageHash
	ReputationContractHash        casper.Hash

	SimpleVoterContractPackageHash casper.ContractPackageHash
	SimpleVoterContractHash        casper.Hash

	RepoVoterContractPackageHash casper.ContractPackageHash
	RepoVoterContractHash        casper.Hash

	ReputationVoterContractPackageHash casper.ContractPackageHash
	ReputationVoterContractHash        casper.Hash

	SlashingVoterContractPackageHash casper.ContractPackageHash
	SlashingVoterContractHash        casper.Hash

	KycVoterContractPackageHash casper.ContractPackageHash
	KycVoterContractHash        casper.Hash

	VANFTContractPackageHash casper.ContractPackageHash
	VANFTContractHash        casper.Hash

	KycNFTContractPackageHash casper.ContractPackageHash
	KycNFTContractHash        casper.Hash

	OnboardingRequestContractPackageHash casper.ContractPackageHash
	OnboardingRequestContractHash        casper.Hash

	AdminContractPackageHash casper.ContractPackageHash
	AdminContractHash        casper.Hash

	BidEscrowContractPackageHash casper.ContractPackageHash
	BidEscrowContractHash        casper.Hash
}

func NewDAOContractsMetadata(contractHashes config.DaoContracts, casperClient casper.RPCClient) (DAOContractsMetadata, error) {
	result := DAOContractsMetadata{}
	stateRootHashRes, err := casperClient.GetStateRootHashLatest(context.Background())
	if err != nil {
		return DAOContractsMetadata{}, err
	}

	stateRootHash := stateRootHashRes.StateRootHash.String()

	for contractName, contractHashHex := range contractHashes.ToMap() {
		stateItemRes, err := casperClient.QueryGlobalStateByStateHash(context.Background(), &stateRootHash, fmt.Sprintf("hash-%s", contractHashHex), []string{})
		if err != nil {
			return DAOContractsMetadata{}, err
		}

		if stateItemRes.StoredValue.Contract == nil {
			return DAOContractsMetadata{}, errors.New("expected Contract StoredValue")
		}

		contractPackageHash := stateItemRes.StoredValue.Contract.ContractPackageHash

		switch contractName {
		case "reputation_contract":
			result.ReputationContractPackageHash = contractPackageHash
			result.ReputationContractHash = contractHashHex
		case "simple_voter_contract":
			result.SimpleVoterContractPackageHash = contractPackageHash
			result.SimpleVoterContractHash = contractHashHex
		case "repo_voter_contract":
			result.RepoVoterContractPackageHash = contractPackageHash
			result.RepoVoterContractHash = contractHashHex
		case "kyc_voter_contract":
			result.KycVoterContractPackageHash = contractPackageHash
			result.KycVoterContractHash = contractHashHex
		case "reputation_voter_contract":
			result.ReputationVoterContractPackageHash = contractPackageHash
			result.ReputationVoterContractHash = contractHashHex
		case "slashing_voter_contract":
			result.SlashingVoterContractPackageHash = contractPackageHash
			result.SlashingVoterContractHash = contractHashHex
		case "va_nft_contract":
			result.VANFTContractPackageHash = contractPackageHash
			result.VANFTContractHash = contractHashHex
		case "kyc_nft_contract":
			result.KycNFTContractPackageHash = contractPackageHash
			result.KycNFTContractHash = contractHashHex
		case "onboarding_request_contract":
			result.OnboardingRequestContractPackageHash = contractPackageHash
			result.OnboardingRequestContractHash = contractHashHex
		case "admin_contract":
			result.AdminContractPackageHash = contractPackageHash
			result.AdminContractHash = contractHashHex
		case "bid_escrow":
			result.BidEscrowContractPackageHash = contractPackageHash
			result.BidEscrowContractHash = contractHashHex
		case "variable_repository_contract":
			result.VariableRepositoryContractPackageHash = contractPackageHash
			result.VariableRepositoryContractHash = contractHashHex
		}

	}

	return result, nil
}

func (d DAOContractsMetadata) ContractHashes() []casper.Hash {
	return []casper.Hash{
		d.ReputationContractHash,
		d.VANFTContractHash,
		d.KycNFTContractHash,
		d.SimpleVoterContractHash,
		d.KycVoterContractHash,
		d.SlashingVoterContractHash,
		d.ReputationVoterContractHash,
		d.RepoVoterContractHash,
		d.VariableRepositoryContractHash,
		d.BidEscrowContractHash,
	}
}
