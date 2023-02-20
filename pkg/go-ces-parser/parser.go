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
		contractsMetadata map[string]contractMetadata
	}

	contractMetadata struct {
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

	contractsMetadata, err := eventParser.loadContractsMetadata(contractHashes)
	if err != nil {
		return nil, err
	}

	return &EventParser{
		casperClient:      casperClient,
		contractsMetadata: contractsMetadata,
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

		clValue, remainder, err := types.ParseCLValueFromBytesWithRemainder(writeCLValue.Bytes)
		if err != nil {
			continue
		}

		if len(clValue.Bytes) < 4 {
			continue
		}

		eventName, eventBody, err := types.ParseBytesWithRemainder(clValue.Bytes[4:])
		if err != nil {
			continue
		}

		if !bytes.HasPrefix(eventName, []byte(eventPrefix)) {
			continue
		}

		urefBytes, _, err := types.ParseBytesWithRemainder(remainder)
		if err != nil {
			continue
		}

		uref, err := types.NewUrefFromBytes(urefBytes)
		if err != nil {
			continue
		}

		contractMetadata, ok := p.contractsMetadata[uref.String()]
		if !ok {
			continue
		}

		parseResult := ParseResult{
			Event: Event{
				Name: strings.TrimPrefix(string(eventName), eventPrefix),
			},
		}

		eventSchema, ok := contractMetadata.schemas[parseResult.Event.Name]
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

		parseResult.Event.ContractHash = contractMetadata.contractHash
		parseResult.Event.ContractPackageHash = contractMetadata.contractPackageHash
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

	schemasURefValue, err := p.casperClient.GetStateItem(stateRootHash.StateRootHash, fmt.Sprintf("hash-%s", contractHash), []string{eventSchemaNamedKey})
	if err != nil {
		return nil, err
	}

	if schemasURefValue.StoredValue.CLValue == nil {
		return nil, ErrExpectCLValueStoredValue
	}

	return schemasURefValue.StoredValue.CLValue.Bytes, nil
}

func (p *EventParser) loadContractsMetadata(contractHashes []types.Hash) (map[string]contractMetadata, error) {
	stateRootHash, err := p.casperClient.GetLatestStateRootHash()
	if err != nil {
		return nil, err
	}

	contractsSchemas := make(map[string]contractMetadata, len(contractHashes))
	for _, hash := range contractHashes {
		contractMetadata, err := p.loadContractMetadataWithoutEventSchemas(stateRootHash.StateRootHash, hash)
		if err != nil {
			return nil, err
		}

		schemas, err := p.loadContractEventSchemas(stateRootHash.StateRootHash, contractMetadata.eventsSchemaURef)
		if err != nil {
			return nil, ErrFailedToParseContractEventSchema
		}

		contractMetadata.schemas = schemas
		contractsSchemas[contractMetadata.eventsURef] = contractMetadata
	}

	return contractsSchemas, nil
}

func (p *EventParser) loadContractMetadataWithoutEventSchemas(stateRootHash string, contractHash types.Hash) (contractMetadata, error) {
	contractResult, err := p.casperClient.GetStateItem(stateRootHash, fmt.Sprintf("hash-%s", contractHash), nil)
	if err != nil {
		return contractMetadata{}, err
	}

	if contractResult.StoredValue.Contract == nil {
		return contractMetadata{}, ErrExpectContractStoredValue
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
		return contractMetadata{}, ErrMissingRequiredNamedKey
	}

	return contractMetadata{
		contractHash:        contractHash,
		contractPackageHash: contractResult.StoredValue.Contract.ContractPackageHash,
		eventsSchemaURef:    eventsSchemaURef,
		eventsURef:          eventsURef,
	}, nil
}

func (p *EventParser) loadContractEventSchemas(stateRootHash, eventSchemaUref string) (Schemas, error) {
	schemasURefValue, err := p.casperClient.GetStateItem(stateRootHash, eventSchemaUref, nil)
	if err != nil {
		return nil, err
	}

	if schemasURefValue.StoredValue.CLValue == nil {
		return nil, ErrExpectCLValueStoredValue
	}

	return NewSchemasFromBytes(schemasURefValue.StoredValue.CLValue.Bytes)
}
