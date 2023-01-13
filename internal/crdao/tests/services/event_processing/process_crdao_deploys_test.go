//go:build integration
// +build integration

package event_processing

import (
	"casper-dao-middleware/pkg/boot"
	"encoding/json"
	"os"
	"testing"
	"time"

	"casper-dao-middleware/internal/crdao/dao_event_parser"
	"casper-dao-middleware/internal/crdao/entities"
	"casper-dao-middleware/internal/crdao/persistence"
	"casper-dao-middleware/internal/crdao/services/event_processing"
	"casper-dao-middleware/pkg/casper"
	"casper-dao-middleware/pkg/casper/types"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ProcessDAODeploysTestSuit struct {
	suite.Suite

	db            *sqlx.DB
	casperClient  casper.RPCClient
	entityManager persistence.EntityManager

	daoContractHashesMap     map[string]string
	daoContractPackageHashes dao_event_parser.DAOContractsMetadata
}

func (suite *ProcessDAODeploysTestSuit) SetupSuite() {
	suite.db = boot.SetUpTestDB()
	suite.casperClient = casper.NewRPCClient(os.Getenv("INTEGRATION_NET_RPC_API_URL"))

	suite.daoContractHashesMap = map[string]string{
		"reputation_contract": "3eaf08d521d01da82151db1f6f5a88c5acea55109b6b84998c8a9073b93b02e8",
		"voter_contract":      "fb3c8edbfffbd3ff7bb294457f9682ab391e7fb181799ee7a8fb8f264f2e19fc",
	}

	var err error
	suite.daoContractPackageHashes, err = dao_event_parser.NewDAOContractsMetadataFromHashesMap(suite.daoContractHashesMap, suite.casperClient)
	assert.NoError(suite.T(), err)
	suite.entityManager = persistence.NewEntityManager(suite.db, suite.daoContractPackageHashes)
}

func (suite *ProcessDAODeploysTestSuit) SetupTest() {
	_, err := suite.db.Exec(`TRUNCATE TABLE reputation_changes`)
	suite.NoError(err)
}

func (suite *ProcessDAODeploysTestSuit) TestParseDaoEvents() {
	var rawProcessedDeploy casper.GetDeployResult
	err := json.Unmarshal([]byte(kycVoterContractDeployJSON), &rawProcessedDeploy)
	assert.NoError(suite.T(), err)

	deployResult := casper.GetDeployResult{}
	err = json.Unmarshal([]byte(kycVoterContractDeployJSON), &deployResult)
	assert.NoError(suite.T(), err)

	daoEventParser, err := dao_event_parser.NewDaoEventParser(suite.casperClient, suite.daoContractHashesMap, 0)
	assert.NoError(suite.T(), err)

	deployHash, _ := types.NewHashFromHexString("9fccaf372ec5f61ac851fcec593d159f928a26df8f2af5aa3522ed9e0b7cbb36")

	processRawDeploy := event_processing.NewProcessRawDeploy()
	processRawDeploy.SetDAOEventParser(daoEventParser)
	processRawDeploy.SetCasperClient(suite.casperClient)
	processRawDeploy.SetEntityManager(suite.entityManager)
	processRawDeploy.SetDAOContractPackageHashes(suite.daoContractPackageHashes)
	processRawDeploy.SetDeployProcessedEvent(&casper.DeployProcessedEvent{
		DeployProcessed: casper.DeployProcessed{
			ExecutionResult: deployResult.ExecutionResults[0].Result,
			DeployHash:      deployHash,
			Timestamp:       time.Now(),
		},
	})
	err = processRawDeploy.Execute()
	suite.NoError(err)

	var reputationChanges = make([]entities.ReputationChange, 0)
	err = suite.db.Select(&reputationChanges, "SELECT * FROM reputation_changes")
	suite.NoError(err)
	suite.True(len(reputationChanges) > 0, true)
	suite.Equal(reputationChanges[0].Reason, entities.ReputationChangeReasonMint)
	suite.Equal(reputationChanges[0].DeployHash, deployHash)
}

func TestSProcessDAODeploysTestSuit(t *testing.T) {
	suite.Run(t, new(ProcessDAODeploysTestSuit))
}

var kycVoterContractDeployJSON = `
{
  "deploy": {
    "hash": "de63d62894d76084166466939e2a057d7c17eaaa45c94d6287d5eb9b9d63e7c0",
    "header": {
      "ttl": "30m",
      "account": "01b2443959f41927cf973fca48e80453d86d15b519be4df08663ff8a2fc3b43f71",
      "body_hash": "bc8cda3acd50b8d58e18c03b71acf1aac80876172d4c232f069c20377c650c08",
      "gas_price": 1,
      "timestamp": "2022-09-14T19:25:11.408Z",
      "chain_name": "integration-test",
      "dependencies": []
    },
    "payment": {
      "ModuleBytes": {
        "args": [
          [
            "amount",
            {
              "bytes": "0400ca9a3b",
              "parsed": "1000000000",
              "cl_type": "U512"
            }
          ]
        ],
        "module_bytes": ""
      }
    },
    "session": {
      "StoredContractByHash": {
        "args": [
          [
            "recipient",
            {
              "bytes": "003b4ffcfb21411ced5fc1560c3f6ffed86f4885e5ea05cde49d90962a48a14d95",
              "parsed": {
                "Account": "account-hash-3b4ffcfb21411ced5fc1560c3f6ffed86f4885e5ea05cde49d90962a48a14d95"
              },
              "cl_type": "Key"
            }
          ],
          [
            "amount",
            {
              "bytes": "0400ca9a3b",
              "parsed": "1000000000",
              "cl_type": "U256"
            }
          ]
        ],
        "hash": "3eaf08d521d01da82151db1f6f5a88c5acea55109b6b84998c8a9073b93b02e8",
        "entry_point": "mint"
      }
    },
    "approvals": [
      {
        "signer": "01b2443959f41927cf973fca48e80453d86d15b519be4df08663ff8a2fc3b43f71",
        "signature": "01e2117105168d652a7c662160675cb043f6ea8a053a60844ce89cbdcf9f94d5f306d6a76c675faf5dd0918800ef239414674e00936313092ddf50abe2d929d903"
      }
    ]
  },
  "api_version": "1.4.11",
  "execution_results": [
    {
      "result": {
        "Success": {
          "cost": "446473770",
          "effect": {
            "operations": [],
            "transforms": [
              {
                "key": "hash-d2dfc9409965993f9e186db762b585274dcafe439fa1321cfca08017262c8e46",
                "transform": "Identity"
              },
              {
                "key": "hash-0a300922655180354a9ee92b808c7b45b08e5b01d9da0bac9a9b3415bcebbf8d",
                "transform": "Identity"
              },
              {
                "key": "hash-f8df015ba26860a7ec8cab4ee99f079325b0bbb9ef0e7810b63d85df39da95fe",
                "transform": "Identity"
              },
              {
                "key": "hash-59c6451dd58463708fa0b122e97114f07fa5f609229c9d67ac9426935416fbeb",
                "transform": "Identity"
              },
              {
                "key": "balance-235e7def62965b17e03eb2068d6cb7035eeb71c57bb7f0562788241f0e91ef7e",
                "transform": "Identity"
              },
              {
                "key": "balance-ea3c9bdcbe57f067a29609d397981b2d0fb39853a0a9f06e444b06404eadcb1a",
                "transform": "Identity"
              },
              {
                "key": "balance-235e7def62965b17e03eb2068d6cb7035eeb71c57bb7f0562788241f0e91ef7e",
                "transform": {
                  "WriteCLValue": {
                    "bytes": "05000e50f781",
                    "parsed": "558200000000",
                    "cl_type": "U512"
                  }
                }
              },
              {
                "key": "balance-ea3c9bdcbe57f067a29609d397981b2d0fb39853a0a9f06e444b06404eadcb1a",
                "transform": {
                  "AddUInt512": "1000000000"
                }
              },
              {
                "key": "hash-3eaf08d521d01da82151db1f6f5a88c5acea55109b6b84998c8a9073b93b02e8",
                "transform": "Identity"
              },
              {
                "key": "hash-4121e052262daeae1daaa745d65527a1212566b91fcf351e7ec735cb0d6a416e",
                "transform": "Identity"
              },
              {
                "key": "hash-add98a905cc09aefcc48934cce3af9a3aed0ad6c7d97f2cd7e4ab86a83c9619b",
                "transform": "Identity"
              },
              {
                "key": "dictionary-53436feb4e73fce2299e4f57969c06d0127a51b273da4725672ac191f61d2552",
                "transform": "Identity"
              },
              {
                "key": "uref-6b7e9c7032211acfa70b82827636e58f34897ab49fca134f51aff979c88ae4e1-000",
                "transform": {
                  "WriteCLValue": {
                    "bytes": "0400ca9a3b",
                    "parsed": "1000000000",
                    "cl_type": "U256"
                  }
                }
              },
              {
                "key": "hash-3eaf08d521d01da82151db1f6f5a88c5acea55109b6b84998c8a9073b93b02e8",
                "transform": {
                  "AddKeys": [
                    {
                      "key": "uref-6b7e9c7032211acfa70b82827636e58f34897ab49fca134f51aff979c88ae4e1-007",
                      "name": "total_supply_token_contract"
                    }
                  ]
                }
              },
              {
                "key": "uref-b0547376181f154df2da85ca451fd866ed9afcf4f2eec887b9a8fdc6be00fb01-000",
                "transform": {
                  "WriteCLValue": {
                    "bytes": "",
                    "parsed": null,
                    "cl_type": "Unit"
                  }
                }
              },
              {
                "key": "hash-3eaf08d521d01da82151db1f6f5a88c5acea55109b6b84998c8a9073b93b02e8",
                "transform": {
                  "AddKeys": [
                    {
                      "key": "uref-b0547376181f154df2da85ca451fd866ed9afcf4f2eec887b9a8fdc6be00fb01-007",
                      "name": "balances_token_contract"
                    }
                  ]
                }
              },
				{
				  "key": "dictionary-d7d98060d372828cb57b58bce26551954ae2a7fab7a123ec325baaaef6b96e79",
				  "transform": {
					"WriteCLValue": {
					  "bytes": "3a0000000135000000100000004164646564546f57686974656c697374010c37ad7d3bfecf08a807122a97a47931fbf14a2b72fe0533abe0ee832277ea7b0d0e032000000006624e86505574d6eeafd021e55f7bb1d2a94a7e95aa7f05cf0df8f02e7101f94000000032366130386534643063353139306630313837316530353639623632393062383637363030383564393966313765623465376536623538666562386436323439",
					  "parsed": null,
					  "cl_type": "Any"
					}
				  }
				},
              {
                "key": "dictionary-4f6c47aa9337c7575b8d59b767b5f127b3304949059031f438df2e9b6dec13a7",
                "transform": {
                  "WriteCLValue": {
                    "bytes": "050000000400ca9a3b0720000000b0547376181f154df2da85ca451fd866ed9afcf4f2eec887b9a8fdc6be00fb014000000065313165663661346333303632663232666463316462633862323130646237343165356632653537333439386233353934326161656637396336373164383262",
                    "parsed": null,
                    "cl_type": "Any"
                  }
                }
              },
              {
                "key": "uref-d1a68e4ae2c8ffe65cafcfc172caf1179bc5fa820d25eb4574a48f89225820a0-000",
                "transform": "Identity"
              },
              {
                "key": "dictionary-cbebf9cd7e1f0ad57cbb166454e66b9bf3c05267779e4e6a91b06e74ffa805f7",
                "transform": {
                  "WriteCLValue": {
                    "bytes": "0500000001030000000d0420000000354a79790db37fb2423976526c327d77c527ec5c60f271119caec84d592f2e394000000031646462303862666362343132656461396366353834643930303263376334363333353835663463666431333165343964336233376662613362326336663766",
                    "parsed": null,
                    "cl_type": "Any"
                  }
                }
              },
              {
                "key": "dictionary-a74290f7c05991181c1548c608da3540078ccc2072de90d5ece50a02f1d0efd1",
                "transform": {
                  "WriteCLValue": {
                    "bytes": "33000000012e000000040000004d696e74003b4ffcfb21411ced5fc1560c3f6ffed86f4885e5ea05cde49d90962a48a14d950400ca9a3b0d0e032000000006624e86505574d6eeafd021e55f7bb1d2a94a7e95aa7f05cf0df8f02e7101f94000000038633033396666376361613137636365626663616463343462643966636536613462363639396334643033646532653333343961613164633131313933636437",
                    "parsed": null,
                    "cl_type": "Any"
                  }
                }
              },
              {
                "key": "uref-d1a68e4ae2c8ffe65cafcfc172caf1179bc5fa820d25eb4574a48f89225820a0-000",
                "transform": {
                  "WriteCLValue": {
                    "bytes": "04000000",
                    "parsed": 4,
                    "cl_type": "U32"
                  }
                }
              },
              {
                "key": "deploy-de63d62894d76084166466939e2a057d7c17eaaa45c94d6287d5eb9b9d63e7c0",
                "transform": {
                  "WriteDeployInfo": {
                    "gas": "446473770",
                    "from": "account-hash-3b4ffcfb21411ced5fc1560c3f6ffed86f4885e5ea05cde49d90962a48a14d95",
                    "source": "uref-235e7def62965b17e03eb2068d6cb7035eeb71c57bb7f0562788241f0e91ef7e-007",
                    "transfers": [],
                    "deploy_hash": "de63d62894d76084166466939e2a057d7c17eaaa45c94d6287d5eb9b9d63e7c0"
                  }
                }
              },
              {
                "key": "hash-d2dfc9409965993f9e186db762b585274dcafe439fa1321cfca08017262c8e46",
                "transform": "Identity"
              },
              {
                "key": "hash-0a300922655180354a9ee92b808c7b45b08e5b01d9da0bac9a9b3415bcebbf8d",
                "transform": "Identity"
              },
              {
                "key": "balance-ea3c9bdcbe57f067a29609d397981b2d0fb39853a0a9f06e444b06404eadcb1a",
                "transform": "Identity"
              },
              {
                "key": "hash-d2dfc9409965993f9e186db762b585274dcafe439fa1321cfca08017262c8e46",
                "transform": "Identity"
              },
              {
                "key": "hash-f8df015ba26860a7ec8cab4ee99f079325b0bbb9ef0e7810b63d85df39da95fe",
                "transform": "Identity"
              },
              {
                "key": "hash-59c6451dd58463708fa0b122e97114f07fa5f609229c9d67ac9426935416fbeb",
                "transform": "Identity"
              },
              {
                "key": "balance-ea3c9bdcbe57f067a29609d397981b2d0fb39853a0a9f06e444b06404eadcb1a",
                "transform": "Identity"
              },
              {
                "key": "balance-23a86caef8f8726deb889e75dae4d6e610c6511551391cb9aeabda2fb6e08f9e",
                "transform": "Identity"
              },
              {
                "key": "balance-ea3c9bdcbe57f067a29609d397981b2d0fb39853a0a9f06e444b06404eadcb1a",
                "transform": {
                  "WriteCLValue": {
                    "bytes": "00",
                    "parsed": "0",
                    "cl_type": "U512"
                  }
                }
              },
              {
                "key": "balance-23a86caef8f8726deb889e75dae4d6e610c6511551391cb9aeabda2fb6e08f9e",
                "transform": {
                  "AddUInt512": "1000000000"
                }
              }
            ]
          },
          "transfers": []
        }
      },
      "block_hash": "7c3349116bce939a7baae62ab47ccda1b70f36ef0de99aaed8bf5b35d1b8a9dc"
    }
  ]
}
`
