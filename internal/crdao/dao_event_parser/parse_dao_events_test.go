package dao_event_parser

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	contract_events "casper-dao-middleware/internal/crdao/dao_event_parser/events"
	"casper-dao-middleware/internal/crdao/dao_event_parser/types"
	"casper-dao-middleware/pkg/casper"

	"github.com/stretchr/testify/assert"
)

func TestParsVotingCreatedDaoEvents(t *testing.T) {
	votingCreatedTransform := `{
		"key": "some-key",
		"transform": {
			"WriteCLValue": {
				"bytes": "50000000014b0000000d000000566f74696e6743726561746564003b4ffcfb21411ced5fc1560c3f6ffed86f4885e5ea05cde49d90962a48a14d9501010101000000ccbf190000000000005c26050000000001640d0e0320000000c7883cef3fa34729510ff09fb822bc39864b21b3a97f3d9815e7ed3e4341ccd94000000038633033396666376361613137636365626663616463343462643966636536613462363639396334643033646532653333343961613164633131313933636437",
				"parsed": "null",
				"cl_type": "Any"
			}
		}
	}`

	var transform casper.Transform
	err := json.Unmarshal([]byte(votingCreatedTransform), &transform)
	assert.NoError(t, err)

	daoEventParser := DaoEventParser{
		daoDictionarySet: map[string]DictionaryKeyMetadata{
			"some-key": {},
		},
	}

	daoEvent, err := daoEventParser.parseDAOEvent(transform)
	assert.NoError(t, err)

	assert.Equal(t, "VotingCreated", daoEvent.EventName)
	assert.NotEmpty(t, daoEvent.EventBody)

	votingCreatedEvent, err := contract_events.ParseVotingCreatedEvent(daoEvent.EventBody)
	if err != nil {
		panic(err)
	}

	assert.NotEmpty(t, votingCreatedEvent.Creator.AccountHash)
	assert.NotEmpty(t, votingCreatedEvent.VotingID)
	assert.NotEmpty(t, votingCreatedEvent.InformalVotingID)
	assert.NotEmpty(t, votingCreatedEvent.ConfigInformalVotingTime)
	assert.NotEmpty(t, votingCreatedEvent.ConfigFormalVotingTime)
}

func TestParseValueUpdatedDaoEvents(t *testing.T) {
	//let event = events::ValueUpdated{
	//	key: "hello".to_string(),
	//	Value: Bytes::from(U256::from(43210301).to_bytes().unwrap()) ,
	//	activation_time: None
	//};
	hexStr := `0c00000056616c7565557064617465640500000068656c6c6f05000000043d56930200`

	res, _ := hex.DecodeString(hexStr)
	valueUpdatedEvent, err := contract_events.ParseValueUpdatedEvent(res[16:])
	assert.NoError(t, err)

	assert.NotEmpty(t, valueUpdatedEvent.Value)
	assert.Equal(t, valueUpdatedEvent.Key, "hello")
	assert.Equal(t, (**valueUpdatedEvent.Value.UValue).Int64(), int64(43210301))
	assert.Nil(t, valueUpdatedEvent.ActivationTime)

	// with u64 encoding ===========================================================

	//let event = events::ValueUpdated{
	//	key: "key".to_string(),
	//	Value: Bytes::from(172800000u64.to_bytes().unwrap()) ,
	//	activation_time: None
	//};
	hexStr = `0c00000056616c756555706461746564030000006b65790800000000b84c0a0000000000`

	res, _ = hex.DecodeString(hexStr)
	valueUpdatedEvent, err = contract_events.ParseValueUpdatedEvent(res[16:])
	assert.NoError(t, err)

	assert.NotEmpty(t, valueUpdatedEvent.Value)
	assert.Equal(t, valueUpdatedEvent.Key, "key")
	assert.Equal(t, valueUpdatedEvent.Value.String(), "172800000")
	assert.Nil(t, valueUpdatedEvent.ActivationTime)

	// with activation_time encoding ===========================================================

	//let event = events::ValueUpdated{
	//	key: "key".to_string(),
	//	Value: Bytes::from(172800000u64.to_bytes().unwrap()) ,
	//	activation_time: Some(123522124),
	//};
	hexStr = `0c00000056616c756555706461746564030000006b65790800000000b84c0a00000000014ccc5c0700000000`

	res, _ = hex.DecodeString(hexStr)
	valueUpdatedEvent, err = contract_events.ParseValueUpdatedEvent(res[16:])
	assert.NoError(t, err)

	assert.NotEmpty(t, valueUpdatedEvent.Value)

	assert.Equal(t, valueUpdatedEvent.Key, "key")
	assert.Equal(t, valueUpdatedEvent.Value.String(), "172800000")
	assert.Equal(t, *valueUpdatedEvent.ActivationTime, uint64(123522124))
}

