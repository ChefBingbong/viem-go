package client

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Client is a JSON-RPC client for Ethereum nodes.
type Client struct {
	transport Transport
	chainID   *big.Int
}

// NewClient creates a new client connected to the given RPC URL.
func NewClient(rpcURL string) (*Client, error) {
	transport := NewHTTPTransport(rpcURL)
	return NewClientWithTransport(transport)
}

// NewClientWithTransport creates a new client with a custom transport.
func NewClientWithTransport(transport Transport) (*Client, error) {
	client := &Client{
		transport: transport,
	}

	// Fetch chain ID
	ctx := context.Background()
	chainID, err := client.GetChainID(ctx)
	if err != nil {
		// Don't fail if chain ID fetch fails - some nodes might not support it
		// Client can still work without it
		client.chainID = nil
	} else {
		client.chainID = chainID
	}

	return client, nil
}

// Transport returns the underlying transport.
func (c *Client) Transport() Transport {
	return c.transport
}

// ChainID returns the chain ID, or nil if not available.
func (c *Client) ChainID() *big.Int {
	return c.chainID
}

// Close closes the client connection.
func (c *Client) Close() error {
	return c.transport.Close()
}

// Call performs an eth_call RPC request.
func (c *Client) Call(ctx context.Context, call CallRequest, block ...BlockNumber) ([]byte, error) {
	blockParam := getBlockParam(block)

	result, err := c.transport.Call(ctx, "eth_call", call, blockParam)
	if err != nil {
		return nil, err
	}

	var hexResult string
	if err := json.Unmarshal(result, &hexResult); err != nil {
		return nil, fmt.Errorf("failed to unmarshal call result: %w", err)
	}

	return hexutil.Decode(hexResult)
}

// SendTransaction sends a transaction and returns the transaction hash.
func (c *Client) SendTransaction(ctx context.Context, tx Transaction) (common.Hash, error) {
	result, err := c.transport.Call(ctx, "eth_sendTransaction", tx)
	if err != nil {
		return common.Hash{}, err
	}

	var hashHex string
	if err := json.Unmarshal(result, &hashHex); err != nil {
		return common.Hash{}, fmt.Errorf("failed to unmarshal transaction hash: %w", err)
	}

	return common.HexToHash(hashHex), nil
}

// SendRawTransaction sends a signed raw transaction.
func (c *Client) SendRawTransaction(ctx context.Context, rawTx []byte) (common.Hash, error) {
	result, err := c.transport.Call(ctx, "eth_sendRawTransaction", hexutil.Encode(rawTx))
	if err != nil {
		return common.Hash{}, err
	}

	var hashHex string
	if err := json.Unmarshal(result, &hashHex); err != nil {
		return common.Hash{}, fmt.Errorf("failed to unmarshal transaction hash: %w", err)
	}

	return common.HexToHash(hashHex), nil
}

// GetBalance returns the balance of an address at the given block.
func (c *Client) GetBalance(ctx context.Context, addr common.Address, block ...BlockNumber) (*big.Int, error) {
	blockParam := getBlockParam(block)

	result, err := c.transport.Call(ctx, "eth_getBalance", addr.Hex(), blockParam)
	if err != nil {
		return nil, err
	}

	var hexBalance string
	if err := json.Unmarshal(result, &hexBalance); err != nil {
		return nil, fmt.Errorf("failed to unmarshal balance: %w", err)
	}

	return hexutil.DecodeBig(hexBalance)
}

// GetBlockNumber returns the current block number.
func (c *Client) GetBlockNumber(ctx context.Context) (uint64, error) {
	result, err := c.transport.Call(ctx, "eth_blockNumber")
	if err != nil {
		return 0, err
	}

	var hexNumber string
	if err := json.Unmarshal(result, &hexNumber); err != nil {
		return 0, fmt.Errorf("failed to unmarshal block number: %w", err)
	}

	return hexutil.DecodeUint64(hexNumber)
}

// GetChainID returns the chain ID.
func (c *Client) GetChainID(ctx context.Context) (*big.Int, error) {
	result, err := c.transport.Call(ctx, "eth_chainId")
	if err != nil {
		return nil, err
	}

	var hexChainID string
	if err := json.Unmarshal(result, &hexChainID); err != nil {
		return nil, fmt.Errorf("failed to unmarshal chain ID: %w", err)
	}

	return hexutil.DecodeBig(hexChainID)
}

// GetTransactionReceipt returns the receipt of a transaction.
func (c *Client) GetTransactionReceipt(ctx context.Context, hash common.Hash) (*Receipt, error) {
	result, err := c.transport.Call(ctx, "eth_getTransactionReceipt", hash.Hex())
	if err != nil {
		return nil, err
	}

	// Check for null result (transaction not found or not yet mined)
	if string(result) == "null" {
		return nil, nil
	}

	var receipt Receipt
	if err := json.Unmarshal(result, &receipt); err != nil {
		return nil, fmt.Errorf("failed to unmarshal receipt: %w", err)
	}

	return &receipt, nil
}

