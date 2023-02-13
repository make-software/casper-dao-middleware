package ces

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"casper-dao-middleware/pkg/casper"
	"casper-dao-middleware/pkg/casper/types"
)

var (
	ErrFailedDeploy                     = errors.New("error: failed deploy, expected successful deploys")
	ErrEventNameNotInSchema             = errors.New("error: event name not found in Schema")
	ErrFailedToParseContractEventSchema = errors.New("error: failed to parse contract event Schema")
	ErrExpectContractStoredValue        = errors.New("error: expect contract stored value")
	ErrExpectCLValueStoredValue         = errors.New("error: expect clValue stored value")
	ErrMissingRequiredNamedKey          = errors.New("error: missing required named key")
	ErrInvalidEventUref                 = errors.New("error: invalid event uref")
	ErrNoEventPrefixInEvent             = errors.New("error: no event_ prefix in event")
	ErrInvalidEventBytes                = errors.New("error: invalid event bytes")
)

const (
	eventSchemaNamedKey = "__events_schema"
	eventNamedKey       = "__events"
	eventPrefix         = "event_"
	dictionaryPrefix    = "dictionary-"
)

type (
	EventParser struct {
		casperClient casper.RPCClient
		// key represent Uref from __events named key
		contractInfos map[string]contractInfo
	}

	contractInfo struct {
		schemas             Schemas
		contractHash        types.Hash
		contractPackageHash types.Hash
		eventsSchemaURef    string
		eventsURef          string
	}
)

func NewParser(casperClient casper.RPCClient, contractHashes []types.Hash) (*EventParser, error) {
	eventParser := EventParser{
		casperClient: casperClient,
	}

	contractInfos, err := eventParser.parseContractsInfos(contractHashes)
	if err != nil {
		return nil, err
	}

	return &EventParser{
		casperClient:  casperClient,
		contractInfos: contractInfos,
	}, nil
}

// ParseExecutionResults accept casper.ExecutionResult analyze its transforms and trying to parse events according to stored contract schema
func (p *EventParser) ParseExecutionResults(executionResult casper.ExecutionResult) ([]ParseResult, error) {
	if executionResult.Success == nil {
		return nil, ErrFailedDeploy
	}

	var results = make([]ParseResult, 0)

	for _, transform := range executionResult.Success.Effect.Transforms {
		if ok := transform.Transform.IsWriteCLValue(); !ok {
			continue
		}

		writeCLValue, err := transform.Transform.ParseAsWriteCLValue()
		if err != nil {
			return nil, err
		}

		if !strings.Contains(transform.Key, dictionaryPrefix) {
			continue
		}

		clValue, reminder, err := types.ParseCLValueFromBytesWithReminder(writeCLValue.Bytes)
		if err != nil {
			continue
		}

		if len(clValue.Bytes) < 4 {
			continue
		}

		eventName, eventBody, err := types.ParseBytesWithReminder(clValue.Bytes[4:])
		if err != nil {
			continue
		}

		if !bytes.HasPrefix(eventName, []byte(eventPrefix)) {
			continue
		}

		urefBytes, _, err := types.ParseBytesWithReminder(reminder)
		if err != nil {
			continue
		}

		parseResult := ParseResult{
			Event: Event{
				Name: strings.TrimPrefix(string(eventName), eventPrefix),
			},
		}

		uref, err := types.NewUrefFromBytes(urefBytes)
		if err != nil {
			parseResult.Error = err
			results = append(results, parseResult)
			continue
		}

		contractSchemas, ok := p.contractInfos[uref.String()]
		if !ok {
			parseResult.Error = ErrInvalidEventUref
			results = append(results, parseResult)
			continue
		}

		eventSchema, ok := contractSchemas.schemas[parseResult.Event.Name]
		if !ok {
			parseResult.Error = ErrEventNameNotInSchema
			results = append(results, parseResult)
			continue
		}

		eventData, err := parseEventDataFromSchemaBytes(eventSchema, eventBody)
		if err != nil {
			parseResult.Error = err
			results = append(results, parseResult)
			continue
		}

		parseResult.Event.ContractHash = contractSchemas.contractHash
		parseResult.Event.ContractPackageHash = contractSchemas.contractPackageHash
		parseResult.Event.Data = eventData
		results = append(results, parseResult)
	}

	return results, nil
}

