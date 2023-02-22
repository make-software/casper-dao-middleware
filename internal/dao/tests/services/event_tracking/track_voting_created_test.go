//go:build integration
// +build integration

package event_tracking

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/persistence"
	"casper-dao-middleware/internal/dao/services/event_processing"
	"casper-dao-middleware/pkg/boot"
	"casper-dao-middleware/pkg/casper"
	"casper-dao-middleware/pkg/casper/mocks"
	"casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/go-ces-parser"
)

type TrackVotingCreatedTestSuit struct {
	suite.Suite
	mockCtrl *gomock.Controller

	db            *sqlx.DB
	casperClient  casper.RPCClient
	entityManager persistence.EntityManager

	daoContractsMetadata utils.DAOContractsMetadata
}

func (suite *TrackVotingCreatedTestSuit) SetupSuite() {
	suite.db = boot.SetUpTestDB()

	suite.mockCtrl = gomock.NewController(suite.T())
	reputationContractHash, err := types.NewHashFromHexString("47f30c5dc923b3ddc0d812b6ac5020bcdabdd42f0c1f99c178d6b7869cdf3251")
	assert.NoError(suite.T(), err)

	simpleVoterContractHash, err := types.NewHashFromHexString("04c6b0a1fa97c5fde9017a9d84513da5c709faef0fd6cd3efddd4ab5040bbc90")
	assert.NoError(suite.T(), err)

	suite.daoContractsMetadata = utils.DAOContractsMetadata{
		ReputationContractPackageHash:  reputationContractHash,
		SimpleVoterContractPackageHash: simpleVoterContractHash,
	}

	suite.entityManager = persistence.NewEntityManager(suite.db, suite.daoContractsMetadata)
}

func (suite *TrackVotingCreatedTestSuit) SetupTest() {
	_, err := suite.db.Exec(`TRUNCATE TABLE votings`)
	suite.NoError(err)

	_, err = suite.db.Exec(`TRUNCATE TABLE reputation_changes`)
	suite.NoError(err)

	_, err = suite.db.Exec(`TRUNCATE TABLE votes`)
	suite.NoError(err)
}

