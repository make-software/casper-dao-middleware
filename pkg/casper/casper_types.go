package casper

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"casper-dao-middleware/pkg/casper/types"

	"github.com/pkg/errors"
)

type RPCRequest struct {
	Version string      `json:"jsonrpc"`
	ID      string      `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

type RPCResponse struct {
	Version string          `json:"jsonrpc"`
	ID      string          `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   *RPCError       `json:"error,omitempty"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type GetBlockResult struct {
	Version string `json:"version"`
	Block   Block  `json:"block"`
}

type Block struct {
	Hash   types.Hash  `json:"hash"`
	Header BlockHeader `json:"header"`
	Body   BlockBody   `json:"body"`
	Proofs []Proof     `json:"proofs"`
}

type EraEnd struct {
	EraReport               EraReport                `json:"era_report"`
	NextEraValidatorWeights []NextEraValidatorWeight `json:"next_era_validator_weights"`
}

type EraReport struct {
	Equivocators       []interface{} `json:"equivocators"`
	Rewards            []Reward      `json:"rewards"`
	InactiveValidators []string      `json:"inactive_validators"`
}

type Reward struct {
	Validator types.PublicKey `json:"validator"`
	Amount    int64           `json:"amount"`
}

type NextEraValidatorWeight struct {
	Validator types.PublicKey `json:"validator"`
	Weight    uint64          `json:"weight,string"`
}

type BlockHeader struct {
	EraEnd          *EraEnd    `json:"era_end"` // @todo
	ParentHash      types.Hash `json:"parent_hash"`
	StateRootHash   types.Hash `json:"state_root_hash"`
	BodyHash        types.Hash `json:"body_hash"`
	AccumulatedSeed string     `json:"accumulated_seed"`
	ProtocolVersion string     `json:"protocol_version"`
	Height          uint64     `json:"height"`
	EraID           uint32     `json:"era_id"`
	RandomBit       bool       `json:"random_bit"`
	Timestamp       time.Time  `json:"timestamp"`
}

type BlockBody struct {
	Proposer       types.PublicKey `json:"proposer"`
	DeployHashes   []string        `json:"deploy_hashes"`
	TransferHashes []string        `json:"transfer_hashes"`
}

type Proof struct {
	PublicKey string `json:"public_key"`
	Signature string `json:"signature"`
}

type StateGetAuctionInfoResult struct {
	Version      string       `json:"api_version"`
	AuctionState AuctionState `json:"auction_state"`
}

type AuctionState struct {
	StateRootHash string          `json:"state_root_hash"`
	BlockHeight   uint64          `json:"block_height"`
	EraValidators []EraValidators `json:"era_validators"`
	Bids          []ValidatorBid  `json:"bids"`
}

type EraValidators struct {
	EraID            uint32            `json:"era_id"`
	ValidatorWeights []ValidatorWeight `json:"validator_weights"`
}

type ValidatorBid struct {
	PublicKey types.PublicKey `json:"public_key"`
	Bid       Bid             `json:"bid"`
}

type Bid struct {
	BondingPurse   string      `json:"bonding_purse"`
	StakedAmount   uint64      `json:"staked_amount,string"`
	DelegationRate float32     `json:"delegation_rate"`
	Inactive       bool        `json:"inactive"`
	Delegators     []Delegator `json:"delegators"`
}

type Delegator struct {
	PublicKey    types.PublicKey `json:"public_key"`
	StakedAmount uint64          `json:"staked_amount,string"`
	BondingPurse string          `json:"bonding_purse"`
	Delegatee    types.PublicKey `json:"delegatee"`
}

type GetBlockTransfersResult struct {
	Version   string     `json:"api_version"`
	BlockHash string     `json:"block_hash"`
	Transfers []Transfer `json:"transfers"`
}

type Transfer struct {
	ID         int64  `json:"id,omitempty"`
	DeployHash string `json:"deploy_hash"`
	From       string `json:"from"`
	To         string `json:"to"`
	Source     string `json:"source"`
	Target     string `json:"target"`
	Amount     string `json:"amount"`
	Gas        string `json:"gas"`
}

type GetDeployResult struct {
	Version          string                  `json:"api_version"`
	Deploy           Deploy                  `json:"deploy"`
	ExecutionResults []DeployExecutionResult `json:"execution_results"`
}

type Deploy struct {
	Hash      types.Hash    `json:"hash"`
	Header    DeployHeader  `json:"header"`
	Payment   DeployPayment `json:"payment"`
	Session   DeploySession `json:"session"`
	Approvals []Approval    `json:"approvals"`
}

type DeployHeader struct {
	Account      types.PublicKey `json:"account"`
	Timestamp    time.Time       `json:"timestamp"`
	TTL          string          `json:"ttl"`
	GasPrice     int             `json:"gas_price"`
	BodyHash     string          `json:"body_hash"`
	Dependencies []string        `json:"dependencies"`
	ChainName    string          `json:"chain_name"`
}

type DeployPayment struct {
	ModuleBytes *DeployModuleBytes `json:"ModuleBytes"`
}

type DeploySession struct {
	ModuleBytes                   *DeployModuleBytes                   `json:"ModuleBytes"`
	StoredContractByHash          *DeployStoredContractByHash          `json:"StoredContractByHash"`
	StoredContractByName          *DeployStoredContractByName          `json:"StoredContractByName"`
	StoredVersionedContractByHash *DeployStoredVersionedContractByHash `json:"StoredVersionedContractByHash"`
	StoredVersionedContractByName *DeployStoredVersionedContractByName `json:"StoredVersionedContractByName"`
	Transfer                      *DeployTransfer                      `json:"Transfer"`
}

type DeployModuleBytes struct {
	Args        Args   `json:"args"`
	ModuleBytes string `json:"module_bytes"`
}

type DeployStoredContractByHash struct {
	Args       Args       `json:"args"`
	Hash       types.Hash `json:"hash"`
	EntryPoint string     `json:"entry_point"`
}

type DeployStoredContractByName struct {
	Args       Args   `json:"args"`
	Name       string `json:"name"`
	EntryPoint string `json:"entry_point"`
}

type DeployStoredVersionedContractByHash struct {
	Args       Args       `json:"args"`
	Hash       types.Hash `json:"hash"`
	Version    uint16     `json:"version"`
	EntryPoint string     `json:"entry_point"`
}

type DeployStoredVersionedContractByName struct {
	Args       Args   `json:"args"`
	Name       string `json:"name"`
	Version    uint16 `json:"version"`
	EntryPoint string `json:"entry_point"`
}

type DeployTransfer struct {
	Args Args `json:"args"`
}

type Args map[string]CLValue

func (a *Args) UnmarshalJSON(data []byte) error {
	var value [][2]interface{}
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	parsedArgs := make(map[string]CLValue, len(value))

	for _, arg := range value {
		name, ok := arg[0].(string)
		if !ok {
			return errors.New("failed to parse arg name")
		}
		values, ok := arg[1].(map[string]interface{})
		if !ok {
			return errors.New("failed to parse arg value")
		}

		clType, ok := values["cl_type"]
		if !ok {
			return errors.New("failed to get cl_type from arg")
		}

		parsed, ok := values["parsed"]
		if !ok {
			return errors.New("failed to get parsed from arg")
		}

		plainBytes, ok := values["bytes"]
		if !ok {
			return errors.New("failed to get bytes from arg")
		}
		bytes, ok := plainBytes.(string)
		if !ok {
			return errors.New("failed to assert bytes from arg")
		}

		clValue := CLValue{
			CLType: clType,
			Parsed: parsed,
			Bytes:  []byte(bytes),
		}

		parsedArgs[name] = clValue
	}

	*a = parsedArgs

	return nil
}

type PutDeployRes struct {
	Hash string `json:"deploy_hash"`
}

type Approval struct {
	Signer    string `json:"signer"`
	Signature string `json:"signature"`
}

type DeployExecutionResult struct {
	BlockHash types.Hash      `json:"block_hash"`
	Result    ExecutionResult `json:"result"`
}

type StateGetItemResult struct {
	StoredValue StoredValue `json:"stored_value"`
}

type StoredValue struct {
	CLValue         *CLValue         `json:"CLValue,omitempty"`
	Account         *Account         `json:"Account,omitempty"`
	Contract        *Contract        `json:"Contract,omitempty"`
	ContractWASM    *string          `json:"ContractWASM,omitempty"`
	ContractPackage *ContractPackage `json:"ContractPackage,omitempty"`
	Transfer        *Transfer        `json:"Transfer,omitempty"`
	DeployInfo      *DeployInfo      `json:"DeployInfo,omitempty"`
}

type CLValue struct {
	CLType interface{} `json:"cl_type"`
	Parsed interface{} `json:"parsed"`
	Bytes  []byte      `json:"bytes"`
}

type Account struct {
	AccountHash      string           `json:"account_hash"`
	NamedKeys        []NamedKey       `json:"named_keys"`
	MainPurse        string           `json:"main_purse"`
	AssociatedKeys   []AssociatedKey  `json:"associated_keys"`
	ActionThresholds ActionThresholds `json:"action_thresholds"`
}

type NamedKey struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

type AssociatedKey struct {
	AccountHash string `json:"account_hash"`
	Weight      uint64 `json:"weight"`
}

type ActionThresholds struct {
	Deployment    uint64 `json:"deployment"`
	KeyManagement uint64 `json:"key_management"`
}

type Contract struct {
	ContractPackageHash types.Hash   `json:"contract_package_hash"`
	ContractWasmHash    string       `json:"contract_wasm_hash"`
	ProtocolVersion     string       `json:"protocol_version"`
	NamedKeys           []NamedKey   `json:"named_keys"`
	Entrypoints         []Entrypoint `json:"entry_points"`
}

type Entrypoint struct {
	Name string `json:"name"`
}

type Version struct {
	ContractHash    string `json:"contract_hash"`
	ContractVersion uint16 `json:"contract_version"`
}

type DisabledVersion struct {
	ContractVersion uint16 `json:"contract_version"`
}

type ContractPackage struct {
	Versions         []Version         `json:"versions"`
	DisabledVersions []DisabledVersion `json:"disabled_versions"`
}

type DeployInfo struct {
	DeployHash string   `json:"deploy_hash"`
	Transfers  []string `json:"transfers"`
	From       string   `json:"from"`
	Source     string   `json:"source"`
	Gas        string   `json:"gas"`
}

type StateGetBalanceResult struct {
	Version      string `json:"api_version"`
	BalanceValue string `json:"balance_value"`
}

type ValidatorWeight struct {
	PublicKey string `json:"public_key"`
	Weight    uint64 `json:"weight,string"`
}

type StatusResult struct {
	Peers                 []Peer `json:"peers"`
	LastAddedBlock        Block  `json:"last_added_block_info"`
	StartingStateRootHash string `json:"starting_state_root_hash"`
	ChainSpecName         string `json:"chainspec_name"`
	BuildVersion          string `json:"build_version"`
	OurPublicSigningKey   string `json:"our_public_signing_key"`
	Uptime                string `json:"uptime"`
}

type Peer struct {
	NodeID  string `json:"node_id"`
	Address string `json:"address"`
}

type PeerResult struct {
	Peers []Peer `json:"peers"`
}

type GetStateRootHashResult struct {
	Version       string `json:"api_version"`
	StateRootHash string `json:"state_root_hash"`
}

type GetEraInfoBySwitchBlockResult struct {
	Version    string     `json:"api_version"`
	EraSummary EraSummary `json:"era_summary"`
}

type EraSummary struct {
	BlockHash     types.Hash            `json:"block_hash"`
	EraID         uint32                `json:"era_id"`
	StoredValue   EraSummaryStoredValue `json:"stored_value"`
	StateRootHash types.Hash            `json:"state_root_hash"`
	MerkleProof   string                `json:"merkle_proof"`
}

type EraSummaryStoredValue struct {
	EraInfo EraInfo `json:"EraInfo"`
}

type EraInfo struct {
	SeigniorageAllocations []SeigniorageAllocation `json:"seigniorage_allocations"`
}

type SeigniorageAllocation struct {
	Validator *ValidatorAllocation `json:"Validator"`
	Delegator *DelegatorAllocation `json:"Delegator"`
}

type ValidatorAllocation struct {
	ValidatorPublicKey types.PublicKey `json:"validator_public_key"`
	Amount             uint64          `json:"amount,string"`
}

type DelegatorAllocation struct {
	DelegatorPublicKey types.PublicKey `json:"delegator_public_key"`
	ValidatorPublicKey types.PublicKey `json:"validator_public_key"`
	Amount             uint64          `json:"amount,string"`
}

func (p *Contract) UnmarshalJSON(data []byte) error {
	var rawContract = struct {
		ContractPackageHash string       `json:"contract_package_hash"`
		ContractWasmHash    string       `json:"contract_wasm_hash"`
		ProtocolVersion     string       `json:"protocol_version"`
		NamedKeys           []NamedKey   `json:"named_keys"`
		Entrypoints         []Entrypoint `json:"entry_points"`
	}{}
	if err := json.Unmarshal(data, &rawContract); err != nil {
		return err
	}

	contractPackageHash, err := types.NewHashFromHexStringWithPrefix(rawContract.ContractPackageHash, "contract-package-wasm")
	if err != nil {
		return err
	}

	p.ContractPackageHash = contractPackageHash
	p.ContractWasmHash = rawContract.ContractWasmHash
	p.ProtocolVersion = rawContract.ProtocolVersion
	p.NamedKeys = rawContract.NamedKeys
	p.Entrypoints = rawContract.Entrypoints
	return nil
}

func (p Contract) MarshalJSON() ([]byte, error) {
	var resp = struct {
		ContractPackageHash string       `json:"contract_package_hash"`
		ContractWasmHash    string       `json:"contract_wasm_hash"`
		ProtocolVersion     string       `json:"protocol_version"`
		NamedKeys           []NamedKey   `json:"named_keys"`
		Entrypoints         []Entrypoint `json:"entry_points"`
	}{
		ContractPackageHash: fmt.Sprintf("contract-package-wasm%s", p.ContractPackageHash.ToHex()),
		ContractWasmHash:    p.ContractWasmHash,
		ProtocolVersion:     p.ProtocolVersion,
		NamedKeys:           p.NamedKeys,
		Entrypoints:         p.Entrypoints,
	}

	return json.Marshal(resp)
}

func (p *CLValue) UnmarshalJSON(data []byte) error {
	var hexValue = struct {
		CLType interface{} `json:"cl_type"`
		Parsed interface{} `json:"parsed"`
		Bytes  string      `json:"bytes"`
	}{}
	if err := json.Unmarshal(data, &hexValue); err != nil {
		return err
	}

	rawBytes, err := hex.DecodeString(hexValue.Bytes)
	if err != nil {
		return err
	}

	p.CLType = hexValue.CLType
	p.Parsed = hexValue.Parsed
	p.Bytes = rawBytes
	return nil
}

func (h CLValue) MarshalJSON() ([]byte, error) {
	var resp = struct {
		CLType interface{} `json:"cl_type"`
		Parsed interface{} `json:"parsed"`
		Bytes  string      `json:"bytes"`
	}{
		CLType: h.CLType,
		Parsed: h.Parsed,
		Bytes:  hex.EncodeToString(h.Bytes),
	}

	return json.Marshal(resp)
}
