package event_processing

import (
	"context"
	"time"

	"go.uber.org/zap"

	"casper-dao-middleware/pkg/casper"
	"casper-dao-middleware/pkg/casper/types"
	"casper-dao-middleware/pkg/go-ces-parser"

	"casper-dao-middleware/internal/crdao/di"
)

var test = `
 {
    "result": {
      "Success": {
        "cost": "8227779160",
        "effect": {
          "operations": [],
          "transforms": [
            {
              "key": "hash-8cf5e4acf51f54eb59291599187838dc3bc234089c46fc6ca8ad17e762ae4401",
              "transform": "Identity"
            },
            {
              "key": "hash-624dbe2395b9d9503fbee82162f1714ebff6b639f96d2084d26d944c354ec4c5",
              "transform": "Identity"
            },
            {
              "key": "hash-010c3fe81b7b862e50c77ef9a958a05bfa98444f26f96f23d37a13c96244cfb7",
              "transform": "Identity"
            },
            {
              "key": "hash-9824d60dc3a5c44a20b9fd260a412437933835b52fc683d8ae36e4ec2114843e",
              "transform": "Identity"
            },
            {
              "key": "balance-23dcda889fb11ec282ad9b409496fef0a2a7dbf5c70e822a7da84aef87f2a590",
              "transform": "Identity"
            },
            {
              "key": "balance-98d945f5324f865243b7c02c0417ab6eac361c5c56602fd42ced834a1ba201b6",
              "transform": "Identity"
            },
            {
              "key": "balance-23dcda889fb11ec282ad9b409496fef0a2a7dbf5c70e822a7da84aef87f2a590",
              "transform": {
                "WriteCLValue": {
                  "bytes": "0500863ba101",
                  "parsed": "7000000000",
                  "cl_type": "U512"
                }
              }
            },
            {
              "key": "balance-98d945f5324f865243b7c02c0417ab6eac361c5c56602fd42ced834a1ba201b6",
              "transform": {
                "AddUInt512": "20000000000"
              }
            },
            {
              "key": "hash-214a0e730e14501d1e3e03504d3a2f940ef32830b13fa47f9d85a40f73b78161",
              "transform": "Identity"
            },
            {
              "key": "hash-690072c9333d9c7a830695a9789154be436ab1a37a303c250d375687fd34890a",
              "transform": "Identity"
            },
            {
              "key": "hash-988326d1b3a99ad6308cd6a0aa09c785ed59222f9d00f6215a0027ceda15a4a8",
              "transform": "Identity"
            },
            {
              "key": "uref-9ccf91c595fc83f2b5cbe081b1a9dc17c8dd5702b565c6c04c202aac05d86499-000",
              "transform": "Identity"
            },
            {
              "key": "uref-4a953092081bbfc209e7315c305e5b32b4ab91dbd35a2d0df4e711d5a931a0c9-000",
              "transform": "Identity"
            },
            {
              "key": "uref-dd1e6c2008e2a032f1537f3ac2c6d9102fd346bd6a8b2b3d83f264e083cf1a4c-000",
              "transform": "Identity"
            },
            {
              "key": "uref-051909875dccb825436c60373476a6b58eed9353a81a78942270d4f8a13bad0b-000",
              "transform": "Identity"
            },
            {
              "key": "uref-c1929a590834162de4d0b67f8d9eed66fcaef7e50fed51df9aa148e524ad5095-000",
              "transform": "Identity"
            },
            {
              "key": "uref-0c6466dcd50b5b1595d9f4ed73663eb63748256c87318812b78011f8dee4d33f-000",
              "transform": "Identity"
            },
            {
              "key": "uref-c1929a590834162de4d0b67f8d9eed66fcaef7e50fed51df9aa148e524ad5095-000",
              "transform": "Identity"
            },
            {
              "key": "uref-66d81c4b7c30cf91a66be5595419b87523a009cefeb819000ea5fb2725407877-000",
              "transform": "Identity"
            },
            {
              "key": "uref-c0db785ca700d4f97a4031ad996c493dafa141b866d9d99509af5c404b7358de-000",
              "transform": "Identity"
            },
            {
              "key": "dictionary-6ac9df7657ad9a03fe9596adfe69323958d4506bde2f4ddc35a0908c406dc3c9",
              "transform": {
                "WriteCLValue": {
                  "bytes": "210000000056befc13a6fd62e18f361700a5e08f966901c34df8041b36ec97d54d605c23de0b20000000e6be0116a38878e6fe7289cc5d554a6826280b97015590c91fbc0ee5c9a05d620100000030",
                  "parsed": null,
                  "cl_type": "Any"
                }
              }
            },
            {
              "key": "dictionary-f8ccd1ba4cdb1fd3cade278226ebdbff5a90f572033845fa9f449d08044839e0",
              "transform": {
                "WriteCLValue": {
                  "bytes": "210000000056befc13a6fd62e18f361700a5e08f966901c34df8041b36ec97d54d605c23de0b200000006afb210443e3fecf8fcc277ec582b732b8017d8ba8633f67d22788308058b0c40100000030",
                  "parsed": null,
                  "cl_type": "Any"
                }
              }
            },
            {
              "key": "dictionary-14c6c21bdbe86790282af284ffbddf3bf9c5c4736579b81eb5e887546cbfd087",
              "transform": {
                "WriteCLValue": {
                  "bytes": "af000000ab0000007b0a2020226e616d65223a20224153434949204172742069732046756e222c0a202022746f6b656e5f757269223a202268747470733a2f2f6d61726974696d652e7365616c73746f726167652e696f2f6d616b652f363439376335396631626366626262632f6265346562633537353139303431313839646366333433313730323862643337222c0a202022636865636b73756d223a20224153434949204172742069732046756e220a7d0a20000000850c5eab27b15cb83a3b208e94bf33e2fe0212bac13640bb420c04e55e2dcc520100000030",
                  "parsed": null,
                  "cl_type": "Any"
                }
              }
            },
            {
              "key": "dictionary-d06ff0c0da4fd64e9be95cf878346d0132f73409a8b15459488d86d458464b1b",
              "transform": {
                "WriteCLValue": {
                  "bytes": "0800000001000000000000000520000000c34ef1ccc4fe5757a5680ac685bc1aa0f06dd09d4f7c1c1025b02bac386568944000000035366265666331336136666436326531386633363137303061356530386639363639303163333464663830343162333665633937643534643630356332336465",
                  "parsed": null,
                  "cl_type": "Any"
                }
              }
            },
            {
              "key": "uref-dd1e6c2008e2a032f1537f3ac2c6d9102fd346bd6a8b2b3d83f264e083cf1a4c-000",
              "transform": {
                "WriteCLValue": {
                  "bytes": "0100000000000000",
                  "parsed": 1,
                  "cl_type": "U64"
                }
              }
            },
            {
              "key": "uref-e97f4c95c44700c38a6828a1aad038fcde95eba1e686026622f1d9bd8301b35b-000",
              "transform": "Identity"
            },
            {
              "key": "dictionary-cc7122247eae1ff47a4a1f4c31742070d3aa842225b72735155d7831c6e64296",
              "transform": {
                "WriteCLValue": {
                  "bytes": "eb000000e70000000a0000006576656e745f4d696e740056befc13a6fd62e18f361700a5e08f966901c34df8041b36ec97d54d605c23de000000000000000000ab0000007b0a2020226e616d65223a20224153434949204172742069732046756e222c0a202022746f6b656e5f757269223a202268747470733a2f2f6d61726974696d652e7365616c73746f726167652e696f2f6d616b652f363439376335396631626366626262632f6265346562633537353139303431313839646366333433313730323862643337222c0a202022636865636b73756d223a20224153434949204172742069732046756e220a7d0e0320000000c21ab84f222b33e354e933fe406446acfa74eea626cf1d40ddb3453a42d561450100000030",
                  "parsed": null,
                  "cl_type": "Any"
                }
              }
            },
            {
              "key": "uref-e97f4c95c44700c38a6828a1aad038fcde95eba1e686026622f1d9bd8301b35b-000",
              "transform": {
                "WriteCLValue": {
                  "bytes": "01000000",
                  "parsed": 1,
                  "cl_type": "U32"
                }
              }
            },
            {
              "key": "uref-75716eb3ccabe5941ee8836f4112a0cb8561959cf2a83a80a510c4c2c069d5c1-000",
              "transform": "Identity"
            },
            {
              "key": "deploy-471f2d7eba2627a486a25cff964f730b7ea747590937519f9fe657be37b045c5",
              "transform": {
                "WriteDeployInfo": {
                  "gas": "8227779160",
                  "from": "account-hash-56befc13a6fd62e18f361700a5e08f966901c34df8041b36ec97d54d605c23de",
                  "source": "uref-23dcda889fb11ec282ad9b409496fef0a2a7dbf5c70e822a7da84aef87f2a590-007",
                  "transfers": [],
                  "deploy_hash": "471f2d7eba2627a486a25cff964f730b7ea747590937519f9fe657be37b045c5"
                }
              }
            },
            {
              "key": "hash-8cf5e4acf51f54eb59291599187838dc3bc234089c46fc6ca8ad17e762ae4401",
              "transform": "Identity"
            },
            {
              "key": "hash-624dbe2395b9d9503fbee82162f1714ebff6b639f96d2084d26d944c354ec4c5",
              "transform": "Identity"
            },
            {
              "key": "balance-98d945f5324f865243b7c02c0417ab6eac361c5c56602fd42ced834a1ba201b6",
              "transform": "Identity"
            },
            {
              "key": "hash-8cf5e4acf51f54eb59291599187838dc3bc234089c46fc6ca8ad17e762ae4401",
              "transform": "Identity"
            },
            {
              "key": "hash-010c3fe81b7b862e50c77ef9a958a05bfa98444f26f96f23d37a13c96244cfb7",
              "transform": "Identity"
            },
            {
              "key": "hash-9824d60dc3a5c44a20b9fd260a412437933835b52fc683d8ae36e4ec2114843e",
              "transform": "Identity"
            },
            {
              "key": "balance-98d945f5324f865243b7c02c0417ab6eac361c5c56602fd42ced834a1ba201b6",
              "transform": "Identity"
            },
            {
              "key": "balance-874289dbe721508e8d2893efd86364ea1ca67a6a2456825259efd6db8efb427c",
              "transform": "Identity"
            },
            {
              "key": "balance-98d945f5324f865243b7c02c0417ab6eac361c5c56602fd42ced834a1ba201b6",
              "transform": {
                "WriteCLValue": {
                  "bytes": "00",
                  "parsed": "0",
                  "cl_type": "U512"
                }
              }
            },
            {
              "key": "balance-874289dbe721508e8d2893efd86364ea1ca67a6a2456825259efd6db8efb427c",
              "transform": {
                "AddUInt512": "20000000000"
              }
            }
          ]
        },
        "transfers": []
      }
    },
    "block_hash": "9174106b233a7c1999c0b5c1565de86322f725e300bea4323cc16baf8703411f"
  }
`

