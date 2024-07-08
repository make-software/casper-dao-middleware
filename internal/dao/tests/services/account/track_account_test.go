//go:build integration
// +build integration

package account

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

type TrackAccountTestSuit struct {
	suite.Suite
	mockCtrl *gomock.Controller

	db            *sqlx.DB
	casperClient  casper.RPCClient
	entityManager persistence.EntityManager

	daoContractsMetadata utils.DAOContractsMetadata
}

func (suite *TrackAccountTestSuit) SetupSuite() {
	suite.db = boot.SetUpTestDB()

	suite.mockCtrl = gomock.NewController(suite.T())
	vaNFTContractHash, err := casper.NewContractPackageHash("df45d1994399b13d88dc9b9740d9e807f989b12345c30c59b1ebbfbfe99d0473")
	assert.NoError(suite.T(), err)

	kycNFTContractHash, err := casper.NewContractPackageHash("fb4cdd53b23a0359a626e65f579fd54e2826315c13d660283e10f126661abdd5")
	assert.NoError(suite.T(), err)

	suite.daoContractsMetadata = utils.DAOContractsMetadata{
		VANFTContractPackageHash:  vaNFTContractHash,
		KycNFTContractPackageHash: kycNFTContractHash,
	}

	suite.entityManager = persistence.NewEntityManager(suite.db, suite.daoContractsMetadata)
}

func (suite *TrackAccountTestSuit) SetupTest() {
	_, err := suite.db.Exec(`TRUNCATE TABLE accounts`)
	suite.NoError(err)
}

func (suite *TrackAccountTestSuit) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite *TrackAccountTestSuit) TestTransferKycNFT() {
	var schemaHex = `0300000008000000417070726f76616c03000000050000006f776e65720b08000000617070726f7665640d0b08000000746f6b656e5f6964070e000000417070726f76616c466f72416c6c03000000050000006f776e65720b080000006f70657261746f720b08000000617070726f76656400080000005472616e73666572030000000400000066726f6d0d0b02000000746f0d0b08000000746f6b656e5f696407`

	mockedClient := mocks.NewMockRPCClient(suite.mockCtrl)

	kycNFTContractHash, err := casper.NewHash("cade1d82f5150c5b6335eb8b48a9f3a18270d83d5dfe652fae8bffd839696da9")
	assert.NoError(suite.T(), err)

	mockedClient.EXPECT().GetStateRootHashLatest(context.Background()).Return(casper.ChainGetStateRootHashResult{}, nil)

	eventUref, err := casper.NewUref("uref-856586b0cd5fc67b9ccd4a2f345d89cba5a3b697c6019d12b6017d8cd3457dfa-007")
	assert.NoError(suite.T(), err)

	eventSchemaUref, err := casper.NewUref("uref-1a891d8006f7b56527505c8066945f7d789a75dd48db82dd556a32b35413e0e3-007")
	assert.NoError(suite.T(), err)

	mockedClient.EXPECT().QueryGlobalStateByStateHash(context.Background(), gomock.Any(), fmt.Sprintf("hash-%s", kycNFTContractHash.ToHex()), nil).Return(rpc.QueryGlobalStateResult{
		StoredValue: casper.StoredValue{
			Contract: &casper.Contract{
				ContractPackageHash: suite.daoContractsMetadata.KycNFTContractPackageHash,
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

	data, err := os.ReadFile("../../fixtures/events/accounts/kyc_nft_mint.json")
	assert.NoError(suite.T(), err)

	err = json.Unmarshal(data, &res)
	assert.NoError(suite.T(), err)

	cesParser, err := ces.NewParser(suite.casperClient, []casper.Hash{kycNFTContractHash})
	assert.NoError(suite.T(), err)

	processRawDeploy := event_processing.NewProcessRawDeploy()
	processRawDeploy.SetEntityManager(suite.entityManager)
	processRawDeploy.SetCESParser(cesParser)
	processRawDeploy.SetDAOContractsMetadata(suite.daoContractsMetadata)
	processRawDeploy.SetDeployProcessedEvent(sse.DeployProcessedEvent{
		DeployProcessed: sse.DeployProcessedPayload{
			DeployHash:      kycNFTContractHash,
			ExecutionResult: res.ExecutionResults[0].Result,
			Timestamp:       time.Now(),
		},
	})
	assert.NoError(suite.T(), processRawDeploy.Execute())

	count, err := suite.entityManager.AccountRepository().Count(nil)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), count, uint64(1))

	var account entities.Account
	err = suite.db.Get(&account, "select * from accounts")
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), account.IsKyc)
}

func (suite *TrackAccountTestSuit) TestTransferVANFT() {
	var schemaHex = `0300000008000000417070726f76616c03000000050000006f776e65720b08000000617070726f7665640d0b08000000746f6b656e5f6964070e000000417070726f76616c466f72416c6c03000000050000006f776e65720b080000006f70657261746f720b08000000617070726f76656400080000005472616e73666572030000000400000066726f6d0d0b02000000746f0d0b08000000746f6b656e5f696407`

	mockedClient := mocks.NewMockRPCClient(suite.mockCtrl)

	vaNFTContractHash, err := casper.NewHash("c4d2c0499ba8fe7546006d77e010fa5c8dfd1b2a982028fa2188439b53a140e0")
	assert.NoError(suite.T(), err)

	mockedClient.EXPECT().GetStateRootHashLatest(context.Background()).Return(casper.ChainGetStateRootHashResult{}, nil)

	eventUref, err := casper.NewUref("uref-ecd658c007ee088f4eea8f88f3a9815f4abc1de412b6777fd410d0e47d707b96-007")
	assert.NoError(suite.T(), err)

	eventSchemaUref, err := casper.NewUref("uref-1a891d8006f7b56527505c8066945f7d789a75dd48db82dd556a32b35413e0e3-007")
	assert.NoError(suite.T(), err)

	mockedClient.EXPECT().QueryGlobalStateByStateHash(context.Background(), gomock.Any(), fmt.Sprintf("hash-%s", vaNFTContractHash.ToHex()), nil).Return(rpc.QueryGlobalStateResult{
		StoredValue: casper.StoredValue{
			Contract: &casper.Contract{
				ContractPackageHash: suite.daoContractsMetadata.VANFTContractPackageHash,
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

	data, err := os.ReadFile("../../fixtures/events/accounts/va_nft_mint.json")
	assert.NoError(suite.T(), err)

	err = json.Unmarshal(data, &res)
	assert.NoError(suite.T(), err)

	cesParser, err := ces.NewParser(suite.casperClient, []casper.Hash{vaNFTContractHash})
	assert.NoError(suite.T(), err)

	processRawDeploy := event_processing.NewProcessRawDeploy()
	processRawDeploy.SetEntityManager(suite.entityManager)
	processRawDeploy.SetCESParser(cesParser)
	processRawDeploy.SetDAOContractsMetadata(suite.daoContractsMetadata)
	processRawDeploy.SetDeployProcessedEvent(sse.DeployProcessedEvent{
		DeployProcessed: sse.DeployProcessedPayload{
			DeployHash:      vaNFTContractHash,
			ExecutionResult: res.ExecutionResults[0].Result,
			Timestamp:       time.Now(),
		},
	})
	assert.NoError(suite.T(), processRawDeploy.Execute())

	count, err := suite.entityManager.AccountRepository().Count(nil)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), count, uint64(1))

	var account entities.Account
	err = suite.db.Get(&account, "select * from accounts")
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), account.IsVA)
}

func TestTrackAccountsEventTestSuit(t *testing.T) {
	suite.Run(t, new(TrackAccountTestSuit))
}
