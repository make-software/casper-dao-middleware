//go:build integration
// +build integration

package event_tracking

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"casper-dao-middleware/internal/crdao/dao_event_parser/events"

	"casper-dao-middleware/internal/crdao/persistence"
	"casper-dao-middleware/internal/crdao/services/event_tracking"
	"casper-dao-middleware/pkg/boot"
	"casper-dao-middleware/pkg/casper"
	"casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/pagination"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TrackVotingCreatedTestSuit struct {
	suite.Suite

	db            *sqlx.DB
	casperClient  casper.RPCClient
	entityManager persistence.EntityManager

	daoContractHashesMap     map[string]string
	daoContractPackageHashes dao.DAOContractsMetadata
}

func (suite *TrackVotingCreatedTestSuit) SetupSuite() {
	suite.db = boot.SetUpTestDB()
	suite.casperClient = casper.NewRPCClient(os.Getenv("INTEGRATION_NET_RPC_API_URL"))
	suite.daoContractHashesMap = map[string]string{
		"reputation_contract": "3eaf08d521d01da82151db1f6f5a88c5acea55109b6b84998c8a9073b93b02e8",
		"voter_contract":      "fb3c8edbfffbd3ff7bb294457f9682ab391e7fb181799ee7a8fb8f264f2e19fc",
	}

	var err error
	suite.daoContractPackageHashes, err = dao.NewDAOContractsMetadataFromHashesMap(suite.daoContractHashesMap, suite.casperClient)
	assert.NoError(suite.T(), err)
	suite.entityManager = persistence.NewEntityManager(suite.db, suite.daoContractPackageHashes)
}

func (suite *TrackVotingCreatedTestSuit) SetupTest() {
	_, err := suite.db.Exec(`TRUNCATE TABLE votings`)
	suite.NoError(err)
}

func (suite *TrackVotingCreatedTestSuit) TestTrackVotingCreated() {
	deployResult := casper.GetDeployResult{}
	err := json.Unmarshal([]byte(createVotingCreatedDeployJSON), &deployResult)
	assert.NoError(suite.T(), err)

	daoEventParser, err := dao.NewDaoEventParser(suite.casperClient, suite.daoContractHashesMap, 0)
	assert.NoError(suite.T(), err)

	deployHash, _ := types.NewHashFromHexString("661a754588fd07ca050b15f7f28f5189449b974f476e678f1ca913104b5644d0")

	daoEvents, err := daoEventParser.Parse(&casper.DeployProcessedEvent{
		DeployProcessed: casper.DeployProcessed{
			ExecutionResult: deployResult.ExecutionResults[0].Result,
			DeployHash:      deployHash,
			Timestamp:       time.Now(),
		},
	})
	assert.NoError(suite.T(), err)

	trackVotingCreated := event_tracking.NewTrackVotingCreated()
	trackVotingCreated.SetEntityManager(suite.entityManager)
	trackVotingCreated.SetDeployProcessed(casper.DeployProcessed{
		ExecutionResult: deployResult.ExecutionResults[0].Result,
		DeployHash:      deployHash,
		Timestamp:       time.Now(),
	})

	for _, event := range daoEvents {
		if event.EventName == events.VotingCreatedEventName {
			trackVotingCreated.SetEventBody(event.EventBody)
		}
	}

	err = trackVotingCreated.Execute()
	assert.NoError(suite.T(), err)

	votings, err := suite.entityManager.VotingRepository().Find(&pagination.Params{
		OrderDirection: pagination.OrderDirectionDESC,
		Page:           1,
		PageSize:       1,
	}, nil)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), len(votings) == 1)
	assert.Equal(suite.T(), votings[0].HasEnded, false)
	assert.Equal(suite.T(), votings[0].VotingTime, uint64(86400000))
	assert.Equal(suite.T(), votings[0].VotingQuorum, uint64(0))

	err = suite.entityManager.VotingRepository().UpdateHasEnded(1, true)
	assert.NoError(suite.T(), err)

	votings, err = suite.entityManager.VotingRepository().Find(&pagination.Params{
		OrderDirection: pagination.OrderDirectionDESC,
		Page:           1,
		PageSize:       1,
	}, nil)
	assert.Equal(suite.T(), votings[0].HasEnded, true)
}