func (suite *TrackVotingCreatedTestSuit) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite *TrackVotingCreatedTestSuit) TestTrackReputationVotingCreated() {
	var schemaHex = `08000000100000004164646564546f57686974656c6973740100000007000000616464726573730b0e00000042616c6c6f7443616e63656c65640500000005000000766f7465720b09000000766f74696e675f6964040b000000766f74696e675f74797065030600000063686f69636503050000007374616b65080a00000042616c6c6f74436173740500000005000000766f7465720b09000000766f74696e675f6964040b000000766f74696e675f74797065030600000063686f69636503050000007374616b65080c0000004f776e65724368616e67656401000000090000006e65775f6f776e65720b1400000052656d6f76656446726f6d57686974656c6973740100000007000000616464726573730b1700000052657075746174696f6e566f74696e67437265617465640f000000070000006163636f756e740b06000000616374696f6e0306000000616d6f756e74080d000000646f63756d656e745f686173680a0700000063726561746f720b050000007374616b650d0809000000766f74696e675f69640416000000636f6e6669675f696e666f726d616c5f71756f72756d041b000000636f6e6669675f696e666f726d616c5f766f74696e675f74696d650514000000636f6e6669675f666f726d616c5f71756f72756d0419000000636f6e6669675f666f726d616c5f766f74696e675f74696d650516000000636f6e6669675f746f74616c5f6f6e626f61726465640822000000636f6e6669675f646f75626c655f74696d655f6265747765656e5f766f74696e6773001d000000636f6e6669675f766f74696e675f636c6561726e6573735f64656c7461082e000000636f6e6669675f74696d655f6265747765656e5f696e666f726d616c5f616e645f666f726d616c5f766f74696e67050e000000566f74696e6743616e63656c65640300000009000000766f74696e675f6964040b000000766f74696e675f747970650308000000756e7374616b6573110b080b000000566f74696e67456e6465640d00000009000000766f74696e675f6964040b000000766f74696e675f74797065030d000000766f74696e675f726573756c74030e0000007374616b655f696e5f6661766f72080d0000007374616b655f616761696e73740816000000756e626f756e645f7374616b655f696e5f6661766f720815000000756e626f756e645f7374616b655f616761696e7374080e000000766f7465735f696e5f6661766f72040d000000766f7465735f616761696e73740408000000756e7374616b657311130b0408060000007374616b657311130b0408050000006275726e7311130b0408050000006d696e747311130b0408`

	mockedClient := mocks.NewMockRPCClient(suite.mockCtrl)

	simpleVoterContractHash, err := types.NewHashFromHexString("ea0c001d969da098fefec42b141db88c74c5682e49333ded78035540a0b4f0bc")
	assert.NoError(suite.T(), err)

	mockedClient.EXPECT().GetLatestStateRootHash().Return(casper.GetStateRootHashResult{}, nil)

	mockedClient.EXPECT().GetStateItem("", fmt.Sprintf("hash-%s", simpleVoterContractHash.ToHex()), nil).Return(casper.StateGetItemResult{
		StoredValue: casper.StoredValue{
			Contract: &casper.Contract{
				ContractPackageHash: simpleVoterContractHash,
				NamedKeys: []casper.NamedKey{
					{
						Name: "__events",
						Key:  "uref-1a891d8006f7b56527505c8066945f7d789a75dd48db82dd556a32b35413e0e3-007",
					}, {
						Name: "__events_schema",
						Key:  "uref-1a891d8006f7b56527505c8066945f7d789a75dd48db82dd556a32b35413e0e3-007",
					}},
			},
		},
	}, nil)

	schemaBytes, err := hex.DecodeString(schemaHex)
	assert.NoError(suite.T(), err)

	mockedClient.EXPECT().GetStateItem("", "uref-1a891d8006f7b56527505c8066945f7d789a75dd48db82dd556a32b35413e0e3-007", nil).Return(
		casper.StateGetItemResult{
			StoredValue: casper.StoredValue{
				CLValue: &casper.CLValue{
					Bytes: schemaBytes,
				},
			},
		}, nil)

	suite.casperClient = mockedClient

	var res casper.GetDeployResult

	data, err := os.ReadFile("../../fixtures/events/voting_created/reputation_voting_created.json")
	assert.NoError(suite.T(), err)

	err = json.Unmarshal(data, &res)
	assert.NoError(suite.T(), err)

	cesParser, err := ces.NewParser(suite.casperClient, []types.Hash{simpleVoterContractHash})
	assert.NoError(suite.T(), err)

	processRawDeploy := event_processing.NewProcessRawDeploy()
	processRawDeploy.SetEntityManager(suite.entityManager)
	processRawDeploy.SetCESEventParser(cesParser)
	processRawDeploy.SetDAOContractsMetadata(suite.daoContractsMetadata)

	processRawDeploy.SetDeployProcessedEvent(casper.DeployProcessedEvent{
		DeployProcessed: casper.DeployProcessed{
			DeployHash:      simpleVoterContractHash,
			ExecutionResult: res.ExecutionResults[0].Result,
			Timestamp:       time.Now(),
		},
	})
	assert.NoError(suite.T(), processRawDeploy.Execute())

	count, err := suite.entityManager.VotingRepository().Count(nil)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), count, uint64(1))

	var voting entities.Voting
	err = suite.db.Get(&voting, "select * from votings")
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), voting.VotingID, uint32(1))
	assert.NotEmpty(suite.T(), voting.Metadata)
	assert.Equal(suite.T(), voting.VotingTypeID, entities.VotingTypeReputation)

	var reputationChangesCount int
	err = suite.db.Get(&reputationChangesCount, "select count(*) from reputation_changes")
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), reputationChangesCount, 2)
}

