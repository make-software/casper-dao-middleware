package casper

type NodeStatus struct {
	APIVersion            string                       `json:"api_version"`
	ChainSpecName         string                       `json:"chainspec_name"`
	StartingStateRootHash string                       `json:"starting_state_root_hash"`
	Peers                 []NodeStatusPeer             `json:"peers"`
	LastAddedBlockInfo    NodeStatusLastAddedBlockInfo `json:"last_added_block_info"`
	OutPublicSigningKey   string                       `json:"our_public_signing_key"`
	RoundLength           string                       `json:"round_length"`
	BuildVersion          string                       `json:"build_version"`
	NextUpgrade           NodeStatusNextUpgrade        `json:"next_upgrade"`
}

type NodeStatusPeer struct {
	NodeID  string `json:"node_id"`
	Address string `json:"address"`
}

type NodeStatusLastAddedBlockInfo struct {
	Hash          string `json:"hash"`
	Timestamp     string `json:"timestamp"`
	EraID         uint64 `json:"era_id"`
	Height        uint64 `json:"height"`
	StateRootHash string `json:"state_root_hash"`
	Creator       string `json:"creator"`
}

type NodeStatusNextUpgrade struct {
	ActivationPoint uint64 `json:"activation_point"`
	ProtocolVersion string `json:"protocol_version"`
}