func TestParseRecord(t *testing.T) {
	//let rec: Record = (Bytes::from(U256::from(172234435535355 as u64).to_bytes().unwrap()), None);
	//let rec_bytes = rec.to_bytes().unwrap();
	//println!("{:?}", hex::encode(rec_bytes));

	hexStr := `0700000006fb215974a59c00`

	res, _ := hex.DecodeString(hexStr)
	record, err := types.NewRecordFromBytes(res)
	assert.NoError(t, err)

	assert.NotEmpty(t, record.Value.UValue)
	assert.Nil(t, record.FutureValue)
	assert.Equal(t, record.Value.String(), "172234435535355")

	//let rec: Record = (Bytes::from(U512::from(978905355u64).to_bytes().unwrap()), Some((Bytes::from(U512::from(1232523u64).to_bytes().unwrap()), 12002)));
	//let rec_bytes = rec.to_bytes().unwrap();
	//println!("{:?}", hex::encode(rec_bytes));

	hexStr = `05000000040be9583a0104000000038bce12e22e000000000000`

	res, _ = hex.DecodeString(hexStr)
	record, err = types.NewRecordFromBytes(res)
	assert.NoError(t, err)

	assert.NotEmpty(t, record.Value.UValue)
	assert.NotNil(t, record.FutureValue)
	assert.Equal(t, record.Value.String(), "978905355")
	assert.Equal(t, record.FutureValue.Value.String(), "1232523")
	assert.Equal(t, record.FutureValue.ActivationTime, uint64(12002))

	// ============== u64 bytes ==================================================

	//let rec: Record = (Bytes::from(172234435535355u64.to_bytes().unwrap()), None);
	//let rec_bytes = rec.to_bytes().unwrap();
	//println!("{:?}", hex::encode(rec_bytes));

	hexStr = `08000000fb215974a59c000000`

	res, _ = hex.DecodeString(hexStr)
	record, err = types.NewRecordFromBytes(res)
	assert.NoError(t, err)

	assert.NotEmpty(t, record.Value.U64Value)
	assert.Nil(t, record.FutureValue)
	assert.Equal(t, record.Value.String(), "172234435535355")

	// ============== u64 bytes with FutureValue ==================================================

	//let rec: Record = (Bytes::from(172234435535355u64.to_bytes().unwrap()), Some((Bytes::from(17223u64.to_bytes().unwrap()), 19000)));
	//let rec_bytes = rec.to_bytes().unwrap();
	//println!("{:?}", hex::encode(rec_bytes));

	hexStr = `08000000fb215974a59c000001080000004743000000000000384a000000000000`

	res, _ = hex.DecodeString(hexStr)
	record, err = types.NewRecordFromBytes(res)
	assert.NoError(t, err)

	assert.NotEmpty(t, record.Value.U64Value)
	assert.NotNil(t, record.FutureValue)
	assert.Equal(t, record.Value.String(), "172234435535355")
	assert.Equal(t, record.FutureValue.Value.String(), "17223")
	assert.Equal(t, record.FutureValue.ActivationTime, uint64(19000))

	// ============== boolean bytes ==================================================

	//let rec: Record = (Bytes::from(true.to_bytes().unwrap()), None);
	//let rec_bytes = rec.to_bytes().unwrap();
	//println!("{:?}", hex::encode(rec_bytes));

	hexStr = `010000000100`

	res, _ = hex.DecodeString(hexStr)
	record, err = types.NewRecordFromBytes(res)
	assert.NoError(t, err)

	assert.NotEmpty(t, record.Value.BoolValue)
	assert.Nil(t, record.FutureValue)
	assert.Equal(t, record.Value.String(), "true")
}

func TestParseBallotCastDaoEvents(t *testing.T) {
	ballotCastTransform := `	{
		"key": "some-key",
		"transform": {
		"WriteCLValue": {
			"bytes": "3d00000001380000000a00000042616c6c6f7443617374003b4ffcfb21411ced5fc1560c3f6ffed86f4885e5ea05cde49d90962a48a14d9501010200000002e8030d0e0320000000c7883cef3fa34729510ff09fb822bc39864b21b3a97f3d9815e7ed3e4341ccd94000000032366130386534643063353139306630313837316530353639623632393062383637363030383564393966313765623465376536623538666562386436323439",
				"parsed": "null",
				"cl_type": "Any"
		}
	}
	}`

	var transform casper.Transform
	err := json.Unmarshal([]byte(ballotCastTransform), &transform)
	assert.NoError(t, err)

	daoEventParser := DaoEventParser{
		daoDictionarySet: map[string]DictionaryKeyMetadata{
			"some-key": {},
		},
	}

	daoEvent, err := daoEventParser.parseDAOEvent(transform)
	assert.NoError(t, err)

	assert.Equal(t, "BallotCast", daoEvent.EventName)
	assert.NotEmpty(t, daoEvent.EventBody)

	ballotCastEvent, err := contract_events.ParseBallotCastEvent(daoEvent.EventBody)
	if err != nil {
		panic(err)
	}

	assert.NotEmpty(t, ballotCastEvent.Address.AccountHash)
	assert.NotEmpty(t, ballotCastEvent.VotingID)
}