func (suite *TrackVotingCreatedTestSuit) TestTrackSimpleVotingCreated() {
	var schemaHex = `08000000100000004164646564546f57686974656c6973740100000007000000616464726573730b0e00000042616c6c6f7443616e63656c65640500000005000000766f7465720b09000000766f74696e675f6964040b000000766f74696e675f74797065030600000063686f69636503050000007374616b65080a00000042616c6c6f74436173740500000005000000766f7465720b09000000766f74696e675f6964040b000000766f74696e675f74797065030600000063686f69636503050000007374616b65080c0000004f776e65724368616e67656401000000090000006e65775f6f776e65720b1400000052656d6f76656446726f6d57686974656c6973740100000007000000616464726573730b1300000053696d706c65566f74696e67437265617465640c0000000d000000646f63756d656e745f686173680a0700000063726561746f720b050000007374616b650d0809000000766f74696e675f69640416000000636f6e6669675f696e666f726d616c5f71756f72756d041b000000636f6e6669675f696e666f726d616c5f766f74696e675f74696d650514000000636f6e6669675f666f726d616c5f71756f72756d0419000000636f6e6669675f666f726d616c5f766f74696e675f74696d650516000000636f6e6669675f746f74616c5f6f6e626f61726465640822000000636f6e6669675f646f75626c655f74696d655f6265747765656e5f766f74696e6773001d000000636f6e6669675f766f74696e675f636c6561726e6573735f64656c7461082e000000636f6e6669675f74696d655f6265747765656e5f696e666f726d616c5f616e645f666f726d616c5f766f74696e67050e000000566f74696e6743616e63656c65640300000009000000766f74696e675f6964040b000000766f74696e675f747970650308000000756e7374616b6573110b080b000000566f74696e67456e6465640d00000009000000766f74696e675f6964040b000000766f74696e675f74797065030d000000766f74696e675f726573756c74030e0000007374616b655f696e5f6661766f72080d0000007374616b655f616761696e73740816000000756e626f756e645f7374616b655f696e5f6661766f720815000000756e626f756e645f7374616b655f616761696e7374080e000000766f7465735f696e5f6661766f72040d000000766f7465735f616761696e73740408000000756e7374616b657311130b0408060000007374616b657311130b0408050000006275726e7311130b0408050000006d696e747311130b0408`

	mockedClient := mocks.NewMockRPCClient(suite.mockCtrl)

	simpleVoterContractHash, err := types.NewHashFromHexString("ea0c001d969da098fefec42b141db88c74c5682e49333ded78035540a0b4f0bc")
	assert.NoError(suite.T(), err)

	mockedClient.EXPECT().GetLatestStateRootHash().Return(casper.GetStateRootHashResult{}, nil)

	mockedClient.EXPECT().GetStateItem("", fmt.Sprintf("hash-%s", simpleVoterContractHash.ToHex()), nil).Return(casper.StateGetItemResult{
		StoredValue: casper.StoredValue{
			Contract: &casper.Contract{
				ContractPackageHash: simpleVoterContractHash,
				NamedKeys: []casper.NamedKey{
					{
						Name: "__events",
						Key:  "uref-d2263e86f497f42e405d5d1390aa3c1a8bfc35f3699fdc3be806a5cfe139dac9-007",
					}, {
						Name: "__events_schema",
						Key:  "uref-d2263e86f497f42e405d5d1390aa3c1a8bfc35f3699fdc3be806a5cfe139dac9-007",
					}},
			},
		},
	}, nil)

	schemaBytes, err := hex.DecodeString(schemaHex)
	assert.NoError(suite.T(), err)

	mockedClient.EXPECT().GetStateItem("", "uref-d2263e86f497f42e405d5d1390aa3c1a8bfc35f3699fdc3be806a5cfe139dac9-007", nil).Return(
		casper.StateGetItemResult{
			StoredValue: casper.StoredValue{
				CLValue: &casper.CLValue{
					Bytes: schemaBytes,
				},
			},
		}, nil)

	suite.casperClient = mockedClient

	var res casper.GetDeployResult

	data, err := os.ReadFile("../../fixtures/events/voting_created/simple_voting_created.json")
	assert.NoError(suite.T(), err)

	err = json.Unmarshal(data, &res)
	assert.NoError(suite.T(), err)

	cesParser, err := ces.NewParser(suite.casperClient, []types.Hash{simpleVoterContractHash})
	assert.NoError(suite.T(), err)

	processRawDeploy := event_processing.NewProcessRawDeploy()
	processRawDeploy.SetEntityManager(suite.entityManager)
	processRawDeploy.SetCESEventParser(cesParser)
	processRawDeploy.SetDAOContractsMetadata(suite.daoContractsMetadata)

	processRawDeploy.SetDeployProcessedEvent(casper.DeployProcessedEvent{
		DeployProcessed: casper.DeployProcessed{
			DeployHash:      simpleVoterContractHash,
			ExecutionResult: res.ExecutionResults[0].Result,
			Timestamp:       time.Now(),
		},
	})
	assert.NoError(suite.T(), processRawDeploy.Execute())

	count, err := suite.entityManager.VotingRepository().Count(nil)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), count, uint64(1))

	var voting entities.Voting
	err = suite.db.Get(&voting, "select * from votings")
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), voting.VotingID, uint32(0))
	assert.NotEmpty(suite.T(), voting.Metadata)
	assert.Equal(suite.T(), voting.VotingTypeID, entities.VotingTypeSimple)

	var reputationChangesCount int
	err = suite.db.Get(&reputationChangesCount, "select count(*) from reputation_changes")
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), reputationChangesCount, 2)
}

func TestTrackVotingCreatedTestSuit(t *testing.T) {
	suite.Run(t, new(TrackVotingCreatedTestSuit))
}
