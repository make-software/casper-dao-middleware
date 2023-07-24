//go:build integration
// +build integration

package bid

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/rpc"
	"github.com/make-software/casper-go-sdk/sse"
	"github.com/make-software/casper-go-sdk/types/key"
	"github.com/make-software/ces-go-parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"casper-dao-middleware/internal/dao/entities"
	"casper-dao-middleware/internal/dao/persistence"
	"casper-dao-middleware/internal/dao/services/event_processing"
	"casper-dao-middleware/internal/dao/tests/mocks"
	"casper-dao-middleware/internal/dao/utils"
	"casper-dao-middleware/pkg/boot"
)

type TrackBidTestSuit struct {
	suite.Suite
	mockCtrl *gomock.Controller

	db            *sqlx.DB
	casperClient  casper.RPCClient
	entityManager persistence.EntityManager

	daoContractsMetadata utils.DAOContractsMetadata
}

func (suite *TrackBidTestSuit) SetupSuite() {
	suite.db = boot.SetUpTestDB()

	suite.mockCtrl = gomock.NewController(suite.T())
	bidEscrowContractHash, err := casper.NewContractPackageHash("24065460907eb6a86b1917fe30c4a3f47bae7c50e951576d53f0f30a6024a865")
	assert.NoError(suite.T(), err)

	suite.daoContractsMetadata = utils.DAOContractsMetadata{
		BidEscrowContractPackageHash: bidEscrowContractHash,
	}

	suite.entityManager = persistence.NewEntityManager(suite.db, suite.daoContractsMetadata)
}

func (suite *TrackBidTestSuit) SetupTest() {
	_, err := suite.db.Exec(`TRUNCATE TABLE bids`)
	suite.NoError(err)
}

