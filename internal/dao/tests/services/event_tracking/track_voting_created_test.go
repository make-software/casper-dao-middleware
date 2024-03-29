//go:build integration
// +build integration

package event_tracking

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/make-software/casper-go-sdk/sse"
	"github.com/make-software/casper-go-sdk/types/key"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/make-software/casper-go-sdk/casper"

	"github.com/make-software/ces-go-parser"

	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/persistence"
	"casper-dao-middleware/internal/dao/services/event_processing"
	"casper-dao-middleware/internal/dao/tests/mocks"
	"casper-dao-middleware/internal/dao/utils"
	"casper-dao-middleware/pkg/boot"
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
	reputationContractHash, err := casper.NewContractPackageHash("ea0c001d969da098fefec42b141db88c74c5682e49333ded78035540a0b4f0bc")
	assert.NoError(suite.T(), err)

	simpleVoterContractHash, err := casper.NewContractPackageHash("954998ff95b0210e994f43f5afb5174b5085fda92d7c63962ef09c17886658c1")
	assert.NoError(suite.T(), err)

	repoVoterContractHash, err := casper.NewContractPackageHash("6a3213fe5db928dd4bb3d1c5ecd3bfbc68656823c9486ef389a3080921d0d3ec")
	assert.NoError(suite.T(), err)

	suite.daoContractsMetadata = utils.DAOContractsMetadata{
		ReputationVoterContractPackageHash: reputationContractHash,
		SimpleVoterContractPackageHash:     simpleVoterContractHash,
		RepoVoterContractPackageHash:       repoVoterContractHash,
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

	mockedClient := mocks.NewMockClient(suite.mockCtrl)

	reputationVoterContractHash, err := casper.NewHash("ea0c001d969da098fefec42b141db88c74c5682e49333ded78035540a0b4f0bc")
	assert.NoError(suite.T(), err)

	reputationVoterContractPackageHash, err := casper.NewContractPackageHash("ea0c001d969da098fefec42b141db88c74c5682e49333ded78035540a0b4f0bc")
	assert.NoError(suite.T(), err)

	mockedClient.EXPECT().GetStateRootHashLatest(context.Background()).Return(casper.ChainGetStateRootHashResult{}, nil)

	eventUref, err := casper.NewUref("uref-1a891d8006f7b56527505c8066945f7d789a75dd48db82dd556a32b35413e0e3-007")
	assert.NoError(suite.T(), err)

	eventSchemaUref, err := casper.NewUref("uref-1a891d8006f7b56527505c8066945f7d789a75dd48db82dd556a32b35413e0e3-007")
	assert.NoError(suite.T(), err)

	mockedClient.EXPECT().GetStateItem(context.Background(), "0000000000000000000000000000000000000000000000000000000000000000", fmt.Sprintf("hash-%s", reputationVoterContractHash.ToHex()), nil).Return(casper.StateGetItemResult{
		StoredValue: casper.StoredValue{
			Contract: &casper.Contract{
				ContractPackageHash: reputationVoterContractPackageHash,
				NamedKeys: []casper.NamedKey{
					{
						Name: "__events",
						Key: key.Key{
							Type: key.TypeIDURef,
							URef: &eventUref,
						},
					}, {
						Name: "__events_schema",
						Key:  key.Key{Type: key.TypeIDURef, URef: &eventSchemaUref},
					}},
			},
		},
	}, nil)

	var arg casper.Argument
	err = json.Unmarshal([]byte(fmt.Sprintf(`{"cl_type": "Any", "bytes": "%s"}`, schemaHex)), &arg)
	require.NoError(suite.T(), err)

	mockedClient.EXPECT().GetStateItem(context.Background(), "0000000000000000000000000000000000000000000000000000000000000000", "uref-1a891d8006f7b56527505c8066945f7d789a75dd48db82dd556a32b35413e0e3-007", nil).Return(
		casper.StateGetItemResult{
			StoredValue: casper.StoredValue{
				CLValue: &arg,
			},
		}, nil)

	suite.casperClient = mockedClient

	var res casper.InfoGetDeployResult

	data, err := os.ReadFile("../../fixtures/events/voting_created/reputation_voting_created.json")
	assert.NoError(suite.T(), err)

	err = json.Unmarshal(data, &res)
	assert.NoError(suite.T(), err)

	cesParser, err := ces.NewParser(suite.casperClient, []casper.Hash{reputationVoterContractHash})
	assert.NoError(suite.T(), err)

	creator, err := casper.NewHash("ea0c001d969da098fefec42b141db88c74c5682e49333ded78035540a0b4f0bc")
	assert.NoError(suite.T(), err)

	deployHash, err := casper.NewHash("ea0c001d969da098fefec42b141db88c74c5682e49333ded78035540a0b4f0bc")
	assert.NoError(suite.T(), err)

	err = suite.entityManager.VotingRepository().Save(&entities.Voting{
		Creator:                                  creator,
		DeployHash:                               deployHash,
		VotingID:                                 1,
		VotingTypeID:                             0,
		InformalVotingQuorum:                     0,
		InformalVotingStartsAt:                   time.Now(),
		InformalVotingEndsAt:                     time.Now(),
		FormalVotingQuorum:                       0,
		FormalVotingTime:                         0,
		FormalVotingStartsAt:                     nil,
		FormalVotingEndsAt:                       nil,
		Metadata:                                 nil,
		IsCanceled:                               false,
		InformalVotingResult:                     nil,
		FormalVotingResult:                       nil,
		ConfigTotalOnboarded:                     0,
		ConfigVotingClearnessDelta:               0,
		ConfigTimeBetweenInformalAndFormalVoting: 0,
	})
	assert.NoError(suite.T(), err)

	processRawDeploy := event_processing.NewProcessRawDeploy()
	processRawDeploy.SetEntityManager(suite.entityManager)
	processRawDeploy.SetCESEventParser(cesParser)
	processRawDeploy.SetDAOContractsMetadata(suite.daoContractsMetadata)

	processRawDeploy.SetDeployProcessedEvent(sse.DeployProcessedEvent{
		DeployProcessed: sse.DeployProcessed{
			DeployHash:      reputationVoterContractHash,
			ExecutionResult: res.ExecutionResults[0].Result,
			Timestamp:       time.Now(),
		},
	})
	assert.NoError(suite.T(), processRawDeploy.Execute())

	count, err := suite.entityManager.VotingRepository().Count(nil)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), count, uint64(2))

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

	mockedClient := mocks.NewMockClient(suite.mockCtrl)

	simpleVoterContractPackageHash, err := casper.NewContractPackageHash("954998ff95b0210e994f43f5afb5174b5085fda92d7c63962ef09c17886658c1")
	assert.NoError(suite.T(), err)

	simpleVoterContractHash, err := casper.NewHash("954998ff95b0210e994f43f5afb5174b5085fda92d7c63962ef09c17886658c1")
	assert.NoError(suite.T(), err)

	mockedClient.EXPECT().GetStateRootHashLatest(context.Background()).Return(casper.ChainGetStateRootHashResult{}, nil)

	eventUref, err := casper.NewUref("uref-d2263e86f497f42e405d5d1390aa3c1a8bfc35f3699fdc3be806a5cfe139dac9-007")
	assert.NoError(suite.T(), err)

	eventSchemaUref, err := casper.NewUref("uref-d2263e86f497f42e405d5d1390aa3c1a8bfc35f3699fdc3be806a5cfe139dac9-007")
	assert.NoError(suite.T(), err)

	mockedClient.EXPECT().GetStateItem(context.Background(), "0000000000000000000000000000000000000000000000000000000000000000", fmt.Sprintf("hash-%s", simpleVoterContractHash.ToHex()), nil).Return(casper.StateGetItemResult{
		StoredValue: casper.StoredValue{
			Contract: &casper.Contract{
				ContractPackageHash: simpleVoterContractPackageHash,
				NamedKeys: []casper.NamedKey{
					{
						Name: "__events",
						Key:  key.Key{Type: key.TypeIDURef, URef: &eventUref},
					}, {
						Name: "__events_schema",
						Key:  key.Key{Type: key.TypeIDURef, URef: &eventSchemaUref},
					}},
			},
		},
	}, nil)

	var arg casper.Argument
	err = json.Unmarshal([]byte(fmt.Sprintf(`{"cl_type": "Any", "bytes": "%s"}`, schemaHex)), &arg)
	require.NoError(suite.T(), err)

	mockedClient.EXPECT().GetStateItem(context.Background(), "0000000000000000000000000000000000000000000000000000000000000000", "uref-d2263e86f497f42e405d5d1390aa3c1a8bfc35f3699fdc3be806a5cfe139dac9-007", nil).Return(
		casper.StateGetItemResult{
			StoredValue: casper.StoredValue{
				CLValue: &arg,
			},
		}, nil)

	suite.casperClient = mockedClient

	var res casper.InfoGetDeployResult

	data, err := os.ReadFile("../../fixtures/events/voting_created/simple_voting_created.json")
	assert.NoError(suite.T(), err)

	err = json.Unmarshal(data, &res)
	assert.NoError(suite.T(), err)

	cesParser, err := ces.NewParser(suite.casperClient, []casper.Hash{simpleVoterContractHash})
	assert.NoError(suite.T(), err)

	creator, err := casper.NewHash("ea0c001d969da098fefec42b141db88c74c5682e49333ded78035540a0b4f0bc")
	assert.NoError(suite.T(), err)

	deployHash, err := casper.NewHash("ea0c001d969da098fefec42b141db88c74c5682e49333ded78035540a0b4f0bc")
	assert.NoError(suite.T(), err)

	err = suite.entityManager.VotingRepository().Save(&entities.Voting{
		Creator:                                  creator,
		DeployHash:                               deployHash,
		VotingID:                                 0,
		VotingTypeID:                             0,
		InformalVotingQuorum:                     0,
		InformalVotingStartsAt:                   time.Now(),
		InformalVotingEndsAt:                     time.Now(),
		FormalVotingQuorum:                       0,
		FormalVotingTime:                         0,
		FormalVotingStartsAt:                     nil,
		FormalVotingEndsAt:                       nil,
		Metadata:                                 nil,
		IsCanceled:                               false,
		InformalVotingResult:                     nil,
		FormalVotingResult:                       nil,
		ConfigTotalOnboarded:                     0,
		ConfigVotingClearnessDelta:               0,
		ConfigTimeBetweenInformalAndFormalVoting: 0,
	})
	assert.NoError(suite.T(), err)

	processRawDeploy := event_processing.NewProcessRawDeploy()
	processRawDeploy.SetEntityManager(suite.entityManager)
	processRawDeploy.SetCESEventParser(cesParser)
	processRawDeploy.SetDAOContractsMetadata(suite.daoContractsMetadata)

	processRawDeploy.SetDeployProcessedEvent(sse.DeployProcessedEvent{
		DeployProcessed: sse.DeployProcessed{
			DeployHash:      simpleVoterContractHash,
			ExecutionResult: res.ExecutionResults[0].Result,
			Timestamp:       time.Now(),
		},
	})
	assert.NoError(suite.T(), processRawDeploy.Execute())

	count, err := suite.entityManager.VotingRepository().Count(nil)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), count, uint64(2))

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

