package casper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"time"

	"github.com/casper-ecosystem/casper-golang-sdk/keypair"
	"github.com/pkg/errors"
)

type blockIdentifierParam struct {
	BlockIdentifier blockIdentifier `json:"block_identifier"`
}

type blockIdentifier struct {
	Hash   string `json:"Hash,omitempty"`
	Height uint64 `json:"Height,omitempty"`
}

//go:generate mockgen -destination=../../internal/crdao/network-store/tests/mocks/rpc_client_mock.go -package=mocks -source=./rpc_client.go RpcClient
type RPCClient interface {
	GetDeploy(hash string) (GetDeployResult, error)
	GetStateItem(stateRootHash, key string, path []string) (StateGetItemResult, error)
	GetDictionaryItem(stateRootHash, uref, key string) (StateGetItemResult, error)
	GetAccountBalance(stateRootHash, balanceUref string) (big.Int, error)
	GetAccountMainPurseURef(accountHash string) (string, error)
	GetEraInfoBySwitchBlockHeight(height uint64) (GetEraInfoBySwitchBlockResult, error)
	GetAccountBalanceByKeypair(stateRootHash string, key keypair.KeyPair) (big.Int, error)
	GetLatestBlock() (GetBlockResult, error)
	GetBlockByHeight(height uint64) (GetBlockResult, error)
	GetBlockByHash(hash string) (GetBlockResult, error)
	GetLatestBlockTransfers() (GetBlockTransfersResult, error)
	GetBlockTransfersByHeight(height uint64) (GetBlockTransfersResult, error)
	GetBlockTransfersByHash(blockHash string) (GetBlockTransfersResult, error)
	GetAuctionState() (StateGetAuctionInfoResult, error)
	GetLatestStateRootHash() (GetStateRootHashResult, error)
	GetStateRootHashByHeight(height uint64) (GetStateRootHashResult, error)
	GetStateRootHashByHash(stateRootHash string) (GetStateRootHashResult, error)
	GetStatus() (StatusResult, error)
	GetNodeStatus(nodeIP string) (NodeStatus, error)
	GetPeers() (PeerResult, error)
}

type rpcClient struct {
	endpoint   string
	httpClient http.Client
}

func NewRPCClient(endpoint string) RPCClient {
	return &rpcClient{
		endpoint: endpoint,
		httpClient: http.Client{
			Timeout: 20 * time.Second,
		},
	}
}

func (c *rpcClient) GetDeploy(hash string) (GetDeployResult, error) {
	resp, err := c.rpcCall("info_get_deploy", map[string]string{
		"deploy_hash": hash,
	})
	if err != nil {
		return GetDeployResult{}, err
	}

	var result GetDeployResult
	err = json.Unmarshal(resp.Result, &result)

	if err != nil {
		return GetDeployResult{}, fmt.Errorf("failed to get result: %w", err)
	}

	return result, nil
}

func (c *rpcClient) GetStateItem(stateRootHash, key string, path []string) (StateGetItemResult, error) {
	params := map[string]interface{}{
		"state_root_hash": stateRootHash,
		"key":             key,
	}
	if len(path) > 0 {
		params["path"] = path
	}
	resp, err := c.rpcCall("state_get_item", params)
	if err != nil {
		return StateGetItemResult{}, newErrorFromRPCError(err)
	}

	var result StateGetItemResult
	err = json.Unmarshal(resp.Result, &result)
	if err != nil {
		return StateGetItemResult{}, fmt.Errorf("failed to get result: %w", err)
	}

	return result, nil
}

func (c *rpcClient) GetDictionaryItem(stateRootHash, uref, key string) (StateGetItemResult, error) {
	params := map[string]interface{}{
		"state_root_hash": stateRootHash,
		"dictionary_identifier": map[string]interface{}{
			"URef": map[string]string{
				"dictionary_item_key": key,
				"seed_uref":           uref,
			},
		},
	}

	resp, err := c.rpcCall("state_get_dictionary_item", params)
	if err != nil {
		return StateGetItemResult{}, newErrorFromRPCError(err)
	}

	var result StateGetItemResult
	err = json.Unmarshal(resp.Result, &result)
	if err != nil {
		return StateGetItemResult{}, fmt.Errorf("failed to get result: %w", err)
	}

	return result, nil
}

func (c *rpcClient) GetAccountBalance(stateRootHash, balanceUref string) (big.Int, error) {
	resp, err := c.rpcCall("state_get_balance", map[string]string{
		"state_root_hash": stateRootHash,
		"purse_uref":      balanceUref,
	})
	if err != nil {
		return big.Int{}, err
	}

	var result StateGetBalanceResult
	err = json.Unmarshal(resp.Result, &result)
	if err != nil {
		return big.Int{}, fmt.Errorf("failed to get result: %w", err)
	}

	balance := big.Int{}
	balance.SetString(result.BalanceValue, 10)
	return balance, nil
}