func TestParseMintDaoEvents(t *testing.T) {
	mintTransform := `	{
		"key": "some-key",
		"transform": {
		"WriteCLValue": {
			"bytes": "33000000012e000000040000004d696e74003b4ffcfb21411ced5fc1560c3f6ffed86f4885e5ea05cde49d90962a48a14d950400ca9a3b0d0e032000000006624e86505574d6eeafd021e55f7bb1d2a94a7e95aa7f05cf0df8f02e7101f94000000038633033396666376361613137636365626663616463343462643966636536613462363639396334643033646532653333343961613164633131313933636437",
				"parsed": "null",
				"cl_type": "Any"
		}
	}
	}`

	var transform casper.Transform
	err := json.Unmarshal([]byte(mintTransform), &transform)
	assert.NoError(t, err)

	daoEventParser := DaoEventParser{
		daoDictionarySet: map[string]DictionaryKeyMetadata{
			"some-key": {},
		},
	}

	daoEvent, err := daoEventParser.parseDAOEvent(transform)
	assert.NoError(t, err)

	assert.Equal(t, "Mint", daoEvent.EventName)
	assert.NotEmpty(t, daoEvent.EventBody)

	mintEvent, err := contract_events.ParseMintEvent(daoEvent.EventBody)
	if err != nil {
		panic(err)
	}

	assert.NotEmpty(t, mintEvent.Address.AccountHash)
	assert.NotEmpty(t, mintEvent.Amount)
}

func TestParseVotingEndedDaoEvents(t *testing.T) {
	// rust code event

	//let addr = AccountHash::from_formatted_str("account-hash-24b6d5aabb8f0ac17d272763a405e9ceca9166b75b745cf200695e172857c2dd").unwrap();
	//let addr = Address::from(addr);
	//
	//let addr2 = AccountHash::from_formatted_str("account-hash-f1075fce3b8cd4eab748b8705ca02444a5e35c0248662649013d8a5cb2b1a87c").unwrap();
	//let addr2 = Address::from(addr2);
	//
	//let mut transfers = BTreeMap::new();
	//transfers.insert(addr, U256::from(100));
	//// transfers.insert(addr2, U256::from(200));
	//
	//let transfersBytes = transfers.to_bytes().unwrap();
	//println!("{:?}", transfersBytes);
	//
	//let mut mints = BTreeMap::new();
	//mints.insert(addr, U256::from(300));
	//
	//let event = VotingEnded {
	//	voting_id: U256::from(1),
	//	informal_voting_id: U256::from(1),
	//	formal_voting_id: Some(U256::from(213)),
	//	result: String::from("passed"),
	//	votes_count: U256::from(2),
	//	stake_in_favor: U256::from(2),
	//	stake_against: U256::from(56),
	//	transfers,
	//	burns: BTreeMap::new(),
	//	mints,
	//};

	votingEndedRawBytes := `010101010101d506000000706173736564010201020138010000000024b6d5aabb8f0ac17d272763a405e9ceca9166b75b745cf200695e172857c2dd016400000000010000000024b6d5aabb8f0ac17d272763a405e9ceca9166b75b745cf200695e172857c2dd022c01`

	eventRawBytes, err := hex.DecodeString(votingEndedRawBytes)
	assert.NoError(t, err)

	votingEnded, err := contract_events.ParseVotingEndedEvent(eventRawBytes)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, (*votingEnded.VotingID).String(), "1")
	assert.Equal(t, (*votingEnded.InformalVotingID).String(), "1")
	assert.NotNil(t, votingEnded.FormalVotingID)
	assert.Equal(t, votingEnded.Result, "passed")
	assert.Equal(t, (*votingEnded.VotesCount).String(), "2")
	assert.Equal(t, (*votingEnded.StakeInFavour).String(), "2")
	assert.Equal(t, (*votingEnded.StakeAgainst).String(), "56")

	assert.True(t, len(votingEnded.Transfers) == 1)
	assert.Nil(t, votingEnded.Burns)
	assert.True(t, len(votingEnded.Mints) == 1)
}
