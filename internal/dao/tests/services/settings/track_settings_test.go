//go:build integration
// +build integration

package account

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/make-software/casper-go-sdk/casper"
	"github.com/make-software/casper-go-sdk/rpc"
	"github.com/make-software/casper-go-sdk/types/key"
	"github.com/make-software/ces-go-parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"casper-dao-middleware/internal/dao/persistence"
	"casper-dao-middleware/internal/dao/services/settings"
	"casper-dao-middleware/internal/dao/tests/mocks"
	"casper-dao-middleware/internal/dao/utils"
	"casper-dao-middleware/pkg/boot"
)

type TrackSettingsTestSuit struct {
	suite.Suite
	mockCtrl *gomock.Controller

	db            *sqlx.DB
	casperClient  casper.RPCClient
	entityManager persistence.EntityManager

	daoContractsMetadata utils.DAOContractsMetadata
}

func (suite *TrackSettingsTestSuit) SetupSuite() {
	suite.db = boot.SetUpTestDB()

	suite.mockCtrl = gomock.NewController(suite.T())
	variableRepoContractHash, err := casper.NewContractPackageHash("df45d1994399b13d88dc9b9740d9e807f989b12345c30c59b1ebbfbfe99d0473")
	assert.NoError(suite.T(), err)

	suite.daoContractsMetadata = utils.DAOContractsMetadata{
		VariableRepositoryContractPackageHash: variableRepoContractHash,
	}

	suite.entityManager = persistence.NewEntityManager(suite.db, suite.daoContractsMetadata)
}

func (suite *TrackSettingsTestSuit) SetupTest() {
	_, err := suite.db.Exec(`TRUNCATE TABLE settings`)
	suite.NoError(err)
}

func (suite *TrackSettingsTestSuit) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite *TrackSettingsTestSuit) TestTrackSettingsFromDeploy() {
	var schemaHex = `010000000c00000056616c75655570646174656403000000030000006b65790a0500000076616c75650e030f00000061637469766174696f6e5f74696d650d05`

	mockedClient := mocks.NewMockRPCClient(suite.mockCtrl)

	variableRepositoryContractHash, err := casper.NewHash("aae41f28e461d01515722169648843b05e09b9aa649d247ea765442a62a876ec")
	assert.NoError(suite.T(), err)

	mockedClient.EXPECT().GetStateRootHashLatest(context.Background()).Return(casper.ChainGetStateRootHashResult{}, nil)

	eventUref, err := casper.NewUref("uref-1d460d63cb7d606c862922159a3a76717616663bb7f0db7f2a43f78e5dfb2ec7-007")
	assert.NoError(suite.T(), err)

	eventSchemaUref, err := casper.NewUref("uref-1a891d8006f7b56527505c8066945f7d789a75dd48db82dd556a32b35413e0e3-007")
	assert.NoError(suite.T(), err)

	mockedClient.EXPECT().QueryGlobalStateByStateHash(context.Background(), gomock.Any(), fmt.Sprintf("hash-%s", variableRepositoryContractHash.ToHex()), nil).Return(rpc.QueryGlobalStateResult{
		StoredValue: casper.StoredValue{
			Contract: &casper.Contract{
				ContractPackageHash: suite.daoContractsMetadata.VariableRepositoryContractPackageHash,
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
	data, err := os.ReadFile("../../fixtures/events/settings/install_variable_repository_deploy.json")
	assert.NoError(suite.T(), err)

	err = json.Unmarshal(data, &res)
	assert.NoError(suite.T(), err)

	variableRepoInstallDeployHash, err := casper.NewHash("67e9507a604ea8a5c31c116e4c3ca9d888a34b573035a7d887fb08bd62962434")
	assert.NoError(suite.T(), err)

	mockedClient.EXPECT().GetDeploy(context.Background(), variableRepoInstallDeployHash.ToHex()).Return(res, nil)

	cesParser, err := ces.NewParser(suite.casperClient, []casper.Hash{variableRepositoryContractHash})
	assert.NoError(suite.T(), err)

	syncDaoSetting := settings.NewSyncInitialDAOSettings()
	syncDaoSetting.SetCasperClient(suite.casperClient)
	syncDaoSetting.SetVariableRepoInstallDeployHash(variableRepoInstallDeployHash)
	syncDaoSetting.SetDAOContractsMetadata(suite.daoContractsMetadata)
	syncDaoSetting.SetEntityManager(suite.entityManager)
	syncDaoSetting.SetCESParser(cesParser)
	assert.NoError(suite.T(), syncDaoSetting.Execute())

	count, err := suite.entityManager.SettingRepository().Count(nil)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), count, uint64(26))
}

func TestTrackSettings(t *testing.T) {
	suite.Run(t, new(TrackSettingsTestSuit))
}
