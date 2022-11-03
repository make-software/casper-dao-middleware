package dao_event_parser

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	contract_events "casper-dao-middleware/internal//dao_event_parser/events"
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
