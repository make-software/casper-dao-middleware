package config

import "casper-dao-middleware/pkg/casper/types"

type DaoContracts struct {
	VariableRepositoryContractHash types.Hash `env:"VARIABLE_REPOSITORY_CONTRACT_HASH,required"`
	ReputationContractHash         types.Hash `env:"REPUTATION_CONTRACT_HASH,required"`
	SimpleVoterContractHash        types.Hash `env:"SIMPLE_VOTER_CONTRACT_HASH,required"`
	RepoVoterContractHash          types.Hash `env:"REPO_VOTER_CONTRACT_HASH,required"`
	ReputationVoterContractHash    types.Hash `env:"REPUTATION_VOTER_CONTRACT_HASH,required"`
	SlashingVoterContractHash      types.Hash `env:"SLASHING_VOTER_CONTRACT_HASH,required"`
	KycVoterContractHash           types.Hash `env:"KYC_VOTER_CONTRACT_HASH,required"`
	VANFTContractHash              types.Hash `env:"VA_NFT_CONTRACT_HASH,required"`
	KycNFTContractHash             types.Hash `env:"KYC_NFT_CONTRACT_HASH,required"`
	OnboardingRequestContractHash  types.Hash `env:"ONBOARDING_REQUEST_CONTRACT_HASH,required"`
	AdminContractHash              types.Hash `env:"ADMIN_CONTRACT_HASH,required"`
}

func (d DaoContracts) ToMap() map[string]types.Hash {
	return map[string]types.Hash{
		"reputation_contract":          d.ReputationContractHash,
		"simple_voter_contract":        d.SimpleVoterContractHash,
		"repo_voter_contract":          d.RepoVoterContractHash,
		"kyc_voter_contract":           d.KycVoterContractHash,
		"reputation_voter_contract":    d.ReputationVoterContractHash,
		"slashing_voter_contract":      d.SlashingVoterContractHash,
		"va_nft_contract":              d.VANFTContractHash,
		"kyc_nft_contract":             d.KycNFTContractHash,
		"variable_repository_contract": d.VariableRepositoryContractHash,
		"onboarding_request_contract":  d.OnboardingRequestContractHash,
		"admin":                        d.AdminContractHash,
	}
}
