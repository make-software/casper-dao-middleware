package config

import (
	"errors"
	"fmt"

	"casper-dao-middleware/pkg/casper"
	"casper-dao-middleware/pkg/casper/types"

	"github.com/go-ozzo/ozzo-validation"
)

const variableRepositoryContractStorageUrefName = "storage__repository__contract"

type DAOContractsMetadata struct {
	VariableRepositoryContractPackageHash types.Hash
	VariableRepositoryContractHash        types.Hash
	VariableRepositoryContractStorageUref string

	ReputationContractPackageHash types.Hash
	ReputationContractHash        types.Hash

	SimpleVoterContractPackageHash types.Hash
	SimpleVoterContractHash        types.Hash

	RepoVoterContractPackageHash types.Hash
	RepoVoterContractHash        types.Hash

	ReputationVoterContractPackageHash types.Hash
	ReputationVoterContractHash        types.Hash

	SlashingVoterContractPackageHash types.Hash
	SlashingVoterContractHash        types.Hash

	KycVoterContractPackageHash types.Hash
	KycVoterContractHash        types.Hash

	VANFTContractPackageHash types.Hash
	VANFTContractHash        types.Hash

	KycNFTContractPackageHash types.Hash
	KycNFTContractHash        types.Hash
}

func NewDAOContractsMetadataFromHashesMap(contractHashes map[string]types.Hash, casperClient casper.RPCClient) (DAOContractsMetadata, error) {
	result := DAOContractsMetadata{}
	stateRootHash, err := casperClient.GetStateRootHashByHash("")
	if err != nil {
		return DAOContractsMetadata{}, err
	}

	for contractName, contractHashHex := range contractHashes {
		stateItemRes, err := casperClient.GetStateItem(stateRootHash.StateRootHash, fmt.Sprintf("hash-%s", contractHashHex), []string{})
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
		case "variable_repository_contract":
			result.VariableRepositoryContractPackageHash = contractPackageHash
			result.VariableRepositoryContractHash = contractHashHex
			for _, namedKey := range stateItemRes.StoredValue.Contract.NamedKeys {
				if namedKey.Name == variableRepositoryContractStorageUrefName {
					result.VariableRepositoryContractStorageUref = namedKey.Key
					break
				}
			}

			if result.VariableRepositoryContractStorageUref == "" {
				return DAOContractsMetadata{}, errors.New("error: missing variable repository contract storage uref in contract")
			}
		}

	}

	return result, result.Validate()
}

func (d DAOContractsMetadata) ContractHashes() []types.Hash {
	return []types.Hash{
		d.ReputationContractHash,
		d.VANFTContractHash,
		d.KycNFTContractHash,
		d.SimpleVoterContractHash,
		d.KycVoterContractHash,
		d.SlashingVoterContractHash,
		d.ReputationVoterContractHash,
		d.RepoVoterContractHash,
		d.VariableRepositoryContractHash,
	}
}

func (d DAOContractsMetadata) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.ReputationContractHash, validation.Required),
		validation.Field(&d.SimpleVoterContractHash, validation.Required),
		validation.Field(&d.RepoVoterContractHash, validation.Required),
		validation.Field(&d.ReputationVoterContractHash, validation.Required),
		validation.Field(&d.SlashingVoterContractHash, validation.Required),
		validation.Field(&d.KycVoterContractHash, validation.Required),
		validation.Field(&d.VANFTContractHash, validation.Required),
		validation.Field(&d.KycNFTContractHash, validation.Required),
		validation.Field(&d.VariableRepositoryContractHash, validation.Required),
		validation.Field(&d.VariableRepositoryContractStorageUref, validation.Required),
	)
}
