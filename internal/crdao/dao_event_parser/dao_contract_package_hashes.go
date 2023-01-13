package dao_event_parser

import (
	"errors"
	"fmt"
	"strings"

	"casper-dao-middleware/pkg/casper"
	"casper-dao-middleware/pkg/casper/types"

	"github.com/go-ozzo/ozzo-validation"
)

const variableRepositoryContractStorageUrefName = "storage_repository_contract"

type DAOContractsMetadata struct {
	ReputationContractPackageHash         types.Hash
	VoterContractPackageHash              types.Hash
	VariableRepositoryContractPackageHash types.Hash
	VariableRepositoryContractStorageUref string
}

func NewDAOContractsMetadataFromHashesMap(contractHashes map[string]string, casperClient casper.RPCClient) (DAOContractsMetadata, error) {
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

		contractPackageHash, err := types.NewHashFromHexString(strings.TrimPrefix(stateItemRes.StoredValue.Contract.ContractPackageHash, "contract-package-wasm"))
		if err != nil {
			return DAOContractsMetadata{}, err
		}

		switch contractName {
		case "reputation_contract":
			result.ReputationContractPackageHash = contractPackageHash
		case "voter_contract":
			result.VoterContractPackageHash = contractPackageHash
		case "variable_repository_contract":
			result.VariableRepositoryContractPackageHash = contractPackageHash
			for _, namedKey := range stateItemRes.StoredValue.Contract.NamedKeys {
				if namedKey.Name == variableRepositoryContractStorageUrefName {
					result.VariableRepositoryContractStorageUref = namedKey.Key
					break
				}
			}
		}

	}

	return result, result.Validate()
}

func (d DAOContractsMetadata) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.ReputationContractPackageHash, validation.Required),
		validation.Field(&d.VoterContractPackageHash, validation.Required),
		validation.Field(&d.VariableRepositoryContractPackageHash, validation.Required),
		validation.Field(&d.VariableRepositoryContractStorageUref, validation.Required),
	)
}
