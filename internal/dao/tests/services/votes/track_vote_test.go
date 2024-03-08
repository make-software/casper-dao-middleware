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
	"github.com/make-software/casper-go-sdk/rpc"
	"github.com/make-software/casper-go-sdk/sse"
	"github.com/make-software/casper-go-sdk/types/key"
	"github.com/make-software/ces-go-parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/make-software/casper-go-sdk/casper"

	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/persistence"
	"casper-dao-middleware/internal/dao/services/event_processing"
	"casper-dao-middleware/internal/dao/tests/mocks"
	"casper-dao-middleware/internal/dao/utils"
	"casper-dao-middleware/pkg/boot"
)

type TrackVoteTestSuit struct {
	suite.Suite
	mockCtrl *gomock.Controller

	db            *sqlx.DB
	casperClient  casper.RPCClient
	entityManager persistence.EntityManager

	daoContractsMetadata utils.DAOContractsMetadata
}

func (suite *TrackVoteTestSuit) SetupSuite() {
	suite.db = boot.SetUpTestDB()

	suite.mockCtrl = gomock.NewController(suite.T())
	simpleVoterContractHash, err := casper.NewContractPackageHash("6701fd5c747c8fe643184dcd2e143a85ab207bc1865a90211c9e2d0549abb5c1")
	assert.NoError(suite.T(), err)
	suite.daoContractsMetadata = utils.DAOContractsMetadata{
		SimpleVoterContractPackageHash: simpleVoterContractHash,
	}

	suite.entityManager = persistence.NewEntityManager(suite.db, suite.daoContractsMetadata)
}

func (suite *TrackVoteTestSuit) SetupTest() {
	_, err := suite.db.Exec(`TRUNCATE TABLE reputation_changes`)
	suite.NoError(err)

	_, err = suite.db.Exec(`TRUNCATE TABLE votes`)
	suite.NoError(err)
}