func (suite *TrackBidTestSuit) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite *TrackBidTestSuit) TestTrackBid() {
	var schemaHex = `0e0000000e00000042616c6c6f7443616e63656c65640500000005000000766f7465720b09000000766f74696e675f6964040b000000766f74696e675f74797065040600000063686f69636504050000007374616b65080a00000042616c6c6f74436173740500000005000000766f7465720b09000000766f74696e675f6964040b000000766f74696e675f74797065040600000063686f69636504050000007374616b65080c00000042696443616e63656c6c656403000000060000006269645f6964040600000063616c6c65720b0c0000006a6f625f6f666665725f69640416000000426964457363726f77566f74696e67437265617465640f000000060000006269645f696404060000006a6f625f6964040c0000006a6f625f6f666665725f6964040a0000006a6f625f706f737465720b06000000776f726b65720b0700000063726561746f720b09000000766f74696e675f69640416000000636f6e6669675f696e666f726d616c5f71756f72756d041b000000636f6e6669675f696e666f726d616c5f766f74696e675f74696d650514000000636f6e6669675f666f726d616c5f71756f72756d0419000000636f6e6669675f666f726d616c5f766f74696e675f74696d650516000000636f6e6669675f746f74616c5f6f6e626f61726465640822000000636f6e6669675f646f75626c655f74696d655f6265747765656e5f766f74696e6773001d000000636f6e6669675f766f74696e675f636c6561726e6573735f64656c7461082e000000636f6e6669675f74696d655f6265747765656e5f696e666f726d616c5f616e645f666f726d616c5f766f74696e67050c0000004269645375626d697474656408000000060000006269645f6964040c0000006a6f625f6f666665725f69640406000000776f726b65720b070000006f6e626f617264001200000070726f706f7365645f74696d656672616d65051000000070726f706f7365645f7061796d656e74081000000072657075746174696f6e5f7374616b650d080a000000637370725f7374616b650d080c0000004a6f6243616e63656c6c656405000000060000006269645f6964040600000063616c6c65720b0a0000006a6f625f706f737465720b06000000776f726b65720b0b000000637370725f616d6f756e74080a0000004a6f624372656174656406000000060000006a6f625f696404060000006269645f6964040a0000006a6f625f706f737465720b06000000776f726b65720b0b00000066696e6973685f74696d6505070000007061796d656e7408070000004a6f62446f6e6505000000060000006269645f6964040600000063616c6c65720b0a0000006a6f625f706f737465720b06000000776f726b65720b0b000000637370725f616d6f756e74080f0000004a6f624f6666657243726561746564040000000c0000006a6f625f6f666665725f6964040a0000006a6f625f706f737465720b0a0000006d61785f627564676574081200000065787065637465645f74696d656672616d65050b0000004a6f6252656a656374656405000000060000006269645f6964040600000063616c6c65720b0a0000006a6f625f706f737465720b06000000776f726b65720b0b000000637370725f616d6f756e74080c0000004a6f625375626d697474656404000000060000006269645f6964040a0000006a6f625f706f737465720b06000000776f726b65720b06000000726573756c740a0e000000566f74696e6743616e63656c65640300000009000000766f74696e675f6964040b000000766f74696e675f747970650408000000756e7374616b6573110b0811000000566f74696e6743726561746564496e666f0b0000000700000063726561746f720b050000007374616b650d0809000000766f74696e675f69640416000000636f6e6669675f696e666f726d616c5f71756f72756d041b000000636f6e6669675f696e666f726d616c5f766f74696e675f74696d650514000000636f6e6669675f666f726d616c5f71756f72756d0419000000636f6e6669675f666f726d616c5f766f74696e675f74696d650516000000636f6e6669675f746f74616c5f6f6e626f61726465640822000000636f6e6669675f646f75626c655f74696d655f6265747765656e5f766f74696e6773001d000000636f6e6669675f766f74696e675f636c6561726e6573735f64656c7461082e000000636f6e6669675f74696d655f6265747765656e5f696e666f726d616c5f616e645f666f726d616c5f766f74696e67050b000000566f74696e67456e6465640d00000009000000766f74696e675f6964040b000000766f74696e675f74797065040d000000766f74696e675f726573756c74040e0000007374616b655f696e5f6661766f72080d0000007374616b655f616761696e73740816000000756e626f756e645f7374616b655f696e5f6661766f720815000000756e626f756e645f7374616b655f616761696e7374080e000000766f7465735f696e5f6661766f72040d000000766f7465735f616761696e73740408000000756e7374616b657311130b0408060000007374616b657311130b0408050000006275726e7311130b0408050000006d696e747311130b0408`

	mockedClient := mocks.NewMockRPCClient(suite.mockCtrl)

	bidEscrowContractHash, err := casper.NewHash("8d710e6a825de784e3fb3a7754061a35381a4cdafedcfd667daf65a8ccc70a25")
	assert.NoError(suite.T(), err)

	mockedClient.EXPECT().GetStateRootHashLatest(context.Background()).Return(casper.ChainGetStateRootHashResult{}, nil)

	eventUref, err := casper.NewUref("uref-dcd03b3df9902b5fd1ec4328dbfa691a2fb39af4df84147f29a8ea6adf0c8f15-007")
	assert.NoError(suite.T(), err)

	eventSchemaUref, err := casper.NewUref("uref-1a891d8006f7b56527505c8066945f7d789a75dd48db82dd556a32b35413e0e3-007")
	assert.NoError(suite.T(), err)

	mockedClient.EXPECT().QueryGlobalStateByStateHash(context.Background(), gomock.Any(), fmt.Sprintf("hash-%s", bidEscrowContractHash.ToHex()), nil).Return(rpc.QueryGlobalStateResult{
		StoredValue: casper.StoredValue{
			Contract: &casper.Contract{
				ContractPackageHash: suite.daoContractsMetadata.BidEscrowContractPackageHash,
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

	data, err := os.ReadFile("../../fixtures/events/bid_escrow/submit_bid.json")
	assert.NoError(suite.T(), err)

	err = json.Unmarshal(data, &res)
	assert.NoError(suite.T(), err)

	cesParser, err := ces.NewParser(suite.casperClient, []casper.Hash{bidEscrowContractHash})
	assert.NoError(suite.T(), err)

	processRawDeploy := event_processing.NewProcessRawDeploy()
	processRawDeploy.SetEntityManager(suite.entityManager)
	processRawDeploy.SetCESParser(cesParser)
	processRawDeploy.SetDAOContractsMetadata(suite.daoContractsMetadata)
	processRawDeploy.SetDeployProcessedEvent(sse.DeployProcessedEvent{
		DeployProcessed: sse.DeployProcessedPayload{
			DeployHash:      bidEscrowContractHash,
			ExecutionResult: res.ExecutionResults[0].Result,
			Timestamp:       time.Now(),
		},
	})
	assert.NoError(suite.T(), processRawDeploy.Execute())

	count, err := suite.entityManager.BidRepository().Count(nil)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), count, uint64(1))

	var bid entities.Bid
	err = suite.db.Get(&bid, "select * from bids")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), bid.BidID, uint32(1))
	assert.NotEmpty(suite.T(), bid.Worker)
	assert.Equal(suite.T(), bid.JobOfferID, uint32(2))
	assert.False(suite.T(), bid.Onboard)
	assert.False(suite.T(), bid.PickedByJobPoster)
}

func TestBidEscrowBidsTestSuit(t *testing.T) {
	suite.Run(t, new(TrackBidTestSuit))
}