// FetchContractSchemasBytes accept contract hash to fetch stored contract schema
func (p *EventParser) FetchContractSchemasBytes(contractHash types.Hash) ([]byte, error) {
	stateRootHash, err := p.casperClient.GetLatestStateRootHash()
	if err != nil {
		return nil, err
	}

	contractInfo, err := p.getContractInfo(stateRootHash.StateRootHash, contractHash)
	if err != nil {
		return nil, err
	}

	schemasURefValue, err := p.casperClient.GetStateItem(stateRootHash.StateRootHash, contractInfo.eventsSchemaURef, nil)
	if err != nil {
		return nil, err
	}

	if schemasURefValue.StoredValue.CLValue == nil {
		return nil, ErrExpectCLValueStoredValue
	}

	return schemasURefValue.StoredValue.CLValue.Bytes, nil
}

func (p *EventParser) parseContractsInfos(contractHashes []types.Hash) (map[string]contractInfo, error) {
	stateRootHash, err := p.casperClient.GetLatestStateRootHash()
	if err != nil {
		return nil, err
	}

	contractsSchemas := make(map[string]contractInfo, len(contractHashes))
	for _, hash := range contractHashes {
		contractInfo, err := p.getContractInfo(stateRootHash.StateRootHash, hash)
		if err != nil {
			return nil, err
		}

		schemas, err := p.parseContractEventsSchemas(stateRootHash.StateRootHash, contractInfo.eventsSchemaURef)
		if err != nil {
			return nil, ErrFailedToParseContractEventSchema
		}

		contractInfo.schemas = schemas
		contractsSchemas[contractInfo.eventsURef] = contractInfo
	}

	return contractsSchemas, nil
}

func (p *EventParser) getContractInfo(stateRootHash string, contractHash types.Hash) (contractInfo, error) {
	contractResult, err := p.casperClient.GetStateItem(stateRootHash, fmt.Sprintf("hash-%s", contractHash.ToHex()), nil)
	if err != nil {
		return contractInfo{}, err
	}

	if contractResult.StoredValue.Contract == nil {
		return contractInfo{}, ErrExpectContractStoredValue
	}

	var (
		eventsURef       string
		eventsSchemaURef string
	)

	for _, namedKey := range contractResult.StoredValue.Contract.NamedKeys {
		switch namedKey.Name {
		case eventNamedKey:
			eventsURef = namedKey.Key
		case eventSchemaNamedKey:
			eventsSchemaURef = namedKey.Key
		}

		if eventsURef != "" && eventsSchemaURef != "" {
			break
		}
	}

	if eventsURef == "" || eventsSchemaURef == "" {
		return contractInfo{}, ErrMissingRequiredNamedKey
	}

	return contractInfo{
		contractHash:        contractHash,
		contractPackageHash: contractResult.StoredValue.Contract.ContractPackageHash,
		eventsSchemaURef:    eventsSchemaURef,
		eventsURef:          eventsURef,
	}, nil
}

func (p *EventParser) parseContractEventsSchemas(stateRootHash, eventSchemaUref string) (Schemas, error) {
	schemasURefValue, err := p.casperClient.GetStateItem(stateRootHash, eventSchemaUref, nil)
	if err != nil {
		return nil, err
	}

	if schemasURefValue.StoredValue.CLValue == nil {
		return nil, ErrExpectCLValueStoredValue
	}

	return NewSchemasFromBytes(schemasURefValue.StoredValue.CLValue.Bytes)
}