func (c *rpcClient) GetAccountMainPurseURef(accountHash string) (string, error) {
	latestBlockResult, err := c.GetLatestBlock()
	if err != nil {
		return "", err
	}

	stateItemResult, err := c.GetStateItem(latestBlockResult.Block.Header.StateRootHash.ToHex(), accountHash, []string{})
	if err != nil {
		return "", err
	}

	return stateItemResult.StoredValue.Account.MainPurse, nil
}

func (c *rpcClient) GetEraInfoBySwitchBlockHeight(height uint64) (GetEraInfoBySwitchBlockResult, error) {
	resp, err := c.rpcCall("chain_get_era_info_by_switch_block",
		blockIdentifierParam{blockIdentifier{
			Height: height,
		}})
	if err != nil {
		return GetEraInfoBySwitchBlockResult{}, err
	}
	var result GetEraInfoBySwitchBlockResult

	err = json.Unmarshal(resp.Result, &result)
	if err != nil {
		return GetEraInfoBySwitchBlockResult{}, fmt.Errorf("failed to get result: %w", err)
	}

	return result, nil
}

func (c *rpcClient) GetAccountBalanceByKeypair(stateRootHash string, key keypair.KeyPair) (big.Int, error) {
	stateItemResult, err := c.GetStateItem(stateRootHash, key.AccountHash(), []string{})
	if err != nil {
		return big.Int{}, err
	}
	return c.GetAccountBalance(stateRootHash, stateItemResult.StoredValue.Account.MainPurse)
}

func (c *rpcClient) GetLatestBlock() (GetBlockResult, error) {
	resp, err := c.rpcCall("chain_get_block", nil)
	if err != nil {
		return GetBlockResult{}, err
	}

	var result GetBlockResult
	err = json.Unmarshal(resp.Result, &result)
	if err != nil {
		return GetBlockResult{}, fmt.Errorf("failed to get result: %w", err)
	}

	return result, nil
}

func (c *rpcClient) GetBlockByHeight(height uint64) (GetBlockResult, error) {
	resp, err := c.rpcCall("chain_get_block",
		blockIdentifierParam{blockIdentifier{
			Height: height,
		}})
	if err != nil {
		return GetBlockResult{}, err
	}

	var result GetBlockResult
	err = json.Unmarshal(resp.Result, &result)
	if err != nil {
		return GetBlockResult{}, fmt.Errorf("failed to get result: %w", err)
	}

	return result, nil
}

func (c *rpcClient) GetBlockByHash(hash string) (GetBlockResult, error) {
	resp, err := c.rpcCall("chain_get_block",
		blockIdentifierParam{blockIdentifier{
			Hash: hash,
		}})
	if err != nil {
		return GetBlockResult{}, err
	}

	var result GetBlockResult
	err = json.Unmarshal(resp.Result, &result)
	if err != nil {
		return GetBlockResult{}, fmt.Errorf("failed to get result: %w", err)
	}

	return result, nil
}

func (c *rpcClient) GetLatestBlockTransfers() (GetBlockTransfersResult, error) {
	resp, err := c.rpcCall("chain_get_block_transfers", nil)
	if err != nil {
		return GetBlockTransfersResult{}, err
	}

	var result GetBlockTransfersResult
	err = json.Unmarshal(resp.Result, &result)
	if err != nil {
		return GetBlockTransfersResult{}, fmt.Errorf("failed to get result: %w", err)
	}

	return result, nil
}

func (c *rpcClient) GetBlockTransfersByHeight(height uint64) (GetBlockTransfersResult, error) {
	resp, err := c.rpcCall("chain_get_block_transfers",
		blockIdentifierParam{blockIdentifier{
			Height: height,
		}})
	if err != nil {
		return GetBlockTransfersResult{}, err
	}

	var result GetBlockTransfersResult
	err = json.Unmarshal(resp.Result, &result)
	if err != nil {
		return GetBlockTransfersResult{}, fmt.Errorf("failed to get result: %w", err)
	}

	return result, nil
}

func (c *rpcClient) GetBlockTransfersByHash(blockHash string) (GetBlockTransfersResult, error) {
	resp, err := c.rpcCall("chain_get_block_transfers",
		blockIdentifierParam{blockIdentifier{
			Hash: blockHash,
		}})
	if err != nil {
		return GetBlockTransfersResult{}, err
	}

	var result GetBlockTransfersResult
	err = json.Unmarshal(resp.Result, &result)
	if err != nil {
		return GetBlockTransfersResult{}, fmt.Errorf("failed to get result: %w", err)
	}

	return result, nil
}