func (suite *TrackVoteTestSuit) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite *TrackVoteTestSuit) TestTrackVote() {
	var schemaHex = `060000000e00000042616c6c6f7443616e63656c65640500000005000000766f7465720b09000000766f74696e675f6964040b000000766f74696e675f74797065040600000063686f69636504050000007374616b65080a00000042616c6c6f74436173740500000005000000766f7465720b09000000766f74696e675f6964040b000000766f74696e675f74797065040600000063686f69636504050000007374616b65081700000052657075746174696f6e566f74696e67437265617465640f000000070000006163636f756e740b06000000616374696f6e0406000000616d6f756e74080d000000646f63756d656e745f686173680a0700000063726561746f720b050000007374616b650d0809000000766f74696e675f69640416000000636f6e6669675f696e666f726d616c5f71756f72756d041b000000636f6e6669675f696e666f726d616c5f766f74696e675f74696d650514000000636f6e6669675f666f726d616c5f71756f72756d0419000000636f6e6669675f666f726d616c5f766f74696e675f74696d650516000000636f6e6669675f746f74616c5f6f6e626f61726465640822000000636f6e6669675f646f75626c655f74696d655f6265747765656e5f766f74696e6773001d000000636f6e6669675f766f74696e675f636c6561726e6573735f64656c7461082e000000636f6e6669675f74696d655f6265747765656e5f696e666f726d616c5f616e645f666f726d616c5f766f74696e67050e000000566f74696e6743616e63656c65640300000009000000766f74696e675f6964040b000000766f74696e675f747970650408000000756e7374616b6573110b0811000000566f74696e6743726561746564496e666f0b0000000700000063726561746f720b050000007374616b650d0809000000766f74696e675f69640416000000636f6e6669675f696e666f726d616c5f71756f72756d041b000000636f6e6669675f696e666f726d616c5f766f74696e675f74696d650514000000636f6e6669675f666f726d616c5f71756f72756d0419000000636f6e6669675f666f726d616c5f766f74696e675f74696d650516000000636f6e6669675f746f74616c5f6f6e626f61726465640822000000636f6e6669675f646f75626c655f74696d655f6265747765656e5f766f74696e6773001d000000636f6e6669675f766f74696e675f636c6561726e6573735f64656c7461082e000000636f6e6669675f74696d655f6265747765656e5f696e666f726d616c5f616e645f666f726d616c5f766f74696e67050b000000566f74696e67456e6465640d00000009000000766f74696e675f6964040b000000766f74696e675f74797065040d000000766f74696e675f726573756c74040e0000007374616b655f696e5f6661766f72080d0000007374616b655f616761696e73740816000000756e626f756e645f7374616b655f696e5f6661766f720815000000756e626f756e645f7374616b655f616761696e7374080e000000766f7465735f696e5f6661766f72040d000000766f7465735f616761696e73740408000000756e7374616b657311130b0408060000007374616b657311130b0408050000006275726e7311130b0408050000006d696e747311130b0408`

	mockedClient := mocks.NewMockRPCClient(suite.mockCtrl)

	simpleVoterContractPackageHash, _ := casper.NewHash("4ca76dbcb92d79bc21363ced9e502a3a44e3696470645b823e6413ab97f88659")
	mockedClient.EXPECT().GetStateRootHashLatest(context.Background()).Return(casper.ChainGetStateRootHashResult{}, nil)

	eventUref, err := casper.NewUref("uref-9f7bb81aad08c82ddc2a02be79d6482e061a7476d435f4fd88ce2cb147143a74-007")
	assert.NoError(suite.T(), err)

	eventSchemaUref, err := casper.NewUref("uref-1a891d8006f7b56527505c8066945f7d789a75dd48db82dd556a32b35413e0e3-007")
	assert.NoError(suite.T(), err)

	mockedClient.EXPECT().QueryGlobalStateByStateHash(context.Background(), gomock.Any(), fmt.Sprintf("hash-%s", simpleVoterContractPackageHash.ToHex()), nil).Return(rpc.QueryGlobalStateResult{
		StoredValue: casper.StoredValue{
			Contract: &casper.Contract{
				ContractPackageHash: suite.daoContractsMetadata.SimpleVoterContractPackageHash,
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

	mockedClient.EXPECT().QueryGlobalStateByStateHash(context.Background(), gomock.Any(), "uref-1a891d8006f7b56527505c8066945f7d789a75dd48db82dd556a32b35413e0e3-007", nil).Return(
		rpc.QueryGlobalStateResult{
			StoredValue: casper.StoredValue{
				CLValue: &arg,
			},
		}, nil)

	suite.casperClient = mockedClient

	var res casper.InfoGetDeployResult

	data, err := os.ReadFile("../../fixtures/events/vote/vote.json")
	assert.NoError(suite.T(), err)

	err = json.Unmarshal(data, &res)
	assert.NoError(suite.T(), err)

	cesParser, err := ces.NewParser(suite.casperClient, []casper.Hash{simpleVoterContractPackageHash})
	assert.NoError(suite.T(), err)

	processRawDeploy := event_processing.NewProcessRawDeploy()
	processRawDeploy.SetEntityManager(suite.entityManager)
	processRawDeploy.SetCESParser(cesParser)
	processRawDeploy.SetDAOContractsMetadata(suite.daoContractsMetadata)
	processRawDeploy.SetDeployProcessedEvent(sse.DeployProcessedEvent{
		DeployProcessed: sse.DeployProcessedPayload{
			DeployHash:      simpleVoterContractPackageHash,
			ExecutionResult: res.ExecutionResults[0].Result,
			Timestamp:       time.Now(),
		},
	})
	assert.NoError(suite.T(), processRawDeploy.Execute())

	count, err := suite.entityManager.VoteRepository().Count(nil)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), count, uint64(1))

	var vote entities.Vote
	err = suite.db.Get(&vote, "select * from votes")
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), vote.VotingID, uint32(5))
	assert.NotEmpty(suite.T(), vote.Address)

	var reputationChangesCount int
	err = suite.db.Get(&reputationChangesCount, "select count(*) from reputation_changes")
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), reputationChangesCount, 2)
}

func TestTrackVoteTestSuit(t *testing.T) {
	suite.Run(t, new(TrackVoteTestSuit))
}