func (suite *TrackVotingCreatedTestSuit) TestTrackRepoVoterVotingCreated() {
	var schemaHex = `08000000100000004164646564546f57686974656c6973740100000007000000616464726573730b0e00000042616c6c6f7443616e63656c65640500000005000000766f7465720b09000000766f74696e675f6964040b000000766f74696e675f74797065030600000063686f69636503050000007374616b65080a00000042616c6c6f74436173740500000005000000766f7465720b09000000766f74696e675f6964040b000000766f74696e675f74797065030600000063686f69636503050000007374616b65080c0000004f776e65724368616e67656401000000090000006e65775f6f776e65720b1400000052656d6f76656446726f6d57686974656c6973740100000007000000616464726573730b110000005265706f566f74696e67437265617465640f000000150000007661726961626c655f7265706f5f746f5f656469740b030000006b65790a0500000076616c75650e030f00000061637469766174696f6e5f74696d650d050700000063726561746f720b050000007374616b650d0809000000766f74696e675f69640416000000636f6e6669675f696e666f726d616c5f71756f72756d041b000000636f6e6669675f696e666f726d616c5f766f74696e675f74696d650514000000636f6e6669675f666f726d616c5f71756f72756d0419000000636f6e6669675f666f726d616c5f766f74696e675f74696d650516000000636f6e6669675f746f74616c5f6f6e626f61726465640822000000636f6e6669675f646f75626c655f74696d655f6265747765656e5f766f74696e6773001d000000636f6e6669675f766f74696e675f636c6561726e6573735f64656c7461082e000000636f6e6669675f74696d655f6265747765656e5f696e666f726d616c5f616e645f666f726d616c5f766f74696e67050e000000566f74696e6743616e63656c65640300000009000000766f74696e675f6964040b000000766f74696e675f747970650308000000756e7374616b6573110b080b000000566f74696e67456e6465640d00000009000000766f74696e675f6964040b000000766f74696e675f74797065030d000000766f74696e675f726573756c74030e0000007374616b655f696e5f6661766f72080d0000007374616b655f616761696e73740816000000756e626f756e645f7374616b655f696e5f6661766f720815000000756e626f756e645f7374616b655f616761696e7374080e000000766f7465735f696e5f6661766f72040d000000766f7465735f616761696e73740408000000756e7374616b657311130b0408060000007374616b657311130b0408050000006275726e7311130b0408050000006d696e747311130b0408`
	mockedClient := mocks.NewMockClient(suite.mockCtrl)

	repoVoterContractHash, err := casper.NewHash("6a3213fe5db928dd4bb3d1c5ecd3bfbc68656823c9486ef389a3080921d0d3ec")
	assert.NoError(suite.T(), err)

	repoVoterContractPackageHash, err := casper.NewContractPackageHash("6a3213fe5db928dd4bb3d1c5ecd3bfbc68656823c9486ef389a3080921d0d3ec")
	assert.NoError(suite.T(), err)

	mockedClient.EXPECT().GetStateRootHashLatest(context.Background()).Return(casper.ChainGetStateRootHashResult{}, nil)

	eventUref, err := casper.NewUref("uref-26babcb0c9924a1f983d1aaf7b32d718c34b9ea85fc12f3b6c46068c86422079-007")
	assert.NoError(suite.T(), err)

	eventSchemaUref, err := casper.NewUref("uref-114af2ad7972ddadf4b92d62170ce14fad21741895cf24415d24ace244a91a9f-007")
	assert.NoError(suite.T(), err)

	mockedClient.EXPECT().GetStateItem(context.Background(), "0000000000000000000000000000000000000000000000000000000000000000", fmt.Sprintf("hash-%s", repoVoterContractHash.ToHex()), nil).Return(casper.StateGetItemResult{
		StoredValue: casper.StoredValue{
			Contract: &casper.Contract{
				ContractPackageHash: repoVoterContractPackageHash,
				NamedKeys: []casper.NamedKey{
					{
						Name: "__events",
						Key:  key.Key{Type: key.TypeIDURef, URef: &eventUref},
					}, {
						Name: "__events_schema",
						Key:  key.Key{Type: key.TypeIDURef, URef: &eventSchemaUref},
					}},
			},
		},
	}, nil)

	var arg casper.Argument
	err = json.Unmarshal([]byte(fmt.Sprintf(`{"cl_type": "Any", "bytes": "%s"}`, schemaHex)), &arg)
	require.NoError(suite.T(), err)

	mockedClient.EXPECT().GetStateItem(context.Background(),
		"0000000000000000000000000000000000000000000000000000000000000000", "uref-114af2ad7972ddadf4b92d62170ce14fad21741895cf24415d24ace244a91a9f-007", nil).Return(
		casper.StateGetItemResult{
			StoredValue: casper.StoredValue{
				CLValue: &arg,
			},
		}, nil)

	suite.casperClient = mockedClient

	var res casper.InfoGetDeployResult

	data, err := os.ReadFile("../../fixtures/events/voting_created/repo_voter_voting_created.json")
	assert.NoError(suite.T(), err)

	err = json.Unmarshal(data, &res)
	assert.NoError(suite.T(), err)

	cesParser, err := ces.NewParser(suite.casperClient, []casper.Hash{repoVoterContractHash})
	assert.NoError(suite.T(), err)

	creator, err := casper.NewHash("ea0c001d969da098fefec42b141db88c74c5682e49333ded78035540a0b4f0bc")
	assert.NoError(suite.T(), err)

	deployHash, err := casper.NewHash("ea0c001d969da098fefec42b141db88c74c5682e49333ded78035540a0b4f0bc")
	assert.NoError(suite.T(), err)

	err = suite.entityManager.VotingRepository().Save(&entities.Voting{
		Creator:                                  creator,
		DeployHash:                               deployHash,
		VotingID:                                 2,
		VotingTypeID:                             4,
		InformalVotingQuorum:                     0,
		InformalVotingStartsAt:                   time.Now(),
		InformalVotingEndsAt:                     time.Now(),
		FormalVotingQuorum:                       0,
		FormalVotingTime:                         0,
		FormalVotingStartsAt:                     nil,
		FormalVotingEndsAt:                       nil,
		Metadata:                                 nil,
		IsCanceled:                               false,
		InformalVotingResult:                     nil,
		FormalVotingResult:                       nil,
		ConfigTotalOnboarded:                     0,
		ConfigVotingClearnessDelta:               0,
		ConfigTimeBetweenInformalAndFormalVoting: 0,
	})
	assert.NoError(suite.T(), err)

	processRawDeploy := event_processing.NewProcessRawDeploy()
	processRawDeploy.SetEntityManager(suite.entityManager)
	processRawDeploy.SetCESEventParser(cesParser)
	processRawDeploy.SetDAOContractsMetadata(suite.daoContractsMetadata)

	processRawDeploy.SetDeployProcessedEvent(sse.DeployProcessedEvent{
		DeployProcessed: sse.DeployProcessed{
			DeployHash:      repoVoterContractHash,
			ExecutionResult: res.ExecutionResults[0].Result,
			Timestamp:       time.Now(),
		},
	})
	assert.NoError(suite.T(), processRawDeploy.Execute())

	count, err := suite.entityManager.VotingRepository().Count(nil)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), count, uint64(2))

	var voting entities.Voting
	err = suite.db.Get(&voting, "select * from votings")
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), voting.VotingID, uint32(2))
	assert.NotEmpty(suite.T(), voting.Metadata)
	assert.Equal(suite.T(), voting.VotingTypeID, entities.VotingTypeRepo)

	var reputationChangesCount int
	err = suite.db.Get(&reputationChangesCount, "select count(*) from reputation_changes")
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), reputationChangesCount, 2)
}

func TestTrackVotingCreatedTestSuit(t *testing.T) {
	suite.Run(t, new(TrackVotingCreatedTestSuit))
}