// ProcessEventStream command start number of concurrent worker to process event from synchronous event stream
type ProcessEventStream struct {
	di.BaseStreamURLAware
	di.CasperClientAware
	di.EntityManagerAware
	di.DAOContractsMetadataAware

	daoContractHashes         map[string]types.Hash
	eventStreamPath           string
	nodeStartFromEventID      uint64
	dictionarySetEventsBuffer uint32
}

func NewProcessEventStream() *ProcessEventStream {
	return &ProcessEventStream{}
}

func (c *ProcessEventStream) SetNodeStartFromEventID(eventID uint64) *ProcessEventStream {
	c.nodeStartFromEventID = eventID
	return c
}

func (c *ProcessEventStream) SetDAOContractHashes(daoContractHashes map[string]types.Hash) *ProcessEventStream {
	c.daoContractHashes = daoContractHashes
	return c
}

func (c *ProcessEventStream) SetDictionarySetEventsBuffer(buffer uint32) *ProcessEventStream {
	c.dictionarySetEventsBuffer = buffer
	return c
}

func (c *ProcessEventStream) SetEventStreamPath(eventPath string) *ProcessEventStream {
	c.eventStreamPath = eventPath
	return c
}

func (c *ProcessEventStream) Execute(ctx context.Context) error {
	//eventListener, err := casper.NewEventListener(c.GetBaseStreamURL(), c.eventStreamPath, &c.nodeStartFromEventID)
	//if err != nil {
	//	return err
	//}

	//daoMetaData := c.GetDAOContractsMetadata()

	//syncDaoSetting := settings.NewSyncDAOSettings()
	//syncDaoSetting.SetCasperClient(c.GetCasperClient())
	//syncDaoSetting.SetVariableRepositoryContractStorageUref(daoMetaData.VariableRepositoryContractStorageUref)
	//syncDaoSetting.SetEntityManager(c.GetEntityManager())
	//syncDaoSetting.SetSettings(settings.VariableRepoSettings)
	//syncDaoSetting.Execute()

	hash, _ := types.NewHashFromHexString("91e93913fdfc4fb2f542c996a9bb337d64d10332b24167466a830846dc1e6410")
	cesParser, err := ces.NewParser(c.GetCasperClient(), []types.Hash{hash})
	if err != nil {
		zap.S().With(zap.Error(err)).Error("Failed to create CES Parser")
		return err
	}

	processRawDeploy := NewProcessRawDeploy()
	processRawDeploy.SetEntityManager(c.GetEntityManager())
	processRawDeploy.SetCESEventParser(cesParser)
	processRawDeploy.SetDAOContractsMetadata(c.GetDAOContractsMetadata())

	res, _ := c.GetCasperClient().GetDeploy("18cd3d13852d6fb75f0eabe09807afb14a9af7edae5a885c07c9b9bce340c7ce")

	processRawDeploy.SetDeployProcessedEvent(&casper.DeployProcessedEvent{
		DeployProcessed: casper.DeployProcessed{
			DeployHash:      hash,
			ExecutionResult: res.ExecutionResults[0].Result,
			Timestamp:       time.Now(),
		},
	})
	if err = processRawDeploy.Execute(); err != nil {
		zap.S().With(zap.Error(err)).Error("Failed to handle DeployProcessedEvent")
	}

	//stopListening := func() {
	//	eventListener.Close()
	//	zap.S().Info("Finish ProcessEvents command successfully")
	//}
	//// in case of blocking on eventListener.ReadEvent(), shutdown will happen on next event/ loop iteration
	//for {
	//	select {
	//	case <-ctx.Done():
	//		stopListening()
	//		return nil
	//	default:
	//		rawEventData, err := eventListener.ReadEvent()
	//		if err != nil {
	//			zap.S().With(zap.Error(err)).Error("Error on event listening")
	//			stopListening()
	//			return err
	//		}
	//
	//		if rawEventData.EventType != casper.DeployProcessedEventType {
	//			zap.S().Info("Skip not supported event type, expect DeployProcessedEvent")
	//			continue
	//		}
	//
	//		deployProcessedEvent, err := rawEventData.Data.ParseAsDeployProcessedEvent()
	//		if err != nil {
	//			zap.S().With(zap.Error(err)).Info("Failed to parse rawEvent as DeployProcessedEvent")
	//			return err
	//		}
	//
	//		if deployProcessedEvent.DeployProcessed.ExecutionResult.Success == nil {
	//			zap.S().With(zap.Error(err)).Info("Failed to parse rawEvent as DeployProcessedEvent")
	//			continue
	//		}
	//
	//		processRawDeploy.SetDeployProcessedEvent(deployProcessedEvent)
	//		if err = processRawDeploy.Execute(); err != nil {
	//			zap.S().With(zap.Error(err)).Error("Failed to handle DeployProcessedEvent")
	//		}
	//	}
	//}
	return nil
}