func TestTrackVotingCreatedTestSuit(t *testing.T) {
	suite.Run(t, new(TrackVotingCreatedTestSuit))
}

var createVotingCreatedDeployJSON = `
{
  "deploy": {
    "hash": "661a754588fd07ca050b15f7f28f5189449b974f476e678f1ca913104b5644d0",
    "header": {
      "ttl": "30m",
      "account": "01b2443959f41927cf973fca48e80453d86d15b519be4df08663ff8a2fc3b43f71",
      "body_hash": "012a7e70b67496a7dcb6f24464701b6b7aeadacf8da85660b15dd0574abf3cbd",
      "gas_price": 1,
      "timestamp": "2022-09-15T13:18:56.779Z",
      "chain_name": "integration-test",
      "dependencies": []
    },
    "payment": {
      "ModuleBytes": {
        "args": [
          [
            "amount",
            {
              "bytes": "0500f2052a01",
              "parsed": "5000000000",
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
            "subject_address",
            {
              "bytes": "006d87e1a98e9122460573b8bc6a4cf93c0fd2736b51d388ab28155f881e5d3c81",
              "parsed": {
                "Account": "account-hash-6d87e1a98e9122460573b8bc6a4cf93c0fd2736b51d388ab28155f881e5d3c81"
              },
              "cl_type": "Key"
            }
          ],
          [
            "document_hash",
            {
              "bytes": "012b",
              "parsed": "43",
              "cl_type": "U256"
            }
          ],
          [
            "stake",
            {
              "bytes": "02e803",
              "parsed": "1000",
              "cl_type": "U256"
            }
          ]
        ],
        "hash": "fb3c8edbfffbd3ff7bb294457f9682ab391e7fb181799ee7a8fb8f264f2e19fc",
        "entry_point": "create_voting"
      }
    },
    "approvals": [
      {
        "signer": "01b2443959f41927cf973fca48e80453d86d15b519be4df08663ff8a2fc3b43f71",
        "signature": "018349553b9ad2e53215f027782ea5d837a12ce31f40c965d43b4c34bd6ddd168b5133266b8a857c35018f58178db95c9cedfe6ea091fc63485f7d435bd800ad07"
      }
    ]
  },
  "api_version": "1.4.11",
  "execution_results": [
    {
      "result": {
        "Success": {
          "cost": "1971251650",
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
                    "bytes": "050074c02de8",
                    "parsed": "997200000000",
                    "cl_type": "U512"
                  }
                }
              },
              {
                "key": "balance-ea3c9bdcbe57f067a29609d397981b2d0fb39853a0a9f06e444b06404eadcb1a",
                "transform": {
                  "AddUInt512": "5000000000"
                }
              },
              {
                "key": "hash-fb3c8edbfffbd3ff7bb294457f9682ab391e7fb181799ee7a8fb8f264f2e19fc",
                "transform": "Identity"
              },
              {
                "key": "hash-0c37ad7d3bfecf08a807122a97a47931fbf14a2b72fe0533abe0ee832277ea7b",
                "transform": "Identity"
              },
              {
                "key": "hash-f4c18623893c09a4209e1c4284178920d2215b87fcb88000cb2e0fba95e92ba3",
                "transform": "Identity"
              },
              {
                "key": "uref-49f20de29269a93db034343a8ee2e6260d77641b7c55322eed1c5fca47210439-000",
                "transform": "Identity"
              },
              {
                "key": "hash-926745ae3b44b7a28ba3b8df5c84a2afd86dcfd16826befca5165e87a8c34087",
                "transform": "Identity"
              },
              {
                "key": "hash-3e7c9f8d67a5729774e663f7740a2d07e6160f5574e0b4124db8fb8ca1fe6632",
                "transform": "Identity"
              },
              {
                "key": "hash-4c40c81253583774c80700272df574cd1d27d39a3d4a147f1b94e7451adc856e",
                "transform": "Identity"
              },
              {
                "key": "uref-375631180ef09cab8ba21b7c2edd93d3ccc61ee44ed93d32ddd65704b84bc07c-000",
                "transform": "Identity"
              },
              {
                "key": "hash-b02d1209e47b06ebaa3b4cdb15f95f807e25da4d202cc1c35e083dc32ba0f9b8",
                "transform": "Identity"
              },
              {
                "key": "hash-a717e476603b2528b39a30914204470a804933bd8950311509eab49ab90c0cee",
                "transform": "Identity"
              },
              {
                "key": "hash-b4f1663f9b5e6beaeb44af89032c27857345c6f23376b7001269d0d95e61d449",
                "transform": "Identity"
              },
              {
                "key": "uref-ee25a0b359f242973c4d25a3dbd632de1f6659cbeae7e14873253f5042e514d3-000",
                "transform": "Identity"
              },
              {
                "key": "uref-010c51bf07fcfc32398bd9f1034438dfd32b6c0412fd61bb0e3f8e0b9a45104e-000",
                "transform": "Identity"
              },
              {
                "key": "hash-b877ae42aca04d5f1249781930d867818a0a4db2b1c37834e01eedd0f9fdf78b",
                "transform": "Identity"
              },
              {
                "key": "hash-0ffbff03f3d8fcbf5aee2adc81fb6682631f5df507e6c124f6599961bed87a3b",
                "transform": "Identity"
              },
              {
                "key": "hash-b969dd5c403cbf580347ba4affbff5a48aeda60f9d97fffb0fe11aef63b9cafc",
                "transform": "Identity"
              },
              {
                "key": "dictionary-b17a3ea9f97552fa1173884e85160da3ed0814b8a1b3131cb40e790299c0b6ef",
                "transform": "Identity"
              },
              {
                "key": "hash-b877ae42aca04d5f1249781930d867818a0a4db2b1c37834e01eedd0f9fdf78b",
                "transform": "Identity"
              },
              {
                "key": "hash-0ffbff03f3d8fcbf5aee2adc81fb6682631f5df507e6c124f6599961bed87a3b",
                "transform": "Identity"
              },
              {
                "key": "hash-b969dd5c403cbf580347ba4affbff5a48aeda60f9d97fffb0fe11aef63b9cafc",
                "transform": "Identity"
              },
              {
                "key": "dictionary-376da7fed16ee3f48e0f0f3c6411780909ee5482ca767758fcc07af8cb66f8da",
                "transform": "Identity"
              },
              {
                "key": "hash-b877ae42aca04d5f1249781930d867818a0a4db2b1c37834e01eedd0f9fdf78b",
                "transform": "Identity"
              },
              {
                "key": "hash-0ffbff03f3d8fcbf5aee2adc81fb6682631f5df507e6c124f6599961bed87a3b",
                "transform": "Identity"
              },
              {
                "key": "hash-b969dd5c403cbf580347ba4affbff5a48aeda60f9d97fffb0fe11aef63b9cafc",
                "transform": "Identity"
              },
              {
                "key": "dictionary-cf7a7f9911bc9ffcee72c7067d0eb2aa2f8773a13b3d7008be84bfc306e26c47",
                "transform": "Identity"
              },
              {
                "key": "hash-b877ae42aca04d5f1249781930d867818a0a4db2b1c37834e01eedd0f9fdf78b",
                "transform": "Identity"
              },
              {
                "key": "hash-0ffbff03f3d8fcbf5aee2adc81fb6682631f5df507e6c124f6599961bed87a3b",
                "transform": "Identity"
              },
              {
                "key": "hash-b969dd5c403cbf580347ba4affbff5a48aeda60f9d97fffb0fe11aef63b9cafc",
                "transform": "Identity"
              },
              {
                "key": "dictionary-6a2714d066e5aa02096b724576b6e61e693b55e5c0f5bdac4416fbc6070d43e6",
                "transform": "Identity"
              },
              {
                "key": "hash-b877ae42aca04d5f1249781930d867818a0a4db2b1c37834e01eedd0f9fdf78b",
                "transform": "Identity"
              },
              {
                "key": "hash-0ffbff03f3d8fcbf5aee2adc81fb6682631f5df507e6c124f6599961bed87a3b",
                "transform": "Identity"
              },
              {
                "key": "hash-b969dd5c403cbf580347ba4affbff5a48aeda60f9d97fffb0fe11aef63b9cafc",
                "transform": "Identity"
              },
              {
                "key": "dictionary-1f52123a6eaa9bef2a9a058c5b824218508aa8e9ec69c7a58e2f24e7e7f4a2f8",
                "transform": "Identity"
              },
              {
                "key": "uref-49f20de29269a93db034343a8ee2e6260d77641b7c55322eed1c5fca47210439-000",
                "transform": "Identity"
              },
              {
                "key": "uref-19d2856ef2da5c86c526d06bbf1a31713ea78b8edd115767ba415a98f8bed0a7-000",
                "transform": "Identity"
              },
              {
                "key": "uref-19d2856ef2da5c86c526d06bbf1a31713ea78b8edd115767ba415a98f8bed0a7-000",
                "transform": {
                  "WriteCLValue": {
                    "bytes": "0102",
                    "parsed": "2",
                    "cl_type": "U256"
                  }
                }
              },
              {
                "key": "uref-a83532e7ecef953619a8d3aa7cf8745713a1fe25a5b170924beebb2cece74d10-000",
                "transform": "Identity"
              },
              {
                "key": "dictionary-0eb9bfe806a4896c6ecfdea288b75685585a44555eb7a8305c448555793c53ec",
                "transform": {
                  "WriteCLValue": {
                    "bytes": "0500000001030000000d042000000007027ae55db9184f5f871373a0ec014c09b0b9bbf7357c9e450e8edb6dcc381d4000000033386239316166356465343761336534323332626561303263616330353532653233393231616137656132663366323263336266633466343433613661313034",
                    "parsed": null,
                    "cl_type": "Any"
                  }
                }
              },
              {
                "key": "dictionary-1c996e3ec04c6f3b5588e720d4278a676a215173e3a99a99b083831c167bd7d9",
                "transform": {
                  "WriteCLValue": {
                    "bytes": "50000000014b0000000d000000566f74696e6743726561746564003b4ffcfb21411ced5fc1560c3f6ffed86f4885e5ea05cde49d90962a48a14d9501010101000000ccbf190000000000005c26050000000001640d0e0320000000c7883cef3fa34729510ff09fb822bc39864b21b3a97f3d9815e7ed3e4341ccd94000000038633033396666376361613137636365626663616463343462643966636536613462363639396334643033646532653333343961613164633131313933636437",
                    "parsed": null,
                    "cl_type": "Any"
                  }
                }
              },
              {
                "key": "uref-a83532e7ecef953619a8d3aa7cf8745713a1fe25a5b170924beebb2cece74d10-000",
                "transform": {
                  "WriteCLValue": {
                    "bytes": "04000000",
                    "parsed": 4,
                    "cl_type": "U32"
                  }
                }
              },
              {
                "key": "dictionary-e3f7d6dda1c80b37936aab26406c7dcb497650eac27b09004a408faf34a27efe",
                "transform": {
                  "WriteCLValue": {
                    "bytes": "9400000001010100000000004e41830100000101000000ccbf190000000000005c260500000000010164000101926745ae3b44b7a28ba3b8df5c84a2afd86dcfd16826befca5165e87a8c34087040000006d696e740200000002000000746f21000000006d87e1a98e9122460573b8bc6a4cf93c0fd2736b51d388ab28155f881e5d3c810b08000000746f6b656e5f696402000000012b070d1520000000be75a303e1e65056971bfbca98c4a9d9f971978daa3f796bc98aeac3f44629554000000062313338636662333535336239623062616461643533616363346630623531626564643132363163396332616339633932633337323735353035353666356161",
                    "parsed": null,
                    "cl_type": "Any"
                  }
                }
              },
              {
                "key": "dictionary-e3f7d6dda1c80b37936aab26406c7dcb497650eac27b09004a408faf34a27efe",
                "transform": "Identity"
              },
              {
                "key": "uref-18df0b8d2563eb7fd0fd036636bfc0d4b83eb0e35cbdd4bfbc57d5d2a8e89055-000",
                "transform": "Identity"
              },
              {
                "key": "hash-4121e052262daeae1daaa745d65527a1212566b91fcf351e7ec735cb0d6a416e",
                "transform": "Identity"
              },
              {
                "key": "hash-3eaf08d521d01da82151db1f6f5a88c5acea55109b6b84998c8a9073b93b02e8",
                "transform": "Identity"
              },
              {
                "key": "hash-add98a905cc09aefcc48934cce3af9a3aed0ad6c7d97f2cd7e4ab86a83c9619b",
                "transform": "Identity"
              },
              {
                "key": "dictionary-e3b19218f8e6023ad63a00a7474ad3332b010ad4db905def5d9a54343153d19a",
                "transform": "Identity"
              },
              {
                "key": "dictionary-4f6c47aa9337c7575b8d59b767b5f127b3304949059031f438df2e9b6dec13a7",
                "transform": "Identity"
              },
              {
                "key": "dictionary-7c4cf88fb210d0ba75c22bbf921847460d43acc58302408c9de3a896b7aa8d2d",
                "transform": "Identity"
              },
              {
                "key": "dictionary-4f6c47aa9337c7575b8d59b767b5f127b3304949059031f438df2e9b6dec13a7",
                "transform": {
                  "WriteCLValue": {
                    "bytes": "050000000430c29a3b0720000000b0547376181f154df2da85ca451fd866ed9afcf4f2eec887b9a8fdc6be00fb014000000065313165663661346333303632663232666463316462633862323130646237343165356632653537333439386233353934326161656637396336373164383262",
                    "parsed": null,
                    "cl_type": "Any"
                  }
                }
              },
              {
                "key": "dictionary-7c4cf88fb210d0ba75c22bbf921847460d43acc58302408c9de3a896b7aa8d2d",
                "transform": {
                  "WriteCLValue": {
                    "bytes": "0300000002d0070720000000b0547376181f154df2da85ca451fd866ed9afcf4f2eec887b9a8fdc6be00fb014000000032623030313465623565386330636364353236393538366161346330323866613165333835333766313034323639663730326162383038313633366632393038",
                    "parsed": null,
                    "cl_type": "Any"
                  }
                }
              },
              {
                "key": "uref-a83532e7ecef953619a8d3aa7cf8745713a1fe25a5b170924beebb2cece74d10-000",
                "transform": "Identity"
              },
              {
                "key": "dictionary-6ff158ecd61337353b9b7d6263dadf2390a017faf55a49113e38598195b5e5b5",
                "transform": {
                  "WriteCLValue": {
                    "bytes": "0500000001040000000d042000000007027ae55db9184f5f871373a0ec014c09b0b9bbf7357c9e450e8edb6dcc381d4000000031633030326331636362366564323666386239636262353735386333363837626361303632666264613735646366396161633265333464633132626237646237",
                    "parsed": null,
                    "cl_type": "Any"
                  }
                }
              },
              {
                "key": "dictionary-577a681a2d0a28e11c69699c5feb63a959bce16aa8ac418e39c0b8cdfdb726c7",
                "transform": {
                  "WriteCLValue": {
                    "bytes": "3d00000001380000000a00000042616c6c6f7443617374003b4ffcfb21411ced5fc1560c3f6ffed86f4885e5ea05cde49d90962a48a14d9501010200000002e8030d0e0320000000c7883cef3fa34729510ff09fb822bc39864b21b3a97f3d9815e7ed3e4341ccd94000000032366130386534643063353139306630313837316530353639623632393062383637363030383564393966313765623465376536623538666562386436323439",
                    "parsed": null,
                    "cl_type": "Any"
                  }
                }
              },
              {
                "key": "uref-a83532e7ecef953619a8d3aa7cf8745713a1fe25a5b170924beebb2cece74d10-000",
                "transform": {
                  "WriteCLValue": {
                    "bytes": "05000000",
                    "parsed": 5,
                    "cl_type": "U32"
                  }
                }
              },
              {
                "key": "dictionary-4464357f500589cf4fda3bb69d5bfac6e0c626cf7dcd3d28faa4e8ce166180f5",
                "transform": {
                  "WriteCLValue": {
                    "bytes": "0400000001000000042000000085c1bc79208097e252334aa50deec28389c868a20aaf5ef5bc6c78ec6a56ee044000000062313338636662333535336239623062616461643533616363346630623531626564643132363163396332616339633932633337323735353035353666356161",
                    "parsed": null,
                    "cl_type": "Any"
                  }
                }
              },
              {
                "key": "dictionary-cc03ce9bbbc8e36a4a9a6a3fd848b6fdfeece8a2a6edad714c11f3e85fee2706",
                "transform": {
                  "WriteCLValue": {
                    "bytes": "21000000003b4ffcfb21411ced5fc1560c3f6ffed86f4885e5ea05cde49d90962a48a14d950b20000000533af9f1465a0dafff9f5845e74dfa60407dbe33d1ea44ed483c07041a8a4bd74000000031323734383334393039613030663138336636316537386562336666633261343161346431346433326332353030636430363665316331386538343836616561",
                    "parsed": null,
                    "cl_type": "Any"
                  }
                }
              },
              {
                "key": "dictionary-cc4c9a63bc49bffb78895e7f56e59ece7f77cc1a50720fc858f46f4e8882d7d2",
                "transform": {
                  "WriteCLValue": {
                    "bytes": "2a000000003b4ffcfb21411ced5fc1560c3f6ffed86f4885e5ea05cde49d90962a48a14d9501010200000002e803152000000075b6877d8a0794893e44db14366fc604037f0026ebd2cf91e51be5387f4832f54000000066326562323664653564306339666164383939666166626430373437343235663661616234373365353864326461393738623162643662616338313866633864",
                    "parsed": null,
                    "cl_type": "Any"
                  }
                }
              },
              {
                "key": "dictionary-e3f7d6dda1c80b37936aab26406c7dcb497650eac27b09004a408faf34a27efe",
                "transform": {
                  "WriteCLValue": {
                    "bytes": "960000000101010002e8030000004e41830100000101000000ccbf190000000000005c260500000000010164000101926745ae3b44b7a28ba3b8df5c84a2afd86dcfd16826befca5165e87a8c34087040000006d696e740200000002000000746f21000000006d87e1a98e9122460573b8bc6a4cf93c0fd2736b51d388ab28155f881e5d3c810b08000000746f6b656e5f696402000000012b070d1520000000be75a303e1e65056971bfbca98c4a9d9f971978daa3f796bc98aeac3f44629554000000062313338636662333535336239623062616461643533616363346630623531626564643132363163396332616339633932633337323735353035353666356161",
                    "parsed": null,
                    "cl_type": "Any"
                  }
                }
              },
              {
                "key": "dictionary-8499a97fca600d506033a47048e04b36a4d53b9ada520fd4fe9acb995f7f4e46",
                "transform": {
                  "WriteCLValue": {
                    "bytes": "0100000001002000000008ac6bae64f5aca3f80e5798b09ec5e61e3720ecca192683952f43b785bfcef44000000062633334336433643866646336626263643536303836336230613636323565633964666439313533356366623764656561363432306339656665353235633466",
                    "parsed": null,
                    "cl_type": "Any"
                  }
                }
              },
              {
                "key": "deploy-661a754588fd07ca050b15f7f28f5189449b974f476e678f1ca913104b5644d0",
                "transform": {
                  "WriteDeployInfo": {
                    "gas": "1971251650",
                    "from": "account-hash-3b4ffcfb21411ced5fc1560c3f6ffed86f4885e5ea05cde49d90962a48a14d95",
                    "source": "uref-235e7def62965b17e03eb2068d6cb7035eeb71c57bb7f0562788241f0e91ef7e-007",
                    "transfers": [],
                    "deploy_hash": "661a754588fd07ca050b15f7f28f5189449b974f476e678f1ca913104b5644d0"
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
                "key": "balance-3200886de7b5c8b487d731ebab4075ec50757274047bd51ed7d84942f17803fb",
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
                "key": "balance-3200886de7b5c8b487d731ebab4075ec50757274047bd51ed7d84942f17803fb",
                "transform": {
                  "AddUInt512": "5000000000"
                }
              }
            ]
          },
          "transfers": []
        }
      },
      "block_hash": "b8b71d417138c98d87591604a692933a18f2f3d3c78faa227ae8b3dec04d1d3e"
    }
  ]
}
`
