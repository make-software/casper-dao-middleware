package config

import "github.com/make-software/casper-go-sdk/casper"

type DaoContracts struct {
	VariableRepositoryContractHash casper.Hash `env:"VARIABLE_REPOSITORY_CONTRACT_HASH,required"`
	ReputationContractHash         casper.Hash `env:"REPUTATION_CONTRACT_HASH,required"`
	SimpleVoterContractHash        casper.Hash `env:"SIMPLE_VOTER_CONTRACT_HASH,required"`
	RepoVoterContractHash          casper.Hash `env:"REPO_VOTER_CONTRACT_HASH,required"`
	ReputationVoterContractHash    casper.Hash `env:"REPUTATION_VOTER_CONTRACT_HASH,required"`
	SlashingVoterContractHash      casper.Hash `env:"SLASHING_VOTER_CONTRACT_HASH,required"`
	KycVoterContractHash           casper.Hash `env:"KYC_VOTER_CONTRACT_HASH,required"`
	VANFTContractHash              casper.Hash `env:"VA_NFT_CONTRACT_HASH,required"`
	KycNFTContractHash             casper.Hash `env:"KYC_NFT_CONTRACT_HASH,required"`
	OnboardingRequestContractHash  casper.Hash `env:"ONBOARDING_REQUEST_CONTRACT_HASH,required"`
	AdminContractHash              casper.Hash `env:"ADMIN_CONTRACT_HASH,required"`
	BidEscrowContractHash          casper.Hash `env:"BID_ESCROW_CONTRACT_HASH,required"`
}

func (d DaoContracts) ToMap() map[string]casper.Hash {
	return map[string]casper.Hash{
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
		"bid_escrow":                   d.BidEscrowContractHash,
	}
}
