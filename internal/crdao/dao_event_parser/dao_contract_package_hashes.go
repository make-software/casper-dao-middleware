package dao_event_parser

import (
	"errors"
	"fmt"
	"strings"

	"casper-dao-middleware/pkg/casper"
	"casper-dao-middleware/pkg/casper/types"

	"github.com/go-ozzo/ozzo-validation"
)

type DAOContractPackageHashes struct {
	ReputationContractPackageHash types.Hash
	VoterContractPackageHash      types.Hash
	VariableRepositoryContract    types.Hash
}

func NewDAOContractPackageHashesFromHashesMap(contractHashes map[string]string, casperClient casper.RPCClient) (DAOContractPackageHashes, error) {
	result := DAOContractPackageHashes{}
	stateRootHash, err := casperClient.GetStateRootHashByHash("")
	if err != nil {
		return DAOContractPackageHashes{}, err
	}

	for contractName, contractHashHex := range contractHashes {
		stateItemRes, err := casperClient.GetStateItem(stateRootHash.StateRootHash, fmt.Sprintf("hash-%s", contractHashHex), []string{})
		if err != nil {
			return DAOContractPackageHashes{}, err
		}

		if stateItemRes.StoredValue.Contract == nil {
			return DAOContractPackageHashes{}, errors.New("expected Contract StoredValue")
		}

		contractPackageHash, err := types.NewHashFromHexString(strings.TrimPrefix(stateItemRes.StoredValue.Contract.ContractPackageHash, "contract-package-wasm"))
		if err != nil {
			return DAOContractPackageHashes{}, err
		}

		switch contractName {
		case "reputation_contract":
			result.ReputationContractPackageHash = contractPackageHash
		case "voter_contract":
			result.VoterContractPackageHash = contractPackageHash
		case "variable_repository_contract":
			result.VariableRepositoryContract = contractPackageHash
		}
	}

	return result, result.Validate()
}

func (d DAOContractPackageHashes) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.ReputationContractPackageHash, validation.Required),
		validation.Field(&d.VoterContractPackageHash, validation.Required),
		validation.Field(&d.VariableRepositoryContract, validation.Required),
	)
}
