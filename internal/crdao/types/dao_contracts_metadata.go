package types

import (
	"errors"
	"fmt"

	"casper-dao-middleware/pkg/casper"
	"casper-dao-middleware/pkg/casper/types"

	"github.com/go-ozzo/ozzo-validation"
)

const variableRepositoryContractStorageUrefName = "storage__repository__contract"

type DAOContractsMetadata struct {
	ReputationContractPackageHash types.Hash
	ReputationContractHash        types.Hash

	// There are many voter contracts
	SimpleVoterContractPackageHash types.Hash
	SimpleVoterContractHash        types.Hash

	VariableRepositoryContractPackageHash types.Hash
	VariableRepositoryContractHash        types.Hash

	VANFTContractPackageHash              types.Hash
	VANFTContractHash                     types.Hash
	VariableRepositoryContractStorageUref string
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
		case "voter_contract":
			result.SimpleVoterContractPackageHash = contractPackageHash
			result.SimpleVoterContractHash = contractHashHex
		case "va_nft_contract":
			result.VANFTContractPackageHash = contractPackageHash
			result.VANFTContractHash = contractHashHex
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

func (d DAOContractsMetadata) CESContracts() []types.Hash {
	return []types.Hash{
		d.ReputationContractHash,
		d.VANFTContractHash,
		d.SimpleVoterContractPackageHash,
		d.VariableRepositoryContractHash,
	}
}

func (d DAOContractsMetadata) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.ReputationContractPackageHash, validation.Required),
		validation.Field(&d.SimpleVoterContractHash, validation.Required),
		validation.Field(&d.VariableRepositoryContractPackageHash, validation.Required),
		validation.Field(&d.VariableRepositoryContractStorageUref, validation.Required),
	)
}
