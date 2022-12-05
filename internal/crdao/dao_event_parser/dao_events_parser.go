package dao_event_parser

import (
	"errors"
	"fmt"

	"casper-dao-middleware/internal/crdao/dao_event_parser/utils"
	"casper-dao-middleware/pkg/casper"
	"casper-dao-middleware/pkg/casper/types"

	"go.uber.org/zap"
)

var (
	ErrKeyNotInDictionary    = errors.New("error: transform key not in dictionary keys set")
	ErrInvalidDAOEventFormat = errors.New("error: invalid DAO event format")
	ErrNotDAOEventTransform  = errors.New("error: not DAO event transform ")
)

type DAOEvent struct {
	EventName string
	EventBody []byte
}

type DictionaryKeyMetadata struct {
	EventLengthUref string
	EventsUref      string
	EventIndex      uint32
}

type DaoEventParser struct {
	casperClient     casper.RPCClient
	daoDictionarySet map[string]DictionaryKeyMetadata
}

func NewDaoEventParser(casperClient casper.RPCClient, daoContractHashes map[string]string, eventsBuffer uint32) (*DaoEventParser, error) {
	daoEventParser := &DaoEventParser{
		casperClient: casperClient,
	}

	daoDictionarySet, err := daoEventParser.calculateDAOEventsDictionarySet(daoContractHashes, eventsBuffer)
	if err != nil {
		return nil, err
	}

	daoEventParser.daoDictionarySet = daoDictionarySet
	return daoEventParser, nil
}

func (c *DaoEventParser) Parse(event *casper.DeployProcessedEvent) ([]DAOEvent, error) {
	daoEvents := make([]DAOEvent, 0)
	for _, transform := range event.DeployProcessed.ExecutionResult.Success.Effect.Transforms {
		daoEvent, err := c.parseDAOEvent(transform)
		if err != nil {
			continue
		}

		daoEvents = append(daoEvents, daoEvent)

		if err := c.actualizeDAODictionarySet(transform.Key); err != nil {
			return nil, err
		}
	}

	return daoEvents, nil
}

func (c *DaoEventParser) parseDAOEvent(transform casper.Transform) (DAOEvent, error) {
	_, ok := c.daoDictionarySet[transform.Key]
	if !ok {
		zap.S().Debug("transform key is not in dictionary set")
		return DAOEvent{}, ErrKeyNotInDictionary
	}

	if ok := transform.Transform.IsWriteCLValue(); !ok {
		zap.S().Debug("transform is not WriteCLValue")
		return DAOEvent{}, ErrNotDAOEventTransform
	}

	writeCLValue, err := transform.Transform.ParseAsWriteCLValue()
	if err != nil {
		zap.S().Debug("failed to parse transform as WriteCLValue")
		return DAOEvent{}, ErrNotDAOEventTransform
	}

	if string(writeCLValue.CLType) != "\"Any\"" {
		zap.S().Debug("expect CLType as Any")
		return DAOEvent{}, ErrInvalidDAOEventFormat
	}

	clValue, err := utils.ParseDAOCLValueFromBytes(writeCLValue.Bytes)
	if err != nil {
		zap.S().With(zap.Error(err)).Debug("failed to parse CLValue from bytes")
		return DAOEvent{}, err
	}

	// expect DAO events in format of Option<Vec<u8>>
	if clValue.Type.ToString() != "Option(List(U8))" {
		zap.S().Debug("invalid DAO event format, expect Option(List(U8))")
		return DAOEvent{}, ErrInvalidDAOEventFormat
	}

	data := clValue.Bytes[5:]
	eventName, body, err := types.ParseBytesWithReminder(data)
	if err != nil {
		zap.S().With(zap.Error(err)).Debug("failed to parse DAO eventName and body")
		return DAOEvent{}, err
	}
	return DAOEvent{
		EventName: string(eventName),
		EventBody: body,
	}, nil
}

func (c *DaoEventParser) actualizeDAODictionarySet(dictionaryKey string) error {
	dictionaryMeta := c.daoDictionarySet[dictionaryKey]

	stateRootHash, err := c.casperClient.GetStateRootHashByHash("")
	if err != nil {
		return err
	}

	actualEventsLength, err := c.getEventLengthFromUref(stateRootHash.StateRootHash, dictionaryMeta.EventLengthUref)
	if err != nil {
		return err
	}

	// delete processed key from set
	delete(c.daoDictionarySet, dictionaryKey)

	if actualEventsLength != dictionaryMeta.EventIndex {
		return nil
	}

	// means that we need to update daoDictionarySet with new dictionary key of next event
	nextEventIndex := actualEventsLength + 1
	dictionaryItem, err := utils.ToDictionaryKey(dictionaryMeta.EventsUref, nextEventIndex)
	if err != nil {
		return err
	}
	dictionaryMeta.EventIndex = nextEventIndex
	c.daoDictionarySet[dictionaryItem] = dictionaryMeta

	return nil
}

func (c *DaoEventParser) calculateDAOEventsDictionarySet(daoContractHashes map[string]string, eventsBuffer uint32) (map[string]DictionaryKeyMetadata, error) {
	dictionarySet := make(map[string]DictionaryKeyMetadata)

	stateRootHash, err := c.casperClient.GetStateRootHashByHash("")
	if err != nil {
		return nil, err
	}

	for _, hash := range daoContractHashes {
		stateItemRes, err := c.casperClient.GetStateItem(stateRootHash.StateRootHash, fmt.Sprintf("hash-%s", hash), []string{})
		if err != nil {
			return nil, err
		}

		if stateItemRes.StoredValue.Contract == nil {
			return nil, errors.New("expected Contract StoredValue")
		}

		var eventsLengthUref, eventsUref string
		for _, namedKey := range stateItemRes.StoredValue.Contract.NamedKeys {
			if namedKey.Name == "events" {
				eventsUref = namedKey.Key
			}

			if namedKey.Name == "events_length" {
				eventsLengthUref = namedKey.Key
			}
		}
		eventsLenght, err := c.getEventLengthFromUref(stateRootHash.StateRootHash, eventsLengthUref)
		if err != nil {
			return nil, err
		}

		startEventIdx := 1
		if eventsBuffer != 0 && eventsLenght > eventsBuffer {
			startEventIdx = int(eventsLenght - eventsBuffer)
		}

		// iterate over all indexes to calculate all dictionary items
		for index := startEventIdx; index <= int(eventsLenght); index++ {
			dictionaryKey, err := utils.ToDictionaryKey(eventsUref, uint32(index))
			if err != nil {
				return nil, err
			}

			dictionarySet[dictionaryKey] = DictionaryKeyMetadata{
				EventLengthUref: eventsLengthUref,
				EventsUref:      eventsUref,
				EventIndex:      uint32(index),
			}
		}
	}

	return dictionarySet, nil
}

func (c *DaoEventParser) getEventLengthFromUref(stateRootHash string, eventsLengthUref string) (uint32, error) {
	stateItemResult, err := c.casperClient.GetStateItem(stateRootHash, eventsLengthUref, nil)
	if err != nil {
		return 0, err
	}

	if stateItemResult.StoredValue.CLValue == nil {
		return 0, errors.New("expect CLValue as StoredValue")
	}

	parsed, ok := stateItemResult.StoredValue.CLValue.Parsed.(float64)
	if !ok {
		return 0, errors.New("CLValue.Parsed should be float64")
	}

	return uint32(parsed), nil
}
