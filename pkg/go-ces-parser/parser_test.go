package ces

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"casper-dao-middleware/pkg/casper"
	"casper-dao-middleware/pkg/casper/mocks"
	"casper-dao-middleware/pkg/casper/types"
)

var schemaHex = `0800000008000000417070726f76616c03000000050000006f776e65720b080000006f70657261746f720b08000000746f6b656e5f6964150e000000417070726f76616c466f72416c6c03000000050000006f776e65720b080000006f70657261746f720d0b09000000746f6b656e5f6964730e15040000004275726e02000000050000006f776e65720b08000000746f6b656e5f6964150f0000004d65746164617461557064617465640200000008000000746f6b656e5f69641504000000646174610a090000004d6967726174696f6e00000000040000004d696e740200000009000000726563697069656e740b08000000746f6b656e5f696415080000005472616e7366657204000000050000006f776e65720b080000006f70657261746f720d0b09000000726563697069656e740b08000000746f6b656e5f6964150c0000005661726961626c657353657400000000`

func TestEventParser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockedClient := mocks.NewMockRPCClient(mockCtrl)

	contractHashToParse, err := types.NewHashFromHexString("ea0c001d969da098fefec42b141db88c74c5682e49333ded78035540a0b4f0bc")
	assert.NoError(t, err)

	contractPackageHash, err := types.NewHashFromHexString("7a5fce1d9ad45c9d71a5e59638602213295a51a6cf92518f8b262cd3e23d6d7e")
	assert.NoError(t, err)

	eventParser := EventParser{
		casperClient: mockedClient,
	}

	mockedClient.EXPECT().GetLatestStateRootHash().Return(casper.GetStateRootHashResult{}, nil)
	mockedClient.EXPECT().GetStateItem("", fmt.Sprintf("hash-%s", contractHashToParse.ToHex()), nil).Return(casper.StateGetItemResult{
		StoredValue: casper.StoredValue{
			Contract: &casper.Contract{
				ContractPackageHash: contractPackageHash,
				NamedKeys: []casper.NamedKey{
					{
						Name: eventNamedKey,
						Key:  "uref-70d95cbeae8ce00c0ca493762cc99aed052adfcb3e279c7440f5241b1bdf27a1-007",
					}, {
						Name: eventSchemaNamedKey,
						Key:  "events-Schema-named-key-uref",
					}},
			},
		},
	}, nil)

	schemaBytes, err := hex.DecodeString(schemaHex)
	assert.NoError(t, err)

	mockedClient.EXPECT().GetStateItem("events-Schema-named-key-uref", "", nil).Return(
		casper.StateGetItemResult{
			StoredValue: casper.StoredValue{
				CLValue: &casper.CLValue{
					Bytes: schemaBytes,
				},
			},
		}, nil)

	contractInfos, err := eventParser.parseContractsInfos([]types.Hash{contractHashToParse})
	assert.NoError(t, err)

	eventParser.contractInfos = contractInfos

	t.Run("Test Mint Event", func(t *testing.T) {
		var res casper.GetDeployResult

		data, err := os.ReadFile("./fixtures/deploys/mint.json")
		assert.NoError(t, err)

		err = json.Unmarshal(data, &res)
		assert.NoError(t, err)

		parseResults, err := eventParser.ParseExecutionResults(res.ExecutionResults[0].Result)
		assert.NoError(t, err)
		assert.True(t, len(parseResults) == 1)

		assert.Equal(t, parseResults[0].Event.Name, "Mint")
		assert.Equal(t, parseResults[0].Event.ContractHash.String(), contractHashToParse.String())
		assert.Equal(t, parseResults[0].Event.ContractPackageHash.String(), contractPackageHash.String())
		assert.True(t, len(parseResults[0].Event.Data) > 0)
	})

	t.Run("Test Transfer Event", func(t *testing.T) {
		var res casper.GetDeployResult

		data, err := os.ReadFile("./fixtures/deploys/transfer.json")
		assert.NoError(t, err)

		err = json.Unmarshal(data, &res)
		assert.NoError(t, err)

		parseResults, err := eventParser.ParseExecutionResults(res.ExecutionResults[0].Result)
		assert.NoError(t, err)
		assert.True(t, len(parseResults) == 1)
		assert.NoError(t, parseResults[0].Error)

		assert.Equal(t, parseResults[0].Event.Name, "Transfer")
		assert.Equal(t, parseResults[0].Event.ContractHash.String(), contractHashToParse.String())
		assert.Equal(t, parseResults[0].Event.ContractPackageHash.String(), contractPackageHash.String())
		assert.True(t, len(parseResults[0].Event.Data) > 0)
	})
}