// GetTransactionCount returns the nonce for an address.
func (c *Client) GetTransactionCount(ctx context.Context, addr common.Address, block ...BlockNumber) (uint64, error) {
	blockParam := getBlockParam(block)

	result, err := c.transport.Call(ctx, "eth_getTransactionCount", addr.Hex(), blockParam)
	if err != nil {
		return 0, err
	}

	var hexNonce string
	if err := json.Unmarshal(result, &hexNonce); err != nil {
		return 0, fmt.Errorf("failed to unmarshal nonce: %w", err)
	}

	return hexutil.DecodeUint64(hexNonce)
}

// GetGasPrice returns the current gas price.
func (c *Client) GetGasPrice(ctx context.Context) (*big.Int, error) {
	result, err := c.transport.Call(ctx, "eth_gasPrice")
	if err != nil {
		return nil, err
	}

	var hexGasPrice string
	if err := json.Unmarshal(result, &hexGasPrice); err != nil {
		return nil, fmt.Errorf("failed to unmarshal gas price: %w", err)
	}

	return hexutil.DecodeBig(hexGasPrice)
}

// EstimateGas estimates the gas required for a call.
func (c *Client) EstimateGas(ctx context.Context, call CallRequest) (uint64, error) {
	result, err := c.transport.Call(ctx, "eth_estimateGas", call)
	if err != nil {
		return 0, err
	}

	var hexGas string
	if err := json.Unmarshal(result, &hexGas); err != nil {
		return 0, fmt.Errorf("failed to unmarshal gas estimate: %w", err)
	}

	return hexutil.DecodeUint64(hexGas)
}

// GetCode returns the code at an address.
func (c *Client) GetCode(ctx context.Context, addr common.Address, block ...BlockNumber) ([]byte, error) {
	blockParam := getBlockParam(block)

	result, err := c.transport.Call(ctx, "eth_getCode", addr.Hex(), blockParam)
	if err != nil {
		return nil, err
	}

	var hexCode string
	if err := json.Unmarshal(result, &hexCode); err != nil {
		return nil, fmt.Errorf("failed to unmarshal code: %w", err)
	}

	return hexutil.Decode(hexCode)
}

// GetStorageAt returns the value from a storage position.
func (c *Client) GetStorageAt(ctx context.Context, addr common.Address, position common.Hash, block ...BlockNumber) ([]byte, error) {
	blockParam := getBlockParam(block)

	result, err := c.transport.Call(ctx, "eth_getStorageAt", addr.Hex(), position.Hex(), blockParam)
	if err != nil {
		return nil, err
	}

	var hexValue string
	if err := json.Unmarshal(result, &hexValue); err != nil {
		return nil, fmt.Errorf("failed to unmarshal storage value: %w", err)
	}

	return hexutil.Decode(hexValue)
}

// GetLogs returns logs matching the filter query.
func (c *Client) GetLogs(ctx context.Context, query FilterQuery) ([]Log, error) {
	result, err := c.transport.Call(ctx, "eth_getLogs", query)
	if err != nil {
		return nil, err
	}

	var logs []Log
	if err := json.Unmarshal(result, &logs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal logs: %w", err)
	}

	return logs, nil
}

// GetBlockByNumber returns a block by number.
func (c *Client) GetBlockByNumber(ctx context.Context, block BlockNumber, fullTx bool) (*Block, error) {
	var blockParam string
	if block == nil {
		blockParam = BlockLatest.String()
	} else {
		blockParam = block.String()
	}

	result, err := c.transport.Call(ctx, "eth_getBlockByNumber", blockParam, fullTx)
	if err != nil {
		return nil, err
	}

	if string(result) == "null" {
		return nil, nil
	}

	var b Block
	if err := json.Unmarshal(result, &b); err != nil {
		return nil, fmt.Errorf("failed to unmarshal block: %w", err)
	}

	return &b, nil
}

// GetBlockByHash returns a block by hash.
func (c *Client) GetBlockByHash(ctx context.Context, hash common.Hash, fullTx bool) (*Block, error) {
	result, err := c.transport.Call(ctx, "eth_getBlockByHash", hash.Hex(), fullTx)
	if err != nil {
		return nil, err
	}

	if string(result) == "null" {
		return nil, nil
	}

	var b Block
	if err := json.Unmarshal(result, &b); err != nil {
		return nil, fmt.Errorf("failed to unmarshal block: %w", err)
	}

	return &b, nil
}

// WaitForTransaction waits for a transaction to be mined and returns the receipt.
// It polls the node until the transaction is found or the context is cancelled.
func (c *Client) WaitForTransaction(ctx context.Context, hash common.Hash) (*Receipt, error) {
	for {
		receipt, err := c.GetTransactionReceipt(ctx, hash)
		if err != nil {
			return nil, err
		}
		if receipt != nil {
			return receipt, nil
		}

		// Check if context is done
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			// Continue polling
		}
	}
}

// getBlockParam returns the block parameter string.
func getBlockParam(block []BlockNumber) string {
	if len(block) == 0 || block[0] == nil {
		return BlockLatest.String()
	}
	return block[0].String()
}

// RawCall performs a raw JSON-RPC call.
func (c *Client) RawCall(ctx context.Context, method string, params ...any) (json.RawMessage, error) {
	return c.transport.Call(ctx, method, params...)
}