func (c *rpcClient) GetAuctionState() (StateGetAuctionInfoResult, error) {
	resp, err := c.rpcCall("state_get_auction_info", nil)
	if err != nil {
		return StateGetAuctionInfoResult{}, err
	}

	var result StateGetAuctionInfoResult
	err = json.Unmarshal(resp.Result, &result)
	if err != nil {
		return StateGetAuctionInfoResult{}, fmt.Errorf("failed to get result: #{err}")
	}

	return result, nil
}

func (c *rpcClient) GetStatus() (StatusResult, error) {
	resp, err := c.rpcCall("info_get_status", nil)
	if err != nil {
		return StatusResult{}, err
	}

	var result StatusResult
	err = json.Unmarshal(resp.Result, &result)
	if err != nil {
		return StatusResult{}, fmt.Errorf("failed to get result: #{err}")
	}

	return result, nil
}

func (c *rpcClient) GetPeers() (PeerResult, error) {
	resp, err := c.rpcCall("info_get_peers", nil)
	if err != nil {
		return PeerResult{}, err
	}

	var result PeerResult
	err = json.Unmarshal(resp.Result, &result)
	if err != nil {
		return PeerResult{}, fmt.Errorf("failed to get result: #{err}")
	}

	return result, nil
}

func (c *rpcClient) GetLatestStateRootHash() (GetStateRootHashResult, error) {
	resp, err := c.rpcCall("chain_get_state_root_hash", nil)
	if err != nil {
		return GetStateRootHashResult{}, err
	}

	var result GetStateRootHashResult
	err = json.Unmarshal(resp.Result, &result)
	if err != nil {
		return GetStateRootHashResult{}, fmt.Errorf("failed to get result: %w", err)
	}

	return result, nil
}

func (c *rpcClient) GetStateRootHashByHeight(height uint64) (GetStateRootHashResult, error) {
	resp, err := c.rpcCall("chain_get_state_root_hash",
		blockIdentifierParam{blockIdentifier{
			Height: height,
		}})
	if err != nil {
		return GetStateRootHashResult{}, err
	}

	var result GetStateRootHashResult
	err = json.Unmarshal(resp.Result, &result)
	if err != nil {
		return GetStateRootHashResult{}, fmt.Errorf("failed to get result: %w", err)
	}

	return result, nil
}

func (c *rpcClient) GetStateRootHashByHash(hash string) (GetStateRootHashResult, error) {
	resp, err := c.rpcCall("chain_get_state_root_hash",
		blockIdentifierParam{blockIdentifier{
			Hash: hash,
		}})
	if err != nil {
		return GetStateRootHashResult{}, err
	}

	var result GetStateRootHashResult
	err = json.Unmarshal(resp.Result, &result)
	if err != nil {
		return GetStateRootHashResult{}, fmt.Errorf("failed to get result: %w", err)
	}

	return result, nil
}

func (c *rpcClient) rpcCall(method string, params interface{}) (RPCResponse, error) {
	body, err := json.Marshal(RPCRequest{
		Version: "2.0",
		Method:  method,
		Params:  params,
	})

	if err != nil {
		return RPCResponse{}, errors.Wrap(err, "failed to marshal json")
	}

	resp, err := c.httpClient.Post(c.endpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		return RPCResponse{}, fmt.Errorf("failed to make request: %w", err)
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return RPCResponse{}, fmt.Errorf("failed to get response body: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return RPCResponse{}, fmt.Errorf("request failed, status code - %d, response - %s", resp.StatusCode, string(b))
	}

	var rpcResponse RPCResponse
	err = json.Unmarshal(b, &rpcResponse)
	if err != nil {
		return RPCResponse{}, fmt.Errorf("failed to parse response body: %w", err)
	}

	if rpcResponse.Error != nil {
		return rpcResponse, fmt.Errorf("rpc call failed, code - %d, message - %s", rpcResponse.Error.Code, rpcResponse.Error.Message)
	}

	return rpcResponse, nil
}

func (c *rpcClient) GetNodeStatus(nodeIP string) (NodeStatus, error) {
	status := NodeStatus{}

	statusURL := fmt.Sprintf("http://%s:8888/status", nodeIP)
	resp, err := c.httpClient.Get(statusURL)
	if err != nil {
		return status, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return status, err
	}

	err = json.Unmarshal(bodyBytes, &status)
	if err != nil {
		return status, err
	}

	return status, nil
}
