package casper

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"casper-dao-middleware/pkg/casper/types"

	"github.com/pkg/errors"
)

type (
	EventType int
	EventData []byte
)

type RawEventData struct {
	EventType EventType
	Data      EventData
	EventID   uint64
}

const (
	APIVersionEventType EventType = iota + 1
	BlockAddedEventType
	DeployProcessedEventType
	DeployAcceptedEventType
	DeployExpiredEventType
	EventIDEventType
	FinalitySignatureType
	StepEventType
	FaultEventType
)

var ErrInvalidEventType = errors.New("error on parse target invalid event")

type APIVersionEvent struct {
	APIVersion string `json:"ApiVersion"`
}

// DeployProcessedEvent definition
type (
	DeployProcessed struct {
		DeployHash      types.Hash      `json:"deploy_hash"`
		Account         string          `json:"account"`
		Timestamp       time.Time       `json:"timestamp"`
		TTL             string          `json:"ttl"`
		BlockHash       string          `json:"block_hash"`
		ExecutionResult ExecutionResult `json:"execution_result"`
	}

	ExecutionResult struct {
		Success *Status `json:"Success"`
		Failure *Status `json:"Failure"`
	}
	Status struct {
		Effect       Effect   `json:"effect"`
		Transfers    []string `json:"transfers"`
		Cost         string   `json:"cost"`
		ErrorMessage string   `json:"error_message"`
	}

	TransformValue json.RawMessage

	WriteTransfer struct {
		ID         *uint64     `json:"id"`
		To         *types.Hash `json:"to"`
		DeployHash types.Hash  `json:"deploy_hash"`
		From       types.Hash  `json:"from"`
		Amount     uint64      `json:"amount,string"`
		Source     string      `json:"source"`
		Target     string      `json:"target"`
		Gas        string      `json:"gas"`
	}

	WriteCLValue struct {
		Bytes  string          `json:"bytes"`
		Parsed interface{}     `json:"parsed"`
		CLType json.RawMessage `json:"cl_type"`
	}

	WriteWithdraw struct {
		Amount             uint64          `json:"amount,string"`
		BondingPurse       string          `json:"bonding_purse"`
		EraOfCreation      uint32          `json:"era_of_creation"`
		UnbonderPublicKey  types.PublicKey `json:"unbonder_public_key"`
		ValidatorPublicKey types.PublicKey `json:"validator_public_key"`
	}

	Effect struct {
		Operations []interface{} `json:"operations"`
		Transforms []Transform   `json:"transforms"`
	}
	Transform struct {
		Key       string         `json:"key"`
		Transform TransformValue `json:"transform"`
	}
	DeployProcessedEvent struct {
		DeployProcessed DeployProcessed `json:"DeployProcessed"`
	}
)

// BlockAddedEvent definition
type (
	BlockAdded struct {
		BlockHash string `json:"block_hash"`
		Block     Block  `json:"block"`
	}
	BlockAddedEvent struct {
		BlockAdded BlockAdded `json:"BlockAdded"`
	}
)

func NewRawEventData(eventType EventType, rawEventID []byte, data EventData) (RawEventData, error) {
	eventID, err := strconv.ParseUint(string(rawEventID), 10, 0)
	if err != nil {
		return RawEventData{}, err
	}

	return RawEventData{
		EventType: eventType,
		EventID:   eventID,
		Data:      data,
	}, nil
}

func (d *EventData) ParseAsAPIVersionEvent() (*APIVersionEvent, error) {
	res := APIVersionEvent{}
	if err := json.Unmarshal(*d, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (d *EventData) ParseAsDeployProcessedEvent() (*DeployProcessedEvent, error) {
	res := DeployProcessedEvent{}
	if err := json.Unmarshal(*d, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (d *EventData) ParseAsBlockAddedEvent() (*BlockAddedEvent, error) {
	res := BlockAddedEvent{}
	if err := json.Unmarshal(*d, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (t *TransformValue) IsWriteTransfer() bool {
	return strings.Contains(string(*t), "WriteTransfer")
}

func (t *TransformValue) ParseAsWriteTransfer() (*WriteTransfer, error) {
	type RawWriteTransferTransform struct {
		WriteTransfer struct {
			DeployHash types.Hash `json:"deploy_hash"`
			From       string     `json:"from"`
			To         string     `json:"to"`
			Source     string     `json:"source"`
			Target     string     `json:"target"`
			Amount     uint64     `json:"amount,string"`
			Gas        string     `json:"gas"`
			ID         *uint64    `json:"id"`
		} `json:"WriteTransfer"`
	}

	jsonRes := RawWriteTransferTransform{}
	if err := json.Unmarshal(*t, &jsonRes); err != nil {
		return nil, err
	}

	var toHash *types.Hash
	if hash, err := types.NewHashFromHexString(strings.TrimPrefix(jsonRes.WriteTransfer.To, "account-hash-")); err == nil {
		toHash = &hash
	}

	fromHash, err := types.NewHashFromHexString(strings.TrimPrefix(jsonRes.WriteTransfer.From, "account-hash-"))
	if err != nil {
		return nil, err
	}

	return &WriteTransfer{
		ID:         jsonRes.WriteTransfer.ID,
		To:         toHash,
		From:       fromHash,
		DeployHash: jsonRes.WriteTransfer.DeployHash,
		Amount:     jsonRes.WriteTransfer.Amount,
		Source:     jsonRes.WriteTransfer.Source,
		Target:     jsonRes.WriteTransfer.Target,
		Gas:        jsonRes.WriteTransfer.Gas,
	}, nil
}

func (t *TransformValue) IsWriteContract() bool {
	return bytes.Equal(*t, []byte("\"WriteContract\""))
}

func (t *TransformValue) IsWriteWithdraw() bool {
	return strings.Contains(string(*t), "WriteWithdraw")
}

func (t *TransformValue) IsWriteCLValue() bool {
	return bytes.Contains(*t, []byte("\"WriteCLValue\""))
}

func (t *TransformValue) ParseAsWriteWithdraws() ([]WriteWithdraw, error) {
	type RawWriteWithdrawals struct {
		Withdraws []WriteWithdraw `json:"WriteWithdraw"`
	}

	jsonRes := RawWriteWithdrawals{}
	if err := json.Unmarshal(*t, &jsonRes); err != nil {
		return nil, err
	}

	return jsonRes.Withdraws, nil
}

func (t *TransformValue) ParseAsWriteCLValue() (*WriteCLValue, error) {
	type RawWriteCLValue struct {
		WriteCLValue WriteCLValue `json:"WriteCLValue"`
	}

	jsonRes := RawWriteCLValue{}
	if err := json.Unmarshal(*t, &jsonRes); err != nil {
		return nil, err
	}

	return &jsonRes.WriteCLValue, nil
}

func (t *TransformValue) UnmarshalJSON(data []byte) error {
	valueBuf := make([]byte, len(data))
	copy(valueBuf, data)

	*t = valueBuf
	return nil
}
